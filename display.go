package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/gdamore/tcell/v2"
)

var TileRunes = map[TileType]rune{
	TileEmpty:    ' ',
	TileWallH:    '-',
	TileWallV:    '|',
	TileWallUL:   '-',
	TileWallUR:   '-',
	TileWallLL:   '-',
	TileWallLR:   '-',
	TileFloor:    '.',
	TilePath:     '#',
	TileDoorOp:   '`',
	TileDoorCl:   '+',
	TileStairsUp: '<',
	TileStairsDn: '>',
}

type Display struct {
	Screen     tcell.Screen
	DefStyle   tcell.Style
	DebugStyle tcell.Style
}

// -----------------------------------------------------------------------------
func (d *Display) Init() {
	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := s.Init(); err != nil {
		log.Fatalf("%+v", err)
	}

	d.DefStyle = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
	d.DebugStyle = tcell.StyleDefault.Foreground(tcell.ColorLightSkyBlue)

	s.SetStyle(d.DefStyle)
	s.Clear()
	d.Screen = s
}

// -----------------------------------------------------------------------------
func (d *Display) Quit() {
	maybePanic := recover()
	d.Screen.Fini()
	if maybePanic != nil {
		panic(maybePanic)
	}
}

// -----------------------------------------------------------------------------
func (d *Display) GetCommand() tcell.Key {

	ev := d.Screen.PollEvent()

	var cmd tcell.Key // game command that the main loop will handle

	// Process event
	switch ev := ev.(type) {
	case *tcell.EventResize:
		d.Screen.Clear()
		d.Screen.Sync()
	case *tcell.EventKey:
		if ev.Key() == tcell.KeyCtrlR {
			d.Screen.Sync()
		} else {
			cmd = ev.Key()
		}
	}
	return cmd
}

// -----------------------------------------------------------------------------
func (d *Display) DrawEntity(e Entity) {
	x, y := e.Pos()
	d.Screen.SetContent(x, y, e.Rune(), nil, d.DefStyle)
}

// -----------------------------------------------------------------------------
func (d *Display) DrawMap(m *DungeonMap) {
	for x, col := range m {
		for y, t := range col {
			r := TileRunes[t.typ]
			d.Screen.SetContent(x, y, r, nil, d.DefStyle)
		}
	}
}

// -----------------------------------------------------------------------------
func (d *Display) DrawMessages(messages []string) {
	if len(messages) > 0 {
		entireStr := strings.Join(messages, " ")
		drawTextWrap(d.Screen, 0, 0, 80, 3, d.DefStyle, entireStr)
		clearMessages()
	}
}

// -----------------------------------------------------------------------------
func (d *Display) DrawText(x1, y1 int, text string) {
	row := y1
	col := x1
	for _, r := range []rune(text) {
		d.Screen.SetContent(col, row, r, nil, d.DefStyle)
		col++
	}
}

// -----------------------------------------------------------------------------
func (d *Display) DrawDebug(p *Player, ml *MonsterList) {
	maxX, maxY := 80, 25
	for x := 0; x < maxX; x++ {
		d.Screen.SetContent(x, maxY, tcell.RuneHLine, nil, d.DebugStyle)
	}
	for y := 0; y < maxY; y++ {
		d.Screen.SetContent(maxX, y, tcell.RuneVLine, nil, d.DebugStyle)
	}
	d.Screen.SetContent(maxX, maxY, tcell.RuneLRCorner, nil, d.DebugStyle)

	//texth := "012345678901234567890123456789012345678901234567890123456789012345678901234567890"
	//drawTextWrap(d.Screen, 0, 26, 81, 26, d.DebugStyle, texth)
	//textv := "01234567890123456789012345"
	//drawTextWrap(d.Screen, 81, 0, 82, 27, d.DebugStyle, textv)

	drawTextWrap(d.Screen, 84, 1, 200, 1, d.DebugStyle, fmt.Sprintf("Moves:  %d", p.moves))
	drawTextWrap(d.Screen, 84, 2, 200, 2, d.DebugStyle, fmt.Sprintf("Player: %d, %d", p.X, p.Y))
	for i, m := range *ml {
		msg := fmt.Sprintf("%d: %v", i, m.DebugString())
		drawTextWrap(d.Screen, 84, 4+i, 200, 4+i, d.DebugStyle, msg)
	}
}

// ============================================================================

func drawTextWrap(s tcell.Screen, x1, y1, x2, y2 int, style tcell.Style, text string) {
	row := y1
	col := x1
	for _, r := range []rune(text) {
		s.SetContent(col, row, r, nil, style)
		col++
		if col >= x2 {
			row++
			col = x1
		}
		if row > y2 {
			break
		}
	}
}
