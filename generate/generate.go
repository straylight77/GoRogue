package main

import "fmt"

type Rect struct {
	X1, Y1 int
	W, H   int
}

func NewRandomRect(x1, y1 int, x2, y2 int) *Rect {
	// create a random rect that fits within the given bounds
	return nil
}

func (r *Rect) SetPos(x, y int) {
	r.X1 = x
	r.Y1 = y
}

var rooms []Rect

func createRandomRooms() {
	roomW := (80 - 3) / 3
	roomH := (23 - 3) / 3

	fmt.Printf("begin: roomW=%d, roomH=%d\n", roomW, roomH)
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			r := Rect{(roomW + 1) * j, (roomH + 1) * i, roomW, roomH}
			fmt.Printf("%d,%d : %v\n", i, j, r)
			rooms = append(rooms, r)
		}
	}

}

// func generateRandomLevel(d *DungeonMap, ml *MonsterList, p *Player) {
func generateRandomLevel() {
	// assuming 3x3 rooms

	roomW := (80 - 3) / 3
	roomH := (24 - 3) / 3

	type XY struct{ X, Y int }
	var centers [3][3]XY
	var centerPairs [9][2]XY
	var x, y int

	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			//x, y = d.CreateRoom(1+(roomW+1)*j, 1+(roomH+1)*i, roomW, roomH)
			//d.SetTile(x, y, TileStairsDn)
			rooms = append(rooms, Rect{(roomW + 1) * j, (roomH + 1) * i, roomW, roomH})
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

func main() {
	fmt.Println("Hi.")
	createRandomRooms()
	fmt.Println(rooms)
}
