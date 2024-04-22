package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/gdamore/tcell/v2"
)

var KeyCmdLookup = map[tcell.Key]GameCommand{
	tcell.KeyCtrlD:  CmdDebug1,
	tcell.KeyCtrlG:  CmdDebug2,
	tcell.KeyCtrlX:  CmdDebug3,
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
	'c': CmdConsume,
	'e': CmdEquip,
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
func (d *Display) DrawActor(a Actor) {
	x, y := a.Pos().XY()
	d.Screen.SetContent(x, y+1, a.Rune(), nil, d.Style("default"))
}

// -----------------------------------------------------------------------------
func (d *Display) DrawItem(pos Coord, item Item) {
	x, y := pos.XY()
	d.Screen.SetContent(x, y+1, item.Rune(), nil, d.Style("default"))
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
	for i, m := range log.Last(22) {
		d.Printf(0, i, "%v", m)
	}
	d.Printf(0, 24, "Press space to continue...")
	d.Screen.HideCursor()
	d.Show()
}

// -----------------------------------------------------------------------------
func (d *Display) InventoryScreen(p *Player) {
	d.Clear()
	d.Print(0, 0, "You are carrying:")
	d.ListInventory(p, 0, false)

	col := 60
	row := 1

	//label := map[string]string{
	//	"weapon": "W",
	//	"armor":  "A",
	//	"left":   "L",
	//	"right":  "R",
	//}
	//order := []string{"weapon", "armor", "left", "right"}
	//d.Printf(col, 0, "Equipment")
	//for _, slot := range order {
	//	if p.equiped[slot] == nil {
	//		d.Printf(col, row, "%s) -none-", label[slot])
	//	} else {
	//		d.Printf(col, row, "%s) %v", label[slot], p.equiped[slot].InvString())
	//	}
	//	row++
	//}

	d.Printf(col, row, "       STATS")
	row += 2
	stats := p.StatsStrings()
	for _, str := range stats {
		d.Printf(col, row, str)
		row++
	}

	d.Printf(0, 23, "Press space to continue...")
	d.Screen.HideCursor()
	d.Show()
}

// -----------------------------------------------------------------------------
func (d *Display) PromptInventory(prompt string, p *Player) int {
	lo := 'a'
	hi := rune(int(lo) + len(p.inventory) - 1)
	str := fmt.Sprintf("%s (%c-%c, ? for list, ESC to cancel):", prompt, lo, hi)

	d.Print(0, 0, strings.Repeat(" ", 80))
	d.Print(0, 0, str)
	d.Screen.ShowCursor(len(str), 0)
	d.Show()

	ch := d.PromptRune()
	if ch == '?' {
		d.ListInventory(p, len(str), false)
		ch = d.PromptRune()
	}
	if ch >= lo && ch <= hi {
		return int(ch - 'a')
	}
	return -1
}

// -----------------------------------------------------------------------------
func (d *Display) ListInventory(p *Player, startWidth int, showWorth bool) {

	height := len(p.inventory)
	if height <= 0 {
		d.Print(0, 1, "Your inventory is empty.")
		d.Show()
		return
	}

	// determine strings to print and largest length
	width := 0
	strList := make([]string, height)
	for i, item := range p.inventory {
		equip := ""
		if item == p.equiped["weapon"] {
			equip = " (weapon in hand)"
		}
		if item == p.equiped["armor"] {
			equip = " (being worn)"
		}
		str := fmt.Sprintf("%c) %c %v%s", 'a'+i, item.Rune(), item.InvString(), equip)

		if check := len(str); check > width {
			width = check
		}
		strList[i] = str
	}

	// make a blank rectangle
	boxWidth := max(startWidth, width+3) + 1
	if showWorth {
	}
	for row := 0; row < height+1; row++ {
		d.Print(0, row+1, strings.Repeat(" ", boxWidth))
	}

	// print the items
	for i, str := range strList {
		d.Print(0, 1+i, str)
	}

	// show worth of each item if showWorth is true
	if showWorth {
		for i, item := range p.inventory {
			d.Printf(width+1, 1+i, "%4d ", item.Worth())
		}
	}
	d.Show()
}

// -----------------------------------------------------------------------------
func (d *Display) PromptRune() rune {
	ev := d.Screen.PollEvent()

	switch ev := ev.(type) {
	case *tcell.EventKey:
		if ev.Key() == tcell.KeyEscape {
			return -1
		} else {
			return ev.Rune()
		}
	}
	return 0
}

// -----------------------------------------------------------------------------
func (d *Display) WaitForKeypress() {
	//d.Screen.PollEvent() // blocks until input from user
	ch := d.PromptRune()
	for ch != ' ' {
		ch = d.PromptRune()
	}
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

// -----------------------------------------------------------------------------
func (d *Display) TombstoneScreen(gs *GameState) {

	tombstone := []string{
		"              __________",
		"             /          \\",
		"            /    REST    \\",
		"           /      IN      \\",
		"          /     PEACE      \\",
		"         /                  \\",
		"         |                  |",
		"         |                  |",
		"         |     killed by    |",
		"         |                  |",
		"         |  with the score  |",
		"         |                  |",
		"         |                  |",
		"         |       1980       |",
		"         |                  |",
		"        *|     *  *  *      | *",
		"________)/\\\\_//(\\/(/\\)/\\//\\/|_)_______",
		"  ",
		"       Press SPACE to continue...",
	}

	length := 0
	for _, str := range tombstone {
		if len(str) > length {
			length = len(str)
		}
	}
	col := 40 - (length / 2)
	row := 24 - len(tombstone)

	d.Clear()
	for _, str := range tombstone {
		d.Print(col, row, str)
		row++
	}
	name := "Nameless Hero"
	d.Print(40-(len(name)/2), 24-13, name)
	killedBy := gs.player.killedBy
	d.Print(40-(len(killedBy)/2), 24-10, killedBy)
	scoreStr := fmt.Sprintf("%d", gs.player.Score())
	d.Print(40-(len(scoreStr)/2), 24-8, scoreStr)
	d.Screen.HideCursor()
	d.Show()
	d.WaitForKeypress()
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
