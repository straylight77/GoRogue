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
)

// -----------------------------------------------------------------------
type Item struct {
	typ  ItemType
	qty  int
	val1 int
	val2 int
}

func (i Item) Rune() rune {
	switch i.typ {
	case 0:
		return '*'
	case 1:
		return '%'
	default:
		return '?'
	}
}

func (i Item) String() string {
	switch i.typ {
	case 0:
		return fmt.Sprintf("%d pieces of gold", i.qty)
	case 1:
		if i.qty == 1 {
			return "a ration of food"
		} else {
			return fmt.Sprintf("%d rations of food", i.qty)
		}
	default:
		return "mysterious artifact"
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
