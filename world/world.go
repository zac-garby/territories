package world

import (
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
}

func NewGen(width, height, npoints int, seed int64) *WorldGen {
	g := &WorldGen{Width: width, Height: height, Seed: seed}
	g.Noise = perlin.NewPerlin(alpha, beta, n, seed)
	g.Rand = rand.New(rand.NewSource(seed))

	g.precomputePerlin()
	g.initPoints(npoints)
	g.voronoi()

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

	//fmt.Println(noise)

	return noise // + 10*noise // * (1 + noise*z)
}

func (g *WorldGen) disjoinRegions() {
	for x := 0; x < g.Width; x++ {
		for y := 0; y < g.Height; y++ {

		}
	}
}
