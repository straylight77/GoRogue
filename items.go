package main

import (
	"fmt"
	"math/rand"
)

// -----------------------------------------------------------------------
type Item interface {
	Rune() rune
	InvString() string
	Qty() int
}

// -----------------------------------------------------------------------
type ItemList map[Coord]Item

func (list *ItemList) Clear() {
	clear(*list)
}

// -----------------------------------------------------------------------
type Gold struct {
	amt int
}

func (g Gold) Rune() rune {
	return '*'
}

func (g Gold) Qty() int {
	return g.amt
}

func (g Gold) String() string {
	return fmt.Sprintf("%d pieces of gold", g.amt)
}

func (g Gold) InvString() string {
	return fmt.Sprintf("%d pieces of gold", g.amt)
}

func randGoldAmt(depth int) int {
	return rand.Intn(50+10*depth) + 2
}

// -----------------------------------------------------------------------
type Food struct {
	typ int
	qty int
}

func (f Food) Rune() rune {
	return '%'
}

func (f Food) Qty() int {
	return f.qty
}

func (f Food) String() string {
	return "a ration"
}

func (f Food) InvString() string {
	return fmt.Sprintf("a ration of food")
}
