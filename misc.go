package main

// -----------------------------------------------------------------------
type Entity interface {
	Pos() (int, int)
	SetPos(int, int)
	Rune() rune
	Attack(Entity) string
	Label() string
	AdjustHP(amt int)
}

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

// -----------------------------------------------------------------------
func abs(val int) int {
	if val < 0 {
		val = -val
	}
	return val
}

// -----------------------------------------------------------------------
func max(val1, val2 int) int {
	if val1 > val2 {
		return val1
	} else {
		return val2
	}
}
