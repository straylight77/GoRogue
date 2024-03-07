package main

import (
	"fmt"
	"math/rand"
)

type Point struct {
	X, Y int
}

// ----------------------------------------------------------------------------
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

// ----------------------------------------------------------------------------
type RoomGrid [9]Room

func (g *RoomGrid) getNeighbours(idx int) []int {
	return []int{}
}

var roomGrid = RoomGrid{}

var nbList [9][]int

// ----------------------------------------------------------------------------
// called from main()
func drawGenerateDebug(disp *Display) {
	disp.DrawHLine(8, 0, 80, disp.DebugStyle)
	disp.DrawHLine(16, 0, 80, disp.DebugStyle)
	disp.DrawVLine(26, 1, 24, disp.DebugStyle)
	disp.DrawVLine(53, 1, 24, disp.DebugStyle)

	for i := 0; i < len(roomGrid); i++ {
		info := fmt.Sprintf("%d: %v", i, roomGrid[i])
		disp.DrawDebug(0, 28+i, info)
	}

	for i, lst := range nbList {
		disp.DrawDebug(20, 28+i, fmt.Sprint(lst))
	}

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

	// drop a few rooms
	for i := 0; i < 2; i++ {
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
		{1, 3, 4}, // room 0 neighbours
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
func generateRandomLevel(d *DungeonMap, ml *MonsterList, p *Player) {
	d.Clear()

	// assume 3x3 rooms on the map
	roomW := (MapMaxX - 2) / 3 // 25
	roomH := (MapMaxY - 2) / 3 // 6

	areas := getAreaBounds(roomW, roomH)
	makeRandomRooms(areas, roomW, roomH)

	// get the list of neighbours for each room
	for idx := range roomGrid {
		nbList[idx] = getNeighbours(idx, 3)
	}

	// create the rooms on the map
	for idx, r := range roomGrid {
		if roomGrid[idx].mark >= 0 {
			d.CreateRoom(r.X, r.Y, r.W, r.H)
		}
	}

	// CONNECT THE ROOMS:
	// 1. pick a random unconnected room
	// 2. connect it to one of its neighbours
	// 3. while there are unconnected rooms
	//    a. pick one at random
	//    b. connect it to a neighbour already connected
	//    c. if there are none, skip

}
