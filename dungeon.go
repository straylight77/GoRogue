package main

import (
	"fmt"
	"math/rand"
)

const (
	MapMaxX, MapMaxY = 80, 23
)

// If set to true, draw corridors without accouting for existing tiles
var IgnoreTiles = false

// -----------------------------------------------------------------------
type TileType int

const (
	TileEmpty TileType = iota
	TileWallH
	TileWallV
	TileWallUL
	TileWallUR
	TileWallLL
	TileWallLR
	TileFloor
	TileCorridor
	TileDoor
	TileStairsDn
	TileStairsUp
)

// -----------------------------------------------------------------------
type Tile struct {
	typ     TileType
	visible bool
	visited bool
}

func (t *Tile) IsWalkable() bool {
	switch t.typ {
	case TileFloor, // consider these tiles as "walkable"
		TileCorridor,
		TileDoor,
		TileStairsDn,
		TileStairsUp:
		return true
	default:
		return false
	}
}

func (t *Tile) IsType(t2 TileType) bool {
	return t.typ == t2
}

/************************************************************************/

type Coord struct {
	X, Y int
}

func (c Coord) String() string {
	return fmt.Sprintf("(%d,%d)", c.X, c.Y)
}

func (c Coord) XY() (int, int) {
	return c.X, c.Y
}

func (from Coord) IsDiagonal(to Coord) bool {
	dx := to.X - from.X
	dy := to.Y - from.Y
	return dx != 0 && dy != 0
}

func (from Coord) Distance(to Coord) int {
	dx := abs(to.X - from.X)
	dy := abs(to.Y - from.Y)
	return max(dx, dy)
}

func (c1 Coord) Sum(c2 Coord) Coord {
	return Coord{c1.X + c2.X, c1.Y + c2.Y}
}

func (c1 Coord) Diff(c2 Coord) Coord {
	return Coord{c1.X - c2.X, c1.Y - c2.Y}
}

/************************************************************************/

type DungeonMap struct {
	tiles [MapMaxX][MapMaxY]Tile
	rooms []Room
}

// -----------------------------------------------------------------------
func (m *DungeonMap) Clear() {
	for x, col := range m.tiles {
		for y := range col {
			m.tiles[x][y] = Tile{typ: TileEmpty}
			m.rooms = nil
		}
	}
}

// -----------------------------------------------------------------------
func (m *DungeonMap) SetTile(pos Coord, t TileType) {
	m.tiles[pos.X][pos.Y] = Tile{typ: t}
}

// -----------------------------------------------------------------------
func (m *DungeonMap) TileAt(pos Coord) Tile {
	return m.tiles[pos.X][pos.Y]
}

// -----------------------------------------------------------------------
func (m *DungeonMap) TileTypeAt(pos Coord) TileType {
	return m.tiles[pos.X][pos.Y].typ
}

// -----------------------------------------------------------------------
func (m *DungeonMap) IsOutOfBounds(pos Coord) bool {
	return pos.X < 0 || pos.X >= MapMaxX || pos.Y < 0 || pos.Y >= MapMaxY
}

// -----------------------------------------------------------------------
func (m *DungeonMap) IsWalkableAt(pos Coord) bool {
	if m.IsOutOfBounds(pos) {
		return false
	}
	return m.tiles[pos.X][pos.Y].IsWalkable()
}

// -----------------------------------------------------------------------
// Prevent diagonal movement through doors and cooridors
func (m *DungeonMap) IsWalkable(from, to Coord) bool {

	walkable := m.IsWalkableAt(to)

	if from.IsDiagonal(to) {
		if m.TileTypeAt(to) == TileDoor ||
			m.TileTypeAt(to) == TileCorridor ||
			m.TileTypeAt(from) == TileDoor ||
			m.TileTypeAt(from) == TileCorridor {
			walkable = false
		}
	}
	return walkable
}

// -----------------------------------------------------------------------
// Returns a random direction, as the delta in coordinates (dx, dy), that
// are always walkable or 0,0 if there are no options available.
func (m *DungeonMap) RandDirectionCoords(orig Coord) Coord {
	var cList []Coord
	for x := -1; x <= 1; x++ {
		for y := -1; y <= 1; y++ {
			dest := orig.Sum(Coord{x, y})
			if dest != orig && m.IsWalkable(orig, dest) {
				cList = append(cList, Coord{x, y})
			}
		}
	}

	if len(cList) > 0 {
		idx := rand.Intn(len(cList))
		return cList[idx]
	} else {
		return Coord{0, 0}
	}
}

