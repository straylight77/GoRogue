package main

import "fmt"

type Player struct {
	X, Y   int
	Symbol rune
	moves  int
	depth  int
}

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
	return fmt.Sprintf("Level: 1  Gold: 4       Hp: 11 (20)  Str: 16(16)  Arm: 4   Exp: 2/14")
}
