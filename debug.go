package main

import (
	"fmt"
)

// ----------------------------------------------------------------------------
type DebugMessageLog struct {
	messages []string
}

func (log *DebugMessageLog) Add(format string, vals ...any) {
	msg := fmt.Sprintf(format, vals...)
	log.messages = append(log.messages, msg)
}

func (log *DebugMessageLog) Clear() {
	log.messages = nil
}

func (log *DebugMessageLog) Draw(disp *Display, startX, startY int) {
	for i, msg := range debug.messages {
		disp.Debug(startX, startY+i, msg)
	}

}

// -----------------------------------------------------------------------------
func drawDebugFrame(d *Display, gs *GameState) {
	maxX, maxY := 80, 25
	d.DrawBox(-1, -1, maxX+1, maxY+1, "debug")

	d.Debugf(84, 1, "Moves: %d, Pos: %v", gs.player.moves, gs.player.Pos())
	d.Debugf(84, 3, "M=%d, H=%d, F=%d, W=%d, SF=%d",
		gs.player.moves,
		gs.player.healCount,
		gs.player.foodCount,
		gs.wander,
		gs.spawnFoodTimer)
	d.Debugf(84, 4, "conf=%d, blind=%d", gs.player.confused, gs.player.blind)
	d.Debugf(84, 5, "path1: %v", path1)
	d.Debugf(84, 6, "path2: %v", path2)
	if gs.dmap != nil {
		d.Debugf(84, 7, "dmap: iter=%d", gs.dmap.iter)
	}

	//for i, m := range *gs.monsters {
	//	d.Debugf(84, 9+i, "%d: %v", i, m.DebugString())
	//}

	//for i, item := range gs.player.inventory {
	//	d.Debugf(84, 9+i, "%c) %v", 'a'+i, item.InvString())
	//}
}

// ----------------------------------------------------------------------------
func debugMapGrid(disp *Display) {
	disp.DrawHLine(8, 0, 79, "debug")
	disp.DrawHLine(16, 0, 79, "debug")
	disp.DrawVLine(26, 1, 24, "debug")
	disp.DrawVLine(53, 1, 24, "debug")

	disp.Debugf(0, 1, "0")
	disp.Debugf(27, 1, "1")
	disp.Debugf(54, 1, "2")
	disp.Debugf(0, 9, "3")
	disp.Debugf(27, 9, "4")
	disp.Debugf(54, 9, "5")
	disp.Debugf(0, 17, "6")
	disp.Debugf(27, 17, "7")
	disp.Debugf(54, 17, "8")
}

// ----------------------------------------------------------------------------
// called from main()
func drawGenerateDebug(disp *Display) {

	debugMapGrid(disp)

	for i, r := range graph.rooms {
		disp.Debugf(0, 28+i, "%d: %v", i, r)
		c := r.Center()
		disp.Debug(c.X, c.Y+1, "X") // Y+1 to convert to map coords
	}

	for i := 0; i < 9; i++ {
		lst := graph.Neighbours(i)
		disp.Debug(20, 28+i, lst)
	}

	for i, p := range graph.corridors {
		disp.Debug(35, 28+i, p)
	}

	cell := 0
	for row := 0; row < 3; row++ {
		for col := 0; col < 3; col++ {
			if graph.rooms[cell].mark != -1 {
				disp.Debug(50+(4*col), 28+(2*row), cell)
			}
			if graph.AreConnected(cell, cell+1) {
				disp.Debugf(50+(4*col)+2, 28+(2*row), "-")
			}
			if graph.AreConnected(cell, cell+3) {
				disp.Debugf(50+(4*col), 28+(2*row)+1, "|")
			}
			cell++
		}
	}
}
