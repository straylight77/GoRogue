package main

const (
	MapMaxX, MapMaxY = 80, 24
)

const (
	//TileWallH  = tcell.RuneHLine
	//TileWallV  = tcell.RuneVLine
	//TileWallUL = tcell.RuneULCorner
	//TileWallUR = tcell.RuneURCorner
	//TileWallLL = tcell.RuneLLCorner
	//TileWallLR = tcell.RuneLRCorner
	//TileFloor  = tcell.RuneBullet
	//TilePath   = tcell.RuneBoard

	TileWallH  = '-'
	TileWallV  = '|'
	TileWallUL = '-'
	TileWallUR = '-'
	TileWallLL = '-'
	TileWallLR = '-'
	TileFloor  = '.'
	TilePath   = '#'

	TileDoorOp = '`'
	TileDoorCl = '+'
)

const (
	NORTH = iota
	EAST
	SOUTH
	WEST
)

type DungeonTile struct {
	sym     rune
	blocks  bool
	visible bool
}

type DungeonMap [MapMaxX][MapMaxY]rune

// -----------------------------------------------------------------------
func (m *DungeonMap) Clear() {
	for x, col := range m {
		for y := range col {
			m[x][y] = ' '
		}
	}
}

// -----------------------------------------------------------------------
func (m *DungeonMap) SetTile(x, y int, r rune) {
	m[x][y] = r
}

// -----------------------------------------------------------------------
func (m *DungeonMap) Tile(x, y int) rune {
	return m[x][y]
}

// -----------------------------------------------------------------------
func (m *DungeonMap) GenerateLevel(lvl int, p *Player) {
	var x, y int

	m.Clear()
	x, y = m.CreateRoom(42, 3, 12, 8)
	x, y = m.CreatePath(x, y, WEST, 15)
	m.CreateRoom(27, 15, 10, 6)
	x, y = m.CreatePath(x, y, SOUTH, 10)

	p.SetPos(45, 5)
	p.depth++
}

// -----------------------------------------------------------------------
func (m *DungeonMap) CreatePath(x1, y1 int, dir int, length int) (int, int) {
	dx, dy := getDirectionCoords(dir)
	x, y := x1, y1
	for i := length; i > 0; i-- {

		switch m[x][y] {
		case TileFloor:
			//ignore floor tiles
		case TileWallH, TileWallV:
			m.SetTile(x, y, '+')
		default:
			m[x][y] = TilePath
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
		m[x][y1] = TileWallH
		m[x][y1+h] = TileWallH
	}

	for y := y1; y < y1+h; y++ {
		m[x1][y] = TileWallV
		m[x1+w][y] = TileWallV
	}

	for x := x1 + 1; x < x1+w; x++ {
		for y := y1 + 1; y < y1+h; y++ {
			m[x][y] = TileFloor
		}
	}

	m[x1][y1] = TileWallUL
	m[x1+w][y1] = TileWallUR
	m[x1][y1+h] = TileWallLL
	m[x1+w][y1+h] = TileWallLR

	// return the coords of the room center
	return x1 + (w / 2), y1 + (h / 2)
}

// ============================================================================

func getDirectionCoords(dir int) (int, int) {
	dx, dy := 0, 0
	switch dir {
	case NORTH:
		dy = -1
	case SOUTH:
		dy = 1
	case EAST:
		dx = 1
	case WEST:
		dx = -1
	}
	return dx, dy
}
