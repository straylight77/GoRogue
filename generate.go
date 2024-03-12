package main

import (
	"math/rand"
)

var debug = DebugMessageLog{}

// ----------------------------------------------------------------------------

var graph *RoomGraph

func generateRandomLevel(dm *DungeonMap, ml *MonsterList, p *Player) {
	debug.Clear()
	dm.Clear()

	graph = newRandomGraph()

	graph.getAreaBounds()
	graph.MakeRandomRooms()
	x, y := buildMap(graph, dm)

	p.SetPos(x, y)
	p.depth++
}

// ----------------------------------------------------------------------------
func buildMap(g *RoomGraph, d *DungeonMap) (int, int) {

	// create the rooms on the dungeon map
	for _, r := range g.rooms {
		if r.mark == 1 {
			d.CreateRoom(r.X, r.Y, r.W, r.H)
		}
	}

	// create the paths on the map
	for _, p := range g.paths {

		if p.mark == 0 { // ignore dropped paths (-1)
			x1, y1 := g.rooms[p.origID].Center()
			x2, y2 := g.rooms[p.destID].Center()
			dir := g.Direction(p.origID, p.destID)
			debug.Add("making path: %d -> %d, dir=%v", p.origID, p.destID, dir)
			d.ConnectRooms(x1, y1, x2, y2, dir)

			// BUG: paths aren't fully connected for dropped rooms
			if g.rooms[p.destID].mark == -1 {
				x3, y3 := g.rooms[p.destID].Center()
				d.SetTile(x3, y3, TilePath)
			}
		}
	}

	// place the player in a random location (as well as the stairs up)
	c1 := g.RandCell(1)
	pX, pY := g.rooms[c1].RandPoint()
	d.SetTile(pX, pY, TileStairsUp)

	// place the stairs down in a random location
	c2 := g.RandCell(1)
	sX, sY := g.rooms[c2].RandPoint()
	d.SetTile(sX, sY, TileStairsDn)

	return pX, pY
}

/*****************************************************************************/
//   0 - 1 - 2
//   |   |   |
//   3 - 4 - 5
//   |   |   |
//   6 - 7 - 8

type RoomGraph struct {
	rooms  [9]Room
	paths  []Path
	bounds [9]Room
}

// ----------------------------------------------------------------------------
func newRandomGraph() *RoomGraph {
	g := RoomGraph{}

	c1 := g.RandCell(0) // connect 2 rooms at random
	c2 := g.RandNeighbour(c1, 0)
	g.Connect(c1, c2)
	debug.Add("First %d -> %d", c1, c2)

	count := 0
	next := g.RandCell(0)          // pick a random unconnected room
	for next != -1 && count < 20 { // while there are unconnected rooms

		nb := g.RandNeighbour(next, 1) // connect it to an already connected neighbour
		if nb != -1 {                  // if there are none, just skip it
			g.Connect(next, nb)
			debug.Add("Connect %d -> %d", next, nb)
		}
		next = g.RandCell(0) // pick the next unconnected room
		count++
	}

	g.DropRandomRooms(2)

	return &g
}

// ----------------------------------------------------------------------------
func (g *RoomGraph) DropRandomRooms(count int) {
	for i := 0; i < count; i++ {
		cell := g.RandCell(1)
		debug.Add("Dropping room %d", cell)
		g.rooms[cell].mark = -1
		g.PruneDeadends(cell)
	}
}

// ----------------------------------------------------------------------------
func (g *RoomGraph) CountPaths(cell int) int {
	count := 0
	for _, p := range g.paths {
		if (p.origID == cell || p.destID == cell) && p.mark != -1 {
			count++
		}
	}
	return count
}

// ----------------------------------------------------------------------------
func (g *RoomGraph) PruneDeadends(cell int) {
	// check if the given cell is a dead end
	if g.rooms[cell].mark == -1 && g.CountPaths(cell) == 1 {

		// if it is, remove all the paths (should be just one)
		for i, p := range g.paths {
			if p.origID == cell || p.destID == cell {
				g.paths[i].mark = -1
				debug.Add("dropping path %v, cell=%d", p, cell)
			}
		}
	}
}

