package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
)

// === WEAPONS ===========================================================

type Weapon struct {
	name    string
	dmgDice int
	dmgSize int
	bonus   int
}

func newWeapon(name string) *Weapon {
	t, ok := WeaponLib[name]
	if !ok {
		panic("No weapon with the name " + name)
	}
	v1, v2 := parseDiceStr(t.melee)

	return &Weapon{
		name:    name,
		dmgDice: v1,
		dmgSize: v2,
	}
}

func randWeapon() *Weapon {
	// Pick a weapon from the list at random
	i := rand.Intn(len(WeaponLib))
	var w *Weapon
	for name := range WeaponLib {
		if i == 0 {
			w = newWeapon(name)
		}
		i--
	}
	//randEnchant(item, 5, 10)
	return w
}

func (w *Weapon) Equip(p *Player, msg *MessageLog) bool {
	if p.equiped["weapon"] != nil {
		msg.Add("You are already wielding %v.", p.equiped["weapon"].GndString())
		return false
	} else {
		p.equiped["weapon"] = w
		//TODO set the player's damage and to hit stats
		msg.Add("You are now wielding the %v.", w)
		return true
	}
}

func (w *Weapon) Unequip(p *Player, msg *MessageLog) bool {
	p.equiped["weapon"] = nil
	msg.Add("You put the %v back into your pack.", w)
	return true
}

func (w *Weapon) Rune() rune {
	return ')'
}

func (w *Weapon) GndString() string {
	return fmt.Sprintf("a %s", w)
}

func (w *Weapon) InvString() string {
	return fmt.Sprintf("%+d %v [%dd%d]", w.bonus, w, w.dmgDice, w.dmgSize)
}

func (w Weapon) String() string {
	return w.name
}

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

// === ARMOR =============================================================

type Armor struct {
	Name  string
	AC    int
	bonus int
}

func newArmor(name string) *Armor {
	t, ok := ArmorLib[name]
	if !ok {
		panic("No armor with the name " + name)
	}

	return &Armor{
		Name: name,
		AC:   t.AC,
	}
}

func randArmor() *Armor {
	// Pick an armor from the list at random
	i := rand.Intn(len(ArmorLib))
	var a *Armor
	for name := range ArmorLib {
		if i == 0 {
			a = newArmor(name)
		}
		i--
	}
	//randEnchant(item, 8, 20)
	return a
}

func (a *Armor) Equip(p *Player, msg *MessageLog) bool {
	p.equiped["armor"] = a
	p.AC = a.AC
	msg.Add("You are now wearing the %v.", a)
	return true
}

func (a *Armor) Unequip(p *Player, msg *MessageLog) bool {
	p.equiped["armor"] = nil
	p.AC = 10
	msg.Add("You take off the %v.", a)
	return true
}

func (a *Armor) Rune() rune {
	return ']'
}

func (a *Armor) GndString() string {
	return fmt.Sprintf("some %s", a)
}

func (a *Armor) InvString() string {
	return fmt.Sprintf("%+d %v [%d]", a.bonus, a, a.AC)
}

func (a Armor) String() string {
	return a.Name
}

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

// -----------------------------------------------------------------------

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

func randEnchant(item Equipable, enchantProb int, cursedProb int) int {
	// 10% chance of a cursed weapon with -1 to -3 penalty, and a 5% chance
	// of an enchanted weapon with a +1 to +3 bonus.
	var ench int
	if rand.Intn(100) < enchantProb { // enchanted
		ench = rand.Intn(2) + 1
	} else if rand.Intn(100) < cursedProb { // cursed
		ench = -1 * (rand.Intn(2) + 1)
	}
	return ench
}
