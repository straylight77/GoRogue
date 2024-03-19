package main

import "fmt"

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
		disp.DrawDebug(startX, startY+i, msg)
	}

}

// ----------------------------------------------------------------------------
func debugMapGrid(disp *Display) {
	disp.DrawHLine(8, 0, 79, disp.DebugStyle)
	disp.DrawHLine(16, 0, 79, disp.DebugStyle)
	disp.DrawVLine(26, 1, 24, disp.DebugStyle)
	disp.DrawVLine(53, 1, 24, disp.DebugStyle)

	disp.DrawDebug(0, 1, "0")
	disp.DrawDebug(27, 1, "1")
	disp.DrawDebug(54, 1, "2")
	disp.DrawDebug(0, 9, "3")
	disp.DrawDebug(27, 9, "4")
	disp.DrawDebug(54, 9, "5")
	disp.DrawDebug(0, 17, "6")
	disp.DrawDebug(27, 17, "7")
	disp.DrawDebug(54, 17, "8")
}

// ----------------------------------------------------------------------------
// called from main()
func drawGenerateDebug(disp *Display) {

	debugMapGrid(disp)

	for i, r := range graph.rooms {
		info := fmt.Sprintf("%d: %v", i, r)
		disp.DrawDebug(0, 28+i, info)
		cX, cY := r.Center()
		disp.DrawDebug(cX, cY+1, "X") // Y+1 to convert to map coords
	}

	for i := 0; i < 9; i++ {
		lst := graph.Neighbours(i)
		disp.DrawDebug(20, 28+i, fmt.Sprint(lst))
	}

	for i, p := range graph.corridors {
		disp.DrawDebug(35, 28+i, fmt.Sprint(p))
	}

	debug.Draw(disp, 84, 5)

	cell := 0
	for row := 0; row < 3; row++ {
		for col := 0; col < 3; col++ {
			if graph.rooms[cell].mark != -1 {
				disp.DrawDebug(50+(4*col), 28+(2*row), fmt.Sprint(cell))
			}
			if graph.AreConnected(cell, cell+1) {
				disp.DrawDebug(50+(4*col)+2, 28+(2*row), "-")
			}
			if graph.AreConnected(cell, cell+3) {
				disp.DrawDebug(50+(4*col), 28+(2*row)+1, "|")
			}
			cell++
		}
	}
}
