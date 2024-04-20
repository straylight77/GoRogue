package main

import "math/rand"

type GameState struct {
	done           bool
	dungeon        *DungeonMap
	player         *Player
	monsters       *MonsterList
	messages       *MessageLog
	dmap           *DMap
	wander         int
	spawnFoodTimer int
	items          ItemList
}

// -----------------------------------------------------------------------
func (gs *GameState) Init() {

	assignPotionColors()

	gs.dungeon = &DungeonMap{}
	gs.player = &Player{}
	gs.monsters = &MonsterList{}
	gs.items = ItemList{}
	gs.messages = &MessageLog{}
	gs.wander = WanderTimer
	gs.spawnFoodTimer = SpawnFood

	gs.player.Init()

	// Set up player's starting equipment
	var item Equipable
	//gs.player.Pickup(newPotion("healing"))
	gs.player.Pickup(newFood("ration"))

	item = newWeapon("mace")
	gs.player.Pickup(item)
	item.Equip(gs.player, gs.messages)

	item = newArmor("ring mail")
	gs.player.Pickup(item)
	item.Equip(gs.player, gs.messages)

	//gs.player.Pickup(newPotion("confusion"))
	//gs.player.Pickup(newPotion("blindness"))
	//gs.player.Pickup(newPotion("paralysis"))

	generateRandomLevel(gs)
	gs.Pathfinding()
	gs.UpdatePlayerFOV()

	gs.messages.Clear() // equipping will create messages which we don't want here
	gs.messages.Add("Welcome to the Dungeons of Doom!")
}

// -----------------------------------------------------------------------
func (gs *GameState) MoveActor(a Actor, delta Coord) bool {

	// Override the direction if the entity is confused
	if a.IsConfused() {
		delta = gs.dungeon.RandDirectionCoords(a.Pos())
	}
	dest := a.Pos().Sum(delta)

	// Check edges of the map
	if gs.dungeon.IsOutOfBounds(dest) {
		gs.messages.Add("As you gaze into the abyss, it begins to gaze back into you...")
		return false
	}

	// Slightly different logic between monsters and the player
	switch a.(type) {

	case *Monster:
		// If player is there attack them
		if dest == gs.player.Pos() {
			gs.messages.Add(a.Attack(gs.player))
			return true
		}

		// If another monster is there, don't move
		m2 := gs.monsters.MonsterAt(dest)
		if m2 != nil {
			return false
		}

	case *Player:
		// Check if player is paralyzed
		if gs.player.Timer("paralyzed") > 0 {
			gs.messages.Add("You remain unable to move.")
			return true
		}

		// If a monster is there, attack it
		m := gs.monsters.MonsterAt(dest)
		if m != nil {
			gs.messages.Add(a.Attack(m))
			m.State = StateChase
			return true
		}
	}

	// Finally, check if the dungeon tile blocks movement or not
	if gs.dungeon.IsWalkable(a.Pos(), dest) {
		a.SetPos(dest)
		return true
	}

	return false
}

// -----------------------------------------------------------------------
func (gs *GameState) GoDownstairs() bool {
	if gs.player.IsParalyzed() {
		gs.messages.Add("You remain unable to move.")
		return true
	}

	if gs.dungeon.TileTypeAt(gs.player.Pos()) == TileStairsDn || debugFlag["main"] {
		gs.messages.Add("You descend the ancient stairs.")
		generateRandomLevel(gs)
		return true
	} else {
		gs.messages.Add("There are no stairs to go down here.")
		return false
	}
}

// -----------------------------------------------------------------------
func (gs *GameState) GoUpstairs() bool {
	if gs.player.IsParalyzed() {
		gs.messages.Add("You remain unable to move.")
		return true
	}

	if gs.dungeon.TileTypeAt(gs.player.Pos()) == TileStairsUp {
		gs.messages.Add("Your way is magically blocked.")
	} else {
		gs.messages.Add("There are no stairs to go up here.")
	}
	return false
}

