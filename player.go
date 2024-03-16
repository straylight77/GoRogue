package main

import "fmt"

type Player struct {
	X, Y   int
	Symbol rune
	moves  int
	depth  int
	HP     int
	maxHP  int
}

func (p *Player) Init() {
	p.HP = 10
	p.maxHP = 10
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
func (p *Player) InfoString() string {
	info := fmt.Sprintf(
		"Level:%2d  Gold: 4       Hp:%2d(%2d)  Str:16(16)  Arm: 4   Exp: 2/14",
		p.depth,
		p.HP,
		p.maxHP,
	)
	return info
}
