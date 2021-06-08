package game

import "github.com/zac-garby/territories/world"

type Game struct {
	World *world.World
}

func NewGame(width, height, npoints int, seed int64) *Game {
	g := &Game{}

	wg := world.NewGen(width, height, npoints, seed)
	g.World = wg.World()

	return g
}
