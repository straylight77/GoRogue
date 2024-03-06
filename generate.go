package main

import "math/rand"

type Point struct {
	X, Y int
}

// ----------------------------------------------------------------------------
type Rect struct {
	X, Y int
	W, H int
}

func (r Rect) Center() Point {
	return Point{r.X + r.W/2, r.Y + r.H/2}
}

func (r Rect) RandPoint() Point {
	return Point{
		r.X + rand.Intn(r.W-2) + 1,
		r.Y + rand.Intn(r.H-2) + 1,
	}
}

// ----------------------------------------------------------------------------
func randomRooms() (rooms []Rect) {
	roomW := (MapMaxX - 2) / 3 // 25
	roomH := (MapMaxY - 2) / 3 // 6

	var areas []Rect

	// split the map into 3x3 areas determine the bounds of each one
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			r := Rect{(roomW + 1) * j, (roomH + 1) * i, roomW, roomH}
			areas = append(areas, r)
		}
	}

	// make a random room within each area
	for _, a := range areas {
		randW := rand.Intn(12) + 8
		randH := rand.Intn(roomH-4) + 4
		dx := rand.Intn(roomW - randW)
		dy := rand.Intn(roomH - randH)
		r := Rect{a.X + dx, a.Y + dy, randW, randH}
		rooms = append(rooms, r)
	}

	// drop a few rooms
	for i := 0; i < 2; i++ {
		idx := rand.Intn(len(rooms))
		rooms = append(rooms[:idx], rooms[idx+1:]...)
	}

	return
}

// ----------------------------------------------------------------------------
func generateRandomLevel(d *DungeonMap, ml *MonsterList, p *Player) {
	rooms := randomRooms()

	d.Clear()
	for _, r := range rooms {
		d.CreateRoom(r.X, r.Y, r.W, r.H)
		//pt := r.Center()
		pt := r.RandPoint()
		d.SetTile(pt.X, pt.Y, TileStairsUp)
	}
}

/*********************************************************************
 *********************************************************************
 */
func generateRandomLevel2(d *DungeonMap, ml *MonsterList, p *Player) {
	p.SetPos(2, 3)

	// assuming 3x3 rooms
	roomW := (MapMaxX - 2) / 3
	roomH := (MapMaxY - 2) / 3

	type XY struct{ X, Y int }
	var centers [3][3]XY
	var centerPairs [9][2]XY
	var x, y int

	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			x, y = d.CreateRoom((roomW+1)*j, (roomH+1)*i, roomW, roomH)
			//d.SetTile(x, y, TileStairsDn)
			centers[i][j] = XY{x, y}
		}
	}

	for i := 0; i < 3; i++ {
		for j := 0; j < 2; j++ {
			c1 := centers[i][j]
			c2 := centers[i][j+1]
			centerPairs[i+j][0] = c1
			centerPairs[i+j][1] = c2

			//d.CreatePath(c1.X, c1.Y, East, c2.X-c1.X)
		}
	}

	for i := 0; i < 2; i++ {
		for j := 0; j < 3; j++ {
			c1 := centers[i][j]
			c2 := centers[i+1][j]
			centerPairs[5+i+j][0] = c1
			centerPairs[5+i+j][1] = c2

			//d.CreatePath(c1.X, c1.Y, South, c2.Y-c1.Y)
		}
	}

	//for _, p := range centerPairs {
	//	d.SetTile(p[0].X, p[0].Y, TileStairsUp)
	//	d.SetTile(p[1].X, p[1].Y, TileStairsUp)
	//}

}
