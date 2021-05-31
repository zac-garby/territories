package world

import (
	"math"
	"math/rand"

	perlin "github.com/aquilax/go-perlin"
)

const (
	alpha float64 = 2
	beta  float64 = 2
	n     int     = 3
	w     float64 = 7
	z     float64 = 0.35
	delta float64 = 1
	nint  int     = 3
)

type Coord struct {
	X, Y float64
}

type WorldGen struct {
	Noise         *perlin.Perlin
	rand          *rand.Rand
	Midpoints     []Coord
	Pixels        [][]int
	width, height int
	seed          int64
}

func NewGen(width, height, npoints int, seed int64) *WorldGen {
	g := &WorldGen{width: width, height: height, seed: seed}
	g.Noise = perlin.NewPerlin(alpha, beta, n, seed)
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
			row[x] = g.closestPoint(Coord{float64(x), float64(y)})
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
		noise  = 0.0 //math.Abs(g.noise.Noise2D(w*(b.X/float64(g.width)), w*(b.Y/float64(g.height))) - g.noise.Noise2D(w*(a.X/float64(g.width)), w*(a.Y/float64(g.height))))
	)

	delta := length / float64(nint)

	for i := 0; i < nint; i++ {
		x := a.X + float64(i)*delta*ndx
		y := a.Y + float64(i)*delta*ndy
		n := delta * math.Abs(g.Noise.Noise2D(w*(x/float64(g.width)), w*(y/float64(g.height))))
		noise += math.Sqrt(delta*ndx*delta*ndx + delta*ndy*delta*ndy + 81*n*n)
	}

	//fmt.Println(noise)

	return noise // + 10*noise // * (1 + noise*z)
}

func (g *WorldGen) disjoinRegions() {
	for x := 0; x < g.width; x++ {
		for y := 0; y < g.height; y++ {

		}
	}
}
