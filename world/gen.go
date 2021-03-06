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

type WorldGen struct {
	Noise         *perlin.Perlin
	Perlin        [][]float64
	Rand          *rand.Rand
	Midpoints     []Coord
	Starts        []Coord
	Pixels        [][]int
	Width, Height int
	Seed          int64
	NumRegions    int
	Polygons      [][]Coord
	Adjacency     [][]bool
	Centroids     []Coord
}

func NewGen(width, height, npoints int, seed int64) *WorldGen {
	g := &WorldGen{Width: width, Height: height, Seed: seed}
	g.Noise = perlin.NewPerlin(alpha, beta, n, seed)
	g.Rand = rand.New(rand.NewSource(seed))

	fmt.Println("generating a world:")
	fmt.Println(" 1. precomputing perlin noise")
	g.precomputePerlin()
	fmt.Println(" 2. initialising points")
	g.initPoints(npoints)
	fmt.Println(" 3. applying voronoi")
	g.voronoi()
	fmt.Println(" 4. finding disjoin regions")
	g.disjoinRegions()
	fmt.Println(" 5. renumbering regions")
	g.NumRegions = g.renumber()

	fmt.Println(" 6. removing tiny regions")
	for g.removeTiny(3000) {
		g.NumRegions = g.renumber()
	}

	fmt.Println(" 7. finding region start points")
	g.Starts = make([]Coord, g.NumRegions)
	for k, v := range g.getRegionStarts() {
		g.Starts[k] = v
	}

	fmt.Println(" 8. finding region polygon vertices")
	g.Polygons = make([][]Coord, g.NumRegions)
	for n := 0; n < g.NumRegions; n++ {
		g.Polygons[n] = g.Polygon(n)
	}

	fmt.Println(" 9. constructing adjacency matrix")
	adj := g.getAdjacencyMatrix() // TODO: this can be called only once
	g.Adjacency = make([][]bool, g.NumRegions)
	for n := 0; n < g.NumRegions; n++ {
		g.Adjacency[n] = make([]bool, g.NumRegions)
		for m := 0; m < g.NumRegions; m++ {
			g.Adjacency[n][m] = adj[n][m] > 0
		}
	}

	fmt.Println(" 10. simplifying region polygons")
	g.Polygons = g.reduceVertices(1.5)
	//g.Polygons = g.reduceVertices(2)

	fmt.Println(" 11. finding centroids")
	g.Centroids = g.getCentroids()

	return g
}

