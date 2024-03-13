package main

// -----------------------------------------------------------------------
type Direction int

const (
	North Direction = iota
	East
	South
	West
)

func (d Direction) String() string {
	switch d {
	case North:
		return "north"
	case East:
		return "east"
	case South:
		return "south"
	case West:
		return "west"
	default:
		return "unknown"
	}
}

func getDirectionCoords(dir Direction) (int, int) {
	dx, dy := 0, 0
	switch dir {
	case North:
		dy = -1
	case South:
		dy = 1
	case East:
		dx = 1
	case West:
		dx = -1
	}
	return dx, dy
}
