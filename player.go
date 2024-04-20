package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
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
type Dice struct {
	num, size, bonus int
}

func (d Dice) String() string {
	return fmt.Sprintf("%dd%d%+d", d.num, d.size, d.bonus)
}

func (d Dice) Min() int {
	return d.num + d.bonus
}

func (d Dice) Max() int {
	return (d.num * d.size) + d.bonus
}

func (d Dice) Add(amt int) Dice {
	return Dice{d.num, d.size, d.bonus + amt}
}

func (d Dice) Roll() int {
	sum := d.bonus
	rolls := make([]int, d.num)
	for i := 0; i < d.num; i++ {
		roll := rand.Intn(d.size) + 1
		rolls[i] = roll
		sum += roll
	}
	debug.Add("Roll: %v, rolls=%v, sum=%d", d, rolls, sum)
	return sum
}

// TODO combine Roll() and rollDice()?
func rollDice(num, size, bonus int) int {
	sum := bonus
	rolls := make([]int, num)
	for i := 0; i < num; i++ {
		roll := rand.Intn(size) + 1
		rolls[i] = roll
		sum += roll
	}
	debug.Add("rollDice: %dd%d+%d, rolls=%v, sum=%d", num, size, bonus, rolls, sum)
	return sum
}

func parseDiceStr(dice string) (int, int) {
	parts := strings.Split(dice, "d")
	v1, err := strconv.Atoi(parts[0])
	if err != nil {
		panic(err)
	}
	v2, err := strconv.Atoi(parts[1])
	if err != nil {
		panic(err)
	}
	return v1, v2
}

func attackHits(toHit int, targetAC int) bool {
	roll := rand.Intn(20) + 1

	check := roll + toHit
	isHit := check >= targetAC
	debug.Add("hit? roll=%d (%d%+d) AC=%d  -> %v", roll, toHit, check, targetAC, isHit)

	return isHit
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

// implement the Actor interface

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
func (p *Player) Attack(m Actor) string {

	var label string

	if p.IsBlind() {
		label = "something"
	} else {
		label = fmt.Sprintf("%v", m.Label())
	}

	if attackHits(p.ToHit(), m.ArmorClass()) {
		dmg := p.RollDamage()
		m.AdjustHP(-dmg)
		p.healCount++ // this shouldn't decrement when fighting
		return fmt.Sprintf("You hit %v for %d damage.", label, dmg)
	}
	return fmt.Sprintf("You miss the %v.", label)

}

func (p *Player) ArmorClass() int {
	return p.AC
}

func (p *Player) RollDamage() int {
	return p.DamageDice().Roll()
}

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
	return p.Level + p.StrAttackBonus()
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
	return fmt.Sprintf(
		"Depth:%-2d  Gold:%-5d  Hp:%2d(%2d)  Str:%-2d  Hit:%+2d  Arm:%-2d  Lvl:%d/%d",
		p.depth,
		p.Gold,
		p.HP,
		p.maxHP,
		p.Str,
		p.ToHit(),
		p.ArmorClass(),
		p.Level,
		p.XP,
	)
}

// -----------------------------------------------------------------------
func (p *Player) StatsStrings() []string {

	savePoison := (7 + p.Level/2) * 5
	saveMagic := (4 + p.Level/2) * 5
	dice := p.DamageDice()

	return []string{
		fmt.Sprintf("Strength: %d (%+d,%+d)",
			p.Str,
			p.StrAttackBonus(),
			p.StrDamageBonus(),
		),
		"",
		fmt.Sprintf("Attack:  %+d", p.ToHit()),
		fmt.Sprintf("Damage:  %d-%d", dice.Min(), dice.Max()),
		fmt.Sprintf("Armor:   %d", p.AC),
		"",
		fmt.Sprintf("Poison:  %d", savePoison),
		fmt.Sprintf("Magic:   %d", saveMagic),
		"",
		fmt.Sprintf("Level:   %d", p.Level),
		fmt.Sprintf("XP:      %d", p.XP),
		fmt.Sprintf("Next:    %d", XPTable[p.Level]),
	}
}
