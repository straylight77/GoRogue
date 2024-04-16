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
func (d *DungeonMap) CanSee(a Actor) bool {
	return d.TileAt(a.Pos()).visible
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
func (m *DungeonMap) ConnectRooms(p1, p2 Coord, startDir Direction) {
	HDir := East
	VDir := South

	if p2.X < p1.X {
		HDir = West
	}
	if p2.Y < p1.Y {
		VDir = North
	}

	dx := p2.X - p1.X
	dy := p2.Y - p1.Y

	var next Coord

	switch startDir {
	case North, South:
		seg1Len := dy / 2
		seg3Len := dy - seg1Len
		next = m.CreateCorridor(p1, VDir, seg1Len)
		next = m.CreateCorridor(next, HDir, dx)
		next = m.CreateCorridor(next, VDir, seg3Len)
	case East, West:
		seg1Len := dx / 2
		seg3Len := dx - seg1Len
		next = m.CreateCorridor(p1, HDir, seg1Len)
		next = m.CreateCorridor(next, VDir, dy)
		next = m.CreateCorridor(next, HDir, seg3Len)
	}
	m.ConvertTile(p2, IgnoreTiles)
}

// -----------------------------------------------------------------------
func (m *DungeonMap) CreateCorridor(pos Coord, dir Direction, length int) Coord {

	//allow length to be given as negative
	if length < 0 {
		length = -1 * length
	}

	delta := getDirectionCoords(dir)
	for i := length; i > 0; i-- {
		m.ConvertTile(pos, IgnoreTiles)
		pos = pos.Sum(delta)
	}
	return pos
}

// -----------------------------------------------------------------------
func (m *DungeonMap) ConvertTile(pos Coord, ignore bool) {
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
func (m *DungeonMap) CreateRoom(pos Coord, w, h int) Coord {
	h -= 1
	w -= 1

	for x := pos.X; x < pos.X+w; x++ {
		m.SetTile(Coord{x, pos.Y}, TileWallH)
		m.SetTile(Coord{x, pos.Y + h}, TileWallH)
	}

	for y := pos.Y; y < pos.Y+h; y++ {
		m.SetTile(Coord{pos.X, y}, TileWallV)
		m.SetTile(Coord{pos.X + w, y}, TileWallV)
	}

	for x := pos.X + 1; x < pos.X+w; x++ {
		for y := pos.Y + 1; y < pos.Y+h; y++ {
			m.SetTile(Coord{x, y}, TileFloor)
		}
	}

	m.SetTile(pos, TileWallUL)
	m.SetTile(pos.Sum(Coord{w, 0}), TileWallUR)
	m.SetTile(pos.Sum(Coord{0, h}), TileWallLL)
	m.SetTile(pos.Sum(Coord{w, h}), TileWallLR)

	m.rooms = append(m.rooms, Room{X: pos.X, Y: pos.Y, W: w, H: h})

	// return the coords of the room center
	return pos.Sum(Coord{w / 2, h / 2})
}
