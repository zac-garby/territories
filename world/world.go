package world

type Coord struct {
	X, Y float64
}

type World struct {
	Width, Height int
	Regions       [][]Coord
}
