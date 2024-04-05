package main

import (
	"fmt"
	"math/rand"
)

// -----------------------------------------------------------------------
type Item interface {
	Rune() rune
	Pickup(*Player)
}

// -----------------------------------------------------------------------
type ItemList map[Coord]Item

// -----------------------------------------------------------------------
type Gold struct {
	amt int
}

func (g Gold) Rune() rune {
	return '*'
}

func (g Gold) String() string {
	return fmt.Sprintf("%d pieces of gold", g.amt)
}

func (g Gold) Pickup(p *Player) {
	p.Gold += g.amt
}

func randGoldAmt(depth int) int {
	return rand.Intn(50+10*depth) + 2
}
