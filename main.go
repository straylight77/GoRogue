package main

import (
	"fmt"
	"strings"
)

// wrap these into GameState?  Will have handleCommand()?
var dungeon DungeonMap
var player Player
var monsters MonsterList

var messages MessageLog

var disp Display

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
	CmdMessages
)

// -----------------------------------------------------------------------
type MessageLog struct {
	messages []string
	idx      int
}

func (log *MessageLog) Add(format string, vals ...any) {
	msg := fmt.Sprintf(format, vals...)
	log.messages = append(log.messages, msg)
}

func (log *MessageLog) Clear() {
	log.messages = nil
}

func (log *MessageLog) HasUnread() bool {
	return log.idx < len(log.messages)
}

func (log *MessageLog) Last(n int) []string {
	if n >= len(log.messages) {
		return log.messages
	} else {
		return log.messages[len(log.messages)-n:]
	}
}

func (log *MessageLog) LatestAsStr() string {
	s := ""
	if len(log.messages[log.idx:]) > 0 {
		s = strings.Join(log.messages[log.idx:], " ")
	}
	return s
}

func (log *MessageLog) ClearUnread() {
	log.idx = len(log.messages)
}

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

	player.moves++
}

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

	debug := false
	done := false
	for !done {

		// draw the world
		disp.Clear()
		disp.DrawMap(&dungeon)
		disp.DrawMessages(&messages)
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
		case CmdMessages:
			disp.DrawMessageHistory(&messages)
			disp.WaitForKeypress()

		// extra debugging and testing stuff
		case CmdDebug:
			debug = !debug
		case CmdGenerate:
			generateRandomLevel(&dungeon, &monsters, &player)

		}
		// do other world updates
		for i, m := range monsters {
			if m.HP <= 0 {
				monsters.Remove(i)
				messages.Add("You defeated the %s!", m.Name)
			}
		}

	}
}
