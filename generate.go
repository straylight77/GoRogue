package main

import (
	"fmt"
	"math/rand"
)

type Point struct {
	X, Y int
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

/*****************************************************************************/

type Path struct {
	origID int
	destID int
	//dir    Direction
}

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

	//logDebugMsg(fmt.Sprintf(" nb: id=%d, nbList=%v", id, nbList[id]))
	for _, nid := range nbList[id] {
		if g[nid].mark == mark {
			ids = append(ids, nid)
			//logDebugMsg(fmt.Sprintf("     nb: mark=%d nid=%d Y", mark, nid))
		} else {
			//logDebugMsg(fmt.Sprintf("     nb: mark=%d nid=%d -", mark, nid))
		}
	}
	if len(ids) == 0 {
		return -1
	}

	idx := rand.Intn(len(ids))
	return ids[idx]
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
// 0  1  2
// 3  4  5
// 6  7  8
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
func generateRandomLevel(d *DungeonMap, ml *MonsterList, p *Player) {
	clearDebugMsg()
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

	// CONNECT THE ROOMS:

	// 1. pick a random room (that's not dropped)
	id1 := roomGrid.getRandomRoom(0)
	logDebugMsg(fmt.Sprintf("first: %d", id1))

	// 2. connect it to one of its neighbours
	id2 := roomGrid.getRandomNeighbour(id1, 0)
	logDebugMsg(fmt.Sprintf("next: %d", id2))
	connectRooms(id1, id2)

	count := 0
	id := 0
	destId := 0

	// 3. while there are unconnected rooms
	// 3a. pick one at random
	id = roomGrid.getRandomRoom(0)
	//logDebugMsg(fmt.Sprintf("next: %d", id))

	for id != -1 && count < 20 {

		// 3b. connect it to a neighbour already connected
		destId = roomGrid.getRandomNeighbour(id, 1)
		if destId != -1 {
			logDebugMsg(fmt.Sprintf("checking: %d ->  %d connected  (count=%d)", id, destId, count))
			connectRooms(id, destId)
		} else {
			//    c. if there are none, skip
			logDebugMsg(fmt.Sprintf("checking: %d -> %d   skipped  (count=%d)", id, destId, count))
		}

		// 3. while there are unconnected rooms
		// 3a. pick one at random
		id = roomGrid.getRandomRoom(0)
		//logDebugMsg(fmt.Sprintf("next: %d", id))
		count++
	}

	dropRandomRooms(2)

	//TODO:
	// - remove deadends (paths that have exactly one dropped room)
	// - set the direction of paths e.g. North-South, not just East

	// actually create the rooms on the map
	for idx, r := range roomGrid {
		if roomGrid[idx].mark == 1 {
			d.CreateRoom(r.X, r.Y, r.W, r.H)
		} else if roomGrid[idx].mark == -1 {
			// bug workaround: for missing rooms, paths aren't connected in some cases
			pt := roomGrid[idx].Center()
			d.SetTile(pt.X, pt.Y, TilePath)
		}
	}

	// actually create the paths on the map
	for _, p := range pathList {
		pt1 := roomGrid[p.origID].Center()
		pt2 := roomGrid[p.destID].Center()
		d.ConnectRooms(pt1.X, pt1.Y, pt2.X, pt2.Y, East)
	}

	//for _, r := range roomGrid {
	//	pt := r.Center()
	//	d.SetTile(pt.X, pt.Y, TileStairsDn)
	//}

}

// ----------------------------------------------------------------------------

var debugMessages []string

func logDebugMsg(msg string) {
	debugMessages = append(debugMessages, msg)
}
func clearDebugMsg() {
	debugMessages = nil
}

// ----------------------------------------------------------------------------
// called from main()
func drawGenerateDebug(disp *Display) {

	debugMapGrid(disp)

	for i := 0; i < len(roomGrid); i++ {
		info := fmt.Sprintf("%d: %v", i, roomGrid[i])
		disp.DrawDebug(0, 28+i, info)
	}

	for i, lst := range nbList {
		disp.DrawDebug(20, 28+i, fmt.Sprint(lst))
	}

	for i, p := range pathList {
		disp.DrawDebug(35, 28+i, fmt.Sprint(p))
	}

	for i, msg := range debugMessages {
		disp.DrawDebug(84, 5+i, msg)
	}

}

// ----------------------------------------------------------------------------
func debugMapGrid(disp *Display) {
	disp.DrawHLine(8, 0, 79, disp.DebugStyle)
	disp.DrawHLine(16, 0, 79, disp.DebugStyle)
	disp.DrawVLine(26, 1, 24, disp.DebugStyle)
	disp.DrawVLine(53, 1, 24, disp.DebugStyle)

	disp.DrawDebug(0, 1, "0")
	disp.DrawDebug(27, 1, "1")
	disp.DrawDebug(54, 1, "2")
	disp.DrawDebug(0, 9, "3")
	disp.DrawDebug(27, 9, "4")
	disp.DrawDebug(54, 9, "5")
	disp.DrawDebug(0, 17, "6")
	disp.DrawDebug(27, 17, "7")
	disp.DrawDebug(54, 17, "8")
}