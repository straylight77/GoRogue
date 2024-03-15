package main

// wrap these into GameState?  Will have handleCommand()?
var dungeon DungeonMap
var player Player
var monsters MonsterList

var messages MessageLog

var disp Display
var debug = DebugMessageLog{}

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
	CmdWait
	CmdTick
	CmdGenerate // for testing
	CmdMessages
)

// -----------------------------------------------------------------------
func movePlayer(dx int, dy int, d *DungeonMap, p *Player, mlist *MonsterList) {
	destX, destY := p.X+dx, p.Y+dy

	// check edges of the map
	if destX < 0 || destX >= MapMaxX || destY < 0 || destY >= MapMaxY {
		messages.Add("That way is blocked.")
		return
	}

	// check for monsters
	m := mlist.MonsterAt(destX, destY)
	if m != nil {
		messages.Add(p.Attack(m))
		return
	}

	// check dungeon tile
	destTile := d.TileAt(destX, destY)
	switch {

	case destTile.IsWalkable():
		p.SetPos(destX, destY)

	case destTile.IsType(TileDoorCl): // open the door
		d.SetTile(destX, destY, TileDoorOp)
		messages.Add("You open the door.")

	default:
		messages.Add("That way is blocked.")
	}

}

// -----------------------------------------------------------------------
func main() {
	var cmd GameCommand

	// Initialization and setup
	disp = Display{}
	disp.Init()
	defer disp.Quit()

	// Create a dungeon level
	dungeon.GenerateLevel(&player, &monsters)
	//generateRandomLevel(&dungeon, &monsters, &player)

	debugFlag := true
	doneFlag := false
	var doUpdate bool

	for !doneFlag {
		doUpdate = true

		// Draw the world
		disp.Clear()
		disp.DrawMap(&dungeon)
		disp.DrawMessages(&messages)
		disp.DrawText(0, 24, player.InfoString())

		for _, m := range monsters {
			disp.DrawEntity(m)
		}
		disp.DrawPlayer(&player)

		if debugFlag {
			disp.DrawDebugFrame(&player, &monsters)
			//drawGenerateDebug(&disp)
			debug.Draw(&disp, 84, 10)
		}

		disp.Show()

		cmd = disp.GetCommand()

		// Handle user's command
		switch cmd {

		// Commands that do not increment time
		case 0: // unkown command, just ignore
			doUpdate = false
		case CmdTick:
			doUpdate = false
			// Do nothing.  Used to redraw, clear recent messages, etc.
		case CmdMessages:
			doUpdate = false
			disp.DrawMessageHistory(&messages)
			disp.WaitForKeypress()
		case CmdQuit:
			doUpdate = false
			doneFlag = true

		// Commands that do increment time
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
				messages.Add("You decend the ancient stairs.")
				generateRandomLevel(&dungeon, &monsters, &player)
			} else {
				messages.Add("There are no stairs to go down here.")
			}
		case CmdUp:
			if dungeon.TileAt(player.X, player.Y).typ == TileStairsUp {
				messages.Add("Your way is magically blocked.")
			} else {
				messages.Add("There are no stairs to go up here.")
			}
		case CmdWait:
			messages.Add("You rest for a moment.")

		// Extra debugging and testing stuff
		case CmdDebug:
			doUpdate = false
			debugFlag = !debugFlag
		case CmdGenerate:
			doUpdate = false
			debug.Clear()
			//generateRandomLevel(&dungeon, &monsters, &player)
			dungeon.GenerateLevel(&player, &monsters)
		}

		// Do world updates
		if doUpdate {
			updateMonsters(&dungeon, &player, &monsters, &messages)
			player.moves++
		}

	}
}

// TODO: move all handling of game objects into a GameState object
func updateMonsters(d *DungeonMap, p *Player, ml *MonsterList, msg *MessageLog) {
	for i, m := range *ml {
		if m.HP <= 0 {
			ml.Remove(i)
			msg.Add("You defeated the %s!", m.Name)
		}
	}

}
