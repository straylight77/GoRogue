package main

import (
	"fmt"
)

// wrap these into GameState?  Will have handleCommand()?
var dungeon DungeonMap
var player Player
var monsters MonsterList

type Entity interface {
	Pos() (int, int)
	SetPos(int, int)
	Rune() rune
}

type GameCommand int

const (
	CmdNop GameCommand = iota
	CmdDebug
	CmdQuit
	CmdNorth
	CmdSouth
	CmdEast
	CmdWest
	CmdUp
	CmdDown
	CmdGenerate // for testing
)

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

	// check edges of the map
	if destX < 0 || destX >= MapMaxX || destY < 0 || destY >= MapMaxY {
		logMessage("That way is blocked.")
		return
	}

	// check for monsters
	m := mlist.MonsterAt(destX, destY)
	if m != nil {
		logMessage(p.Attack(m))
		return
	}

	// check dungeon tile
	destTile := d.TileAt(destX, destY)
	switch {

	case destTile.IsWalkable():
		p.SetPos(destX, destY)

	case destTile.IsType(TileDoorCl): // open the door
		d.SetTile(destX, destY, TileDoorOp)
		logMessage("You open the door.")

	default:
		logMessage("That way is blocked.")
	}

	player.moves++
}

var disp Display

// -----------------------------------------------------------------------
func main() {
	var cmd GameCommand

	// initialization and setup
	disp = Display{}
	disp.Init()
	defer disp.Quit()

	// create a dungeon level
	//dungeon.GenerateLevel(player.depth, &player, &monsters)
	generateRandomLevel(&dungeon, &monsters, &player)

	debug := true
	done := false
	for !done {

		// draw the world
		disp.Clear()
		disp.DrawMap(&dungeon)
		disp.DrawMessages(messages)
		disp.DrawText(0, 24, player.InfoString())

		for _, m := range monsters {
			disp.DrawEntity(m)
		}
		disp.DrawPlayer(&player)
		if debug {
			disp.DrawDebugFrame(&player, &monsters)
			drawGenerateDebug(&disp)
		}

		disp.Show()

		cmd = disp.GetCommand()

		// handle user's command
		switch cmd {
		case 0: //ignore
		case CmdQuit:
			done = true
		case CmdWest:
			movePlayer(-1, 0, &dungeon, &player, &monsters)
		case CmdEast:
			movePlayer(1, 0, &dungeon, &player, &monsters)
		case CmdNorth:
			movePlayer(0, -1, &dungeon, &player, &monsters)
		case CmdSouth:
			movePlayer(0, 1, &dungeon, &player, &monsters)
		case CmdDown:
			if dungeon.TileAt(player.X, player.Y).typ == TileStairsDn {
				logMessage("You decend the ancient stairs.")
				generateRandomLevel(&dungeon, &monsters, &player)
			} else {
				logMessage("There are no stairs to go down here.")
			}
		case CmdUp:
			if dungeon.TileAt(player.X, player.Y).typ == TileStairsUp {
				logMessage("Your way is magically blocked.")
			} else {
				logMessage("There are no stairs to go up here.")
			}
		case CmdDebug:
			debug = !debug
		case CmdGenerate:
			generateRandomLevel(&dungeon, &monsters, &player)
		}

		// do other world updates
		for i, m := range monsters {
			if m.HP <= 0 {
				monsters.Remove(i)
				logMessage(fmt.Sprintf("You defeated the %s!", m.Name))
			}
		}

	}
}
