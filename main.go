package main

import "math/rand"

// wrap these into GameState?  Will have handleCommand()?
var dungeon DungeonMap
var player Player
var monsters MonsterList

var messages MessageLog

var disp Display
var debug = DebugMessageLog{}

var path1 Path
var path2 Path
var dmap *DMap

var debugFlag = map[string]bool{
	"main":     true,
	"generate": false,
	"dmap":     false,
	"path":     false,
}

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
	CmdSouth
	CmdEast
	CmdWest
	CmdUp
	CmdDown
	CmdWait
	CmdTick     // redundant with CmdNop?
	CmdGenerate // for testing
	CmdMessages
)

// -----------------------------------------------------------------------
func movePlayer(p *Player, dx int, dy int, d *DungeonMap, mlist *MonsterList) {
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
		m.State = StateChase
		return
	}

	// check dungeon tile
	destTile := d.TileAt(destX, destY)
	if destTile.IsWalkable() {
		p.SetPos(destX, destY)
	}

}

// -----------------------------------------------------------------------
func main() {
	var cmd GameCommand

	// Initialization and setup
	disp = Display{}
	disp.Init()
	defer disp.Quit()

	player.Init()

	// Create a dungeon level
	//GenerateTestLevel(&dungeon, &player, &monsters)
	generateRandomLevel(&dungeon, &monsters, &player)

	doneFlag := false
	var doUpdate bool

	for !doneFlag {
		doUpdate = true

		// ===== Test some pathfinding stuff ====
		pathX, pathY := dungeon.rooms[0].Center()
		path1 = findPathBFS(&dungeon, player.X, player.Y, pathX, pathY)

		// ==== Testing Dijkstra Maps ====
		dmap = newDMap()
		dmap.AddTarget(Coord{player.X, player.Y})
		dmap.Calculate(&dungeon)
		path2 = dmap.PathFrom(Coord{pathX, pathY})

		// Draw the world
		disp.Clear()
		disp.DrawMap(&dungeon, debugFlag["main"])
		disp.DrawMessages(&messages)
		disp.Print(0, 24, player.InfoString())

		for _, m := range monsters {
			mx, my := m.Pos()
			if dungeon.TileAt(mx, my).visible || debugFlag["main"] {
				disp.DrawEntity(m)
			}
		}
		disp.DrawPlayer(&player)

		if debugFlag["main"] {
			drawDebugFrame(&disp, &player, &monsters)
			debug.Draw(&disp, 84, 15)
		}
		if debugFlag["generate"] {
			drawGenerateDebug(&disp)
		}
		if debugFlag["dmap"] {
			dmap.Draw(&disp)
		}
		if debugFlag["path"] {
			//drawPathDebug(&disp, path1, 'x')
			drawPathDebug(&disp, path2, '*')
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
			movePlayer(&player, -1, 0, &dungeon, &monsters)
		case CmdEast:
			movePlayer(&player, 1, 0, &dungeon, &monsters)
		case CmdNorth:
			movePlayer(&player, 0, -1, &dungeon, &monsters)
		case CmdSouth:
			movePlayer(&player, 0, 1, &dungeon, &monsters)
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
			//messages.Add("You rest for a moment.")

		// Extra debugging and testing stuff
		case CmdDebug1:
			debugFlag["main"] = !debugFlag["main"]
			doUpdate = false
		case CmdDebug2:
			debugFlag["generate"] = !debugFlag["generate"]
			doUpdate = false
		case CmdDebug3:
			debugFlag["dmap"] = !debugFlag["dmap"]
			doUpdate = false
		case CmdDebug4:
			debugFlag["path"] = !debugFlag["path"]
			doUpdate = false

		case CmdGenerate:
			doUpdate = false
			debug.Clear()
			generateRandomLevel(&dungeon, &monsters, &player)
			//GenerateTestLevel(&dungeon, &player, &monsters)
		}

		// Update the player's field of view and visited tiles
		dungeon.SetVisible(0, 0, MapMaxX, MapMaxY, false)
		dungeon.playerFOV(&player)

		// If the player is in a room, light it up
		for _, r := range dungeon.rooms {
			if r.InRoom(player.X, player.Y) {
				dungeon.SetVisible(r.X, r.Y, r.W+1, r.H+1, true)
			}
		}

		// Do world updates
		if doUpdate {
			updateMonsters(&dungeon, &player, &monsters, &messages)
			player.Update()
		}

	}
}

