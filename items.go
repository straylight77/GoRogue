package main

import (
	"fmt"
	"math/rand"
)

type Item interface {
	Rune() rune
	InvString() string
	GndString() string
}

type Consumable interface {
	Rune() rune
	InvString() string
	GndString() string

	Consume(*GameState) bool
	//Identify()
}

type Equipable interface {
	Rune() rune
	InvString() string
	GndString() string

	Equip(*Player, *MessageLog) bool
	Unequip(*Player, *MessageLog) bool
	//Identify()
}

// -----------------------------------------------------------------------
type ItemList map[Coord]Item

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

func randItem() Item {
	roll := rand.Intn(100) + 1
	//debug.Add("rand item: roll=%d", roll)
	switch {
	case roll <= 27:
		return randPotion()
	case roll <= 54:
		return randPotion()
	case roll <= 72:
		return newFood("ration")
	case roll <= 81:
		return randWeapon()
	case roll <= 90:
		return randArmor()
	case roll <= 95:
		return randPotion()
	case roll <= 100:
		return randPotion()
	default:
		return newFood("slime mold")
	}

}

// === GOLD ==============================================================
type Gold struct {
	qty int
}

func (g *Gold) Rune() rune {
	return '*'
}

func (g *Gold) InvString() string {
	return g.GndString()
}

func (g *Gold) GndString() string {
	if g.qty == 1 {
		return fmt.Sprintf("1 piece of gold")
	} else {
		return fmt.Sprintf("%d pieces of gold", g.qty)
	}
}

func newGold(qty int) *Gold {
	return &Gold{qty: qty}
}

func randGoldAmt(depth int) int {
	return rand.Intn(50+10*depth) + 2
}

// === EFFECTS ===========================================================
const (
	E_Nothing = iota
	E_Healing
	E_ExtraHealing
	E_Strength
	E_Poison
	E_Confusion
	E_Blindness
	E_Restore
	E_DetMagic
	E_DetMonsters
	E_LevelUp
	E_Paralyze
	E_Haste
	E_Truesight
)

func doEffect(effect int, gs *GameState) {
	if effect == -1 {
		panic("Unkown effect id")
	}

	switch effect {
	case E_Nothing:
		//do nothing
	case E_Healing:
		gs.player.AdjustHP(gs.player.Level * 3)
		gs.player.SetTimer("blind", 0)
		gs.player.SetTimer("confusion", 0)
	case E_ExtraHealing:
		gs.player.AdjustHP(gs.player.Level * 5)
		gs.player.SetTimer("blind", 0)
		gs.player.SetTimer("confusion", 0)
	case E_Strength:
		gs.player.Str += 1
		gs.player.maxStr += 1
	case E_Poison:
		gs.player.Str -= rand.Intn(3) + 1
	case E_Restore:
		gs.player.Str = gs.player.maxStr
	case E_Blindness:
		gs.player.SetTimer("blind", 850)
	case E_Confusion:
		gs.player.SetTimer("confused", 20+rand.Intn(8))
	case E_DetMonsters:
		gs.player.SetTimer("detMonsters", 850)
	case E_DetMagic:
		gs.player.SetTimer("detMagic", 850)
	case E_LevelUp:
		gs.player.XP = XPTable[gs.player.Level]
	case E_Paralyze:
		gs.player.SetTimer("paralyzed", 3)
	case E_Haste:
		// if already hasted, faint for 0-7 turns
		gs.player.SetTimer("haste", rand.Intn(5)+10)
	case E_Truesight:
		gs.player.SetTimer("truesight", 850)
		gs.player.SetTimer("blind", 0)
	default:
		gs.messages.Add("This effect (%d) has not been implemented.", effect)
	}
}
