package main

import "github.com/bennicholls/burl-E/burl"

type Player struct {
	Person

	Credit int //money to player's name.

	MissionLog []*Mission

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

	p.MissionLog = make([]*Mission, 0, 0)
	return
}

//Adds a mission to the player's missionlog.
func (p *Player) AddMission(m *Mission) {
	p.MissionLog = append(p.MissionLog, m)
	burl.LogInfo("Added Mission: ", m.name)
	burl.PushEvent(burl.NewEvent(burl.EV_UPDATE_UI, "missions"))
}
