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

// Returns a string that describes the item as it appears on the ground
func (item Item) GndString() string {
	switch item.typ {
	case Gold:
		if item.val1 == 1 {
			return fmt.Sprintf("%d piece of gold", item.val1)
		} else {
			return fmt.Sprintf("%d pieces of gold", item.val1)
		}
	case Potion:
		if item.IsIdentified() {
			return fmt.Sprintf("a %s", item.name)
		} else {
			return fmt.Sprintf("a %s potion", item.Descriptor())
		}
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

// Returns a string that describes the item in a player's inventory
func (item Item) InvString() string {
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
		dice := fmt.Sprintf("%dd%d", item.val1, item.val2)
		if item.ench != 0 {
			dice = fmt.Sprintf("%dd%d%+d", item.val1, item.val2, item.ench)
		}
		return fmt.Sprintf("a %+d %s [%s]%s", item.ench, item.name, dice, cursed)
	case Armor:
		prot := item.val1 - item.ench
		return fmt.Sprintf("a %+d %s [%d]%s", item.ench, item.name, prot, cursed)
	case Ring:
		return fmt.Sprintf("a %s ring", item.name)
	case Potion:
		if item.IsIdentified() {
			return fmt.Sprintf("a %s", item.name)
		} else {
			return fmt.Sprintf("a %s potion", item.Descriptor())
		}
	case Scroll:
		return fmt.Sprintf("a scroll titled '%s'", item.name)
	case Stick:
		return fmt.Sprintf("a %s wand", item.name)
	default:
		return fmt.Sprintf("a %s", item.name)
	}
}

// Return the color, material, stone or title of the item
// Used for listing unidentified consumables
func (item Item) Descriptor() string {
	switch item.typ {
	case Potion:
		i := PotionLib[item.val1].color
		return PotionColors[i]
	default:
		return "mysterious"
	}
}

// Return the id of the effect the item has when triggered
func (item Item) Effect() int {
	switch item.typ {
	case Potion:
		return PotionLib[item.val1].effect
	}
	return -1
}

// If item is consumable, return if that type has been discovered. For
// equipment, return if this particular instance has been identified.
func (item Item) IsIdentified() bool {
	switch item.typ {
	case Potion:
		return PotionLib[item.val1].discovered
	default:
		return item.identified
	}
}

// See IsIdentified().  Set the appropriate flag.  We will only ever need
// setting to true, never to false.
func (item *Item) Identify() {
	switch item.typ {
	case Potion:
		PotionLib[item.val1].discovered = true
	default:
		item.identified = true
	}
}

// Returns the string to use as a game message after consuming the item.
func (item Item) ConsumeMsg() string {
	switch item.typ {
	case Potion:
		return PotionLib[item.val1].message
	default:
		return "Yum!"
	}
}

func (item Item) IsMagical() bool {
	switch item.typ {
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

// === POTIONS ==========================================================
// name: full name of the potion once it's been identified
// val1: index of PotionLib for this potion

type PotionTemplate struct {
	name       string
	effect     int
	color      int
	discovered bool
	message    string
}

var PotionLib = []PotionTemplate{
	{"thirst quenching", E_Nothing, 0, false, "Meh, tastes pretty dull."},
	{"healing", E_Healing, 0, false, "You begin to feel better."},
	{"extra healing", E_ExtraHealing, 0, false, "You begin to feel much better."},
	{"strength", E_Strength, 0, false, "You feel stronger, what bulging muscles!"},
	{"poison", E_Poison, 0, false, "You feel very sick now."},
	{"blindness", E_Blindness, 0, false, "A cloak of darkness falls around you."},
	{"confusion", E_Confusion, 0, false, "Wait, what's going on here. Huh? What? Who?"},
	{"restore strength", E_Restore, 0, false, "Hey, this tastes great, it make you feel warm all over."},
	{"detect magic", E_DetMagic, 0, false, "You sense the presence of magic."},
	{"monster detection", E_DetMonsters, 0, false, "You feel like you are not alone."},
	{"raise level", E_LevelUp, 0, false, "You feel more experienced."},
	{"paralysis", E_Paralyze, 0, false, "You feel your body seizing up, you can't move!"},
	{"haste", E_Haste, 0, false, "Tastes like coffee, everything seems to slow down."},
	{"truesight", E_Truesight, 0, false, "Tastes like slime-mold juice."},
	//see invisible
	//haste
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
	roll := rand.Intn(len(PotionLib))
	return newPotion(PotionLib[roll].name)
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
