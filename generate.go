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
	//mark int // 0=unconnected, 1=connected, -1=dropped
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
func getAreaBounds(roomW, roomH int) (areas []Room) {

	// split the map into 3x3 areas determine the bounds of each one
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			r := Room{(roomW + 1) * j, (roomH + 1) * i, roomW, roomH}
			areas = append(areas, r)
		}
	}
	return
}

var mark [9]int

// ----------------------------------------------------------------------------
func getRandomRooms(bounds []Room, roomW, roomH int) (rooms []Room) {

	// make a random room within each area
	for _, a := range bounds {
		randW := rand.Intn(12) + 8
		randH := rand.Intn(roomH-4) + 4
		dx := rand.Intn(roomW - randW)
		dy := rand.Intn(roomH - randH)
		r := Room{a.X + dx, a.Y + dy, randW, randH}
		rooms = append(rooms, r)
	}

	// drop a few rooms
	mark = [9]int{}
	for i := 0; i < 2; i++ {
		idx := rand.Intn(len(rooms))
		if mark[idx] >= 0 {
			mark[idx] = -1
		} else {
			i-- // if room has already been dropped, choose again
		}
	}

	disp.DrawDebug(0, 28, fmt.Sprint(mark))

	return
}

// ----------------------------------------------------------------------------
func generateRandomLevel(d *DungeonMap, ml *MonsterList, p *Player) {
	d.Clear()

	// assume 3x3 rooms on the map
	roomW := (MapMaxX - 2) / 3 // 25
	roomH := (MapMaxY - 2) / 3 // 6

	areas := getAreaBounds(roomW, roomH)
	rooms := getRandomRooms(areas, roomW, roomH)

	// create the rooms on the map
	for idx, r := range rooms {
		if mark[idx] >= 0 {
			d.CreateRoom(r.X, r.Y, r.W, r.H)
			//pt := r.Center()
			pt := r.RandPoint()
			d.SetTile(pt.X, pt.Y, TileStairsUp)
		}
	}

	disp.Show()
}