// -----------------------------------------------------------------------
func (gs *GameState) CheckItems() {
	for pos, item := range gs.items {
		if pos == gs.player.Pos() {
			if gs.player.Pickup(item) {
				gs.messages.Add("You pick up %v.", item.GndString())
				delete(gs.items, pos)
			}
		}
	}
}

// -----------------------------------------------------------------------
func (gs *GameState) PruneMonsters() {
	for i, m := range *gs.monsters {
		if m.HP <= 0 {
			gs.monsters.Remove(i)
			if gs.player.IsBlind() {
				gs.messages.Add("You defeated something!")
			} else {
				gs.messages.Add("You defeated the %s!", m.Name)
			}
			gs.player.AddXP(m.XP)
		}
	}
	// This is the only place XP is awarded so check player level
	msg := gs.player.CheckLevel()
	gs.messages.Add(msg)
}

// -----------------------------------------------------------------------
func (gs *GameState) MonstersAct() {

	for _, m := range *gs.monsters {
		switch m.State {

		case StateDormant:
			if m.isMean && gs.dungeon.CanSee(m) && rand.Intn(100) < 67 {
				m.State = StateChase
			}

		case StateChase:

			if m.randMove > rand.Intn(100) {
				// Move randomly randMove% of the time (e.g. bats)
				delta := gs.dungeon.RandDirectionCoords(m.Pos())
				gs.MoveActor(m, delta)

			} else {
				// Pathfinding to the player is already calculated with the dmap
				m.nextStep = gs.dmap.NextStep(m.Pos())
				delta := m.DirectionCoordsTo(m.nextStep)
				gs.MoveActor(m, delta)

				// For testing, store the next step
				m.nextStep = gs.dmap.NextStep(Coord{m.X, m.Y})
			}
		}
	}
}

// -----------------------------------------------------------------------
func (gs *GameState) Pathfinding() {
	// Recalculate the DMap for monsters to use to find the player
	gs.dmap = newDMap(gs.dungeon, gs.player.Pos())
}

// -----------------------------------------------------------------------
// Update the player's field of view and visited tiles
func (gs *GameState) UpdatePlayerFOV() {
	gs.dungeon.SetVisible(Coord{0, 0}, MapMaxX, MapMaxY, false)
	gs.dungeon.playerFOV(gs.player.Pos())

	// If the player is in a room, light it up
	for _, r := range gs.dungeon.rooms {
		if r.InRoom(gs.player.Pos()) {
			gs.dungeon.SetVisible(r.TopLeft(), r.W+1, r.H+1, true)
		}
	}

}

// -----------------------------------------------------------------------
// After 70 turns, a “wander” daemon activates. When activated, every fourth
// move has a 1/6th chance that a monster will spawn and deactivate the daemon.
// Monsters spawned this way will immediately be hostile toward the player.
func (gs *GameState) WanderingMonsters() {
	if gs.wander > 0 {
		gs.wander--
	} else {
		if gs.player.moves%4 == 0 && rand.Intn(100) < 16 {

			// Find a random room that the player is not in
			r := rand.Intn(len(gs.dungeon.rooms))
			rm := gs.dungeon.rooms[r]
			for rm.InRoom(gs.player.Pos()) {
				r = rand.Intn(len(gs.dungeon.rooms))
				rm = gs.dungeon.rooms[r]
			}

			// Spawn a new wandering monster that is hostile
			m := randomMonster(gs.player.depth)
			for m.noWander {
				m = randomMonster(gs.player.depth)
			}

			m.State = StateChase
			gs.monsters.Add(m, rm.RandPoint())
			debug.Add("spawned: %v", m)

			// Reset the countdown
			gs.wander = WanderTimer
		}
	}
}

func (gs *GameState) IsBonusMove() bool {
	return gs.player.IsHasted() && gs.player.moves%2 == 0
}
