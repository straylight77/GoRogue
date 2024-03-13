package main

const (
	MapMaxX, MapMaxY = 80, 23
)

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
	TilePath
	TileDoorCl
	TileDoorOp
	TileStairsDn
	TileStairsUp
)

// -----------------------------------------------------------------------
type Tile struct {
	typ     TileType
	visible bool
}

func (t *Tile) IsWalkable() bool {
	switch t.typ {
	case TileFloor, // consider these tiles as "walkable"
		TilePath,
		TileDoorOp,
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

type DungeonMap [MapMaxX][MapMaxY]Tile

// -----------------------------------------------------------------------
func (m *DungeonMap) Clear() {
	for x, col := range m {
		for y := range col {
			m[x][y] = Tile{typ: TileEmpty}
		}
	}
}

// -----------------------------------------------------------------------
func (m *DungeonMap) SetTile(x, y int, t TileType) {
	m[x][y] = Tile{typ: t}
}

// -----------------------------------------------------------------------
func (m *DungeonMap) TileAt(x, y int) Tile {
	return m[x][y]
}

// -----------------------------------------------------------------------
func (m *DungeonMap) IsWalkableAt(x, y int) bool {
	return m[x][y].IsWalkable()
}

// -----------------------------------------------------------------------
// assume x1 < x2 and y1 < y2
func (m *DungeonMap) ConnectRooms(x1, y1 int, x2, y2 int, startDir Direction) {
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
		x, y = m.CreatePath(x1, y1, VDir, seg1Len)
		x, y = m.CreatePath(x, y, HDir, dx)
		x, y = m.CreatePath(x, y, VDir, seg3Len)
	case East, West:
		seg1Len := dx / 2
		seg3Len := dx - seg1Len
		x, y = m.CreatePath(x1, y1, HDir, seg1Len)
		x, y = m.CreatePath(x, y, VDir, dy)
		x, y = m.CreatePath(x, y, HDir, seg3Len)
	}
	m.ConvertTile(x2, y2, false)
}

// -----------------------------------------------------------------------
func (m *DungeonMap) CreatePath(x1, y1 int, dir Direction, length int) (int, int) {

	//debug.Add("path: x1=%d, y1=%d, len=%d", x1, y1, length)

	//allow length to be given as negative
	if length < 0 {
		length = -1 * length
	}

	dx, dy := getDirectionCoords(dir)
	x, y := x1, y1
	for i := length; i > 0; i-- {
		m.ConvertTile(x, y, false)
		x += dx
		y += dy
	}
	return x, y
}

// -----------------------------------------------------------------------
func (m *DungeonMap) ConvertTile(x, y int, ignore bool) {
	if ignore {
		m.SetTile(x, y, TilePath)
	} else {
		switch m.TileAt(x, y).typ {
		case TileFloor: //don't overwrite floor tiles
		case TileWallH, TileWallV:
			m.SetTile(x, y, TileDoorCl)
		default:
			m.SetTile(x, y, TilePath)
		}
	}

}

// -----------------------------------------------------------------------
func (m *DungeonMap) CreateRoom(x1, y1 int, w, h int) (int, int) {
	h -= 1
	w -= 1

	for x := x1; x < x1+w; x++ {
		m.SetTile(x, y1, TileWallH)
		m.SetTile(x, y1+h, TileWallH)
	}

	for y := y1; y < y1+h; y++ {
		m.SetTile(x1, y, TileWallV)
		m.SetTile(x1+w, y, TileWallV)
	}

	for x := x1 + 1; x < x1+w; x++ {
		for y := y1 + 1; y < y1+h; y++ {
			m.SetTile(x, y, TileFloor)
		}
	}

	m.SetTile(x1, y1, TileWallUL)
	m.SetTile(x1+w, y1, TileWallUR)
	m.SetTile(x1, y1+h, TileWallLL)
	m.SetTile(x1+w, y1+h, TileWallLR)

	// return the coords of the room center
	return x1 + (w / 2), y1 + (h / 2)
}

// -----------------------------------------------------------------------
func (m *DungeonMap) GenerateLevel(lvl int, p *Player, ml *MonsterList) {

	m.Clear()
	x1, y1 := m.CreateRoom(42, 16, 13, 5)
	x2, y2 := m.CreateRoom(7, 1, 11, 7)

	//debug.Add("gen: x1=%d, y1=%d, x2=%d, y2=%d", x1, y1, x2, y2)
	m.ConnectRooms(x1, y1, x2, y2, North)

	//m.SetTile(45, 5, TileStairsUp)
	//m.SetTile(31, 18, TileStairsDn)
	//monsters.Add(NewMonster(0), 50, 8)
	//monsters.Add(NewMonster(1), 29, 17)

	p.SetPos(9, 3)
	p.depth++
}
