package main

type Player struct {
	Name string //players name, also name of their save file
	Credit int  //money to player's name. 

	MissionLog []Mission
	SpaceShip *Ship
}

func NewPlayer(n string) (p *Player) {
	p = new(Player)
	p.Name = n

	p.MissionLog = make([]Mission, 0, 20) //20 max missions? Should be fine.
	return
}