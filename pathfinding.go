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
//type Path struct { // Need to rename all other Path (room connections)
//	steps []Coord
//	count int
//}

// A simple Breadth First Seach pathfinding algorithm.  Using A* would be
// more optimal but the complexity is low for this game (small map, only
// a few monsters chasing at any given time.)
// https://www.redblobgames.com/pathfinding/a-star/introduction.html
func findPathBFS(dm *DungeonMap, x1, y1 int, x2, y2 int) []Coord {
	// Declarations
	start := Coord{x1, y1}
	end := Coord{x2, y2}
	frontier := CoordQueue{}
	cameFrom := map[Coord]Coord{}
	pathCount := 0

	// Initialize
	frontier.Add(start)
	cameFrom[start] = Coord{start.X, start.Y}

	// While path not found yet or no more explorable areas
	_, foundPath := cameFrom[end]
	for !frontier.IsEmpty() && !foundPath {
		current := frontier.Next()

		nb := dm.getWalkableNeighbours(current)
		for _, next := range nb {
			_, reached := cameFrom[next]
			if !reached {
				frontier.Add(next)
				cameFrom[next] = current
			}
		}
		_, foundPath = cameFrom[end]
		pathCount++
	}

	//debug.Add("path found: %d steps", pathCount)

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
		disp.Screen.SetContent(pos.X, pos.Y+1, ch, nil, disp.Style("debug2"))
	}
}
