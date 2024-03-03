package main

/***
 * Entity
 * BaseEntity
 * - contains common logic for all game entities e.g. Player, Monster
 * - Entity interface allows use of generic functions e.g. draw, combat
 */

type Entity interface {
	Pos() (int, int)
	SetPos(int, int)
	Rune() rune
}

type BaseEntity struct {
	X, Y   int
	Symbol rune
}

func (e *BaseEntity) SetPos(newX, newY int) {
	e.X = newX
	e.Y = newY
}

func (e *BaseEntity) Pos() (int, int) {
	return e.X, e.Y
}

func (e *BaseEntity) Rune() rune {
	return e.Symbol
}
