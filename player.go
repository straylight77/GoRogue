package main

import "fmt"

type Player struct {
	BaseEntity
	moves int
	depth int
}

// -----------------------------------------------------------------------
func (p *Player) InfoString() string {
	return fmt.Sprintf("Level: 1  Gold: 4       Hp: 11 (20)  Str: 16(16)  Arm: 4   Exp: 2/14")
}
