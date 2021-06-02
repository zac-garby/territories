package world

import (
	"fmt"
	"math"
	"math/rand"

	perlin "github.com/aquilax/go-perlin"
)

const (
	alpha     float64 = 2
	beta      float64 = 2
	n         int     = 3
	w         float64 = 6
	nint      int     = 4
	steepness float64 = 32
)

type Coord struct {
	X, Y float64
}

type WorldGen struct {
	Noise         *perlin.Perlin
	Perlin        [][]float64
	Rand          *rand.Rand
	Midpoints     []Coord
	Pixels        [][]int
	Width, Height int
	Seed          int64
	NumRegions    int
}

func NewGen(width, height, npoints int, seed int64) *WorldGen {
	g := &WorldGen{Width: width, Height: height, Seed: seed}
	g.Noise = perlin.NewPerlin(alpha, beta, n, seed)
	g.Rand = rand.New(rand.NewSource(seed))

	g.precomputePerlin()
	g.initPoints(npoints)
	g.voronoi()
	g.disjoinRegions()
	g.NumRegions = g.renumber()
	fmt.Printf("there are %d regions\n", g.NumRegions)

	for g.removeTiny(3000) {
		g.NumRegions = g.renumber()
	}

	return g
}

func (g *WorldGen) initPoints(n int) {
	g.Midpoints = make([]Coord, n)

	for i := 0; i < n; i++ {
		g.Midpoints[i].X = g.Rand.Float64() * float64(g.Width)
		g.Midpoints[i].Y = g.Rand.Float64() * float64(g.Height)
	}
}

func (g *WorldGen) precomputePerlin() {
	width := float64(g.Width)
	height := float64(g.Height)

	g.Perlin = make([][]float64, g.Height+1)

	for y := 0; y <= g.Height; y++ {
		row := make([]float64, g.Width+1)
		for x := 0; x <= g.Width; x++ {
			row[x] = g.Noise.Noise2D(w*(float64(x)/width), w*(float64(y)/height))
		}
		g.Perlin[y] = row
	}
}

func (g *WorldGen) Interperlin(x, y float64) float64 {
	var (
		ix_, fx = math.Modf(x)
		iy_, fy = math.Modf(y)
		ix, iy  = int(ix_), int(iy_)
		l       = (1-fy)*g.Perlin[iy][ix] + fy*g.Perlin[iy+1][ix]
		r       = (1-fy)*g.Perlin[iy][ix+1] + fy*g.Perlin[iy+1][ix+1]
	)

	return (1-fx)*l + fx*r
}

func (g *WorldGen) voronoi() {
	g.Pixels = make([][]int, g.Height)
	for y := 0; y < g.Height; y++ {
		row := make([]int, g.Width)
		for x := 0; x < g.Width; x++ {
			row[x] = g.closestPoint(Coord{float64(x), float64(y)})
		}
		g.Pixels[y] = row
	}
}

func (g *WorldGen) closestPoint(c Coord) int {
	var (
		cd = float64(g.Width * g.Height)
		ci = -1
	)

	for i, p := range g.Midpoints {
		dist := g.Distance(c, p)
		if dist < cd {
			cd = dist
			ci = i
		}
	}

	return ci
}

func (g *WorldGen) Distance(a, b Coord) float64 {
	var (
		//mx    = (a.X + b.X) / 2
		//my    = (a.Y + b.Y) / 2
		dx     = b.X - a.X
		dy     = b.Y - a.Y
		length = math.Sqrt(dx*dx + dy*dy)
		ndx    = dx / length
		ndy    = dy / length
		noise  = 0.0
	)

	delta := length / float64(nint)

	for i := 0; i < nint; i++ {
		x := a.X + float64(i)*delta*ndx
		y := a.Y + float64(i)*delta*ndy
		n := delta * math.Abs(g.Interperlin(x, y))
		noise += math.Sqrt(delta*ndx*delta*ndx + delta*ndy*delta*ndy + steepness*n*n)
	}

	return noise
}

