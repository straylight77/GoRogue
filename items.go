package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
)

type ItemType int

const (
	Gold ItemType = iota
	Food
	Potion
	Scroll
	Ring
	Stick
	Weapon
	Armor
	Amulet
)

var itemRunes = map[ItemType]rune{
	Gold:   '*',
	Food:   '%',
	Potion: '!',
	Scroll: '?',
	Ring:   '=',
	Stick:  '/',
	Weapon: ')',
	Armor:  ']',
	Amulet: '&',
}

// -----------------------------------------------------------------------
type Item struct {
	typ     ItemType
	name    string
	val1    int
	val2    int
	val3    int
	val4    int
	ench    int
	magical bool
	cursed  bool
}

func (item Item) Rune() rune {
	ch, ok := itemRunes[item.typ]
	if !ok {
		ch = '0' // shouldn't see this but here's a default just in case
	}
	return ch
}

func (item Item) GndString() string {
	switch item.typ {
	case Gold:
		return fmt.Sprintf("%d pieces of gold", item.val1)
	default:
		return fmt.Sprintf("a %v", item.name)
	}
}

func (item Item) String() string {
	switch item.typ {
	case Gold:
		return fmt.Sprintf("%d pieces of gold", item.val1)
	case Weapon:
		minDmg := item.val1 + item.ench
		maxDmg := item.val2 + item.ench
		cursed := ""
		if item.IsCursed() {
			cursed = " {cursed}"
		}
		return fmt.Sprintf("%+d %s [%d-%d]%s", item.ench, item.name, minDmg, maxDmg, cursed)
	case Armor:
		return fmt.Sprintf("+%d %s [%d]", item.ench, item.name, item.val1+item.ench)
	default:
		return item.name
	}
}

func (item Item) GoldQty() int {
	return item.val1
}

func (item Item) Nutrition() int {
	return item.val1
}

func (item Item) MeleeDamage() int {
	sum := 0
	for i := 0; i < item.val1; i++ {
		sum += rand.Intn(item.val2)
	}
	return sum
}

func (item Item) IsMagical() bool {
	return item.magical
}

func (item Item) IsCursed() bool {
	return item.cursed
}

// -----------------------------------------------------------------------
func newGold(qty int) *Item {
	return &Item{typ: Gold, val1: qty}
}

func randGoldAmt(depth int) int {
	return rand.Intn(50+10*depth) + 2
}

// -----------------------------------------------------------------------
func newRation() *Item {
	return &Item{
		typ:  Food,
		name: "ration",
		val1: NutritionTime,
	}
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

// -----------------------------------------------------------------------
type WeaponTemplate struct {
	melee  string
	thrown string
	worth  int
}

var WeaponLib = map[string]WeaponTemplate{
	"mace":             {"2d4", "1d3", 9},
	"long sword":       {"1d10", "1d2", 15},
	"dagger":           {"1d6", "1d4", 2},
	"two-handed sword": {"3d6", "1d2", 30},
	"spear":            {"1d8", "1d6", 2},
}

func newWeapon(name string) *Item {
	t, ok := WeaponLib[name]
	if !ok {
		panic(name)
	}
	v1, v2 := parseDiceStr(t.melee)

	return &Item{
		typ:  Weapon,
		name: name,
		val1: v1,
		val2: v2,
	}
}

func randWeapon() *Item {
	// Pick a weapon from the list at random
	i := rand.Intn(len(WeaponLib))
	var item *Item
	for name := range WeaponLib {
		if i == 0 {
			item = newWeapon(name)
		}
		i--
	}

	// 10% chance of a cursed weapon with -1 to -3 penalty, and a 5% chance
	// of an enchanted weapon with a +1 to +3 bonus.
	if rand.Intn(100) < 5 { // enchanted
		item.magical = true
		item.ench = rand.Intn(2) + 1
	} else if rand.Intn(100) < 10 { // cursed
		item.cursed = true
		item.ench = -1 * (rand.Intn(2) + 1)
	}

	return item
}

// -----------------------------------------------------------------------
func newArmor() *Item {
	return &Item{
		typ:  Armor,
		name: "leather armor",
		val1: 8,
	}
}

// -----------------------------------------------------------------------
type ItemList map[Coord]*Item

func (list *ItemList) Clear() {
	clear(*list)
}
