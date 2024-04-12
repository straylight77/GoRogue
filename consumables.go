package main

import (
	"fmt"
	"math/rand"
)

type Consumable struct {
	typ    ItemType
	name   string
	qty    int
	effect int
	val1   int
}

// === FOOD ==============================================================
func newRation() *Item {
	return &Item{
		typ:  Food,
		name: "ration",
		val1: NutritionTime,
	}
}

func (item Item) Nutrition() int {
	return item.val1
}

// === POTIONS ==========================================================
// name: full name of the potion once it's been identified
// val1: index of PotionLib for this potion

type PotionTemplate struct {
	pct        int // probability of this potion being randomly generated
	chance     int // cumulative probability
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

func newPotion(name string) *Item {
	ok := false
	var templ PotionTemplate
	var idx int
	for i, t := range PotionLib {
		if t.name == name {
			ok = true
			templ = t
			idx = i
			break
		}
	}
	if !ok {
		panic("No potion with the name " + name)
	}

	return &Item{
		typ:  Potion,
		name: fmt.Sprintf("potion of %s", templ.name),
		val1: idx,
	}
}

func randPotion() *Item {
	roll := rand.Intn(100) + 1 //1-100
	name := ""
	for _, t := range PotionLib {
		//debug.Add("rand potion: (%d) chance=%d", roll, t.chance)
		if roll <= t.chance {
			name = t.name
			break
		}
	}
	return newPotion(name)
}

// === SCROLLS ===========================================================
func newScroll() *Item {
	return &Item{
		typ:  Scroll,
		name: "ryfay in the airchay",
	}
}

func randScroll() *Item {
	return newScroll()
}

// === STICKS ============================================================
func newStick() *Item {
	return &Item{
		typ:  Stick,
		name: "bamboo",
	}
}

func randStick() *Item {
	return newStick()
}
