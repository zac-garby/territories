package main

import (
	"github.com/fogleman/gg"
	"github.com/zac-garby/territories/world"
)

func main() {
	g := world.NewGen(512, 512, 15, 3)
	c := gg.NewContext(512, 512)
	for x := 0; x < 512; x++ {
		for y := 0; y < 512; y++ {
			p := float64(g.Pixels[y][x]) / 20
			c.SetRGB(p, p, p)
			c.SetPixel(x, y)
		}
	}
	c.SavePNG("out.png")
}
