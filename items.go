package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
)

type ItemType int

const (
	Gold ItemType = iota
	Food
	Weapon
	Armor
	Ring
	Potion
	Scroll
	Stick
	Amulet
)

var ItemRunes = map[ItemType]rune{
	Gold:   '*',
	Food:   '%',
	Weapon: ')',
	Armor:  ']',
	Ring:   '=',
	Potion: '!',
	Scroll: '?',
	Stick:  '/',
	Amulet: '&',
}

// -----------------------------------------------------------------------
type Item struct {
	typ     ItemType
	name    string
	val1    int
	val2    int
	val3    int
	val4    int
	ench    int
	magical bool
	cursed  bool
}

func (item Item) Rune() rune {
	ch, ok := ItemRunes[item.typ]
	if !ok {
		ch = '0' // shouldn't see this but here's a default just in case
	}
	return ch
}

func (item Item) GndString() string {
	switch item.typ {
	case Gold:
		if item.val1 == 1 {
			return fmt.Sprintf("%d piece of gold", item.val1)
		} else {
			return fmt.Sprintf("%d pieces of gold", item.val1)
		}
	case Potion:
		return fmt.Sprintf("a %s potion", item.name)
	case Scroll:
		return fmt.Sprintf("a scroll titled '%s'", item.name)
	case Ring:
		return fmt.Sprintf("a %s ring", item.name)
	case Stick:
		return fmt.Sprintf("a %s wand", item.name)
	default:
		return fmt.Sprintf("a %s", item.name)
	}
}

func (item Item) String() string {
	cursed := ""
	if item.IsCursed() {
		cursed = " {cursed}"
	}

	switch item.typ {
	case Gold:
		if item.val1 == 1 {
			return fmt.Sprintf("%d piece of gold", item.val1)
		} else {
			return fmt.Sprintf("%d pieces of gold", item.val1)
		}
	case Weapon:
		minDmg := item.val1 + item.ench
		maxDmg := item.val2 + item.ench
		return fmt.Sprintf("a %+d %s [%d-%d]%s", item.ench, item.name, minDmg, maxDmg, cursed)
	case Armor:
		prot := item.val1 - item.ench
		return fmt.Sprintf("a %+d %s [%d]%s", item.ench, item.name, prot, cursed)
	case Ring:
		return fmt.Sprintf("a %s ring", item.name)
	case Potion:
		return fmt.Sprintf("a %s potion", item.name)
	case Scroll:
		return fmt.Sprintf("a scroll titled '%s'", item.name)
	case Stick:
		return fmt.Sprintf("a %s wand", item.name)
	default:
		return fmt.Sprintf("a %s", item.name)
	}
}

func (item Item) IsMagical() bool {
	return item.magical
}

func (item Item) IsCursed() bool {
	return item.cursed
}

// -----------------------------------------------------------------------
func newGold(qty int) *Item {
	return &Item{typ: Gold, val1: qty}
}

func (item Item) GoldQty() int {
	return item.val1
}

func randGoldAmt(depth int) int {
	return rand.Intn(50+10*depth) + 2
}

// -----------------------------------------------------------------------
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

// -----------------------------------------------------------------------
type WeaponTemplate struct {
	melee  string
	thrown string
	worth  int
}

var WeaponLib = map[string]WeaponTemplate{
	"mace":             {"2d4", "1d3", 9},
	"long sword":       {"1d10", "1d2", 15},
	"dagger":           {"1d6", "1d4", 2},
	"two-handed sword": {"3d6", "1d2", 30},
	"spear":            {"1d8", "1d6", 2},
}

func newWeapon(name string) *Item {
	t, ok := WeaponLib[name]
	if !ok {
		panic(name)
	}
	v1, v2 := parseDiceStr(t.melee)

	return &Item{
		typ:  Weapon,
		name: name,
		val1: v1,
		val2: v2,
	}
}

func randWeapon() *Item {
	// Pick a weapon from the list at random
	i := rand.Intn(len(WeaponLib))
	var item *Item
	for name := range WeaponLib {
		if i == 0 {
			item = newWeapon(name)
		}
		i--
	}
	randEnchant(item, 5, 10)
	return item
}

func randEnchant(item *Item, enchantProb int, cursedProb int) {
	// 10% chance of a cursed weapon with -1 to -3 penalty, and a 5% chance
	// of an enchanted weapon with a +1 to +3 bonus.
	if rand.Intn(100) < enchantProb { // enchanted
		item.magical = true
		item.ench = rand.Intn(2) + 1
	} else if rand.Intn(100) < cursedProb { // cursed
		item.cursed = true
		item.ench = -1 * (rand.Intn(2) + 1)
	}
}

func (item Item) MeleeDamage() int {
	sum := 0
	for i := 0; i < item.val1; i++ {
		sum += rand.Intn(item.val2)
	}
	return sum
}

// -----------------------------------------------------------------------
type ArmorTemplate struct {
	AC    int
	worth int
}

var ArmorLib = map[string]ArmorTemplate{
	"leather armor": {8, 0},
	"ring mail":     {7, 0},
	"scale mail":    {6, 3},
	"chain mail":    {5, 75},
	"splint mail":   {4, 80},
	"banded mail":   {3, 90},
	"plate armor":   {2, 440},
}

func newArmor(name string) *Item {
	t, ok := ArmorLib[name]
	if !ok {
		panic(name)
	}

	return &Item{
		typ:  Armor,
		name: name,
		val1: t.AC,
	}
}

func randArmor() *Item {
	// Pick an armor from the list at random
	i := rand.Intn(len(ArmorLib))
	var item *Item
	for name := range ArmorLib {
		if i == 0 {
			item = newArmor(name)
		}
		i--
	}
	randEnchant(item, 8, 20)
	return item
}

// -----------------------------------------------------------------------
func newPotion() *Item {
	return &Item{
		typ:  Potion,
		name: "red",
	}
}

func randPotion() *Item {
	return newPotion()
}

// -----------------------------------------------------------------------
func newScroll() *Item {
	return &Item{
		typ:  Scroll,
		name: "ryfay in the airchay",
	}
}

func randScroll() *Item {
	return newScroll()
}

// -----------------------------------------------------------------------
func newStick() *Item {
	return &Item{
		typ:  Stick,
		name: "bamboo",
	}
}

func randStick() *Item {
	return newStick()
}

// -----------------------------------------------------------------------
func newRing() *Item {
	return &Item{
		typ:  Ring,
		name: "ruby",
	}
}

func randRing() *Item {
	return newRing()
}

// -----------------------------------------------------------------------
type ItemList map[Coord]*Item

func (list *ItemList) Clear() {
	clear(*list)
}

// -----------------------------------------------------------------------
// ITEM   PCT  CUMUL
// Potion  27     27
// Scroll  27     54
// Food    18     72
// Weapon   9     81
// Armor    9     90
// Ring     5     95
// Stick    5    100
func randItemType() ItemType {
	roll := rand.Intn(100)
	//debug.Add("rand item: roll=%d", roll)
	switch {
	case roll < 27:
		return Potion
	case roll < 54:
		return Scroll
	case roll < 72:
		return Food
	case roll < 81:
		return Weapon
	case roll < 90:
		return Armor
	case roll < 95:
		return Ring
	case roll < 100:
		return Stick
	}
	return Food
}

func parseDiceStr(dice string) (int, int) {
	parts := strings.Split(dice, "d")
	v1, err := strconv.Atoi(parts[0])
	if err != nil {
		panic(err)
	}
	v2, err := strconv.Atoi(parts[1])
	if err != nil {
		panic(err)
	}
	return v1, v2
}
