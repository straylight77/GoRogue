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
		logMessage(p.Attack(m))
		return
	}

	// check dungeon tile
	destTile := d[destX][destY]
	switch destTile {

	case TileFloor, // consider these tiles as "walkable"
		TilePath,
		TileDoorOp:
		p.SetPos(destX, destY)

	case TileDoorCl: // open the door
		d.SetTile(destX, destY, TileDoorOp)
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
		disp.Screen.Clear()
		disp.DrawMap(&dungeon)
		disp.DrawMessages(messages)
		disp.DrawText(0, 24, player.InfoString())

		for _, m := range monsters {
			disp.DrawEntity(m)
		}
		disp.DrawEntity(&player)
		disp.DrawDebug(&player, &monsters)

		disp.Screen.Show()

		// get the Game Command from user (blocks until input)
		cmd = disp.GetCommand()

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

		// do other world updates
		for i, m := range monsters {
			if m.HP <= 0 {
				monsters.Remove(i)
				logMessage(fmt.Sprintf("You defeated the %s!", m.Name))
			}
		}

	}
}
