package main

import "math/rand"

type GameState struct {
	done     bool
	dungeon  *DungeonMap
	player   *Player
	monsters *MonsterList
	messages *MessageLog
	dmap     *DMap
}

// -----------------------------------------------------------------------
func (gs *GameState) Init() {

	gs.dungeon = &DungeonMap{}
	gs.player = &Player{}
	gs.monsters = &MonsterList{}
	gs.messages = &MessageLog{}

	gs.player.Init()
	generateRandomLevel(gs)
}

// -----------------------------------------------------------------------
func (gs *GameState) MoveEntity(e Entity, delta Coord) bool {
	dest := e.Pos().Sum(delta)

	// Check edges of the map
	if gs.dungeon.IsOutOfBounds(dest) {
		gs.messages.Add("As you gaze into the abyss, it begins to gaze back into you...")
		return false
	}

	// Slightly different logic between monsters and the player
	switch e.(type) {

	case *Monster:
		// If player is there attack them
		if dest == gs.player.Pos() {
			gs.messages.Add(e.Attack(gs.player))
			return true
		}

		// If another monster is there, don't move
		m2 := gs.monsters.MonsterAt(dest)
		if m2 != nil {
			return false
		}

	case *Player:
		// If a monster is there, attack it
		m := gs.monsters.MonsterAt(dest)
		if m != nil {
			gs.messages.Add(e.Attack(m))
			m.State = StateChase
			return true
		}
	}

	// Finally, check if the dungeon tile blocks movement or not
	if gs.dungeon.IsWalkable(e.Pos(), dest) {
		e.SetPos(dest)
		return true
	}

	return false
}

// -----------------------------------------------------------------------
func (gs *GameState) GoDownstairs() bool {
	if gs.dungeon.TileTypeAt(gs.player.Pos()) == TileStairsDn {
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
	if gs.dungeon.TileTypeAt(gs.player.Pos()) == TileStairsUp {
		gs.messages.Add("Your way is magically blocked.")
	} else {
		gs.messages.Add("There are no stairs to go up here.")
	}
	return false
}

// -----------------------------------------------------------------------
func (gs *GameState) UpdateMonsters() {

	for i, m := range *gs.monsters {

		// Remove any slain monsters
		// TODO move this into a separate function under MonsterList
		if m.HP <= 0 {
			gs.monsters.Remove(i)
			gs.messages.Add("You defeated the %s!", m.Name)
			m := gs.player.AddXP(m.XP)
			if m != "" {
				gs.messages.Add(m)
			}
			continue
		}

		switch m.State {

		case StateDormant:
			if (m.isMean && gs.dungeon.CanSee(m)) || m.Pos().Distance(gs.player.Pos()) <= 2 {
				m.State = StateChase
				//gs.messages.Add("The %s wakes up.", m.Name)
			}

		case StateChase:

			if !m.isMean && m.Pos().Distance(gs.player.Pos()) > 6 {
				// For non-mean monsters, go dormant when far away
				m.State = StateDormant

			} else if m.randMove > rand.Intn(100) {
				// Move randomly randMove% of the time
				delta := gs.dungeon.RandDirectionCoords(m.Pos())
				gs.MoveEntity(m, delta)

			} else {
				// Pathfinding to the player is already calculated with the dmap
				m.nextStep = gs.dmap.NextStep(m.Pos())
				delta := m.DirectionCoordsTo(m.nextStep)
				gs.MoveEntity(m, delta)

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