package main

import (
	"fmt"
	"math/rand"
)

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
	X, Y      int
	Symbol    rune
	moves     int
	depth     int
	HP        int
	maxHP     int
	Level     int
	XP        int
	healCount int
	foodCount int
}

func (p *Player) Init() {
	p.HP = 10
	p.maxHP = 10
	p.Level = 1
	p.foodCount = 2000
	p.ResetHealCount()

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

func (p *Player) Label() string {
	return "you"
}

func (p *Player) AdjustHP(amt int) {
	p.HP += amt
	if p.HP > p.maxHP {
		p.HP = p.maxHP
	}
}

func (p *Player) Attack(m Entity) string {
	dmg := 1
	m.AdjustHP(-dmg)
	msg := fmt.Sprintf("You hit %v for %d damage.", m.Label(), dmg)
	p.healCount++ // this shouldn't decrement when fighting
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
func (p *Player) ResetHealCount() {
	if p.Level < 8 {
		p.healCount = 21 - p.Level*2
	} else {
		p.healCount = 3
	}
}

// -----------------------------------------------------------------------
func (p *Player) Update() {

	// At 300 start being hungry, at 150 weak
	// At 0, every turn 20% chance you faint which paralyzes for 4-11 turns
	p.foodCount--

	// Levels 1-7, heal one point every [21-LVL*2] turns without fighting.
	// Levels 8+, heal between 1 and [LVL-7] points every three turns without fighting.
	// Note: Also see Attack()
	p.healCount--
	if p.healCount == 0 {
		if p.Level < 8 {
			p.AdjustHP(1)
		} else {
			amt := rand.Intn(p.Level - 7)
			p.AdjustHP(amt)
		}
		p.ResetHealCount()
	}

	p.moves++
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
