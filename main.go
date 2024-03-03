package main

import (
	"github.com/gdamore/tcell/v2"
)

type Entity interface {
	Pos() (int, int)
	SetPos(int, int)
	Rune() rune
}

// wrap these into GameState?  Will have handleCommand()?
var dungeon DungeonMap
var player Player

var messages []string

// -----------------------------------------------------------------------
func logMessage(s string) {
	messages = append(messages, s)
}

// -----------------------------------------------------------------------
func clearMessages() {
	messages = nil
}

// -----------------------------------------------------------------------
func movePlayer(dx int, dy int, d *DungeonMap, p *Player) {
	destX, destY := p.X+dx, p.Y+dy

	// check for monsters

	// check dungeon tile
	destTile := d[destX][destY]
	switch destTile {
	case '.', '#', '`', tcell.RuneBullet:
		p.SetPos(destX, destY)
		break
	case '+':
		// open the door
		d.SetTile(destX, destY, '`')
		logMessage("You open the door.")
		break
	default:
		logMessage("That way is blocked.")
		break
	}
}

// -----------------------------------------------------------------------
func main() {
	var cmd tcell.Key

	// initialization and setup
	disp := Display{}
	disp.Init()
	defer disp.Quit()

	// create a dungeon level
	dungeon.Clear()
	dungeon.CreateRoom(7, 7, 8, 5)
	dungeon.CreateRoom(27, 15, 10, 6)
	dungeon.CreateRoom(42, 3, 12, 8)
	dungeon.SetTile(42, 5, '+')

	player.SetPos(45, 5)

	done := false
	for !done {

		// Draw
		disp.DrawDebug(&player)
		disp.DrawText(0, 24, player.InfoString())
		disp.DrawMap(&dungeon)
		disp.DrawEntity(&player)

		// draw messages if we have any
		if len(messages) > 0 {
			for i, msg := range messages {
				disp.DrawText(0, i, msg)
			}
			clearMessages()
		}

		disp.Screen.Show()

		// get the Game Command from user (blocks until input)
		cmd = disp.Command()

		// handle user's command
		switch cmd {
		case 'X':
			done = true
			break
		case tcell.KeyLeft:
			movePlayer(-1, 0, &dungeon, &player)
			break
		case tcell.KeyRight:
			movePlayer(1, 0, &dungeon, &player)
			break
		case tcell.KeyUp:
			movePlayer(0, -1, &dungeon, &player)
			break
		case tcell.KeyDown:
			movePlayer(0, 1, &dungeon, &player)
			break

		}

		player.moves++
	}
}
