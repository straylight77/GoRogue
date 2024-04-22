package main

import (
	"fmt"
	"math/rand"
	"strings"
)

/*************************************************************************
 * MonsterLib
 *
 */
type MonsterTemplate struct {
	Symbol      rune
	Carry       int
	XP          int
	Level       int
	AC          int
	Attacks     string
	AttackVerbs string
	Name        string
	isMean      bool
	noWander    bool
	randMove    int // chance that it will move randomly (percentage)
}

// Index is used as difficulty of the monsters
//
//	min = depth - 6
//	max = depth + 3
//
// (apparently called "vorpalness" in original Rogue source code)
// https://datadrivengamer.blogspot.com/2019/05/identifying-mechanics-of-rogue.html
var MonsterLib = []MonsterTemplate{
	{'K', 0, 2, 1, 7, "1d4", "swings at", "kobold", true, false, 0},
	{'J', 0, 2, 1, 7, "1d2", "bites", "jackal", true, false, 0},
	{'B', 0, 1, 1, 3, "1d2", "bites", "bat", false, false, 50}, // 50% chance to move randomly
	{'S', 0, 3, 1, 5, "1d3", "bites", "snake", true, false, 0},
	{'H', 0, 3, 1, 5, "1d8", "swings at", "hobgoblin", true, false, 0},
	{'E', 0, 5, 1, 9, "0d0", "gazes at", "floating eye", false, true, 0}, // paralyzes 2-3 turns
	{'A', 0, 10, 2, 3, "1d6", "stings", "giant ant", true, false, 0},     // decrease str
	{'O', 15, 5, 1, 6, "1d7", "attacks", "orc", true, false, 0},
	{'Z', 0, 7, 2, 8, "1d8", "slams", "zombie", true, false, 0},
	{'G', 10, 8, 1, 5, "1d6", "attacks", "gnome", false, false, 0},
	{'L', 0, 10, 3, 8, "1d1", "pickpockets", "leprechaun", false, true, 0}, // steal gold unless save vs magic
	{'C', 15, 15, 4, 4, "1d6/1d6", "kicks/kicks", "centaur", false, false, 0},
	{'R', 0, 25, 5, 2, "0d0/0d0", "bites/bites", "rust monster", true, false, 0}, // -1 to armor being worn
	{'Q', 30, 35, 3, 2, "1d2/1d2/1d4", "claws/claws/bites", "quasit", true, false, 0},
	{'N', 100, 40, 3, 9, "0d0", "pickpockets", "nymph", false, true, 0}, // steals random magic item from inventory
	{'Y', 30, 50, 4, 6, "1d6/1d6", "swings/swings", "yeti", false, false, 0},
	{'T', 50, 55, 6, 4, "1d8/1d8/2d6", "claws/claws/bites", "troll", true, true, 0},
	{'W', 0, 55, 5, 4, "1d6", "touches", "wraith", true, false, 0},                // 15% chance to drain level and 1d10 max hp
	{'F', 0, 85, 8, 3, "0d0", "sqeezes", "violet fungi", true, true, 0},           // grapple, damage is 1 then 2 then 3 etc.
	{'I', 0, 120, 8, 3, "4d4", "swings at", "invisible stalker", true, false, 20}, // 20% chance to move randomly
	{'X', 0, 120, 7, -2, "1d3/1d3/1d3/4d6", "claws/claws/claws/bites", "xorn", true, false, 0},
	{'U', 40, 130, 8, 2, "3d4/3d4/2d5", "claws/claws/bites", "umber hulk", true, false, 0}, // confuses for 20-39 turns, only once
	{'M', 30, 140, 7, 7, "3d4", "bites", "mimic", false, true, 0},
	{'V', 30, 380, 8, 1, "1d10", "bites", "vampire", true, false, 0},
	{'D', 100, 9000, 10, -1, "1d8/1d8/3d10", "claws/claws/bites", "dragon", false, true, 0},
	{'P', 70, 7000, 15, 6, "2d12/2d4", "bites/stings", "purple worm", false, true, 0},
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
 * implements Actor interface
 */

type Monster struct {
	X, Y        int
	Symbol      rune
	Name        string
	Level       int
	HP          int
	AC          int
	Attacks     []Dice
	AttackVerbs []string
	XP          int
	State       int
	isMean      bool // once visible, start chasing player
	isGreedy    bool // move towards any nearby gold
	noWander    bool
	randMove    int
	nextStep    Coord
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
		Name:        mt.Name,
		Level:       mt.Level,
		HP:          mt.Level * (rand.Intn(8) + 1),
		AC:          mt.AC,
		Attacks:     parseDice(mt.Attacks),
		AttackVerbs: strings.Split(mt.AttackVerbs, "/"),
		XP:          mt.XP,
		Symbol:      mt.Symbol,
		isMean:      mt.isMean,
		noWander:    mt.noWander,
		randMove:    mt.randMove,
	}
	return m
}

func (m *Monster) DebugString() string {
	return fmt.Sprintf(
		"%c (%2d,%2d) hp=%-2d ac=%-2d thac0=%-2d s=%d step=%v",
		m.Symbol,
		m.X, m.Y,
		m.HP,
		m.AC,
		m.ToHit(),
		m.State,
		m.nextStep,
	)
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

// ----------------------------------------------------------------------
// Implement the Actor interface

func (m *Monster) Pos() Coord {
	return Coord{m.X, m.Y}
}

func (m *Monster) SetPos(newPos Coord) {
	m.X = newPos.X
	m.Y = newPos.Y
}

func (m *Monster) Rune() rune {
	return m.Symbol
}

func (m *Monster) AdjustHP(amt int) {
	m.HP += amt
}

func (m *Monster) Attack(a Actor, msg *MessageLog) {

	var label string
	if a.IsBlind() {
		label = "Something"
	} else {
		label = fmt.Sprintf("The %v", m)
	}

	//debug.Add("attack: %v", m.Attacks)
	for i, atk := range m.Attacks {
		if attackHits(m.ToHit(), a.ArmorClass()) {
			dmg := atk.Roll()
			a.AdjustHP(-dmg)
			msg.Add("%v %s you for %d damage.", label, m.AttackVerbs[i], dmg)
		} else {
			msg.Add("%v misses you.", label)
		}
	}
}

func (m *Monster) ArmorClass() int {
	return m.AC
}

func (m *Monster) IsConfused() bool {
	//TODO implement monster confusion
	return false
}

func (m *Monster) IsBlind() bool {
	//TODO implement monster blindness
	return false
}

// ----------------------------------------------------------------------

func (m *Monster) ToHit() int {
	return 21 - m.Level
}

func (m *Monster) RollDamage() int {
	//return m.Dmg.Roll()
	return 1
}
