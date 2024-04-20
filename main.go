package main

import "fmt"

var debug DebugMessageLog

var debugFlag = map[string]bool{
	"main":     false,
	"generate": false,
	"dmap":     false,
	"path":     false,
}

// For testing
var path1 Path
var path2 Path
var RoomID int

const (
	WanderTimer   = 70 // For spawning wandering monsters
	NutritionTime = 1300
	HungerLimit   = 300
	WeakLimit     = 150
	SpawnFood     = 3 // Guarantee food spawns every 3 levels
)

type GameCommand int

const (
	CmdNop GameCommand = iota
	CmdDebug1
	CmdDebug2
	CmdDebug3
	CmdDebug4
	CmdDebug5
	CmdQuit

	CmdWait
	CmdNorth
	CmdNorthEast
	CmdEast
	CmdSouthEast
	CmdSouth
	CmdSouthWest
	CmdWest
	CmdNorthWest
	CmdUp
	CmdDown
	CmdConsume
	CmdEquip

	CmdTick
	CmdGenerate // for testing
	CmdMessages
	CmdInventory
)

// -----------------------------------------------------------------------
func main() {

	// Initialization
	var display Display
	display.Init()
	defer display.Quit()

	// Set up the initial game state
	var state GameState
	state.Init()

	var doUpdate bool   // If game time has passed this iteration
	var cmd GameCommand // Determined from user's input

	// Main Game Loop
	done := false
	for !done {

		// DEBUG: For testing pathfinding
		dest := state.dungeon.rooms[RoomID].Center()
		path1 = findPathBFS(state.dungeon, state.player.Pos(), dest)
		path2 = state.dmap.PathFrom(dest)

		// Draw the game world and refresh the display
		display.Clear()
		draw(&display, &state)
		drawDebug(&display, &state)
		display.Show()

		// Get user's command (this blocks until we get a key event)
		cmd = display.GetCommand(state.messages)

		// Handle user's command
		doUpdate = false
		switch cmd {

		// Commands that do not increment time
		case 0:
			// unknown command, just ignore
		case CmdTick:
			// Do nothing.  Used to redraw, clear recent messages, etc.
		case CmdMessages:
			display.DrawMessageHistory(state.messages)
			display.WaitForKeypress()
		case CmdInventory:
			display.InventoryScreen(state.player)
			display.WaitForKeypress()
		case CmdQuit:
			done = true

		// Commands that do increment time
		case CmdNorth:
			doUpdate = state.MoveActor(state.player, Coord{0, -1})
		case CmdNorthEast:
			doUpdate = state.MoveActor(state.player, Coord{1, -1})
		case CmdEast:
			doUpdate = state.MoveActor(state.player, Coord{1, 0})
		case CmdSouthEast:
			doUpdate = state.MoveActor(state.player, Coord{1, 1})
		case CmdSouth:
			doUpdate = state.MoveActor(state.player, Coord{0, 1})
		case CmdSouthWest:
			doUpdate = state.MoveActor(state.player, Coord{-1, 1})
		case CmdWest:
			doUpdate = state.MoveActor(state.player, Coord{-1, 0})
		case CmdNorthWest:
			doUpdate = state.MoveActor(state.player, Coord{-1, -1})
		case CmdDown:
			doUpdate = state.GoDownstairs()
		case CmdUp:
			doUpdate = state.GoUpstairs()
		case CmdWait:
			doUpdate = true
			//messages.Add("You rest for a moment.")

		case CmdConsume:
			if state.player.IsParalyzed() {
				state.messages.Add("You cannot consume anything while paralyzed.")
			} else {
				idx := display.PromptInventory("Consume what?", state.player)
				if idx != -1 {
					item := state.player.inventory[idx]
					switch item.(type) {
					case Consumable:
						doUpdate = item.(Consumable).Consume(&state)
						state.player.RemoveItem(idx)
					default:
						state.messages.Add("You cannot consume that item.")
					}
				}
			}

		case CmdEquip:
			if state.player.IsParalyzed() {
				state.messages.Add("You cannot equip anything while paralyzed.")
			} else {
				idx := display.PromptInventory("Equip or unequip what?", state.player)
				if idx != -1 {
					item := state.player.inventory[idx]
					switch item.(type) {
					case Equipable:
						doUpdate = item.(Equipable).Equip(state.player, state.messages)
					default:
						state.messages.Add("You cannot equip that item.")
					}
				}
			}

		// Extra debugging and testing stuff
		case CmdDebug1:
			debugFlag["main"] = !debugFlag["main"]
		case CmdDebug2:
			debugFlag["generate"] = !debugFlag["generate"]
		case CmdDebug3:
			debugFlag["dmap"] = !debugFlag["dmap"]
		case CmdDebug4:
			debugFlag["path"] = !debugFlag["path"]
		case CmdDebug5:
			RoomID++
			if RoomID >= len(state.dungeon.rooms) {
				RoomID = 0
			}
		case CmdGenerate:
			debug.Clear()
			GenerateTestLevel(&state)
		default:
			state.messages.Add("Unknown command.")
		}

		// Check for objects on the ground
		state.CheckItems()

		// Do updates of the game world
		state.Pathfinding()
		state.UpdatePlayerFOV()

		if doUpdate {
			state.PruneMonsters()
			state.player.Update(state.messages)
			if !state.IsBonusMove() {
				state.MonstersAct()
				state.WanderingMonsters()
			}
		}
	}
	display.Quit()
	fmt.Println("Thanks for playing!")
}

