package main

import (
	"fmt"
	"math/rand"
)

// === FOOD ==============================================================

type Food struct {
	name string
	amt  int
}

func newFood(name string) *Food {
	return &Food{name, NutritionTime}
}

func (f *Food) Rune() rune {
	return '%'
}

func (f *Food) InvString() string {
	return f.GndString()
}

func (f *Food) GndString() string {
	return fmt.Sprintf("a %s", f.name)
}

func (f *Food) Worth() int {
	return 2
}

func (f *Food) Consume(gs *GameState) bool {
	gs.messages.Add("You eat %v.", f.InvString())
	gs.player.AdjustFoodCount(f.amt)
	return true
}

// === POTIONS ==========================================================

type Potion struct {
	id int
}

func newPotion(name string) *Potion {
	ok := false
	var idx int
	for i, t := range PotionLib {
		if t.name == name {
			ok = true
			idx = i
			break
		}
	}
	if !ok {
		panic("No potion with the name " + name)
	}

	return &Potion{
		id: idx,
	}
}

func randPotion() *Potion {
	roll := rand.Intn(100) + 1 //1-100
	name := ""
	for _, t := range PotionLib {
		//debug.Add("rand potion: (%d) chance=%d", roll, t.chance)
		if roll <= t.cumPct {
			name = t.name
			break
		}
	}
	return newPotion(name)
}

func (p *Potion) Rune() rune {
	return '!'
}

func (p *Potion) InvString() string {
	return p.GndString()
}

func (p *Potion) GndString() string {
	templ := PotionLib[p.id]
	if templ.discovered {
		return fmt.Sprintf("a potion of %s", templ.name)
	} else {
		color := PotionColors[templ.color]
		return fmt.Sprintf("a %s potion", color)
	}
}

func (p *Potion) Worth() int {
	return 2
}

func (p Potion) String() string {
	return p.GndString()
}

func (p *Potion) Consume(gs *GameState) bool {
	templ := PotionLib[p.id]
	doEffect(templ.effect, gs)
	gs.messages.Add(templ.message)
	p.Identify()
	return true
}

func (p *Potion) IsIdentified() bool {
	return PotionLib[p.id].discovered
}
func (p *Potion) Identify() {
	PotionLib[p.id].discovered = true
}

type PotionTemplate struct {
	pct        int // probability of this potion being randomly generated
	cumPct     int // cumulative probability
	name       string
	effect     int
	color      int
	discovered bool
	message    string
}

var PotionLib = []PotionTemplate{
	{15, 15, "healing", E_Healing, 0, false, "You begin to feel better."},
	{15, 30, "strength", E_Strength, 0, false, "You feel stronger, what bulging muscles!"},
	{14, 44, "restore strength", E_Restore, 0, false, "Hey, this tastes great, it make you feel warm all over."},
	{10, 54, "paralysis", E_Paralyze, 0, false, "You feel your body seizing up, you can't move!"},
	{8, 62, "confusion", E_Confusion, 0, false, "Wait, what's going on here. Huh? What? Who?"},
	{8, 70, "poison", E_Poison, 0, false, "You feel very sick now."},
	{6, 76, "monster detection", E_DetMonsters, 0, false, "You feel like you are not alone."},
	{6, 82, "detect magic", E_DetMagic, 0, false, "You sense the presence of magic."},
	{5, 87, "extra healing", E_ExtraHealing, 0, false, "You begin to feel much better."},
	{4, 91, "haste", E_Haste, 0, false, "Tastes like coffee, everything seems to slow down."},
	{4, 95, "blindness", E_Blindness, 0, false, "A cloak of darkness falls around you."},
	{2, 97, "raise level", E_LevelUp, 0, false, "You feel more experienced."},
	{2, 99, "truesight", E_Truesight, 0, false, "Tastes like slime-mold juice."},
	{1, 100, "thirst quenching", E_Nothing, 0, false, "Meh, tastes pretty dull."},
}

var PotionColors = []string{
	"black",
	"blue",
	"brown",
	"clear",
	"crimson",
	"cyan",
	"gold",
	"green",
	"grey",
	"magenta",
	"pink",
	"plaid",
	"purple",
	"red",
	"silver",
	"tan",
	"tangerine",
	"topaz",
	"turquoise",
	"vermilion",
	"violet",
	"white",
	"yellow",
}

func assignPotionColors() {
	if len(PotionColors) < len(PotionLib) {
		panic("Not enough potion colors to assign")
	}
	used := make(map[int]bool)
	for pid := range PotionLib {
		cid := rand.Intn(len(PotionColors))
		for used[cid] {
			cid = rand.Intn(len(PotionColors))
		}
		used[cid] = true
		PotionLib[pid].color = cid
		//debug.Add("assign %s -> %s", PotionLib[pid].name, PotionColors[cid])
	}
}

// === SCROLLS ===========================================================

// === STICKS ============================================================
