package main

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
)

// wrap these into GameState?  Will have handleCommand()?
var dungeon DungeonMap
var player Player
var monsters MonsterList

// -----------------------------------------------------------------------
var messages []string

func logMessage(s string) {
	messages = append(messages, s)
}

func clearMessages() {
	messages = nil
}

// -----------------------------------------------------------------------
func movePlayer(dx int, dy int, d *DungeonMap, p *Player, mlist *MonsterList) {
	destX, destY := p.X+dx, p.Y+dy

	// check for monsters
	m := mlist.MonsterAt(destX, destY)
	if m != nil {
		m.HP-- // placeholder combat for now
		msg := fmt.Sprintf("You attack the %v (hp=%d).", m.Name, m.HP)
		logMessage(msg)
		return
	}

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
	player.Symbol = '@'

	// create a dungeon level
	dungeon.GenerateLevel(1, &player)

	m1 := NewMonster("bat", 'B', 3)
	monsters.Add(m1, 50, 8)

	done := false
	for !done {

		// draw the world
		disp.DrawMap(&dungeon)
		disp.DrawMessages(messages)
		disp.DrawText(0, 24, player.InfoString())

		for _, m := range monsters {
			disp.DrawEntity(m)
		}
		disp.DrawEntity(&player)
		disp.DrawDebug(&player)

		disp.Screen.Show()

		// get the Game Command from user (blocks until input)
		cmd = disp.Command()

		// handle user's command
		switch cmd {
		case tcell.KeyEscape,
			tcell.KeyCtrlC:
			done = true
		case tcell.KeyLeft:
			movePlayer(-1, 0, &dungeon, &player, &monsters)
		case tcell.KeyRight:
			movePlayer(1, 0, &dungeon, &player, &monsters)
		case tcell.KeyUp:
			movePlayer(0, -1, &dungeon, &player, &monsters)
		case tcell.KeyDown:
			movePlayer(0, 1, &dungeon, &player, &monsters)
		default:
			logMessage(fmt.Sprintf("I don't know that command (%v)", cmd))
		}

		player.moves++
	}
}
