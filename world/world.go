package world

import (
	"math"
	"math/rand"

	perlin "github.com/aquilax/go-perlin"
)

const (
	alpha float64 = 2
	beta  float64 = 2
	n     int     = 1
	w     float64 = 18
	z     float64 = 0.25
)

type Coord struct {
	X, Y float64
}

type WorldGen struct {
	noise         *perlin.Perlin
	rand          *rand.Rand
	Midpoints     []Coord
	Pixels        [][]int
	width, height int
	seed          int64
}

func NewGen(width, height, npoints int, seed int64) *WorldGen {
	g := &WorldGen{width: width, height: height, seed: seed}
	g.noise = perlin.NewPerlin(alpha, beta, n, seed)
	g.rand = rand.New(rand.NewSource(seed))

	g.initPoints(npoints)
	g.voronoi()

	return g
}

func (g *WorldGen) initPoints(n int) {
	g.Midpoints = make([]Coord, n)

	for i := 0; i < n; i++ {
		g.Midpoints[i].X = g.rand.Float64() * float64(g.width)
		g.Midpoints[i].Y = g.rand.Float64() * float64(g.height)
	}
}

func (g *WorldGen) voronoi() {
	g.Pixels = make([][]int, g.height)
	for y := 0; y < g.height; y++ {
		row := make([]int, g.width)
		for x := 0; x < g.width; x++ {
			closest := g.closestPoint(Coord{float64(x), float64(y)})
			row[x] = closest
		}
		g.Pixels[y] = row
	}
}

func (g *WorldGen) closestPoint(c Coord) int {
	var (
		cd = float64(g.width * g.height)
		ci = -1
	)

	for i, p := range g.Midpoints {
		dist := g.distance(c, p)
		if dist < cd {
			cd = dist
			ci = i
		}
	}

	return ci
}

func (g *WorldGen) distance(a, b Coord) float64 {
	var (
		mx    = (a.X + b.X) / 2
		my    = (a.Y + b.Y) / 2
		dx    = b.X - a.X
		dy    = b.Y - a.Y
		noise = g.noise.Noise2D(w*(mx/float64(g.width)), w*(my/float64(g.height)))
	)

	return (math.Abs(dx) + math.Abs(dy)) * (1 + noise*z)
}