// -----------------------------------------------------------------------
func draw(display *Display, state *GameState) {

	if !state.player.IsBlind() {

		display.DrawMap(state.dungeon, debugFlag["main"])

		for pos, item := range state.items {
			if state.dungeon.TileAt(pos).visible || debugFlag["main"] {
				display.DrawItem(pos, item)
			}
		}

		for _, m := range *state.monsters {
			if state.dungeon.TileAt(m.Pos()).visible || debugFlag["main"] {
				display.DrawActor(m)
			}
		}
	}

	// monster detection should work even if blind
	if state.player.Timer("detMonsters") > 0 {
		for _, m := range *state.monsters {
			display.DrawActor(m)
		}
	}

	// if detect magic should work even if blind
	if state.player.Timer("detMagic") > 0 {
		for pos, item := range state.items {
			//if item.IsMagical() { //TODO
			display.DrawItem(pos, item)
			//}
		}
	}

	display.DrawMessages(state.messages)
	display.Print(0, 24, state.player.InfoString())

	display.DrawPlayer(state.player)
}

// -----------------------------------------------------------------------
func drawDebug(display *Display, state *GameState) {
	if debugFlag["main"] {
		drawDebugFrame(display, state)
		debug.Draw(display, 84, 15)
	}
	if debugFlag["generate"] {
		drawGenerateDebug(display)
	}
	if debugFlag["dmap"] {
		state.dmap.Draw(display)
	}
	if debugFlag["path"] {
		drawPathDebugIdx(display, path2)
	}
}

// -----------------------------------------------------------------------
func GenerateTestLevel(gs *GameState) {

	gs.dungeon.Clear()
	gs.monsters.Clear()
	gs.items.Clear()

	p1 := gs.dungeon.CreateRoom(Coord{44, 6}, 13, 7)
	p2 := gs.dungeon.CreateRoom(Coord{25, 15}, 11, 7)
	p3 := gs.dungeon.CreateRoom(Coord{18, 2}, 20, 7)
	gs.dungeon.ConnectRooms(p1, p3, East)
	gs.dungeon.ConnectRooms(p2, p3, South)

	gs.dungeon.SetTile(p2, TileStairsDn)
	gs.monsters.Add(randomMonster(gs.player.depth), Coord{20, 4})
	gs.monsters.Add(randomMonster(gs.player.depth), p2)
	gs.monsters.Add(randomMonster(gs.player.depth), p3)
	gs.monsters.Add(randomMonster(gs.player.depth), Coord{29, 17})

	gs.player.SetPos(p1)

	c := Coord{2, 1}
	gs.items[p1.Sum(c)] = newGold(randGoldAmt(gs.player.depth))
	//gs.items[p3.Sum(c)] = randWeapon()
	//gs.items[p2.Sum(c)] = randWeapon()
	//gs.player.depth++
}
