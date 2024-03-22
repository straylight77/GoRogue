package main

import (
	"fmt"
	"slices"
)

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
type Path struct { // Need to rename all other Path (room connections)
	steps []Coord
	algo  string
	iter  int
}

func (p Path) String() string {
	return fmt.Sprintf("len=%d, algo=%s, iter=%d", len(p.steps), p.algo, p.iter)
}

// A simple Breadth First Seach pathfinding algorithm.  Using A* would be
// more optimal but the complexity is low for this game (small map, only
// a few monsters chasing at any given time.)
// https://www.redblobgames.com/pathfinding/a-star/introduction.html
func findPathBFS(dm *DungeonMap, x1, y1 int, x2, y2 int) Path {
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
	path := Path{
		algo: "bfs",
		iter: pathCount,
	}
	var ok bool
	current := end
	for current != start {
		path.steps = append(path.steps, current)
		current, ok = cameFrom[current]
		if !ok {
			break
		}
	}
	slices.Reverse(path.steps)
	return path
}

// -----------------------------------------------------------------------
func drawPathDebug(disp *Display, path Path, ch rune) {
	for _, pos := range path.steps {
		disp.Screen.SetContent(pos.X, pos.Y+1, ch, nil, disp.Style("debug2"))
	}
}

// -----------------------------------------------------------------------
func drawPathDebugIdx(disp *Display, path Path) {
	for i, pos := range path.steps {
		ch := rune('1' + i%10 - 1)
		disp.Screen.SetContent(pos.X, pos.Y+1, ch, nil, disp.Style("debug2"))
	}
}

/******************************************************************************
* Dijkstra Map or Distance Transform
* https://www.roguebasin.com/index.php/Dijkstra_Maps_Visualized
* https://www.redblobgames.com/pathfinding/tower-defense/
 */

type DMap struct {
	targets  []Coord
	distance map[Coord]int
	iter     int
}

// In most cases we want to give some targets and calculate right away
func newDMap(dng *DungeonMap, targets ...Coord) *DMap {
	m := &DMap{
		make([]Coord, 0),
		make(map[Coord]int),
		0,
	}

	m.AddTargets(targets...)
	m.Calculate(dng)
	return m
}

func (m *DMap) Reset(dng *DungeonMap, targets ...Coord) {
	m.Clear()
	m.AddTargets(targets...)
	m.Calculate(dng)
}

func (m *DMap) AddTargets(c ...Coord) {
	m.targets = append(m.targets, c...)
}

func (m *DMap) RemoveTarget(c Coord) {
}

func (m *DMap) Clear() {
	m.targets = make([]Coord, 0)
	m.distance = make(map[Coord]int)
	m.iter = 0

}

func (m *DMap) Calculate(dng *DungeonMap) {

	frontier := CoordQueue{}
	iterations := 0

	// Initialize
	for _, start := range m.targets {
		frontier.Add(start)
		m.distance[start] = 0
	}

	for !frontier.IsEmpty() {
		current := frontier.Next()
		nb := m.neighbours(current)
		for _, next := range nb {
			_, reached := m.distance[next]
			if !reached && dng.IsWalkable(current, next) {
				frontier.Add(next)
				m.distance[next] = m.distance[current] + 1
			}
		}
		iterations++
	}
	m.iter = iterations
}

func (m *DMap) neighbours(pos Coord) []Coord {
	// The order here determines how we traverse the graph
	return []Coord{
		// Cardinal directions first
		{pos.X - 1, pos.Y},
		{pos.X, pos.Y + 1},
		{pos.X + 1, pos.Y},
		{pos.X, pos.Y - 1},
		// Then the diagonals
		{pos.X - 1, pos.Y - 1},
		{pos.X - 1, pos.Y + 1},
		{pos.X + 1, pos.Y - 1},
		{pos.X + 1, pos.Y + 1},
	}
}

func (m *DMap) PathFrom(pos Coord) Path {
	path := Path{
		algo: "dmap",
		iter: 0,
	}
	current := pos
	for m.distance[current] != 0 {
		path.steps = append(path.steps, current)
		current = m.NextStep(current)
	}

	return path
}

func (m *DMap) NextStep(pos Coord) Coord {

	// Reverse the search order since we're 'going downhill' compared to how
	// the dmap was created
	toCheck := m.neighbours(pos)
	slices.Reverse(toCheck)

	// We don't want just any neighbour that has less distance than our
	// current one, we need to follow the sequence given (ie. the 'next'
	// distance, not simple the minumum distance).  This will handle
	// restricted diagonal movement (e.g. through doors, cooridors)

	nextDist := m.distance[pos] - 1
	nextCoord := Coord{-1, -1}
	for _, check := range toCheck {
		dist, ok := m.distance[check]
		if ok && dist == nextDist {
			nextDist = m.distance[check]
			nextCoord = check
		}
	}

	return nextCoord
}

func (m *DMap) Draw(disp *Display) {

	styleName := []string{
		"yellow",
		"orange",
		"red",
		"purple",
		"darkblue",
		"blue",
		"bluegreen",
		"green",
		"darkgreen",
	}

	for pos, dist := range m.distance {
		if dist != 0 {
			color := dist / 10
			style := disp.Style("default")
			if color < len(styleName) {
				style = disp.Style(styleName[color])
			}
			ch := rune('1' + (dist % 10) - 1)
			disp.Screen.SetContent(pos.X, pos.Y+1, ch, nil, style)
		}
	}
}