// -----------------------------------------------------------------------
func (dm *DungeonMap) getWalkableNeighbours(pos Coord) []Coord {
	toCheck := []Coord{
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
	var final []Coord
	for _, c := range toCheck {
		if dm.IsWalkable(pos, c) {
			final = append(final, c)
		}
	}
	return final
}

// -----------------------------------------------------------------------
func (d *DungeonMap) CanSee(e Entity) bool {
	return d.TileAt(e.Pos()).visible
}

// -----------------------------------------------------------------------
// TODO: implement an actual line-of-sight,raytacing algorithm here
func (d *DungeonMap) playerFOV(pos Coord) {
	radius := 1
	for x := pos.X - radius; x <= pos.X+radius; x++ {
		for y := pos.Y - radius; y <= pos.Y+radius; y++ {

			// Check what the player is currently standing on
			switch d.TileTypeAt(pos) {

			// If the player is not in a room...
			case TileCorridor, TileDoor:

				switch d.TileTypeAt(Coord{x, y}) {
				//... only light up corridors, doors and floors
				case TileCorridor, TileDoor, TileFloor:
					d.tiles[x][y].visible = true
					d.tiles[x][y].visited = true
				}

			default:
				// Otherwise, light up everything
				d.tiles[x][y].visible = true
				d.tiles[x][y].visited = true
			}
		}
	}
}

// -----------------------------------------------------------------------
func (d *DungeonMap) SetVisible(start Coord, w, h int, val bool) {
	for x := start.X; x < start.X+w; x++ {
		for y := start.Y; y < start.Y+h; y++ {
			d.tiles[x][y].visible = val
			if val {
				// Any time we set a tile visible consider it visited
				d.tiles[x][y].visited = true
			}
		}
	}
}

// -----------------------------------------------------------------------
func (m *DungeonMap) ConnectRooms(x1, y1 int, x2, y2 int, startDir Direction) { //TODO
	HDir := East
	VDir := South

	if x2 < x1 {
		HDir = West
	}
	if y2 < y1 {
		VDir = North
	}

	dx := x2 - x1
	dy := y2 - y1

	var x, y int
	switch startDir {
	case North, South:
		seg1Len := dy / 2
		seg3Len := dy - seg1Len
		x, y = m.CreateCorridor(x1, y1, VDir, seg1Len)
		x, y = m.CreateCorridor(x, y, HDir, dx)
		x, y = m.CreateCorridor(x, y, VDir, seg3Len)
	case East, West:
		seg1Len := dx / 2
		seg3Len := dx - seg1Len
		x, y = m.CreateCorridor(x1, y1, HDir, seg1Len)
		x, y = m.CreateCorridor(x, y, VDir, dy)
		x, y = m.CreateCorridor(x, y, HDir, seg3Len)
	}
	m.ConvertTile(x2, y2, IgnoreTiles)
}

// -----------------------------------------------------------------------
func (m *DungeonMap) CreateCorridor(x1, y1 int, dir Direction, length int) (int, int) {

	//allow length to be given as negative
	if length < 0 {
		length = -1 * length
	}

	dx, dy := getDirectionCoords(dir)
	x, y := x1, y1
	for i := length; i > 0; i-- {
		m.ConvertTile(x, y, IgnoreTiles)
		x += dx
		y += dy
	}
	return x, y
}

// -----------------------------------------------------------------------
func (m *DungeonMap) ConvertTile(x, y int, ignore bool) {
	pos := Coord{x, y}
	if ignore {
		m.SetTile(pos, TileCorridor)
	} else {
		switch m.TileTypeAt(pos) {
		case TileFloor: //don't overwrite floor tiles
		case TileWallH, TileWallV:
			m.SetTile(pos, TileDoor)
		default:
			m.SetTile(pos, TileCorridor)
		}
	}

}

// -----------------------------------------------------------------------
func (m *DungeonMap) CreateRoom(x1, y1 int, w, h int) (int, int) {
	h -= 1
	w -= 1

	for x := x1; x < x1+w; x++ {
		m.SetTile(Coord{x, y1}, TileWallH)
		m.SetTile(Coord{x, y1 + h}, TileWallH)
	}

	for y := y1; y < y1+h; y++ {
		m.SetTile(Coord{x1, y}, TileWallV)
		m.SetTile(Coord{x1 + w, y}, TileWallV)
	}

	for x := x1 + 1; x < x1+w; x++ {
		for y := y1 + 1; y < y1+h; y++ {
			m.SetTile(Coord{x, y}, TileFloor)
		}
	}

	m.SetTile(Coord{x1, y1}, TileWallUL)
	m.SetTile(Coord{x1 + w, y1}, TileWallUR)
	m.SetTile(Coord{x1, y1 + h}, TileWallLL)
	m.SetTile(Coord{x1 + w, y1 + h}, TileWallLR)

	m.rooms = append(m.rooms, Room{X: x1, Y: y1, W: w, H: h})

	// return the coords of the room center
	return x1 + (w / 2), y1 + (h / 2)
}
