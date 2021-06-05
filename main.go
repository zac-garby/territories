package main

import (
	"github.com/fogleman/gg"
	"github.com/zac-garby/territories/world"
)

func main() {
	num := 25
	g := world.NewGen(512, 512, num, 4)
	c := gg.NewContext(g.Width, g.Height)

	for x := 0; x < g.Width; x++ {
		for y := 0; y < g.Height; y++ {
			p := float64(g.Pixels[y][x]) / float64(g.NumRegions)
			c.SetRGB(p, p, p)
			c.SetPixel(x, y)
		}
	}

	c.SavePNG("out.png")
}
