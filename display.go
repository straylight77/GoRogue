package main

import (
	"log"

	"github.com/gdamore/tcell/v2"
)

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
	// You have to catch panics in a defer, clean up, and
	// re-raise them - otherwise your application can
	// die without leaving any diagnostic trace.
	maybePanic := recover()
	d.Screen.Fini()
	if maybePanic != nil {
		panic(maybePanic)
	}
}

// -----------------------------------------------------------------------------
func (d *Display) Command() rune {

	ev := d.Screen.PollEvent()

	var cmd rune // game command that the main loop will handle

	// Process event
	switch ev := ev.(type) {
	case *tcell.EventResize:
		d.Screen.Clear()
		d.Screen.Sync()
	case *tcell.EventKey:
		if ev.Key() == tcell.KeyEscape {
			cmd = 'X'
		} else if ev.Key() == tcell.KeyCtrlR {
			d.Screen.Sync()
		} else {
			cmd = ev.Rune()
		}
	}
	return cmd
}

// -----------------------------------------------------------------------------
func (d *Display) DrawMap(m *DungeonMap) {
	for x, col := range m {
		for y, ch := range col {
			d.Screen.SetContent(x, y, ch, nil, d.DefStyle)
		}
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
func (d *Display) DrawDebug() {
	maxX, maxY := 80, 25
	for x := 0; x < maxX; x++ {
		d.Screen.SetContent(x, maxY, tcell.RuneHLine, nil, d.DebugStyle)
	}
	for y := 0; y < maxY; y++ {
		d.Screen.SetContent(maxX, y, tcell.RuneVLine, nil, d.DebugStyle)
	}
	d.Screen.SetContent(maxX, maxY, tcell.RuneLRCorner, nil, d.DebugStyle)

	texth := "012345678901234567890123456789012345678901234567890123456789012345678901234567890"
	drawTextWrap(d.Screen, 0, 26, 81, 26, d.DebugStyle, texth)
	textv := "01234567890123456789012345"
	drawTextWrap(d.Screen, 81, 0, 82, 27, d.DebugStyle, textv)

}

// -----------------------------------------------------------------------------
func (d *Display) DrawBox(x1, y1 int, w, h int) {
	h -= 1
	w -= 1

	for x := x1; x < x1+w; x++ {
		d.Screen.SetContent(x, y1, tcell.RuneHLine, nil, d.DefStyle)
		d.Screen.SetContent(x, y1+h, tcell.RuneHLine, nil, d.DefStyle)
	}

	for y := y1; y < y1+h; y++ {
		d.Screen.SetContent(x1, y, tcell.RuneVLine, nil, d.DefStyle)
		d.Screen.SetContent(x1+w, y, tcell.RuneVLine, nil, d.DefStyle)
	}

	d.Screen.SetContent(x1, y1, tcell.RuneULCorner, nil, d.DefStyle)
	d.Screen.SetContent(x1+w, y1, tcell.RuneURCorner, nil, d.DefStyle)
	d.Screen.SetContent(x1, y1+h, tcell.RuneLLCorner, nil, d.DefStyle)
	d.Screen.SetContent(x1+w, y1+h, tcell.RuneLRCorner, nil, d.DefStyle)
}

// -----------------------------------------------------------------------------
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
