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
	AC        int
	HP        int
	maxHP     int
	Level     int
	XP        int
	Gold      int
	healCount int
	foodCount int
	inventory []*Item
	weapon    *Item
	armor     *Item
}

func (p *Player) Init() {
	p.HP = 10
	p.maxHP = 10
	p.AC = 0
	p.Level = 1
	p.foodCount = NutritionTime
	p.ResetHealCount()
}

// implement the Entity interface

func (p *Player) SetPos(newPos Coord) {
	p.X = newPos.X
	p.Y = newPos.Y
}

func (p *Player) Pos() Coord {
	return Coord{p.X, p.Y}
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

func (p *Player) AdjustFoodCount(amt int) {
	p.foodCount += amt
	if p.foodCount > NutritionTime {
		p.foodCount = NutritionTime
	}
}

// -----------------------------------------------------------------------
func (p *Player) Pickup(item *Item) bool {
	switch item.typ {
	case Gold:
		p.Gold += item.GoldQty()
		return true
	default:
		p.inventory = append(p.inventory, item)
		return true
	}
}

// -----------------------------------------------------------------------
func (p *Player) RemoveItem(idx int) {
	p.inventory = append(p.inventory[:idx], p.inventory[idx+1:]...)
}

// -----------------------------------------------------------------------
func (p *Player) AddXP(amt int) {
	p.XP += amt
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
	//debug.Add("level: xp=%d, ply=%d level=%d", p.XP, p.Level, level)
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
func (p *Player) Update(msg *MessageLog) {

	// At 300 start being hungry, at 150 weak
	// At 0, every turn 20% chance you faint which paralyzes for 4-11 turns
	f1 := p.foodCount
	p.foodCount--
	if f1 > HungerLimit && p.foodCount <= HungerLimit {
		msg.Add("You are starting to get hungry.")
	}
	if f1 > WeakLimit && p.foodCount <= WeakLimit {
		msg.Add("You are starting to feel weak.")
		//TODO: handle feinting from hunger
	}

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
		"Level: %-2d  Gold: %-5d  Hp: %2d(%2d)  Str: 16(16)  Arm: %-2d   Exp: %d/%d",
		p.depth,
		p.Gold,
		p.HP,
		p.maxHP,
		p.AC,
		p.Level,
		p.XP,
	)
	return info
}
