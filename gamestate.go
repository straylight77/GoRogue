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
func (gs *GameState) MoveEntity(e Entity, dx, dy int) bool {
	destX, destY := e.Pos()
	destX += dx
	destY += dy

	// Check edges of the map
	if gs.dungeon.IsOutOfBounds(destX, destY) {
		gs.messages.Add("As you gaze into the abyss, it begins to gaze back into you...")
		return false
	}

	// Slightly different logic between monsters and the player
	switch e.(type) {

	case *Monster:
		// If player is there attack them
		if destX == gs.player.X && destY == gs.player.Y {
			gs.messages.Add(e.Attack(gs.player))
			return true
		}

		// If another monster is there, don't move
		m2 := gs.monsters.MonsterAt(destX, destY)
		if m2 != nil {
			return false
		}

	case *Player:
		// If a monster is there, attack it
		m := gs.monsters.MonsterAt(destX, destY)
		if m != nil {
			gs.messages.Add(e.Attack(m))
			m.State = StateChase
			return true
		}
	}

	// Finally, check if the dungeon tile blocks movement or not
	origX, origY := e.Pos()
	if gs.dungeon.IsWalkable(Coord{origX, origY}, Coord{destX, destY}) {
		e.SetPos(destX, destY)
		return true
	}

	return false
}

// -----------------------------------------------------------------------
func (gs *GameState) GoDownstairs() bool {
	if gs.dungeon.TileTypeAt(gs.player.X, gs.player.Y) == TileStairsDn {
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
	if gs.dungeon.TileTypeAt(gs.player.X, gs.player.Y) == TileStairsUp {
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
			if (m.isMean && gs.dungeon.CanSee(m)) || m.DistanceFrom(gs.player) <= 2 {
				m.State = StateChase
			}

		case StateChase:

			if !m.isMean && m.DistanceFrom(gs.player) > 6 {
				// For non-mean monsters, go dormant when far away
				m.State = StateDormant

			} else if m.randMove > rand.Intn(100) {
				// Move randomly randMove% of the time
				dx, dy := gs.dungeon.RandDirectionCoords(m.X, m.Y)
				gs.MoveEntity(m, dx, dy)

			} else {
				// Pathfinding to the player is already calculated with the dmap
				m.nextStep = gs.dmap.NextStep(Coord{m.X, m.Y})
				dx, dy := m.DirectionCoordsTo(m.nextStep.X, m.nextStep.Y)
				gs.MoveEntity(m, dx, dy)

				// For testing, store the next step
				m.nextStep = gs.dmap.NextStep(Coord{m.X, m.Y})
			}
		}
	}
}

// -----------------------------------------------------------------------
func (gs *GameState) Pathfinding() {
	// Recalculate the DMap for monsters to use to find the player
	gs.dmap = newDMap(gs.dungeon, Coord{gs.player.X, gs.player.Y})
}

// -----------------------------------------------------------------------
// Update the player's field of view and visited tiles
func (gs *GameState) UpdatePlayerFOV() {
	gs.dungeon.SetVisible(0, 0, MapMaxX, MapMaxY, false)
	gs.dungeon.playerFOV(gs.player)

	// If the player is in a room, light it up
	for _, r := range gs.dungeon.rooms {
		if r.InRoom(gs.player.X, gs.player.Y) {
			gs.dungeon.SetVisible(r.X, r.Y, r.W+1, r.H+1, true)
		}
	}

}
