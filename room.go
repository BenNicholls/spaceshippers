package main

import "strconv"

//room indexes
const (
	BRIDGE int = iota
	ENGINEERING
	MESSHALL
	MEDBAY
	QUARTERS
	HALLWAY
	MAX_ROOMS
)

type Room struct {
	name string

	x, y int
	w, h int //rectange rooms for now, use a little bitmap later

	state            Stat //state of repair.
	upkeep           int  //periodic decay of repair state.
	repairDifficulty int  //default time to repair by 1 unit.
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
