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
	typ        ItemType
	name       string
	identified bool
	val1       int
	val2       int
	val3       int
	val4       int
	ench       int
	magical    bool
	cursed     bool
}

func (item Item) Rune() rune {
	ch, ok := ItemRunes[item.typ]
	if !ok {
		ch = '0' // shouldn't see this but here's a default just in case
	}
	return ch
}

func (item Item) Type() ItemType {
	return item.typ
}

func (item Item) TypeString() string {
	switch item.Type() {
	case Potion:
		return "potion"
	case Scroll:
		return "scroll"
	case Ring:
		return "ring"
	case Stick:
		return "wand"
	default:
		return item.name
	}
}

// Returns a string that describes the item as it appears on the ground
func (item Item) GndString() string {
	switch item.Type() {
	case Gold:
		if item.val1 == 1 {
			return fmt.Sprintf("%d piece of gold", item.val1)
		} else {
			return fmt.Sprintf("%d pieces of gold", item.val1)
		}
	case Potion, Ring, Stick:
		if item.IsIdentified() {
			return fmt.Sprintf("a %s", item.name)
		} else {
			return fmt.Sprintf("a %s %s", item.Descriptor(), item.TypeString())
		}
	case Scroll:
		return fmt.Sprintf("a scroll titled '%s'", item.name)
	default:
		return fmt.Sprintf("a %s", item.name)
	}
}

// Returns a string that describes the item in a player's inventory
func (item Item) InvString() string {
	cursed := ""
	if item.IsCursed() {
		cursed = " {cursed}"
	}

	switch item.Type() {
	case Gold:
		return item.GndString()
	case Weapon:
		dice := fmt.Sprintf("%dd%d", item.val1, item.val2)
		if item.ench != 0 {
			dice = fmt.Sprintf("%dd%d%+d", item.val1, item.val2, item.ench)
		}
		return fmt.Sprintf("a %+d %s [%s]%s", item.ench, item.name, dice, cursed)
	case Armor:
		prot := item.val1 - item.ench
		return fmt.Sprintf("a %+d %s [%d]%s", item.ench, item.name, prot, cursed)
	default:
		return item.GndString()
	}
}

// Return the color, material, stone or title of the item
// Used for listing unidentified consumables
func (item Item) Descriptor() string {
	switch item.Type() {
	case Potion:
		i := PotionLib[item.val1].color
		return PotionColors[i]
	default:
		return "mysterious"
	}
}

// Return the id of the effect the item has when triggered
func (item Item) Effect() int {
	switch item.Type() {
	case Potion:
		return PotionLib[item.val1].effect
	}
	return -1
}

// If item is consumable, return if that type has been discovered. For
// equipment, return if this particular instance has been identified.
func (item Item) IsIdentified() bool {
	switch item.Type() {
	case Potion:
		return PotionLib[item.val1].discovered
	default:
		return item.identified
	}
}

// See IsIdentified().  Set the appropriate flag.  We will only ever need
// setting to true, never to false.
func (item *Item) Identify() {
	switch item.Type() {
	case Potion:
		PotionLib[item.val1].discovered = true
	default:
		item.identified = true
	}
}

// Returns the string to use as a game message after consuming the item.
func (item Item) ConsumeMsg() string {
	switch item.Type() {
	case Potion:
		return PotionLib[item.val1].message
	default:
		return "Yum!"
	}
}

func (item Item) IsMagical() bool {
	switch item.Type() {
	case Potion, Scroll, Stick, Ring:
		return true
	case Weapon, Armor:
		return item.ench > 0
	default:
		return false
	}
}

func (item Item) IsCursed() bool {
	return item.cursed
}

// === GOLD ==============================================================
func newGold(qty int) *Item {
	return &Item{typ: Gold, val1: qty}
}

func (item Item) GoldQty() int {
	return item.val1
}

func randGoldAmt(depth int) int {
	return rand.Intn(50+10*depth) + 2
}

// === EFFECTS ===========================================================
// (move this to its own file effects.go)
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
	roll := rand.Intn(100) + 1
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
