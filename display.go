package main

import (
	"fmt"
	"log"

	"github.com/gdamore/tcell/v2"
)

var KeyCmdLookup = map[tcell.Key]GameCommand{
	tcell.KeyCtrlD:  CmdDebug1,
	tcell.KeyCtrlG:  CmdDebug2,
	tcell.KeyCtrlM:  CmdDebug3,
	tcell.KeyCtrlP:  CmdDebug4,
	tcell.KeyCtrlR:  CmdDebug5,
	tcell.KeyEscape: CmdQuit,
	tcell.KeyCtrlC:  CmdQuit,
	tcell.KeyLeft:   CmdWest,
	tcell.KeyRight:  CmdEast,
	tcell.KeyUp:     CmdNorth,
	tcell.KeyDown:   CmdSouth,
}

var RuneCmdLookup = map[rune]GameCommand{
	'G': CmdGenerate,
	'.': CmdWait,
	' ': CmdTick,
	'Q': CmdQuit,
	'M': CmdMessages,
	'>': CmdDown,
	'<': CmdUp,
	'1': CmdSouthWest,
	'2': CmdSouth,
	'3': CmdSouthEast,
	'4': CmdWest,
	'5': CmdWait,
	'6': CmdEast,
	'7': CmdNorthWest,
	'8': CmdNorth,
	'9': CmdNorthEast,
	'i': CmdInventory,
	'e': CmdEat,
}

var TileRunes = map[TileType]rune{
	TileEmpty:    ' ',
	TileWallH:    '-',
	TileWallV:    '|',
	TileWallUL:   '-',
	TileWallUR:   '-',
	TileWallLL:   '-',
	TileWallLR:   '-',
	TileFloor:    '.',
	TileCorridor: '#',
	TileDoor:     '+',
	TileStairsUp: '<',
	TileStairsDn: '>',
}

type Display struct {
	Screen tcell.Screen
	styles map[string]tcell.Style
}

// -----------------------------------------------------------------------------
func (d *Display) Init() {
	scr, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := scr.Init(); err != nil {
		log.Fatalf("%+v", err)
	}

	d.styles = make(map[string]tcell.Style)
	d.styles["default"] = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
	d.styles["debug"] = tcell.StyleDefault.Foreground(tcell.ColorLightSkyBlue)
	d.styles["debug2"] = tcell.StyleDefault.Foreground(tcell.ColorRed).Background(tcell.ColorDarkRed)
	d.styles["yellow"] = tcell.StyleDefault.Foreground(tcell.ColorYellow)
	d.styles["orange"] = tcell.StyleDefault.Foreground(tcell.ColorOrange)
	d.styles["red"] = tcell.StyleDefault.Foreground(tcell.ColorRed)
	d.styles["purple"] = tcell.StyleDefault.Foreground(tcell.ColorPurple)
	d.styles["darkblue"] = tcell.StyleDefault.Foreground(tcell.ColorNavy)
	d.styles["blue"] = tcell.StyleDefault.Foreground(tcell.ColorBlue)
	d.styles["bluegreen"] = tcell.StyleDefault.Foreground(tcell.ColorAquaMarine)
	d.styles["green"] = tcell.StyleDefault.Foreground(tcell.ColorGreen)
	d.styles["darkgreen"] = tcell.StyleDefault.Foreground(tcell.ColorDarkGreen)

	scr.SetStyle(d.styles["default"])
	scr.SetCursorStyle(tcell.CursorStyleSteadyBlock)
	scr.Clear()
	d.Screen = scr
}

