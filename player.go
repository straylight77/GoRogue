package main

import "fmt"

var XPTable = [21]int{
	0,
	10,
	20,
	40,
	80,
	160,
	320,
	640,
	1300,
	2600,
	5200,
	13000,
	26000,
	50000,
	100000,
	200000,
	400000,
	800000,
	2000000,
	4000000,
	8000000,
}

type Player struct {
	X, Y   int
	Symbol rune
	moves  int
	depth  int
	HP     int
	maxHP  int
	Level  int
	XP     int
}

func (p *Player) Init() {
	p.HP = 10
	p.maxHP = 10
	p.Level = 1
}

// implement the Entity interface

func (p *Player) SetPos(newX, newY int) {
	p.X = newX
	p.Y = newY
}

func (p *Player) Pos() (int, int) {
	return p.X, p.Y
}

func (p *Player) Rune() rune {
	return p.Symbol
}

// -----------------------------------------------------------------------
func (p *Player) Attack(m *Monster) string {
	dmg := 1
	m.HP -= dmg
	msg := fmt.Sprintf("You hit the %v for %d damage.", m.Name, dmg)
	return msg
}

// -----------------------------------------------------------------------
func (p *Player) AddXP(amt int) string {
	p.XP += amt
	return p.CheckLevel()
}

// -----------------------------------------------------------------------
func (p *Player) CheckLevel() string {
	msg := ""
	level := 0
	for _, xp := range XPTable {
		if p.XP < xp {
			break
		}
		level++
	}
	if p.Level < level {
		msg = fmt.Sprintf("Welcome to level %d!", level)
	}
	p.Level = level
	return msg
}

// -----------------------------------------------------------------------
func (p *Player) InfoString() string {
	info := fmt.Sprintf(
		"Level:%2d  Gold: 4       Hp:%2d(%2d)  Str:16(16)  Arm: 4   Exp: %d/%d",
		p.depth,
		p.HP,
		p.maxHP,
		p.Level,
		p.XP,
	)
	return info
}
