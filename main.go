package main

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

type GameCommand int

const (
	CmdNop GameCommand = iota
	CmdDebug1
	CmdDebug2
	CmdDebug3
	CmdDebug4
	CmdDebug5
	CmdQuit

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
	CmdWait
	CmdTick
	CmdGenerate // for testing
	CmdMessages
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

		// Draw the world
		draw(&display, &state)
		drawDebug(&display, &state)
		display.Show()

		cmd = display.GetCommand(state.messages)
		doUpdate = false

		// Handle user's command
		switch cmd {

		// Commands that do not increment time
		case 0: // unknown command, just ignore
		case CmdTick:
			// Do nothing.  Used to redraw, clear recent messages, etc.
		case CmdMessages:
			display.DrawMessageHistory(state.messages)
			display.WaitForKeypress()
		case CmdQuit:
			done = true

		// Commands that do increment time
		case CmdNorth:
			doUpdate = state.MoveEntity(state.player, Coord{0, -1})
		case CmdNorthEast:
			doUpdate = state.MoveEntity(state.player, Coord{1, -1})
		case CmdEast:
			doUpdate = state.MoveEntity(state.player, Coord{1, 0})
		case CmdSouthEast:
			doUpdate = state.MoveEntity(state.player, Coord{1, 1})
		case CmdSouth:
			doUpdate = state.MoveEntity(state.player, Coord{0, 1})
		case CmdSouthWest:
			doUpdate = state.MoveEntity(state.player, Coord{-1, 1})
		case CmdWest:
			doUpdate = state.MoveEntity(state.player, Coord{-1, 0})
		case CmdNorthWest:
			doUpdate = state.MoveEntity(state.player, Coord{-1, -1})
		case CmdDown:
			doUpdate = state.GoDownstairs()
		case CmdUp:
			doUpdate = state.GoUpstairs()
		case CmdWait:
			doUpdate = true
			//messages.Add("You rest for a moment.")

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
			//generateRandomLevel(&state)
			GenerateTestLevel(&state)
		}

		// Do updates that happen regardless of game time
		state.Pathfinding()
		state.UpdatePlayerFOV()

		// Do updates of the game world
		if doUpdate {
			state.UpdateMonsters()
			state.player.Update()
		}

		// For testing pathfinding
		dest := state.dungeon.rooms[RoomID].Center()
		path1 = findPathBFS(state.dungeon, state.player.X, state.player.Y, dest.X, dest.Y)
		path2 = state.dmap.PathFrom(dest)
	}
}

// -----------------------------------------------------------------------
func draw(display *Display, state *GameState) {
	display.Clear()
	display.DrawMap(state.dungeon, debugFlag["main"])
	display.DrawMessages(state.messages)
	display.Print(0, 24, state.player.InfoString())

	for _, m := range *state.monsters {
		if state.dungeon.TileAt(m.Pos()).visible || debugFlag["main"] {
			display.DrawEntity(m)
		}
	}
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

	p1 := gs.dungeon.CreateRoom(Coord{44, 6}, 13, 7)
	p2 := gs.dungeon.CreateRoom(Coord{25, 15}, 11, 7)
	p3 := gs.dungeon.CreateRoom(Coord{18, 2}, 20, 7)
	gs.dungeon.ConnectRooms(p1, p3, East)
	gs.dungeon.ConnectRooms(p2, p3, South)

	//gs.dungeon.SetTile(x1, y1, TileStairsUp)
	gs.dungeon.SetTile(p2, TileStairsDn)
	gs.monsters.Add(randomMonster(gs.player.depth), Coord{20, 4})
	gs.monsters.Add(randomMonster(gs.player.depth), p2)
	gs.monsters.Add(randomMonster(gs.player.depth), p3)
	gs.monsters.Add(randomMonster(gs.player.depth), Coord{29, 17})
	//gs.monsters.Add(newMonster(2), 44, 5)

	gs.player.SetPos(p1)
	//gs.player.depth++

}