func (g *WorldGen) disjoinRegions() {
	var (
		res      = make([][]int, g.Height)
		i        = 1
		mappings = make(map[int]int, 0)
		equiv    = make(map[int]int, 0)
	)

	for y := 0; y < g.Height; y++ {
		row := make([]int, g.Width)
		res[y] = row

		for x := 0; x < g.Width; x++ {
			var (
				val       = g.Pixels[y][x]
				top, left = -1, -1
			)

			if x > 0 {
				if m, ok := mappings[res[y][x-1]]; ok && m == val {
					left = res[y][x-1]
				}
			}

			if y > 0 {
				if m, ok := mappings[res[y-1][x]]; ok && m == val {
					top = res[y-1][x]
				}
			}

			if left == -1 && top == -1 {
				mappings[i] = val
				res[y][x] = i
				i++
			} else if left > -1 && top > -1 && top != left {
				min, max := left, top
				if top < left {
					min, max = top, left
				}

				if _, exists := equiv[max]; !exists { // i think something could be done in the else here
					equiv[min] = max
				}

				res[y][x] = max
			} else if left == -1 {
				res[y][x] = top
			} else {
				res[y][x] = left
			}
		}
	}

	for y := 0; y < g.Height; y++ {
		for x := 0; x < g.Width; x++ {
			val := res[y][x]
			for repl, ok := equiv[val]; ok; repl, ok = equiv[repl] {
				res[y][x] = repl
			}
		}
	}

	g.Pixels = res
}

func (g *WorldGen) measureRegions() map[int]int {
	sizes := make(map[int]int, 32)

	for x := 0; x < g.Width; x++ {
		for y := 0; y < g.Height; y++ {
			val := g.Pixels[y][x]

			if _, ok := sizes[val]; ok {
				sizes[val]++
			} else {
				sizes[val] = 1
			}
		}
	}

	return sizes
}

func (g *WorldGen) renumber() int {
	var (
		mappings = make(map[int]int, 0)
		i        = 0
	)

	for y := 0; y < g.Height; y++ {
		for x := 0; x < g.Width; x++ {
			val := g.Pixels[y][x]

			if m, ok := mappings[val]; ok {
				g.Pixels[y][x] = m
			} else {
				mappings[val] = i
				g.Pixels[y][x] = i
				i++
			}
		}
	}

	return i
}

func (g *WorldGen) removeTiny(threshold int) bool {
	var (
		sizes   = g.measureRegions()
		adj     = g.getAdjacencyMatrix()
		repl    = make([]int, g.NumRegions)
		changed = false
	)

	for from, row := range adj {
		if sizes[from] < threshold {
			fmt.Printf("%d is too small (%d)\n", from, sizes[from])
			best := 0

			for to, a := range row {
				if to != from && a != 0 && a > best {
					best = a
					repl[from] = to
				}
			}

			changed = true
		} else {
			repl[from] = from
		}
	}

	for y := 0; y < g.Height; y++ {
		for x := 0; x < g.Width; x++ {
			g.Pixels[y][x] = repl[g.Pixels[y][x]]
		}
	}

	return changed
}

func (g *WorldGen) getAdjacencyMatrix() [][]int {
	adj := make([][]int, g.NumRegions)

	for i := 0; i < g.NumRegions; i++ {
		adj[i] = make([]int, g.NumRegions)
	}

	for y := 0; y < g.Height; y++ {
		for x := 0; x < g.Width; x++ {
			val := g.Pixels[y][x]

			if x > 0 {
				adj[val][g.Pixels[y][x-1]]++
			}

			if x+1 < g.Width {
				adj[val][g.Pixels[y][x+1]]++
			}

			if y > 0 {
				adj[val][g.Pixels[y-1][x]]++
			}

			if y+1 < g.Height {
				adj[val][g.Pixels[y+1][x]]++
			}
		}
	}

	return adj
}
