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

	p.SetPos(3, 3)
	p.depth++
}

/*****************************************************************************/
type Point struct {
	X, Y int
}

/*****************************************************************************/
//   0 - 1 - 2
//   |   |   |
//   3 - 4 - 5
//   |   |   |
//   6 - 7 - 8

type RoomGraph struct {
	rooms [9]Room
	paths []Path
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

	g.DropRandRooms(2)

	// add one more connection
	//count = 0
	//found := false
	//for !found && count < 10 {
	//	c1 = g.RandCell(1)
	//	c2 = g.RandNeighbour(c1, 1)
	//	if !g.AreConnected(c1, c2) {
	//		g.Connect(c1, c2)
	//		debug.Add("Last %d -> %d", c1, c2)
	//		found = true
	//	} else {
	//		debug.Add("Already connected %d -> %d", c1, c2)
	//	}
	//	count++
	//}

	return &g
}

// ----------------------------------------------------------------------------
func (g *RoomGraph) DropRandRooms(count int) {
	for i := 0; i < count; i++ {
		cell := g.RandCell(1)
		debug.Add("Dropping room %d", cell)
		g.rooms[cell].Drop()
		g.PruneDeadends(cell, 2)
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
func (g *RoomGraph) PruneDeadends(cell int, depth int) {
	if depth == 0 { // handle recursion
		return
	}
	// check if the given cell is a dead end
	if g.rooms[cell].mark == -1 && g.CountPaths(cell) == 1 {

		// if it is, remove all the paths (should be just one)
		for i, p := range g.paths {
			if p.origID == cell || p.destID == cell {
				g.paths[i].Drop()
				debug.Add("dropping path %v, cell=%d", p, cell)
			}
		}

		// TODO handle the case where this creates a new deadend
	}
}

// func (g *RoomGraph) RoomAt(col, row int) Room       {}
// func (g *RoomGraph) IsDeadend(cell int) bool        {}
// func (g *RoomGraph) Path(c1, c2 int)                {}

// ----------------------------------------------------------------------------
func (g *RoomGraph) Room(cell int) Room {
	return g.rooms[cell]
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

/*****************************************************************************/

type Room struct {
	X, Y int
	W, H int
	mark int // 0=unconnected, 1=connected, -1=dropped
}

func (r Room) Center() Point {
	return Point{r.X + r.W/2, r.Y + r.H/2}
}

func (r Room) RandPoint() Point {
	return Point{
		r.X + rand.Intn(r.W-2) + 1,
		r.Y + rand.Intn(r.H-2) + 1,
	}
}

func (r *Room) Drop() {
	r.mark = -1
}

/*****************************************************************************/

type Path struct {
	origID int
	destID int
	mark   int // 0=normal, -1=dropped
}

func (p *Path) Drop() {
	p.mark = -1
}

/*****************************************************************************/
/*****************************************************************************/
/*
 *   0 - 1 - 2
 *   |   |   |
 *   3 - 4 - 5
 *   |   |   |
 *   6 - 7 - 8
 *
 */

type RoomGrid [9]Room

// ----------------------------------------------------------------------------
func (g *RoomGrid) getRandomRoom(mark int) int {
	ids := []int{}

	for i, r := range g {
		if r.mark == mark {
			ids = append(ids, i)
		}
	}
	if len(ids) == 0 {
		return -1
	}

	idx := rand.Intn(len(ids))
	return ids[idx]
}

// ----------------------------------------------------------------------------
func (g *RoomGrid) getRandomNeighbour(id int, mark int) int {
	ids := []int{}

	for _, nid := range nbList[id] {
		if g[nid].mark == mark {
			ids = append(ids, nid)
		} else {
		}
	}
	if len(ids) == 0 {
		return -1
	}

	idx := rand.Intn(len(ids))
	return ids[idx]
}

// ----------------------------------------------------------------------------
func (g *RoomGrid) Direction(id1, id2 int) Direction {
	col1, row1 := id1%3, id1/3
	col2, row2 := id2%3, id2/3

	dx := col2 - col1
	dy := row2 - row1

	switch {
	case dx != 0 && dy != 0:
		// shouldn't happen but let's catch it
		return East
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
func (g *RoomGrid) IsDeadend(id int) bool {

	count := 0
	for _, p := range pathList {
		if (p.origID == id || p.destID == id) && p.mark != -1 && roomGrid[id].mark == -1 {
			count++
		}
	}

	deadend := count == 1
	debug.Add("deadend? %d: %v", id, deadend)

	return deadend
}

// ----------------------------------------------------------------------------
func (g *RoomGrid) DropPaths(id int, maxDepth int) {
	if maxDepth < 0 {
		return
	}

	for i, p := range pathList {
		if p.origID == id || p.destID == id {
			debug.Add("%d dropping path %d: %v", maxDepth, id, p)
			pathList[i].Drop()

			// should check orig and dest to see if they become deadends
			if p.origID != id && g.IsDeadend(p.origID) {
				g.DropPaths(p.origID, maxDepth-1)
			}
			if p.destID != id && g.IsDeadend(p.destID) {
				g.DropPaths(p.destID, maxDepth-1)
			}
		}
	}
}

/*****************************************************************************/

var roomGrid = RoomGrid{}

var pathList = []Path{}

var nbList [9][]int

// ----------------------------------------------------------------------------
func getAreaBounds(roomW, roomH int) (areas []Room) {

	// split the map into 3x3 areas determine the bounds of each one
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			r := Room{X: (roomW + 1) * j, Y: (roomH + 1) * i, W: roomW, H: roomH}
			areas = append(areas, r)
		}
	}
	return
}

// ----------------------------------------------------------------------------
func makeRandomRooms(bounds []Room, roomW, roomH int) {

	// make a random room within each area
	for i, a := range bounds {
		randW := rand.Intn(12) + 8
		randH := rand.Intn(roomH-4) + 4
		dx := rand.Intn(roomW - randW) // position within the boundary area
		dy := rand.Intn(roomH - randH)
		r := Room{X: a.X + dx, Y: a.Y + dy, W: randW, H: randH}
		roomGrid[i] = r
	}

}

// ----------------------------------------------------------------------------
func dropRandomRooms(count int) {
	for i := 0; i < count; i++ {
		idx := rand.Intn(len(roomGrid))
		if roomGrid[idx].mark >= 0 {
			roomGrid[idx].mark = -1
		} else {
			i-- // if room has already been dropped, choose again
		}
	}
}

// ----------------------------------------------------------------------------
func getNeighbours(id int, depth int) []int {

	// handle recursion
	if depth <= 0 {
		return []int{}
	}

	ref := [9][]int{
		{1, 3},    // room 0 neighbours
		{0, 4, 2}, // room 1
		{1, 5},    // etc
		{0, 4, 6},
		{1, 5, 7, 3},
		{2, 4, 8},
		{3, 7},
		{6, 4, 8},
		{7, 5},
	}
	shortList := ref[id]

	finalList := map[int]bool{}
	for _, shortID := range shortList {

		// if this room has been dropped, find its neighbours instead
		if roomGrid[shortID].mark == -1 {
			n2 := getNeighbours(shortID, depth-1)
			for _, innerID := range n2 {
				finalList[innerID] = true
			}
		} else {
			finalList[shortID] = true
		}
	}

	// the original room id will appear after recursion cases
	delete(finalList, id)

	//convert finalList into an array to return
	finalArr := make([]int, len(finalList))
	idx := 0
	for k := range finalList {
		finalArr[idx] = k
		idx++
	}
	return finalArr
}

// ----------------------------------------------------------------------------
func connectRooms(origID int, destID int) {
	p := Path{origID: origID, destID: destID}
	pathList = append(pathList, p)
	roomGrid[origID].mark = 1
	roomGrid[destID].mark = 1
}

// ----------------------------------------------------------------------------
func generateRandomLevel_old(d *DungeonMap, ml *MonsterList, p *Player) {
	debug.Clear()
	d.Clear()

	pathList = []Path{}

	// assume 3x3 rooms on the map
	roomW := (MapMaxX - 2) / 3 // 25
	roomH := (MapMaxY - 2) / 3 // 6

	areas := getAreaBounds(roomW, roomH)
	makeRandomRooms(areas, roomW, roomH)

	// get the list of neighbours for each room
	for idx := range roomGrid {
		nbList[idx] = getNeighbours(idx, 1)
	}

	// pick a random room (that's not dropped)
	id1 := roomGrid.getRandomRoom(0)
	debug.Add("first: %d", id1)

	// connect it to one of its neighbours
	id2 := roomGrid.getRandomNeighbour(id1, 0)
	debug.Add("next: %d", id2)
	connectRooms(id1, id2)

	count := 0
	id := 0
	destId := 0

	// pick an another unconnected room at random
	id = roomGrid.getRandomRoom(0)

	// if it exists (and we haven't hit the limit)
	for id != -1 && count < 20 {

		// connect the chosen room to a neighbour that's already connected
		destId = roomGrid.getRandomNeighbour(id, 1)
		if destId != -1 {
			debug.Add("checking: %d ->  %d connected  (count=%d)", id, destId, count)
			connectRooms(id, destId)
		} else {
			// if there are none, just skip it
			debug.Add("checking: %d -> %d   skipped  (count=%d)", id, destId, count)
		}

		// pick another random unconnected room
		id = roomGrid.getRandomRoom(0)
		count++
	}

	dropRandomRooms(2)

	for i := range roomGrid {
		if roomGrid.IsDeadend(i) {
			roomGrid.DropPaths(i, 2)
		}
	}

	x, y := buildMap(d)
	p.SetPos(x, y)
	p.depth++
}

// ----------------------------------------------------------------------------
func buildMap(d *DungeonMap) (int, int) {
	// actually create the rooms on the map
	for idx, r := range roomGrid {
		if roomGrid[idx].mark == 1 {
			d.CreateRoom(r.X, r.Y, r.W, r.H)
		}
	}

	// actually create the paths on the map
	for _, p := range pathList {

		if p.mark == 0 { // ignore dropped paths (-1)
			pt1 := roomGrid[p.origID].Center()
			pt2 := roomGrid[p.destID].Center()
			dir := roomGrid.Direction(p.origID, p.destID)
			//debug.Add("making path: %d -> %d, dir=%v", p.origID, p.destID, dir))
			d.ConnectRooms(pt1.X, pt1.Y, pt2.X, pt2.Y, dir)

			// bug workaround: for missing rooms, some paths aren't fully connected
			if roomGrid[p.destID].mark == -1 {
				pt := roomGrid[p.destID].Center()
				d.SetTile(pt.X, pt.Y, TilePath)
			}
		}
	}

	// place the player in a random location (as well as the stairs up)
	playerID := roomGrid.getRandomRoom(1)
	playerPt := roomGrid[playerID].RandPoint()
	d.SetTile(playerPt.X, playerPt.Y, TileStairsUp)

	// place the stairs down in a random location
	stairsID := roomGrid.getRandomRoom(1)
	stairsPt := roomGrid[stairsID].RandPoint()
	d.SetTile(stairsPt.X, stairsPt.Y, TileStairsDn)

	return playerPt.X, playerPt.Y
}
