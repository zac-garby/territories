package game

import "github.com/zac-garby/territories/world"

type Game struct {
	world *world.World
}

func NewGame() *Game {
	return &Game{}
}
