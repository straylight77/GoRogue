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
	typ  ItemType
	qty  int
	val1 int
	val2 int
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
		return fmt.Sprintf("%d pieces of gold", item.qty)
	case Food:
		if item.qty == 1 {
			return "a ration of food"
		} else {
			return fmt.Sprintf("%d rations of food", item.qty)
		}
	default:
		return "mysterious artifact" // default value just in case
	}
}

// -----------------------------------------------------------------------
func newGold(qty int) *Item {
	return &Item{typ: Gold, qty: qty}
}

func newRation() *Item {
	return &Item{typ: Food, qty: 1}
}

func randGoldAmt(depth int) int {
	return rand.Intn(50+10*depth) + 2
}

// -----------------------------------------------------------------------
type ItemList map[Coord]*Item

func (list *ItemList) Clear() {
	clear(*list)
}
