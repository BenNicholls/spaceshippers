package main

type Player struct {
	Person

	Credit int //money to player's name.

	MissionLog []Mission

	//EventLog tracks which events (referenced by the int event.id) the player has handled.
	//Naturally does not track repeating events, just unique and story-based ones.
	EventLog map[int]bool

	//The most important thing in the game.
	SpaceShip *Ship
}

func NewPlayer(n string) (p *Player) {
	p = new(Player)
	p.Name = n
	p.Ptype = PERSON_PLAYER

	p.MissionLog = make([]Mission, 0, 20) //20 max missions? Should be fine.
	return
}
