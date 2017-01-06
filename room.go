package main

import "strconv"
import "github.com/bennicholls/delvetown/console"

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

	x,y int
	w,h int //rectange rooms for now, use a little bitmap later

	state            Stat //state of repair.
	upkeep           int //periodic decay of repair state.
	repairDifficulty int //default time to repair by 1 unit.
}

func (r Room) PrintStatus() {
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

	output.Append(roomstatus)
}

func (r Room) Draw() {

	//get colour (interp over state for now, red -> green)
	b := console.MakeColour(255 - (255*r.state.GetPct()/100), 255*r.state.GetPct()/100, 0)


	var left, right, up, down bool
	var g int
	for j := 0; j < r.h; j++ {
		up = (j == 0)
		down = (j == r.h-1)

		for i := 0; i < r.w; i++ {
			left = (i==0)
			right = (i==r.w-1)
			g = 0
			if up {
				g += 1
			}
			if right {
				g += 2
			}
			if down {
				g+=4
			}
			if left {
				g += 8
			}

			shipdisplay.Draw(r.x + i, r.y + j, 0x80 + g, 0xFFFFFFFF, b)
		}
	}
}

func (r *Room) ApplyUpkeep() {
	if r.upkeep == 0 {
		return
	} else if SpaceTime%r.upkeep == 0 {
		r.state.Mod(-1)
	}
}

func (r *Room) Update() {
	r.ApplyUpkeep()
}