// ----------------------------------------------------------------------------
func (g *RoomGraph) Neighbours(cell int) []int {
	ref := [9][]int{
		{1, 3},       // cell 0
		{0, 4, 2},    // cell 1
		{1, 5},       // cell 2
		{0, 4, 6},    // cell 3
		{1, 5, 7, 3}, // cell 4
		{2, 4, 8},    // cell 5
		{3, 7},       // cell 6
		{6, 4, 8},    // cell 7
		{7, 5},       // cell 8
	}
	return ref[cell]
}

// ----------------------------------------------------------------------------
func (g *RoomGraph) RandCell(mark int) int {
	cells := []int{}

	for i, r := range g.rooms {
		if r.mark == mark {
			cells = append(cells, i)
		}
	}
	if len(cells) == 0 {
		return -1
	}

	idx := rand.Intn(len(cells))
	return cells[idx]
}

// ----------------------------------------------------------------------------
func (g *RoomGraph) RandNeighbour(cell int, mark int) int {
	nbList := []int{}

	for _, nb := range g.Neighbours(cell) {
		if g.rooms[nb].mark == mark {
			nbList = append(nbList, nb)
		}
	}

	if len(nbList) == 0 {
		return -1
	}

	idx := rand.Intn(len(nbList))
	return nbList[idx]
}

// ----------------------------------------------------------------------------
func (g *RoomGraph) Direction(c1, c2 int) Direction {
	col1, row1 := c1%3, c1/3
	col2, row2 := c2%3, c2/3

	dx := col2 - col1
	dy := row2 - row1

	switch {
	case dx != 0 && dy != 0:
		return East // shouldn't happen but let's catch it
	case dx > 0:
		return East
	case dx < 0:
		return West
	case dy > 0:
		return South
	case dy < 0:
		return North
	default:
		return East
	}
}

// ----------------------------------------------------------------------------
func (g *RoomGraph) Connect(c1, c2 int) {
	p := Path{origID: c1, destID: c2}
	g.paths = append(g.paths, p)
	g.rooms[c1].mark = 1
	g.rooms[c2].mark = 1
	//debug.Add("Connect(%d, %d) %v %v", c1, c2, g.rooms[c1], g.rooms[c2])
}

// ----------------------------------------------------------------------------
func (g *RoomGraph) AreConnected(c1, c2 int) bool {
	for _, p := range g.paths {
		if p.mark != -1 &&
			(p.origID == c1 || p.destID == c1) &&
			(p.origID == c2 || p.destID == c2) {
			//debug.Add("AreConnected: %d, %d, yes!", c1, c2)
			return true
		}
	}
	return false
}

// ----------------------------------------------------------------------------
func (g *RoomGraph) getAreaBounds() (areas []Room) {

	// assume 3x3 rooms on the map
	roomW := (MapMaxX - 2) / 3 // 25
	roomH := (MapMaxY - 2) / 3 // 6

	// split the map into 3x3 areas determine the bounds of each one
	idx := 0
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			r := Room{X: (roomW + 1) * j, Y: (roomH + 1) * i, W: roomW, H: roomH}
			g.bounds[idx] = r
			idx++
		}
	}
	return
}

// ----------------------------------------------------------------------------
func (g *RoomGraph) MakeRandomRooms() {

	// make a random room within each area
	for i, a := range g.bounds {
		randW := rand.Intn(12) + 8    // between 8 and 20
		randH := rand.Intn(a.H-4) + 4 // between 4 and height of area
		dx := rand.Intn(a.W - randW)  // position within the boundary area
		dy := rand.Intn(a.H - randH)
		g.rooms[i].SetSize(a.X+dx, a.Y+dy, randW, randH)
	}

}

/*****************************************************************************/

type Room struct {
	X, Y int
	W, H int
	mark int // 0=unconnected, 1=connected, -1=dropped
}

func (r *Room) SetSize(x, y, w, h int) {
	r.X = x
	r.Y = y
	r.W = w
	r.H = h
}

func (r Room) Center() (x int, y int) {
	x = r.X + r.W/2
	y = r.Y + r.H/2
	return
}

func (r Room) RandPoint() (x int, y int) {
	x = r.X + rand.Intn(r.W-2) + 1
	y = r.Y + rand.Intn(r.H-2) + 1
	return
}

/*****************************************************************************/

type Path struct {
	origID int
	destID int
	mark   int // 0=normal, -1=dropped
}
