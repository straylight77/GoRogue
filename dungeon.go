package main

import "github.com/gdamore/tcell/v2"

const (
	MapMaxX, MapMaxY = 80, 24
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
func (m *DungeonMap) CreateRoom(x1, y1 int, w, h int) {
	h -= 1
	w -= 1

	for x := x1; x < x1+w; x++ {
		m[x][y1] = tcell.RuneHLine
		m[x][y1+h] = tcell.RuneHLine
	}

	for y := y1; y < y1+h; y++ {
		m[x1][y] = tcell.RuneVLine
		m[x1+w][y] = tcell.RuneVLine
	}

	for x := x1 + 1; x < x1+w; x++ {
		for y := y1 + 1; y < y1+h; y++ {
			m[x][y] = tcell.RuneBullet
		}
	}

	m[x1][y1] = tcell.RuneULCorner
	m[x1+w][y1] = tcell.RuneURCorner
	m[x1][y1+h] = tcell.RuneLLCorner
	m[x1+w][y1+h] = tcell.RuneLRCorner

}
