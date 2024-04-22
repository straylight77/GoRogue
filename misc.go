package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
)

// -----------------------------------------------------------------------
type Actor interface {
	Pos() Coord
	SetPos(Coord)
	Rune() rune
	AdjustHP(amt int)
	Attack(Actor, *MessageLog)
	ArmorClass() int
	IsConfused() bool
	IsBlind() bool
}

// --- COMBAT ------------------------------------------------------------
func attackHits(toHit int, targetAC int) bool {
	roll := rand.Intn(20) + 1
	target := toHit - targetAC
	isHit := roll >= target
	//debug.Add("hit? roll=%d target=%d (%d-%d)  -> %v", roll, target, toHit, targetAC, isHit)
	return isHit
}

// -----------------------------------------------------------------------
type Dice struct {
	num, size, bonus int
}

// Dice rolls given in the format: "1d8/1d6/1d6"
func parseDice(fullStr string) []Dice {

	attacks := strings.Split(fullStr, "/")
	dice := []Dice{}
	for _, str := range attacks {

		parts := strings.Split(str, "d")
		num, err := strconv.Atoi(parts[0])
		if err != nil {
			panic(err)
		}
		size, err := strconv.Atoi(parts[1])
		if err != nil {
			panic(err)
		}

		//TODO handle bonus part
		dice = append(dice, Dice{num, size, 0})
		//debug.Add("parse:   %v", dice)
	}

	return dice
}

func (d Dice) String() string {
	if d.bonus == 0 {
		return fmt.Sprintf("%dd%d", d.num, d.size)
	} else {
		return fmt.Sprintf("%dd%d%+d", d.num, d.size, d.bonus)
	}
}

func (d Dice) Min() int {
	return d.num + d.bonus
}

func (d Dice) Max() int {
	return (d.num * d.size) + d.bonus
}

func (d Dice) Add(amt int) Dice {
	return Dice{d.num, d.size, d.bonus + amt}
}

func (d Dice) Roll() int {
	sum := d.bonus
	rolls := make([]int, d.num)
	for i := 0; i < d.num; i++ {
		roll := rand.Intn(d.size) + 1
		rolls[i] = roll
		sum += roll
	}
	debug.Add("Roll: %v, rolls=%v, sum=%d", d, rolls, sum)
	return sum
}

// -----------------------------------------------------------------------
type Direction int

const (
	North Direction = iota
	East
	South
	West
)

func (d Direction) String() string {
	switch d {
	case North:
		return "north"
	case East:
		return "east"
	case South:
		return "south"
	case West:
		return "west"
	default:
		return "unknown"
	}
}

func getDirectionCoords(dir Direction) Coord {
	dx, dy := 0, 0
	switch dir {
	case North:
		dy = -1
	case South:
		dy = 1
	case East:
		dx = 1
	case West:
		dx = -1
	}
	return Coord{dx, dy}
}

// -----------------------------------------------------------------------

func abs(val int) int {
	if val < 0 {
		val = -val
	}
	return val
}

func max(val1, val2 int) int {
	if val1 > val2 {
		return val1
	} else {
		return val2
	}
}
