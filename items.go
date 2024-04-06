package main

import (
	"fmt"
	"math/rand"
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

func (item Item) String() string {
	switch item.typ {
	case Gold:
		return fmt.Sprintf("%d pieces of gold", item.val1)
	default:
		return fmt.Sprintf("a %v", item.name)
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

func (item Item) isMagical() bool {
	return item.magical
}

func (item Item) isCursed() bool {
	return item.cursed
}

// -----------------------------------------------------------------------
func newGold(qty int) *Item {
	return &Item{typ: Gold, val1: qty}
}

func newRation() *Item {
	return &Item{
		typ:  Food,
		name: "ration",
		val1: NutritionTime,
	}
}

func newWeapon() *Item {
	return &Item{
		typ:  Weapon,
		name: "longsword",
		val1: 1,
		val2: 8,
	}
}

func randGoldAmt(depth int) int {
	return rand.Intn(50+10*depth) + 2
}

// -----------------------------------------------------------------------
type ItemList map[Coord]*Item

func (list *ItemList) Clear() {
	clear(*list)
}