func (g *WorldGen) World() *World {
	w := &World{
		Width:     g.Width,
		Height:    g.Height,
		Regions:   g.Polygons,
		Adjacency: g.Adjacency,
		Centroids: g.Centroids,
	}

	return w
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

func (g *WorldGen) getRegionStarts() map[int]Coord {
	coords := make(map[int]Coord, g.NumRegions)

	for y := 0; y < g.Height; y++ {
		for x := 0; x < g.Width; x++ {
			val := g.Pixels[y][x]
			if _, ok := coords[val]; !ok {
				coords[val] = Coord{X: float64(x), Y: float64(y)}
			}
		}
	}

	return coords
}

func (g *WorldGen) Polygon(r int) []Coord {
	var (
		startC = g.Starts[r]
		px, py = int(startC.X), int(startC.Y)
		sx, sy = px, py
		dx     = 1
		dy     = 0
		coords = make([]Coord, 0, 1000)
	)

	for {
		coords = append(coords, Coord{X: float64(px), Y: float64(py)})
		px += dx
		py += dy

		if px == sx && py == sy {
			break
		}

		if g.hasEdge(px, py, r, dx, dy) {
			dx = dx
			dy = dy
		} else if g.hasEdge(px, py, r, dy, -dx) {
			dx, dy = dy, -dx
		} else if g.hasEdge(px, py, r, -dy, dx) {
			dx, dy = -dy, dx
		} else {
			panic("no edges to follow")
		}
	}

	return coords
}

func (g *WorldGen) hasEdge(px, py int, region int, dx, dy int) bool {
	if dy == 0 {
		if dx == 1 {
			dx = 0
		}
		if px+dx < 0 || px+dx >= g.Width {
			return false
		}
		outAbove := py-1 < 0
		outBelow := py >= g.Height

		if outAbove && !outBelow {
			return g.Pixels[py][px+dx] == region
		} else if outBelow && !outAbove {
			return g.Pixels[py-1][px+dx] == region
		}

		above := g.Pixels[py-1][px+dx]
		below := g.Pixels[py][px+dx]

		return (above == region && below != region) || (above != region && below == region)
	} else if dx == 0 {
		if dy == 1 {
			dy = 0
		}
		if py+dy < 0 || py+dy >= g.Height {
			return false
		}
		outLeft := px-1 < 0
		outRight := px >= g.Width

		if outLeft && !outRight {
			return g.Pixels[py+dy][px] == region
		} else if outRight && !outLeft {
			return g.Pixels[py+dy][px-1] == region
		}

		left := g.Pixels[py+dy][px-1]
		right := g.Pixels[py+dy][px]

		return (left == region && right != region) || (left != region && right == region)
	}

	panic("this shouldn't happen")
}

func (g *WorldGen) getCentroids() []Coord {
	cs := make([]Coord, g.NumRegions)
	ns := make([]int, g.NumRegions)

	// calculate centroids as a moving average
	for y := 0; y < g.Height; y++ {
		for x := 0; x < g.Width; x++ {
			r := g.Pixels[y][x]

			if ns[r] == 0 {
				cs[r] = Coord{}
			} else {
				n := float64(ns[r])
				cs[r].X = (n*cs[r].X + float64(x)) / (n + 1)
				cs[r].Y = (n*cs[r].Y + float64(y)) / (n + 1)
			}

			ns[r] += 1
		}
	}

	return cs
}

func (g *WorldGen) reduceVertices(errThresh float64) [][]Coord {
	var (
		polys     = make([][]Coord, g.NumRegions)
		threshSqr = errThresh * errThresh
		w         = float64(g.Width)
		h         = float64(g.Height)
	)

	for i, poly := range g.Polygons {
		polys[i] = make([]Coord, 0, len(poly)/25)

		v := 0 // v is the start vertex
		for v < len(poly) {
			pv := poly[v]

			// add the current vertex to the new polygon
			polys[i] = append(polys[i], pv)

			u := v + 1 // u is the potential end point of the new line

			var maxErr float64

			for u < len(poly) {
				pu := poly[u]
				maxErr = 0.0

				if u+1 < len(poly) {
					npu := poly[u+1]
					ppu := poly[u-1]

					if pu.X == 0 && (npu.X != 0 || ppu.X != 0) ||
						pu.Y == 0 && (npu.Y != 0 || ppu.Y != 0) ||
						pu.X == w && (npu.X != w || ppu.X != w) ||
						pu.Y == h && (npu.Y != h || ppu.Y != h) {
						break
					}
				}

				for z := v + 1; z < u; z++ {
					pz := poly[z]
					err := lineDistSqr(pv.X, pv.Y, pu.X, pu.Y, pz.X, pz.Y)
					if err > maxErr {
						maxErr = err
					}
				}

				if maxErr > threshSqr {
					break
				} else {
					u += 1
				}
			}

			if u < len(poly) {
				v = u
			} else {
				break
			}
		}

		fmt.Printf("reduced %d vertices to %d\n", len(poly), len(polys[i]))
	}

	return polys
}

func lineDistSqr(x1, y1, x2, y2, px, py float64) float64 {
	return math.Pow((x2-x1)*(y1-py)-(x1-px)*(y2-y1), 2) /
		// --------------------------------------------
		(math.Pow(x2-x1, 2) + math.Pow(y2-y1, 2))
}
