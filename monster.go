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
	{'K', 0, 2, 1, 7, "1d4", "kobold", true, 0},
	{'J', 0, 2, 1, 7, "1d2", "jackal", true, 0},
	{'B', 0, 1, 1, 3, "1d2", "bat", false, 50}, // 50% chance to move randomly
	{'S', 0, 3, 1, 5, "1d3", "snake", true, 0},
	{'H', 0, 3, 1, 5, "1d8", "hobgoblin", true, 0},
	{'E', 0, 5, 1, 9, "0d0", "floating eye", false, 0}, // paralyzes 2-3 turns
	{'A', 0, 10, 2, 3, "1d6", "giant ant", true, 0},    // decrease str
	{'O', 15, 5, 1, 6, "1d7", "orc", true, 0},
	{'Z', 0, 7, 2, 8, "1d8", "zombie", true, 0},
	{'G', 10, 8, 1, 5, "1d6", "gnome", false, 0},
	{'L', 0, 10, 3, 8, "1d1", "leprechaun", false, 0}, // steal gold unless save vs magic
	{'C', 15, 15, 4, 4, "1d6/1d6", "centaur", false, 0},
	{'R', 0, 25, 5, 2, "0d0/0d0", "rust monster", true, 0}, // -1 to armor being worn
	{'Q', 30, 35, 3, 2, "1d2/1d2/1d4", "quasit", true, 0},
	{'N', 100, 40, 3, 9, "0d0", "nymph", false, 0}, // steals random magic item from inventory
	{'Y', 30, 50, 4, 6, "1d6/1d6", "yeti", false, 0},
	{'T', 50, 55, 6, 4, "1d8/1d8/2d6", "troll", true, 0},
	{'W', 0, 55, 5, 4, "1d6", "wraith", true, 0},              // 15% chance to drain level and 1d10 max hp
	{'F', 0, 85, 8, 3, "0d0", "violet fungi", true, 0},        // grapple, damage is 1 then 2 then 3 etc.
	{'I', 0, 120, 8, 3, "4d4", "invisible stalker", true, 20}, // 20% chance to move randomly
	{'X', 0, 120, 7, -2, "1d3/1d3/1d3/4d6", "xorn", true, 0},
	{'U', 40, 130, 8, 2, "3d4/3d4/2d5", "umber hulk", true, 0}, // confuses for 20-39 turns, only once
	{'M', 30, 140, 7, 7, "3d4", "mimic", false, 0},
	{'V', 30, 380, 8, 1, "1d10", "vampire", true, 0},
	{'D', 100, 9000, 10, -1, "1d8/1d8/3d10", "dragon", false, 0},
	{'P', 70, 7000, 15, 6, "2d12/2d4", "purple worm", false, 0},
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
	debug.Add("monster: len=%d, min=%d, max=%d, idx=%d", len(MonsterLib), min, max, idx)
	return newMonster(idx)
}

/*************************************************************************
 * MonsterList
 *
 */

type MonsterList []*Monster

func (ml *MonsterList) Add(m *Monster, x, y int) {
	m.X, m.Y = x, y
	*ml = append(*ml, m)
}

func (ml *MonsterList) Remove(idx int) {
	*ml = append((*ml)[:idx], (*ml)[idx+1:]...)
}

func (ml *MonsterList) Clear() {
	*ml = nil
}

func (ml MonsterList) MonsterAt(x, y int) *Monster {
	for _, m := range ml {
		if m.X == x && m.Y == y {
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
	State    int
	isMean   bool // once visible, start chasing player
	isGreedy bool // move towards any nearby gold
	randMove int
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
		HP:       2,
		Symbol:   mt.Symbol,
		isMean:   mt.isMean,
		randMove: mt.randMove,
	}
	return m
}

func (m *Monster) DebugString() string {
	return fmt.Sprintf("%c (%d,%d) hp=%d state=%d", m.Symbol, m.X, m.Y, m.HP, m.State)
}

// Returns the Chebyshev Distance from the given Entity
func (m *Monster) DistanceFrom(e Entity) int {
	x2, y2 := e.Pos()
	dx := abs(x2 - m.X)
	dy := abs(y2 - m.Y)
	return max(dx, dy)
}

func (m *Monster) DirectionCoordsTo(e Entity) (dx int, dy int) {
	eX, eY := e.Pos()

	dx = 0
	if eX < m.X {
		dx = -1
	} else if eX > m.X {
		dx = 1
	}
	dy = 0
	if eY < m.Y {
		dy = -1
	} else if eY > m.Y {
		dy = 1
	}

	return dx, dy
}

func (m *Monster) Attack(p *Player) string {
	p.HP += -1
	return fmt.Sprintf("The %v attacks.", m)
}

func (m Monster) String() string {
	return m.Name
}

// Implement the Entity interface

func (m *Monster) SetPos(newX, newY int) {
	m.X = newX
	m.Y = newY
}

func (m *Monster) Pos() (int, int) {
	return m.X, m.Y
}

func (m *Monster) Rune() rune {
	return m.Symbol
}
