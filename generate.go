package main

import (
	"math/rand"
)

var graph *RoomGraph = &RoomGraph{}

// ----------------------------------------------------------------------------
// func generateRandomLevel(dm *DungeonMap, ml *MonsterList, p *Player) {
func generateRandomLevel(gs *GameState) {
	debug.Clear()
	gs.dungeon.Clear()
	gs.monsters.Clear()

	graph = newRandomGraph()

	graph.MakeCellBounds()
	graph.MakeRandomRooms()
	pos := buildMap(graph, gs.dungeon)

	// Populate with monsters
	N := 5
	for i := 0; i < N; i++ {
		r := graph.RandCell(1)
		// TODO: use Coord
		pos := graph.rooms[r].RandPoint()
		//if mX != gs.player.X && mY != gs.player.Y && gs.monsters.MonsterAt(mX, mY) == nil {
		if pos != gs.player.Pos() && gs.monsters.MonsterAt(pos.X, pos.Y) == nil {
			gs.monsters.Add(randomMonster(gs.player.depth), pos.X, pos.Y)
		} else {
			i--
		}
	}

	gs.player.SetPos(pos)
	gs.player.depth++

}

// ----------------------------------------------------------------------------
// Takes a completed RoomGraph and changes the tiles in DungeonMap appropriately
// Returns the position of the Stairs Up (in order to set the Player's position)
func buildMap(g *RoomGraph, d *DungeonMap) Coord {

	// create the rooms on the dungeon map
	for _, r := range g.rooms {
		if r.mark == 1 {
			d.CreateRoom(r.X, r.Y, r.W, r.H)
		}
	}

	// create the corridors on the map
	for _, p := range g.corridors {
		if p.mark == 0 { // ignore dropped corridors (-1)

			var x1, y1, x2, y2 int

			// If the room has been dropped use, its center. Otherwise use a
			// random point on the wall closest to the destination cell.
			dir1 := g.Direction(p.origID, p.destID)
			if g.rooms[p.origID].mark == 1 {
				x1, y1 = g.rooms[p.origID].RandWallPoint(dir1)
			} else {
				x1, y1 = g.rooms[p.origID].Center()
			}

			// Same logic as above for the destination room
			dir2 := g.Direction(p.destID, p.origID)
			if g.rooms[p.destID].mark == 1 {
				x2, y2 = g.rooms[p.destID].RandWallPoint(dir2)
			} else {
				x2, y2 = g.rooms[p.destID].Center()
			}

			//debug.Add("making corridor: %d -> %d, dir=%v", p.origID, p.destID, dir)
			d.ConnectRooms(x1, y1, x2, y2, dir1)
		}
	}

	// place the player in a random location (as well as the stairs up)
	c1 := g.RandCell(1)
	pos1 := g.rooms[c1].RandPoint()
	d.SetTile(pos1, TileStairsUp)

	// place the stairs down in a random location
	c2 := g.RandCell(1)
	pos2 := g.rooms[c2].RandPoint()
	d.SetTile(pos2, TileStairsDn)

	return pos1
}

/*****************************************************************************/
//   0 - 1 - 2
//   |   |   |
//   3 - 4 - 5
//   |   |   |
//   6 - 7 - 8

type RoomGraph struct {
	rooms     [9]Room
	corridors []Corridor
	bounds    [9]Room
}

// ----------------------------------------------------------------------------
// Creates a new randomized room graph.  Follows the following steps:
//  1. Pick a random room
//  2. Connect it to a random neighbour
//  3. While there are unconnected rooms remaining:
//     a. Pick a random unconnected room
//     b. Connect it to a random connected neighbour
//     (This guarantees all rooms will end up connected together.)
//     c. If there are no connected neighbours, skip it and continue looping.
//     d. Limit the nubmer of loop to 20 as safeguard.
//  4. Drop 2 of the rooms
//  5. Check for dead ends and prune them (done by DropRandomRooms())
func newRandomGraph() *RoomGraph {
	g := RoomGraph{}

	c1 := g.RandCell(0) // Connect 2 rooms at random
	c2 := g.RandNeighbour(c1, 0)
	g.Connect(c1, c2)
	//debug.Add("First %d -> %d", c1, c2)

	count := 0
	next := g.RandCell(0)          // Pick a random unconnected room
	for next != -1 && count < 20 { // While there are unconnected rooms

		nb := g.RandNeighbour(next, 1) // Connect it to an already connected neighbour
		if nb != -1 {                  // If there are none, just skip it
			g.Connect(next, nb)
			//debug.Add("Connect %d -> %d", next, nb)
		}
		next = g.RandCell(0) // Pick the next unconnected room
		count++
	}

	// Add a few more connections to keep it interesting
	n := rand.Intn(2) + 1 // 1-2
	for i := 0; i < n; i++ {
		found := false
		count = 0
		for !found && count < 10 {
			c1 = g.RandCell(1)
			c2 = g.RandNeighbour(c1, 1)
			if !g.AreConnected(c1, c2) {
				g.Connect(c1, c2)
				//debug.Add("Last %d -> %d", c1, c2)
				found = true
			}
		}
	}

	g.DropRandomRooms(2)

	return &g
}

