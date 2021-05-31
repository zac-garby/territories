package main

import (
	"github.com/fogleman/gg"
	"github.com/zac-garby/territories/world"
)

func main() {
	num := 25
	g := world.NewGen(512, 512, num, 4)
	c := gg.NewContext(512, 512)

	for x := 0; x < 512; x++ {
		for y := 0; y < 512; y++ {
			n := g.Interperlin(float64(x), float64(y))
			//p := float64(g.Pixels[y][x]) / float64(num)
			c.SetRGB(0.5+0.5*n, 0.5+0.5*n, 0.5+0.5*n)
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