// -----------------------------------------------------------------------------
func (d *Display) Quit() {
	maybePanic := recover()
	d.Screen.Fini()
	if maybePanic != nil {
		//log.Fatalf("%+v", maybePanic)
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
func (d *Display) Style(styleName string) tcell.Style {
	style, ok := d.styles[styleName]
	if !ok {
		style = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
	}
	return style
}

// -----------------------------------------------------------------------------
func (d *Display) Printf(x, y int, format string, vals ...any) {
	text := fmt.Sprintf(format, vals...)
	d.DrawText(x, y, "default", text)
}

// -----------------------------------------------------------------------------
func (d *Display) Print(x, y int, vals ...any) {
	text := fmt.Sprint(vals...)
	d.DrawText(x, y, "default", text)
}

// -----------------------------------------------------------------------------
func (d *Display) Debugf(x, y int, format string, vals ...any) {
	text := fmt.Sprintf(format, vals...)
	d.DrawText(x, y, "debug", text)
}

// -----------------------------------------------------------------------------
func (d *Display) Debug(x, y int, vals ...any) {
	text := fmt.Sprint(vals...)
	d.DrawText(x, y, "debug", text)
}

// -----------------------------------------------------------------------------
func (d *Display) DrawText(x1, y1 int, styleName string, text string) {
	style := d.Style(styleName)
	row := y1
	col := x1
	for _, r := range []rune(text) {
		d.Screen.SetContent(col, row, r, nil, style)
		col++
	}
}

// -----------------------------------------------------------------------------
func (d *Display) DrawBox(x, y int, w, h int, styleName string) {
	style := d.Style(styleName)
	for col := x; col <= x+w; col++ {
		d.Screen.SetContent(col, y, tcell.RuneHLine, nil, style)
		d.Screen.SetContent(col, y+h, tcell.RuneHLine, nil, style)
	}
	for row := y; row <= y+h; row++ {
		d.Screen.SetContent(x, row, tcell.RuneVLine, nil, style)
		d.Screen.SetContent(x+w, row, tcell.RuneVLine, nil, style)
	}
	d.Screen.SetContent(x, y, tcell.RuneULCorner, nil, style)
	d.Screen.SetContent(x+w, y, tcell.RuneURCorner, nil, style)
	d.Screen.SetContent(x, y+h, tcell.RuneLLCorner, nil, style)
	d.Screen.SetContent(x+w, y+h, tcell.RuneLRCorner, nil, style)
}

// -----------------------------------------------------------------------------
func (d *Display) DrawHLine(row int, from, to int, styleName string) {
	style := d.Style(styleName)
	for i := from; i <= to; i++ {
		d.Screen.SetContent(i, row, tcell.RuneHLine, nil, style)
	}
}

// -----------------------------------------------------------------------------
func (d *Display) DrawVLine(col int, from, to int, styleName string) {
	style := d.Style(styleName)
	for i := from; i < to; i++ {
		d.Screen.SetContent(col, i, tcell.RuneVLine, nil, style)
	}
}

// -----------------------------------------------------------------------------
func (d *Display) DrawEntity(e Entity) {
	x, y := e.Pos().XY()
	d.Screen.SetContent(x, y+1, e.Rune(), nil, d.Style("default"))
}

// -----------------------------------------------------------------------------
func (d *Display) DrawItem(pos Coord, i *Item) {
	x, y := pos.XY()
	d.Screen.SetContent(x, y+1, i.Rune(), nil, d.Style("default"))
}

// -----------------------------------------------------------------------------
func (d *Display) DrawPlayer(p *Player) {
	x, y := p.Pos().XY()
	d.Screen.SetContent(x, y+1, '@', nil, d.Style("default"))
	d.Screen.ShowCursor(x, y+1)
}

// -----------------------------------------------------------------------------
func (d *Display) DrawMap(m *DungeonMap, showAll bool) {
	for x, col := range m.tiles {
		for y, t := range col {
			r := TileRunes[t.typ]

			if showAll {
				d.Screen.SetContent(x, y+1, r, nil, d.Style("debug"))
			}

			if t.visible {
				// y+1 because first line is the message line
				d.Screen.SetContent(x, y+1, r, nil, d.Style("default"))
			} else if t.visited && t.typ != TileFloor {
				// have the option to use a different style here
				d.Screen.SetContent(x, y+1, r, nil, d.Style("default"))
			}
		}
	}
}

// -----------------------------------------------------------------------------
func (d *Display) DrawMessages(log *MessageLog) {
	if log.HasUnread() {
		s := log.LatestAsStr()
		drawTextWrap(d.Screen, 0, 0, 80, 3, d.Style("default"), s)
		//d.Print(0, 0, s)
		log.ClearUnread()
	}
}

// -----------------------------------------------------------------------------
func (d *Display) DrawMessageHistory(log *MessageLog) {
	d.Clear()
	d.Print(0, 0, "MESSAGE HISTORY:")
	for i, m := range log.Last(20) {
		d.Printf(0, i+1, "%d: %v", i, m)
	}
	d.Printf(0, 22, "Press any key to continue...")
	d.Screen.HideCursor()
	d.Show()
}

// -----------------------------------------------------------------------------
func (d *Display) DrawInventory(p *Player) {
	d.Clear()
	d.Print(0, 0, "YOUR INVENTORY:")
	for i, item := range p.inventory {
		d.Printf(0, i+1, "%c) %v", 'a'+i, item)
	}
	d.Printf(0, 22, "Press any key to continue...")
	d.Screen.HideCursor()
	d.Show()
}

// -----------------------------------------------------------------------------
func (d *Display) WaitForKeypress() {
	d.Screen.PollEvent() // blocks until input from user
}

// -----------------------------------------------------------------------------
// Handles all events appropriateley (e.g. resizing) but this functions will only
// return when a key event is received.  Will return 0 if the command is not
// recognized along with creating a game message.
func (d *Display) GetCommand(msg *MessageLog) (cmd GameCommand) {

	gotEventKey := false
	for !gotEventKey {

		ev := d.Screen.PollEvent()
		//debug.Add("event: %T", ev)

		switch ev := ev.(type) {
		case *tcell.EventResize:
			//d.Screen.Clear()
			d.Screen.Sync()

		case *tcell.EventKey:
			gotEventKey = true
			key := ev.Key()
			rn := ev.Rune()

			var ok bool
			switch key {
			case tcell.KeyRune:
				if cmd, ok = RuneCmdLookup[rn]; !ok {
					msg.Add("I don't know that command (%c)", rn)
				}
			default:
				if cmd, ok = KeyCmdLookup[key]; !ok {
					msg.Add("I don't know that command (%v)", tcell.KeyNames[key])
				}
			}
		}
	}
	return cmd
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