// ----------------------------------------------------------------------------
// Checks if there's a non-dropped corridor between the given cells
func (g *RoomGraph) AreConnected(c1, c2 int) bool {
	for _, p := range g.corridors {
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
// Creates a corridor between the given cells
func (g *RoomGraph) Connect(c1, c2 int) {
	p := Corridor{origID: c1, destID: c2}
	g.corridors = append(g.corridors, p)
	g.rooms[c1].mark = 1
	g.rooms[c2].mark = 1
	//debug.Add("Connect(%d, %d) %v %v", c1, c2, g.rooms[c1], g.rooms[c2])
}

// ----------------------------------------------------------------------------
// Returns the total number corridors coming in or going out of the given cell
func (g *RoomGraph) CountCorridors(cell int) int {
	count := 0
	for _, p := range g.corridors {
		if (p.origID == cell || p.destID == cell) && p.mark != -1 {
			count++
		}
	}
	return count
}

// ----------------------------------------------------------------------------
// Gives the relative Direction when going from c1 to c2.  Used to build corridors.
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
// Marks the given number of cells as dropped.  Assume they are already connected.
func (g *RoomGraph) DropRandomRooms(count int) {
	for i := 0; i < count; i++ {
		cell := g.RandCell(1)
		//debug.Add("Dropping room %d", cell)
		g.rooms[cell].mark = -1
		g.PruneDeadends(cell, 2)
	}
}

// ----------------------------------------------------------------------------
// Determines the boundaries of the area within each cell.  Used to randomize the rooms.
func (g *RoomGraph) MakeCellBounds() (areas []Room) {

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
// Create rooms in each cell with random size and location within the cell bounds.
// Assumes the bounds have already been created.
func (g *RoomGraph) MakeRandomRooms() {

	// make a random room within each area
	for i, a := range g.bounds {
		//randW := rand.Intn(12) + 8    // between 8 and 20
		randW := rand.Intn(a.W-5) + 5
		randH := rand.Intn(a.H-4) + 4 // between 4 and max height of area
		dx := rand.Intn(a.W - randW)  // position within the boundary area
		dy := rand.Intn(a.H - randH)
		g.rooms[i].SetSize(a.X+dx, a.Y+dy, randW, randH)
	}
}

// ----------------------------------------------------------------------------
// Gives a slice of the neighbouring cells of a given cell.
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
// For a given cell, checks if it is a deadend (just one corridor with no room attached)
// and removes the corridor if it is.
func (g *RoomGraph) PruneDeadends(cell int, depth int) {

	// check if the given cell is a dead end
	if g.rooms[cell].mark == -1 && g.CountCorridors(cell) == 1 {

		// if it is, remove all the corridors (should be just one)
		for i, p := range g.corridors {
			if (p.origID == cell || p.destID == cell) && p.mark != -1 {
				g.corridors[i].mark = -1
				//debug.Add("%d dropping corridor %v, cell=%d", depth, p, cell)

				// check to see if we just created another deadend
				if p.origID == cell {
					g.PruneDeadends(p.destID, depth-1)
				} else {
					g.PruneDeadends(p.origID, depth-1)
				}
			}
		}
	}
}

// ----------------------------------------------------------------------------
// Returns a cell chosen at random with the given mark or -1 if none are available.
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
// Returns a randomly chosen neighbour of the given cell with the given mark or
// -1 if none are avilable.
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

/*****************************************************************************/

type Room struct {
	X, Y int
	W, H int
	mark int // 0=unconnected, 1=connected, -1=dropped
}

// Returns the screen coord of the room's center
func (r Room) Center() (x int, y int) { //TODO
	x = r.X + r.W/2
	y = r.Y + r.H/2
	return
}

func (r Room) TopLeft() Coord {
	return Coord{r.X, r.Y}
}

// Returns a random point within the room ensuring it's not on a wall
func (r Room) RandPoint() Coord {
	x := r.X + rand.Intn(r.W-2) + 1
	y := r.Y + rand.Intn(r.H-2) + 1
	return Coord{x, y}
}

// Returns the coord of a random point on the wall of the given direction
func (r Room) RandWallPoint(dir Direction) (x, y int) {
	x, y = r.RandPoint().XY()
	switch dir {
	case North:
		y = r.Y
	case South:
		y = r.Y + r.H - 1
	case East:
		x = r.X + r.W - 1
	case West:
		x = r.X
	}
	return
}

// Updates the dimensions of the room
func (r *Room) SetSize(x, y, w, h int) {
	r.X = x
	r.Y = y
	r.W = w
	r.H = h
}

// Returns true if the given x,y coord in within the bounds of the room
func (r *Room) InRoom(x, y int) bool {
	return r.X-1 < x &&
		x < r.X+r.W+1 &&
		r.Y-1 < y &&
		y < r.Y+r.H+1
}

/*****************************************************************************/

type Corridor struct {
	origID int
	destID int
	mark   int // 0=normal, -1=dropped
}
