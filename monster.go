package main

/*************************************************************************
 * MonsterLib
 *
 */

/*************************************************************************
 * MonsterList
 *
 */

type MonsterList []*Monster

func (ml *MonsterList) Add(m *Monster, x, y int) {
	m.X, m.Y = x, y
	*ml = append(*ml, m)
}

func (ml *MonsterList) Remove(idx int) {
	*ml = append((*ml)[:idx], (*ml)[idx+1:]...)
}

func (ml *MonsterList) Clear() {
	*ml = nil
}

func (ml MonsterList) MonsterAt(x, y int) *Monster {
	for _, m := range ml {
		if m.X == x && m.Y == y {
			return m
		}
	}
	return nil
}

/*************************************************************************
 * Monster
 * implements Entity interface
 */

type Monster struct {
	//X, Y   int
	//Symbol rune
	BaseEntity
	Name string
	HP   int
}

func NewMonster(n string, sym rune, hp int) *Monster {
	newMonster := &Monster{
		//Symbol: sym,
		Name: n,
		HP:   hp,
	}
	newMonster.Symbol = sym
	return newMonster
}
