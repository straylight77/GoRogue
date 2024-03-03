package main

type Player struct {
	X, Y int
}

// -----------------------------------------------------------------------
func (p *Player) SetPos(newX, newY int) {
	p.X = newX
	p.Y = newY
}

// -----------------------------------------------------------------------
func (p *Player) Pos() (int, int) {
	return p.X, p.Y
}

// -----------------------------------------------------------------------
func (p *Player) Rune() rune {
	return rune('@')
}
