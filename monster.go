package main

import (
	"fmt"
	"math/rand"
)

/*************************************************************************
 * MonsterLib
 *
 */
type MonsterTemplate struct {
	Symbol   rune
	Carry    int
	XP       int
	Level    int
	AC       int
	Dmg      string
	Name     string
	isMean   bool
	noWander bool
	randMove int // chance that it will move randomly (percentage)
}

// Index is used as difficulty of the monsters
//
//	min = depth - 6
//	max = depth + 3
//
// (apparently called "vorpalness" in original Rogue source code)
// https://datadrivengamer.blogspot.com/2019/05/identifying-mechanics-of-rogue.html
var MonsterLib = []MonsterTemplate{
	{'K', 0, 2, 1, 7, "1d4", "kobold", true, false, 0},
	{'J', 0, 2, 1, 7, "1d2", "jackal", true, false, 0},
	{'B', 0, 1, 1, 3, "1d2", "bat", false, false, 50}, // 50% chance to move randomly
	{'S', 0, 3, 1, 5, "1d3", "snake", true, false, 0},
	{'H', 0, 3, 1, 5, "1d8", "hobgoblin", true, false, 0},
	{'E', 0, 5, 1, 9, "0d0", "floating eye", false, true, 0}, // paralyzes 2-3 turns
	{'A', 0, 10, 2, 3, "1d6", "giant ant", true, false, 0},   // decrease str
	{'O', 15, 5, 1, 6, "1d7", "orc", true, false, 0},
	{'Z', 0, 7, 2, 8, "1d8", "zombie", true, false, 0},
	{'G', 10, 8, 1, 5, "1d6", "gnome", false, false, 0},
	{'L', 0, 10, 3, 8, "1d1", "leprechaun", false, true, 0}, // steal gold unless save vs magic
	{'C', 15, 15, 4, 4, "1d6/1d6", "centaur", false, false, 0},
	{'R', 0, 25, 5, 2, "0d0/0d0", "rust monster", true, false, 0}, // -1 to armor being worn
	{'Q', 30, 35, 3, 2, "1d2/1d2/1d4", "quasit", true, false, 0},
	{'N', 100, 40, 3, 9, "0d0", "nymph", false, true, 0}, // steals random magic item from inventory
	{'Y', 30, 50, 4, 6, "1d6/1d6", "yeti", false, false, 0},
	{'T', 50, 55, 6, 4, "1d8/1d8/2d6", "troll", true, true, 0},
	{'W', 0, 55, 5, 4, "1d6", "wraith", true, false, 0},              // 15% chance to drain level and 1d10 max hp
	{'F', 0, 85, 8, 3, "0d0", "violet fungi", true, true, 0},         // grapple, damage is 1 then 2 then 3 etc.
	{'I', 0, 120, 8, 3, "4d4", "invisible stalker", true, false, 20}, // 20% chance to move randomly
	{'X', 0, 120, 7, -2, "1d3/1d3/1d3/4d6", "xorn", true, false, 0},
	{'U', 40, 130, 8, 2, "3d4/3d4/2d5", "umber hulk", true, false, 0}, // confuses for 20-39 turns, only once
	{'M', 30, 140, 7, 7, "3d4", "mimic", false, true, 0},
	{'V', 30, 380, 8, 1, "1d10", "vampire", true, false, 0},
	{'D', 100, 9000, 10, -1, "1d8/1d8/3d10", "dragon", false, true, 0},
	{'P', 70, 7000, 15, 6, "2d12/2d4", "purple worm", false, true, 0},
}

// Uses public variable MonsterLib
func randomMonster(depth int) *Monster {
	min := depth - 6
	max := depth + 3
	if min < 0 {
		min = 0
	}
	if max > len(MonsterLib) {
		max = len(MonsterLib)
	}

	idx := len(MonsterLib) - 1 // Default to most difficult monster
	if min < len(MonsterLib) { // Ensure we don't go out of bounds
		idx = rand.Intn(max-min) + min
	}
	//debug.Add("monster: len=%d, min=%d, max=%d, idx=%d", len(MonsterLib), min, max, idx)
	return newMonster(idx)
}

/*************************************************************************
 * MonsterList
 *
 */

type MonsterList []*Monster

func (ml *MonsterList) Add(m *Monster, pos Coord) {
	m.X, m.Y = pos.XY()
	*ml = append(*ml, m)
}

func (ml *MonsterList) Remove(idx int) {
	*ml = append((*ml)[:idx], (*ml)[idx+1:]...)
}

func (ml *MonsterList) Clear() {
	*ml = nil
}

func (ml MonsterList) MonsterAt(pos Coord) *Monster {
	for _, m := range ml {
		if m.Pos() == pos {
			return m
		}
	}
	return nil
}

/*************************************************************************
 * Monster
 * implements Entity interface
 */

type Monster struct {
	X, Y     int
	Symbol   rune
	Name     string
	HP       int
	XP       int
	State    int
	isMean   bool // once visible, start chasing player
	isGreedy bool // move towards any nearby gold
	noWander bool
	randMove int
	nextStep Coord
}

const (
	StateDormant = iota
	StateActive
	StateChase
	StateWander
)

func newMonster(id int) *Monster {
	mt := MonsterLib[id]
	m := &Monster{
		Name:     mt.Name,
		HP:       mt.Level + 1,
		XP:       mt.XP,
		Symbol:   mt.Symbol,
		isMean:   mt.isMean,
		noWander: mt.noWander,
		randMove: mt.randMove,
	}
	return m
}

func (m *Monster) DebugString() string {
	return fmt.Sprintf(
		"%c (%d,%d) hp=%d state=%d, step=%v",
		m.Symbol,
		m.X,
		m.Y,
		m.HP,
		m.State,
		m.nextStep)
}

func (m *Monster) DirectionCoordsTo(pos Coord) Coord {
	dx := 0
	if pos.X < m.X {
		dx = -1
	} else if pos.X > m.X {
		dx = 1
	}
	dy := 0
	if pos.Y < m.Y {
		dy = -1
	} else if pos.Y > m.Y {
		dy = 1
	}

	return Coord{dx, dy}
}

func (m Monster) String() string {
	return m.Name
}

// Implement the Entity interface

func (m *Monster) SetPos(newPos Coord) {
	m.X = newPos.X
	m.Y = newPos.Y
}

func (m *Monster) Pos() Coord {
	return Coord{m.X, m.Y}
}

func (m *Monster) Rune() rune {
	return m.Symbol
}

func (m *Monster) Label() string {
	return "the " + m.Name
}

func (m *Monster) AdjustHP(amt int) {
	m.HP += amt
}

func (m *Monster) Attack(e Entity) string {
	e.AdjustHP(-1)
	return fmt.Sprintf("The %v attacks %s.", m, e.Label())
}
