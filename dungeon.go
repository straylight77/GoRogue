package main

const (
	MapMaxX, MapMaxY = 80, 23
)

// If set to true, draw paths without accouting for existing tiles
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
	TilePath
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
		TilePath,
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

type DungeonMap struct {
	tiles [MapMaxX][MapMaxY]Tile
	rooms []Room
}

// -----------------------------------------------------------------------
func (m *DungeonMap) Clear() {
	for x, col := range m.tiles {
		for y := range col {
			m.tiles[x][y] = Tile{typ: TileEmpty}
		}
	}
}

// -----------------------------------------------------------------------
func (m *DungeonMap) SetTile(x, y int, t TileType) {
	m.tiles[x][y] = Tile{typ: t}
}

// -----------------------------------------------------------------------
func (m *DungeonMap) TileAt(x, y int) Tile {
	return m.tiles[x][y]
}

// -----------------------------------------------------------------------
func (m *DungeonMap) TileTypeAt(x, y int) TileType {
	return m.tiles[x][y].typ
}

// -----------------------------------------------------------------------
func (m *DungeonMap) IsWalkableAt(x, y int) bool {
	return m.tiles[x][y].IsWalkable()
}

// -----------------------------------------------------------------------
func (d *DungeonMap) playerFOV(p *Player) {
	for x := p.X - 1; x <= p.X+1; x++ {
		for y := p.Y - 1; y <= p.Y+1; y++ {

			// Check what the player is currently standing on
			switch d.TileTypeAt(p.X, p.Y) {
			case TilePath, TileFloor, TileDoor:

				// If player is in a hallway, only light up paths or doors
				switch d.TileTypeAt(x, y) {
				case TilePath, TileDoor:
					d.tiles[x][y].visible = true
					d.tiles[x][y].visited = true
				}

			default:
				d.tiles[x][y].visible = true
				d.tiles[x][y].visited = true
			}
		}
	}
}

// -----------------------------------------------------------------------
func (d *DungeonMap) SetVisible(x1, y1, w, h int, val bool) {
	for x := x1; x < x1+w; x++ {
		for y := y1; y < y1+h; y++ {
			d.tiles[x][y].visible = val
			if val {
				// Any time we set the visible (on or off) we consider is visited
				d.tiles[x][y].visited = true
			}
		}
	}
}

// -----------------------------------------------------------------------
// Returns the Chebyshev Distance from the given Entity
func (d *DungeonMap) Distance(x1, y1, x2, y2 int) int {
	dx := abs(x2 - x1)
	dy := abs(y2 - y1)
	return max(dx, dy)
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
	m.ConvertTile(x2, y2, IgnoreTiles)
}

// -----------------------------------------------------------------------
func (m *DungeonMap) CreatePath(x1, y1 int, dir Direction, length int) (int, int) {

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
	if ignore {
		m.SetTile(x, y, TilePath)
	} else {
		switch m.TileAt(x, y).typ {
		case TileFloor: //don't overwrite floor tiles
		case TileWallH, TileWallV:
			m.SetTile(x, y, TileDoor)
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

	m.rooms = append(m.rooms, Room{X: x1, Y: y1, W: w, H: h})

	// return the coords of the room center
	return x1 + (w / 2), y1 + (h / 2)
}

// -----------------------------------------------------------------------
func (m *DungeonMap) GenerateLevel(p *Player, ml *MonsterList) {

	m.Clear()
	ml.Clear()

	x1, y1 := m.CreateRoom(42, 3, 13, 5)
	x2, y2 := m.CreateRoom(25, 15, 11, 7)
	m.ConnectRooms(x1, y1, x2, y2, North)

	x3, y3 := m.CreateRoom(18, 2, 20, 7)
	m.ConnectRooms(x1, y1, x3, y3, East)

	m.SetTile(45, 5, TileStairsUp)
	m.SetTile(31, 18, TileStairsDn)
	//monsters.Add(randomMonster(player.depth), 20, 4)
	monsters.Add(randomMonster(player.depth), 50, 6)
	//monsters.Add(randomMonster(player.depth), 29, 17)

	p.SetPos(45, 5)
	p.depth++
}
