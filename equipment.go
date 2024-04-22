package main

import (
	"fmt"
	"math/rand"
)

// === WEAPONS ===========================================================

type Weapon struct {
	name   string
	damage Dice
	ench   int
	cursed bool
	worth  int
}

// -----------------------------------------------------------------------
func newWeapon(name string) *Weapon {
	t, ok := WeaponLib[name]
	if !ok {
		panic("No weapon with the name " + name)
	}

	return &Weapon{
		name:   name,
		damage: parseDice(t.melee),
		worth:  t.worth,
	}
}

// -----------------------------------------------------------------------
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
	w.ench, w.cursed = randEnchant(5, 10)

	return w
}

// -----------------------------------------------------------------------
func (w *Weapon) Equip(p *Player, msg *MessageLog) bool {
	if p.equiped["weapon"] == w {
		w.Unequip(p, msg)
		return false
	}
	if p.equiped["weapon"] != nil {
		msg.Add("You need to put away the %v first.", p.equiped["weapon"])
		return false
	}
	p.equiped["weapon"] = w
	p.Melee = w.damage
	msg.Add("You are now wielding the %v.", w)
	return true
}

// -----------------------------------------------------------------------
func (w *Weapon) Unequip(p *Player, msg *MessageLog) bool {
	if p.equiped["weapon"] == nil {
		msg.Add("You aren't wielding the %v.", w)
		return false
	}
	if w.cursed {
		msg.Add("You cannot put away the %v, it's cursed!", w)
		return false
	}
	p.equiped["weapon"] = nil
	msg.Add("You put away the %v.", w)
	return true
}

// -----------------------------------------------------------------------
func (w *Weapon) Rune() rune {
	return ')'
}

func (w *Weapon) GndString() string {
	return fmt.Sprintf("a %s", w)
}

func (w *Weapon) InvString() string {
	cursed := ""
	if w.cursed {
		cursed = " {cursed}"
	}
	return fmt.Sprintf("%+d %v [%s]%s", w.ench, w, w.damage, cursed)
}

func (w *Weapon) Worth() int {
	if w.ench < 0 {
		return 0
	} else {
		return (1 + (10 * w.ench)) * w.worth
	}
}

func (w Weapon) String() string {
	return w.name
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

// === ARMOR =============================================================

type Armor struct {
	Name   string
	AC     int
	ench   int
	cursed bool
	worth  int
}

// -----------------------------------------------------------------------
func newArmor(name string) *Armor {
	t, ok := ArmorLib[name]
	if !ok {
		panic("No armor with the name " + name)
	}

	return &Armor{
		Name:  name,
		AC:    t.AC,
		worth: t.worth,
	}
}

// -----------------------------------------------------------------------
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
	a.ench, a.cursed = randEnchant(8, 20)
	return a
}

// -----------------------------------------------------------------------
func (a *Armor) Equip(p *Player, msg *MessageLog) bool {
	if p.equiped["armor"] == a {
		a.Unequip(p, msg)
		return false
	}
	if p.equiped["armor"] != nil {
		msg.Add("You need to take off the %v first.", p.equiped["armor"])
		return false
	}
	p.equiped["armor"] = a
	p.AC = a.AC - a.ench
	msg.Add("You are now wearing the %v.", a)
	return true
}

// -----------------------------------------------------------------------
func (a *Armor) Unequip(p *Player, msg *MessageLog) bool {
	if p.equiped["armor"] == nil {
		msg.Add("You aren't wearing the %v.", a)
		return false
	}
	if a.cursed {
		msg.Add("You cannot take off the %v, it's cursed!", a)
		return false
	}
	p.equiped["armor"] = nil
	p.AC = 10
	msg.Add("You take off the %v.", a)
	return true
}

// -----------------------------------------------------------------------
func (a *Armor) Rune() rune {
	return ']'
}

func (a *Armor) GndString() string {
	return fmt.Sprintf("some %s", a)
}

func (a *Armor) InvString() string {
	cursed := ""
	if a.cursed {
		cursed = " {cursed}"
	}
	return fmt.Sprintf("%+d %v [%d]%s", a.ench, a, a.AC-a.ench, cursed)
}

func (a *Armor) Worth() int {
	if a.ench < 0 {
		return 0
	} else {
		return (1 + (10 * a.ench)) * a.worth
	}
}

func (a Armor) String() string {
	return a.Name
}

// -----------------------------------------------------------------------
type ArmorTemplate struct {
	AC    int
	worth int
}

var ArmorLib = map[string]ArmorTemplate{
	"leather armor": {8, 5},
	"ring mail":     {7, 30},
	"scale mail":    {6, 3},
	"chain mail":    {5, 75},
	"banded mail":   {4, 90},
	"plate mail":    {3, 440},
}

// =======================================================================

func randEnchant(enchantProb int, cursedProb int) (int, bool) {
	// 10% chance of a cursed weapon with -1 to -3 penalty, and a 5% chance
	// of an enchanted weapon with a +1 to +3 bonus.
	var ench int
	if rand.Intn(100) < enchantProb { // enchanted
		ench = rand.Intn(2) + 1
	} else if rand.Intn(100) < cursedProb { // cursed
		ench = -1 * (rand.Intn(2) + 1)
	}
	cursed := false
	if ench < 0 {
		cursed = true
	}
	return ench, cursed
}
