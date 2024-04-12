package main

import "math/rand"

type Equipment struct {
	typ        ItemType
	name       string
	identified bool
	cursed     bool
	bonus      int
	meleeRoll  int
	meleeSize  int
}

// === WEAPONS ===========================================================
// name: the sub-type of weapons e.g. dagger, long sword
// val1: the number of dice to roll for melee damage
// val2: the die size to roll for melee damage

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
		panic("No weapon with the name " + name)
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

// === ARMOR =============================================================
// name: the sub-type of armor e.g. chain mail
// val1: the base armor class before enchantments

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
		panic("No armor with the name " + name)
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

// === RINGS =============================================================
func newRing() *Item {
	return &Item{
		typ:  Ring,
		name: "ruby",
	}
}

func randRing() *Item {
	return newRing()
}
