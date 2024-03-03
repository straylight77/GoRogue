package main

import "fmt"

func main() {
	var cmd rune
	var moves int

	disp := Display{}
	disp.Init()
	defer disp.Quit()

	done := false
	for !done {

		// Draw
		drawBorder(disp.Screen, disp.DefStyle)
		disp.drawText(2, 2, fmt.Sprintf("frames: %d", moves))
		//drawText(disp.Screen, 2, 3, 30, 30, defStyle, fmt.Sprintf("cmd: %v", cmd))
		disp.drawText(5, 5, "Hello world! This is a long string.")
		disp.drawText(135, 7, "This one goes off the screen but it's okay.")
		disp.Screen.Show()

		// Get the Game Command from user (blocks until input)
		cmd = disp.Command()

		// handle user's command
		switch cmd {
		case 'X':
			done = true
			break
		}

		moves++
	}
}
