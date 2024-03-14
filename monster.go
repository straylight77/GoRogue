package main

import "fmt"

/*************************************************************************
 * MonsterLib
 *
 */
type MonsterTemplate struct {
	Name     string
	Symbol   rune
	HP       int
	AC       int
	Dmg      int
	XP       int
	MinDepth int
	MaxDepth int
}

var MonsterLib = []MonsterTemplate{
	{"bat", 'B', 3, 17, 1, 1, 1, 7},
	{"kestrel", 'K', 4, 13, 2, 1, 1, 5},
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
	X, Y   int
	Symbol rune
	Name   string
	HP     int
}

func NewMonster(id int) *Monster {
	mt := MonsterLib[id]
	m := &Monster{
		Name: mt.Name,
		HP:   mt.HP,
	}
	m.Symbol = mt.Symbol
	return m
}

func CreateMonster(n string, sym rune, hp int) *Monster {
	newMonster := &Monster{
		//Symbol: sym,
		Name: n,
		HP:   hp,
	}
	newMonster.Symbol = sym
	return newMonster
}

func (m *Monster) DebugString() string {
	return fmt.Sprintf("%s x=%d y=%d hp=%d", m.Name, m.X, m.Y, m.HP)
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
