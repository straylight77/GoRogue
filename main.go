package main

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
)

var dungeon DungeonMap

type Entity interface {
	Pos() (int, int)
	SetPos(int, int)
	Rune() rune
}

var player Player

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
		break
	default:
		// log message: that way is blocked
		break
	}
}

// -----------------------------------------------------------------------
func main() {
	var cmd tcell.Key
	var moves int

	disp := Display{}
	disp.Init()
	defer disp.Quit()

	dungeon.Clear()
	dungeon.CreateRoom(7, 7, 8, 5)
	dungeon.CreateRoom(27, 15, 10, 6)
	dungeon.CreateRoom(42, 3, 12, 8)
	dungeon.SetTile(42, 5, '+')

	player.SetPos(45, 5)

	done := false
	for !done {

		// Draw
		disp.DrawDebug()
		info := fmt.Sprintf("Level: 1  Gold: 4       Hp: 11 (20)  Str: 16(16)  Arm: 4   Exp: 2/14")
		disp.DrawText(0, 24, info)
		disp.DrawMap(&dungeon)
		disp.DrawEntity(&player)

		disp.Screen.Show()

		// Get the Game Command from user (blocks until input)
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

		moves++
	}
}
