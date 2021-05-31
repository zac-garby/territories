package main

import (
	"math"

	"github.com/fogleman/gg"
	"github.com/zac-garby/territories/world"
)

func main() {
	num := 25
	g := world.NewGen(512, 512, num, 4)
	c := gg.NewContext(512, 512)

	for x := 0; x < 512; x++ {
		for y := 0; y < 512; y++ {
			//n := g.Noise.Noise2D(10*(float64(x)/float64(512)), 10*(float64(y)/float64(512)))
			p := float64(g.Pixels[y][x]) / float64(num)
			c.SetRGB(p, math.Sqrt(1-p), 0)
			c.SetPixel(x, y)
		}
	}

	c.SetRGB(1, 1, 1)
	for _, p := range g.Midpoints {
		c.DrawCircle(p.X, p.Y, 2)
		c.Fill()
	}

	c.SavePNG("out.png")
}
