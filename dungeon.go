package main

const (
	MapMaxX, MapMaxY = 80, 24
)

// -----------------------------------------------------------------------
type Direction int

const (
	North Direction = iota
	East
	South
	West
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

/***
 * Dungeon Map
 *
 */
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
func (m *DungeonMap) GenerateLevel(lvl int, p *Player, ml *MonsterList) {
	var x, y int

	m.Clear()
	x, y = m.CreateRoom(42, 3, 12, 8)
	x, y = m.CreatePath(x, y, West, 15)
	m.CreateRoom(27, 15, 10, 6)
	x, y = m.CreatePath(x, y, South, 10)

	m.SetTile(45, 5, TileStairsUp)
	m.SetTile(31, 18, TileStairsDn)

	monsters.Add(NewMonster(0), 50, 8)
	monsters.Add(NewMonster(1), 29, 17)

	p.SetPos(45, 5)
	p.depth++
}

// -----------------------------------------------------------------------
func (m *DungeonMap) CreatePath(x1, y1 int, dir Direction, length int) (int, int) {
	dx, dy := getDirectionCoords(dir)
	x, y := x1, y1
	for i := length; i > 0; i-- {

		switch m.TileAt(x, y).typ {
		case TileFloor:
			//ignore floor tiles
		case TileWallH, TileWallV:
			m.SetTile(x, y, TileDoorCl)
		default:
			m.SetTile(x, y, TilePath)
		}
		x += dx
		y += dy
	}
	return x, y
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

// ============================================================================

func getDirectionCoords(dir Direction) (int, int) {
	dx, dy := 0, 0
	switch dir {
	case North:
		dy = -1
	case South:
		dy = 1
	case East:
		dx = 1
	case West:
		dx = -1
	}
	return dx, dy
}
