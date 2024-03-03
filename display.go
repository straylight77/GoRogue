package main

import (
	"log"

	"github.com/gdamore/tcell/v2"
)

type Display struct {
	Screen   tcell.Screen
	DefStyle tcell.Style
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
func (d *Display) drawText(x1, y1 int, text string) {
	row := y1
	col := x1
	for _, r := range []rune(text) {
		d.Screen.SetContent(col, row, r, nil, d.DefStyle)
		col++
	}
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

// -----------------------------------------------------------------------------
func drawBorder(s tcell.Screen, style tcell.Style) {
	xmax, ymax := s.Size()

	for x := 0; x < xmax; x++ {
		s.SetContent(x, 0, tcell.RuneHLine, nil, style)
		s.SetContent(x, ymax-7, tcell.RuneHLine, nil, style)
		s.SetContent(x, ymax-1, tcell.RuneHLine, nil, style)
	}

	for y := 0; y < ymax; y++ {
		s.SetContent(0, y, tcell.RuneVLine, nil, style)
		s.SetContent(xmax-1, y, tcell.RuneVLine, nil, style)
	}

	for y := 0; y < ymax-7; y++ {
		s.SetContent(xmax-30, y, tcell.RuneVLine, nil, style)
	}

	s.SetContent(xmax-30, 0, tcell.RuneTTee, nil, style)
	s.SetContent(xmax-30, ymax-7, tcell.RuneBTee, nil, style)

	s.SetContent(0, 0, tcell.RuneULCorner, nil, style)
	s.SetContent(xmax-1, 0, tcell.RuneURCorner, nil, style)
	s.SetContent(0, ymax-1, tcell.RuneLLCorner, nil, style)
	s.SetContent(xmax-1, ymax-1, tcell.RuneLRCorner, nil, style)
	s.SetContent(0, ymax-7, tcell.RuneLTee, nil, style)
	s.SetContent(xmax-1, ymax-7, tcell.RuneRTee, nil, style)
}
