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

// -----------------------------------------------------------------------
type Player struct {
	X, Y      int
	Symbol    rune
	moves     int
	depth     int
	HP        int
	maxHP     int
	Str       int
	maxStr    int
	Level     int
	XP        int
	AC        int
	Melee     Dice
	Gold      int
	healCount int
	foodCount int
	inventory []Item
	equiped   map[string]Equipable
	timer     map[string]int
}

// -----------------------------------------------------------------------
func (p *Player) Init() {
	p.Str = 16
	p.maxStr = 16
	p.HP = 12
	p.maxHP = 12
	p.AC = 10
	p.Level = 1
	p.foodCount = NutritionTime
	p.timer = make(map[string]int)
	p.equiped = map[string]Equipable{
		"weapon": nil,
		"armor":  nil,
		"left":   nil,
		"right":  nil,
	}
	p.ResetHealCount()
}

// -----------------------------------------------------------------------
// implement the Actor interface

func (p *Player) Pos() Coord {
	return Coord{p.X, p.Y}
}

func (p *Player) SetPos(newPos Coord) {
	p.X = newPos.X
	p.Y = newPos.Y
}

func (p *Player) Rune() rune {
	return p.Symbol
}

func (p *Player) AdjustHP(amt int) {
	p.HP += amt
	if p.HP > p.maxHP {
		p.HP = p.maxHP
	}
}

func (p *Player) Attack(m Actor) string {

	var label string
	if p.IsBlind() {
		label = "something"
	} else {
		label = fmt.Sprintf("the %v", m)
	}

	if attackHits(p.ToHit(), m.ArmorClass()) {
		dmg := p.RollDamage()
		m.AdjustHP(-dmg)
		p.healCount++ // this shouldn't decrement when fighting
		return fmt.Sprintf("You hit %v for %d damage.", label, dmg)
	}
	return fmt.Sprintf("You miss %v.", label)

}

func (p *Player) ArmorClass() int {
	return p.AC
}

func (p *Player) IsConfused() bool {
	return p.timer["confused"] > 0
}

func (p *Player) IsBlind() bool {
	return p.timer["blind"] > 0
}

// -----------------------------------------------------------------------

func (p *Player) IsParalyzed() bool {
	return p.timer["paralyzed"] > 0
}

func (p *Player) IsHasted() bool {
	return p.timer["haste"] > 0
}

// -----------------------------------------------------------------------

func (p *Player) StrAttackBonus() int {
	switch {
	case p.Str <= 6:
		return p.Str - 7
	case p.Str <= 16:
		return 0
	case p.Str <= 19:
		return 1
	case p.Str <= 20: // 18/[51-75]
		return 2
	case p.Str >= 22: // 18/[91-100]
		return 3
	}
	return 0
}

func (p *Player) StrDamageBonus() int {
	switch {
	case p.Str <= 6:
		return p.Str - 7
	case p.Str <= 15:
		return 0
	case p.Str <= 17:
		return 1
	case p.Str == 18:
		return 2
	case p.Str == 19: // 18/[1-50]
		return 3
	case p.Str == 20: // 18/[51-75]
		return 4
	case p.Str == 21: // 18/[76-90]
		return 5
	case p.Str >= 22: // 18/[91-100]
		return 6
	default:
		return 0
	}
}

func (p *Player) ToHit() int {
	return 21 - p.Level - p.StrAttackBonus()
}

func (p *Player) RollDamage() int {
	return p.DamageDice().Roll()
}

func (p *Player) DamageDice() Dice {
	return p.Melee.Add(p.StrDamageBonus())
}

// -----------------------------------------------------------------------

func (p *Player) Pickup(item Item) bool {
	switch item.(type) {
	case *Gold:
		p.Gold += item.(*Gold).qty
		return true
	default:
		p.inventory = append(p.inventory, item)
		return true
	}
}

func (p *Player) RemoveItem(idx int) {
	p.inventory = append(p.inventory[:idx], p.inventory[idx+1:]...)
}

// -----------------------------------------------------------------------

func (p *Player) AddXP(amt int) {
	p.XP += amt
}

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
		// Level Up!
		hp := rand.Intn(12) + 1
		p.HP += hp
		p.maxHP += hp
		msg = fmt.Sprintf("Welcome to level %d! [%+d HP]", level, hp)
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

func (p *Player) AdjustFoodCount(amt int) {
	p.foodCount += amt
	if p.foodCount > NutritionTime {
		p.foodCount = NutritionTime
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

// -----------------------------------------------------------------------

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
	condition := ""
	switch {
	case p.IsParalyzed():
		condition = "Paralyzed"
	case p.foodCount <= HungerLimit:
		condition = "Hungry"
	case p.IsConfused():
		condition = "Confused"
	case p.IsBlind():
		condition = "Blind"
	case p.IsHasted():
		condition = "Haste"
	}

	return fmt.Sprintf(
		"Depth:%-2d  Gold:%-5d  Hp:%2d(%2d)  Str:%-2d  Hit:%-2d  Arm:%-2d  Lvl:%d/%-6d %s",
		p.depth,
		p.Gold,
		p.HP,
		p.maxHP,
		p.Str,
		p.ToHit(),
		p.ArmorClass(),
		p.Level,
		p.XP,
		condition,
	)
}

// -----------------------------------------------------------------------
func (p *Player) StatsStrings() []string {

	//savePoison := (7 + p.Level/2) * 5
	//saveMagic := (4 + p.Level/2) * 5

	savePoison := (7 + p.Level/2)
	saveMagic := (4 + p.Level/2)
	dice := p.DamageDice()

	return []string{

		fmt.Sprintf("Level:  %d", p.Level),
		"",
		fmt.Sprintf("Strength:   %d / %d", p.Str, p.maxStr),
		//fmt.Sprintf("(%+d hit, %d dmg)", p.StrAttackBonus(), p.StrDamageBonus()),
		fmt.Sprintf(" %+d hit", p.StrAttackBonus()),
		fmt.Sprintf(" %+d dmg", p.StrDamageBonus()),
		"",
		fmt.Sprintf("Hit Points: %d / %d", p.HP, p.maxHP),
		"",
		fmt.Sprintf("THAC0:  %d    (%+d)", p.ToHit(), p.StrAttackBonus()),
		fmt.Sprintf("Damage: %d-%-2d  (%+d)", dice.Min(), dice.Max(), p.StrDamageBonus()),
		fmt.Sprintf("Armor:  %d", p.AC),
		"",
		fmt.Sprintf("Poison: %d", savePoison),
		fmt.Sprintf("Magic:  %d", saveMagic),
		"",
		fmt.Sprintf("XP:     %d", p.XP),
		fmt.Sprintf("Next:   %d", XPTable[p.Level]),
	}
}

func (p *Player) Score() int {
	sum := p.Gold
	for _, item := range p.inventory {
		sum += item.Worth()
	}
	return sum
}
