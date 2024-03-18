package main

import (
	"fmt"
	"log"

	"github.com/gdamore/tcell/v2"
)

var KeyCmdLookup = map[tcell.Key]GameCommand{
	tcell.KeyEscape: CmdQuit,
	tcell.KeyCtrlC:  CmdQuit,
	tcell.KeyLeft:   CmdWest,
	tcell.KeyRight:  CmdEast,
	tcell.KeyUp:     CmdNorth,
	tcell.KeyDown:   CmdSouth,
}

var RuneCmdLookup = map[rune]GameCommand{
	'.': CmdWait,
	' ': CmdTick,
	'Q': CmdQuit,
	'D': CmdDebug,
	'P': CmdDebug2,
	'p': CmdTest1,
	'G': CmdGenerate,
	'M': CmdMessages,
	'>': CmdDown,
	'<': CmdUp,
}

var TileRunes = map[TileType]rune{
	//TileEmpty:    tcell.RuneCkBoard, // for testing
	TileEmpty:    ' ',
	TileWallH:    '-',
	TileWallV:    '|',
	TileWallUL:   '-',
	TileWallUR:   '-',
	TileWallLL:   '-',
	TileWallLR:   '-',
	TileFloor:    '.',
	TilePath:     '#',
	TileDoor:     '+',
	TileStairsUp: '<',
	TileStairsDn: '>',
}

type Display struct {
	Screen      tcell.Screen
	DefStyle    tcell.Style
	DebugStyle  tcell.Style
	Debug2Style tcell.Style
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
	d.Debug2Style = tcell.StyleDefault.Foreground(tcell.ColorRed)

	s.SetStyle(d.DefStyle)
	s.SetCursorStyle(tcell.CursorStyleSteadyBlock)
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
func (d *Display) Clear() {
	d.Screen.Clear()
}

// -----------------------------------------------------------------------------
func (d *Display) Show() {
	d.Screen.Show()
}

// -----------------------------------------------------------------------------
func (d *Display) Printf(x, y int, format string, vals ...any) {
	text := fmt.Sprintf(format, vals...)
	col, row := x, y
	for _, r := range []rune(text) {
		d.Screen.SetContent(col, row, r, nil, d.DefStyle)
		col++
	}
}

// -----------------------------------------------------------------------------
func (d *Display) DrawEntity(e Entity) {
	x, y := e.Pos()
	d.Screen.SetContent(x, y+1, e.Rune(), nil, d.DefStyle)
}

// -----------------------------------------------------------------------------
func (d *Display) DrawPlayer(p *Player) {
	x, y := p.Pos()
	d.Screen.SetContent(x, y+1, '@', nil, d.DefStyle)
	d.Screen.ShowCursor(x, y+1)
}

// -----------------------------------------------------------------------------
func (d *Display) DrawMap(m *DungeonMap, showAll bool) {
	for x, col := range m.tiles {
		for y, t := range col {
			r := TileRunes[t.typ]
			if showAll {
				d.Screen.SetContent(x, y+1, r, nil, d.DebugStyle)
			}

			if t.visible {
				// y+1 because first line is the message line
				d.Screen.SetContent(x, y+1, r, nil, d.DefStyle)
			} else if t.visited && t.typ != TileFloor {
				// have the option to use a different style here
				d.Screen.SetContent(x, y+1, r, nil, d.DefStyle)
			}
		}
	}
}

// -----------------------------------------------------------------------------
func (d *Display) DrawMessages(log *MessageLog) {
	if log.HasUnread() {
		s := log.LatestAsStr()
		drawTextWrap(d.Screen, 0, 0, 80, 3, d.DefStyle, s)
		log.ClearUnread()
	}
}

// -----------------------------------------------------------------------------
func (d *Display) DrawMessageHistory(log *MessageLog) {
	d.Clear()
	d.DrawText(0, 0, "MESSAGE HISTORY:")
	for i, m := range log.Last(20) {
		d.Printf(0, i+1, "%d: %v", i, m)
	}
	d.Printf(0, 22, "Press any key to continue...")
	d.Screen.HideCursor()
	d.Show()
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
func (d *Display) DrawDebug(x1, y1 int, text string) {
	row := y1
	col := x1
	for _, r := range []rune(text) {
		d.Screen.SetContent(col, row, r, nil, d.DebugStyle)
		col++
	}
}

// -----------------------------------------------------------------------------
func (d *Display) DrawHLine(row int, from, to int, style tcell.Style) {
	for i := from; i <= to; i++ {
		d.Screen.SetContent(i, row, tcell.RuneHLine, nil, style)
	}
}

// -----------------------------------------------------------------------------
func (d *Display) DrawVLine(col int, from, to int, style tcell.Style) {
	for i := from; i < to; i++ {
		d.Screen.SetContent(col, i, tcell.RuneVLine, nil, style)
	}
}

// -----------------------------------------------------------------------------
func (d *Display) DrawDebugFrame(p *Player, ml *MonsterList) {
	maxX, maxY := 80, 25
	d.DrawHLine(maxY, 0, maxX, d.DebugStyle)
	d.DrawVLine(maxX, 0, maxY, d.DebugStyle)
	d.Screen.SetContent(maxX, maxY, tcell.RuneLRCorner, nil, d.DebugStyle)

	d.DrawDebug(84, 1, fmt.Sprintf("Moves:  %d", p.moves))
	d.DrawDebug(84, 2, fmt.Sprintf("Pos: %d, %d", p.X, p.Y))
	d.DrawDebug(84, 3, fmt.Sprintf("Heal: %d,  Food: %d", p.healCount, p.foodCount))
	for i, m := range *ml {
		msg := fmt.Sprintf("%d: %v", i, m.DebugString())
		d.DrawDebug(84, 5+i, msg)
	}
}

// -----------------------------------------------------------------------------
func (d *Display) GetCommand() (cmd GameCommand) {

	var ok bool
	ev := d.Screen.PollEvent() // blocks until input from user

	// do 'display level' processing, otherwise find the GameCommand to return
	switch ev := ev.(type) {
	case *tcell.EventResize:
		d.Screen.Clear()
		d.Screen.Sync()

	case *tcell.EventKey:
		key := ev.Key()
		rn := ev.Rune()

		switch key {
		case tcell.KeyRune:
			if cmd, ok = RuneCmdLookup[rn]; !ok {
				messages.Add("I don't know that command (%c)", rn)
			}
		default:
			if cmd, ok = KeyCmdLookup[key]; !ok {
				messages.Add("I don't know that command (%v)", tcell.KeyNames[key])
			}
		}
	}
	return cmd
}

// -----------------------------------------------------------------------------
func (d *Display) WaitForKeypress() {
	d.Screen.PollEvent() // blocks until input from user
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
