# GoRogue
Back to basics!  A rewrite from scratch of the original game [Rogue](https://en.wikipedia.org/wiki/Rogue_(video_game)) as it was released circa 1980, with a few modern features thrown in, written in Go using [Tcell](https://github.com/gdamore/tcell). The goal is to recreate the playing experience of the original game as much as possible.  

![Screenshot](screenshot.png)

Why?  To teach myself Go, to develop my own framework for future roguelike games, to limit my scope so I actually produce a playable game and, of course, to have a bit of fun.

## Feature Roadmap
```
[X] Basic dungeon map and walking around 
[X] Random level generation 
    [X] Dungeon layout (3X3 Rogue-style)
    [X] Randomized monsters
    [X] Randomized gold
    [X] Randomized items
    [ ] Traps
    [ ] Dark rooms
    [ ] Hidden doors
[X] Monsters
    [X] Stats for all monsters
    [X] Basic states (dormant, chasing player)
    [X] Random movement (e.g. bats or confusion) 
    [X] Chasing the player (pathfinding)
    [X] Spawning wandering monsters 
[X] Player
    [X] Awarding XP and leveling up
    [X] Natural healing
    [X] Hunger
    [X] Inventory
    [X] Combat (AD&D 2nd edition rules)
[ ] Items
    [X] Gold
    [X] Food
    [/] Weapons
    [X] Armor
    [X] Potions
    [ ] Scrolls
    [ ] Rings
    [ ] Sticks
    [/] Cursed items and identification
[ ] Gameplay
    [X] Player score
    [ ] Title screen
    [X] End game screen
    [ ] Tracking high scores
    [ ] Amulet of Yendor
```

## Contributing 
Not yet.  Not until I hit version 1.0 but feedback is welcome!
