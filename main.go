package main

import (
	"math"

	"github.com/fogleman/gg"
	"github.com/zac-garby/territories/world"
)

func main() {
	g := world.NewGen(512, 512, 100, 3)
	c := gg.NewContext(512, 512)

	for x := 0; x < 512; x++ {
		for y := 0; y < 512; y++ {
			p := float64(g.Pixels[y][x]) / 100
			c.SetRGB(math.Sqrt(p), math.Sqrt(1-p), math.Pow(p, p))
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
