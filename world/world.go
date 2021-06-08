package world

type Coord struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type World struct {
	Width, Height int
	Regions       [][]Coord
	Adjacency     [][]bool
	Centroids     []Coord
}
