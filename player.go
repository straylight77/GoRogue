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
	Str       int
	maxStr    int
	Level     int
	XP        int
	Gold      int
	healCount int
	foodCount int
	inventory []Object
	equiped   map[string]Object
	timer     map[string]int
}

func (p *Player) Init() {
	p.Str = 16
	p.maxStr = 16
	p.HP = 10
	p.maxHP = 10
	p.AC = 10
	p.Level = 1
	p.foodCount = NutritionTime
	p.timer = make(map[string]int)
	p.equiped = map[string]Object{
		"weapon": nil,
		"armor":  nil,
		"left":   nil,
		"right":  nil,
	}
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
	p.healCount++ // this shouldn't decrement when fighting
	if p.IsBlind() {
		return fmt.Sprintf("You hit something for %d damage.", dmg)
	} else {
		return fmt.Sprintf("You hit %v for %d damage.", m.Label(), dmg)
	}
}

func (p *Player) AdjustFoodCount(amt int) {
	p.foodCount += amt
	if p.foodCount > NutritionTime {
		p.foodCount = NutritionTime
	}
}

func (p *Player) IsConfused() bool {
	return p.timer["confused"] > 0
}

func (p *Player) IsBlind() bool {
	return p.timer["blind"] > 0
}

func (p *Player) IsParalyzed() bool {
	return p.timer["paralyzed"] > 0
}

func (p *Player) IsHasted() bool {
	return p.timer["haste"] > 0
}

// -----------------------------------------------------------------------
func (p *Player) Pickup(item Object) bool {
	switch item.(type) {
	default:
		p.inventory = append(p.inventory, item)
		return true
	}
}

//func (p *Player) Pickup(item *Item) bool {
//	switch item.Type() {
//	case Gold:
//		p.Gold += item.GoldQty()
//		return true
//	default:
//		p.inventory = append(p.inventory, item)
//		return true
//	}
//}

// -----------------------------------------------------------------------
func (p *Player) RemoveItem(idx int) {
	p.inventory = append(p.inventory[:idx], p.inventory[idx+1:]...)
}

// -----------------------------------------------------------------------
func (p *Player) Equip(item *Item, msg *MessageLog) bool {
	//switch item.Type() {
	//case Weapon:
	//	if p.weapon != nil {
	//		msg.Add("You return %v to your pack.", p.weapon.GndString())
	//		p.weapon = nil
	//	}
	//	msg.Add("You are now wielding %v.", item.InvString())
	//	p.weapon = item
	//case Armor:
	//	if p.armor != nil {
	//		msg.Add("You take off %v.", p.armor.GndString())
	//		p.armor = nil
	//	}

	//	msg.Add("You are now wearing %v.", item.InvString())
	//	p.armor = item
	//	p.AC = item.val1
	//default:
	//	msg.Add("You cannot equip that item.")
	//	return false
	//}
	return true
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

	// Decrement and timers that are set
	for k := range p.timer {
		p.timer[k]--
		if p.timer[k] < 0 {
			delete(p.timer, k)
		}
	}

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

func (p *Player) Timer(name string) int {
	return p.timer[name]
}

func (p *Player) SetTimer(name string, val int) {
	if val == 0 {
		delete(p.timer, name)
	} else {
		p.timer[name] = val
	}
}

// -----------------------------------------------------------------------
func (p *Player) InfoString() string {
	info := fmt.Sprintf(
		"Level: %-2d  Gold: %-5d  Hp: %2d(%2d)  Str: %2d(%2d)  Arm: %-2d   Exp: %d/%d",
		p.depth,
		p.Gold,
		p.HP,
		p.maxHP,
		p.Str,
		p.maxStr,
		p.AC,
		p.Level,
		p.XP,
	)
	return info
}
