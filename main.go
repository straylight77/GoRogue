package main

import "fmt"

var dungeon DungeonMap

func main() {
	var cmd rune
	var moves int

	disp := Display{}
	disp.Init()
	defer disp.Quit()

	dungeon.Clear()
	dungeon.CreateRoom(7, 7, 8, 5)
	dungeon.CreateRoom(27, 15, 10, 6)
	dungeon.CreateRoom(42, 3, 12, 8)

	done := false
	for !done {

		// Draw
		disp.DrawDebug()
		info := fmt.Sprintf("Level: 1  Gold: 4       Hp: 11 (20)  Str: 16(16)  Arm: 4   Exp: 2/14")
		disp.DrawText(0, 24, info)
		disp.DrawMap(&dungeon)

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
