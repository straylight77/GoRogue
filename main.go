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

// -----------------------------------------------------------------------
var messages []string

func logMessage(s string) {
	messages = append(messages, s)
}

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

	case '.', // consider these tiles as "walkable"
		'#',
		'`',
		tcell.RuneBullet,
		tcell.RuneBoard,
		tcell.RuneCkBoard:
		p.SetPos(destX, destY)

	case '+': // open the door
		d.SetTile(destX, destY, '`')
		logMessage("You open the door.")

	default:
		logMessage("That way is blocked.")
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
	dungeon.GenerateLevel(1, &player)

	done := false
	for !done {

		// draw the world
		disp.DrawMap(&dungeon)
		disp.DrawEntity(&player)
		disp.DrawMessages(messages)
		disp.DrawText(0, 24, player.InfoString())
		disp.DrawDebug(&player)

		disp.Screen.Show()

		// get the Game Command from user (blocks until input)
		cmd = disp.Command()

		// handle user's command
		switch cmd {
		case 'X':
			done = true
		case tcell.KeyLeft:
			movePlayer(-1, 0, &dungeon, &player)
		case tcell.KeyRight:
			movePlayer(1, 0, &dungeon, &player)
		case tcell.KeyUp:
			movePlayer(0, -1, &dungeon, &player)
		case tcell.KeyDown:
			movePlayer(0, 1, &dungeon, &player)
		}

		player.moves++
	}
}