// -----------------------------------------------------------------------
func updateMonsters(d *DungeonMap, p *Player, ml *MonsterList, msg *MessageLog) {
	for i, m := range *ml {

		// Remove any slain monsters
		if m.HP <= 0 {
			ml.Remove(i)
			msg.Add("You defeated the %s!", m.Name)
			m := p.AddXP(m.XP)
			if m != "" {
				msg.Add(m)
			}
			continue
		}

		switch m.State {

		case StateDormant, StateActive:
			if (m.isMean && playerCanSee(m, &dungeon)) || m.DistanceFrom(&player) <= 1 {
				m.State = StateChase
			}

		case StateChase:

			if !m.isMean && m.DistanceFrom(&player) > 6 {
				// For non-mean monsters, go dormant when far away
				m.State = StateDormant

			} else if m.randMove > rand.Intn(100) {
				// Move randomly
				moved := false
				count := 0
				for !moved && count < 8 {
					dx, dy := randDirectionCoords()
					moved = moveMonster(m, dx, dy, &dungeon, &player, &monsters)
					count++
				}

			} else {
				// Do pathfinding to the player and take the first step
				// What happens when player goes out of sight?
				m.path = findPathBFS(&dungeon, m.X, m.Y, player.X, player.Y)
				dx, dy := m.DirectionCoordsTo(m.path.steps[0].X, m.path.steps[0].Y)
				moveMonster(m, dx, dy, &dungeon, &player, &monsters)
			}
		}

	}
}

// -----------------------------------------------------------------------
func randDirectionCoords() (x, y int) {
	count := 0
	for (x == 0 && y == 0) && count < 10 {
		x = rand.Intn(3) - 1 // between -1 and 1
		y = rand.Intn(3) - 1 // between -1 and 1
		count++
	}
	return
}

// -----------------------------------------------------------------------
// Returns weather the monster actually did something (moved or attacked)
func moveMonster(m *Monster, dx, dy int, d *DungeonMap, p *Player, mlist *MonsterList) bool {
	destX, destY := m.X+dx, m.Y+dy
	//debug.Add("move: to %d, %d", destX, destY)

	// Check edges of the map
	if destX < 0 || destX >= MapMaxX || destY < 0 || destY >= MapMaxY {
		return false
	}

	// Check if player is there
	if destX == p.X && destY == p.Y {
		messages.Add(m.Attack(p))
		return true
	}

	// Check for other monsters
	m2 := mlist.MonsterAt(destX, destY)
	if m2 != nil {
		return false
	}

	// Check dungeon tile
	destTile := d.TileAt(destX, destY)
	if destTile.IsWalkable() {
		m.SetPos(destX, destY)
		return true
	} else {
		return false
	}
}

// -----------------------------------------------------------------------
func playerCanSee(e Entity, d *DungeonMap) bool {
	eX, eY := e.Pos()
	t := d.TileAt(eX, eY)
	return t.visible
}

// -----------------------------------------------------------------------
func GenerateTestLevel(m *DungeonMap, p *Player, ml *MonsterList) {

	m.Clear()
	ml.Clear()

	x1, y1 := m.CreateRoom(44, 6, 13, 7)
	x2, y2 := m.CreateRoom(25, 15, 11, 7)
	x3, y3 := m.CreateRoom(18, 2, 20, 7)
	m.ConnectRooms(x1, y1, x3, y3, East)
	m.ConnectRooms(x2, y2, x3, y3, South)
	//m.ConnectRooms(x1, y1, x2, y2, South)

	//m.SetTile(x1, y1, TileStairsUp)
	m.SetTile(x2, y2, TileStairsDn)
	monsters.Add(randomMonster(player.depth), 20, 4)
	monsters.Add(randomMonster(player.depth), x2, y2)
	monsters.Add(randomMonster(player.depth), x3, y3)
	monsters.Add(randomMonster(player.depth), 29, 17)
	//monsters.Add(newMonster(2), 44, 5)

	p.SetPos(x1, y1)
	p.depth++

}
