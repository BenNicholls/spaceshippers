package main

import "strconv"
import "github.com/bennicholls/burl/core"

type Room struct {
	name string

	X, Y int
	Width, Height int

	connected []*Room

	state            core.Stat //state of repair.
	upkeep           int  //periodic decay of repair state.
	repairDifficulty int  //default time to repair by 1 unit.
}

func NewRoom(name string, x, y, w, h, upkeep, repair int) *Room {
	r := new(Room)
	r.X, r.Y, r.Width, r.Height = x, y, w, h
	r.upkeep = upkeep
	r.repairDifficulty = repair
	r.state = core.NewStat(100)
	r.connected = make([]*Room, 10)

	return r
}

func (r *Room) AddConnection(c *Room) {
	//check if room is already connected
	for _, room := range r.connected {
		if room == c {
			return
		}
	}
	
	r.connected = append(r.connected, c)
}

func (r *Room) RemoveConnection(c *Room) {
	for i, room := range r.connected {
		if room == c {
			r.connected = append(r.connected[:i], r.connected[i+1:]...)
		}
	}
}

func (r Room) GetStatus() string {
	roomstatus := r.name + ": Status "
	if r.state.Get() > 80 {
		roomstatus += "NOMINAL."
	} else if r.state.Get() > 50 {
		roomstatus += "FINE."
	} else if r.state.Get() > 20 {
		roomstatus += "NEEDS REPAIR."
	} else if r.state.Get() > 0 {
		roomstatus += "CRITICAL."
	} else {
		roomstatus += "DESTROYED."
	}

	roomstatus += " (" + strconv.Itoa(r.state.Get()) + "/100)"

	return roomstatus
}

func (r *Room) ApplyUpkeep(spaceTime int) {
	if r.upkeep == 0 {
		return
	} else if spaceTime%r.upkeep == 0 {
		r.state.Mod(-1)
	}
}

func (r *Room) Update(spaceTime int) {
	r.ApplyUpkeep(spaceTime)
}

func (r Room) Rect() (int, int, int, int) {
	return r.X, r.Y, r.Width, r.Height
}