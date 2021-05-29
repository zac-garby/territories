package world

import (
	"fmt"
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
	x, y float64
}

type WorldGen struct {
	noise         *perlin.Perlin
	rand          *rand.Rand
	midpoints     []Coord
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
	g.midpoints = make([]Coord, n)

	for i := 0; i < n; i++ {
		g.midpoints[i].x = g.rand.Float64() * float64(g.width)
		g.midpoints[i].y = g.rand.Float64() * float64(g.height)
		fmt.Println(g.midpoints[i])
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

	for i, p := range g.midpoints {
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
		mx    = (a.x + b.x) / 2
		my    = (a.y + b.y) / 2
		dx    = b.x - a.x
		dy    = b.y - a.y
		noise = g.noise.Noise2D(w*(mx/float64(g.width)), w*(my/float64(g.height)))
	)

	return (math.Abs(dx) + math.Abs(dy)) * (1 + noise*z)
}
