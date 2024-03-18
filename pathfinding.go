package main

import (
	"slices"
)

type Coord struct {
	X, Y int
}

// -----------------------------------------------------------------------
type CoordQueue struct {
	items []Coord
	idx   int
}

func (q *CoordQueue) Add(item Coord) {
	q.items = append(q.items, item)

}
func (q *CoordQueue) Next() Coord {
	if !q.IsEmpty() {
		q.idx++
		return q.items[q.idx-1]
	} else {
		return Coord{-1, -1}
	}
}

func (q *CoordQueue) IsEmpty() bool {
	return q.idx >= len(q.items)
}

// -----------------------------------------------------------------------
// https://www.redblobgames.com/pathfinding/a-star/introduction.html
func findPathBFS(dm *DungeonMap, x1, y1 int, x2, y2 int) []Coord {
	// Declarations
	start := Coord{x1, y1}
	end := Coord{x2, y2}
	frontier := CoordQueue{}
	reached := map[Coord]bool{}
	cameFrom := map[Coord]Coord{}
	pathCount := 0

	// Initialize
	frontier.Add(start)
	reached[start] = true
	cameFrom[start] = Coord{-1, -1}

	// While path not found yet or no more explorable areas
	_, foundPath := cameFrom[end]
	for !frontier.IsEmpty() && !foundPath {
		current := frontier.Next()

		nb := dm.getWalkableNeighbours(current)
		//debug.Add("nb: %v", nb)
		for _, next := range nb {
			if !reached[next] {
				frontier.Add(next)
				reached[next] = true
				cameFrom[next] = current
			}
		}
		_, foundPath = cameFrom[end]
		pathCount++
	}

	// Build a slice to hold the path we found
	var ok bool
	path := []Coord{}
	current := end
	for current != start {
		path = append(path, current)
		current, ok = cameFrom[current]
		if !ok {
			break
		}
	}
	slices.Reverse(path)
	return path
}

// -----------------------------------------------------------------------
func drawPathDebug(disp *Display, path []Coord, ch rune) {
	for _, pos := range path {
		disp.Screen.SetContent(pos.X, pos.Y+1, ch, nil, disp.Debug2Style)
	}
}
