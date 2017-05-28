package main

import "strconv"
import "github.com/bennicholls/burl/core"
import "github.com/bennicholls/burl/util"

type Room struct {
	Name string

	X, Y          int
	Width, Height int

	RoomMap *core.TileMap

	connected []*Room

	state            core.Stat //state of repair.
	upkeep           int       //periodic decay of repair state.
	repairDifficulty int       //default time to repair by 1 unit.
}

func NewRoom(name string, x, y, w, h, upkeep, repair int) *Room {
	r := new(Room)
	r.X, r.Y, r.Width, r.Height = x, y, w, h
	r.Name = name

	r.RoomMap = core.NewMap(r.Width, r.Height)
	for i := 0; i < r.Width*r.Height; i++ {
		if i < r.Width || i%r.Width == 0 || i%r.Width == r.Width-1 || i/r.Width == r.Height-1 {
			r.RoomMap.ChangeTileType(i%r.Width, i/r.Width, TILE_WALL)
		} else {
			r.RoomMap.ChangeTileType(i%r.Width, i/r.Width, TILE_FLOOR)
		}
	}

	r.upkeep = upkeep
	r.repairDifficulty = repair
	r.state = core.NewStat(100)
	r.connected = make([]*Room, 0, 10)

	return r
}

//Tries to connect room to another. Finds the intersection of the two rooms and puts doors there!
//If rooms not properly lined up, does nothing.
func (r *Room) AddConnection(c *Room) {
	//check if room is already connected
	for _, room := range r.connected {
		if room == c {
			return
		}
	}

	//check if rooms intersect properly
	x, y, w, h := util.FindIntersectionRect(c, r)
	if w != 1 && h != 1 {
		return
	}

	//translate coords from shipspace to roomspace
	x, y = x-r.X, y-r.Y

	if w == 1 && h >= 3 {
		if h > 3 {
			r.RoomMap.ChangeTileType(x, y+h/2-1, TILE_DOOR)
		}
		r.RoomMap.ChangeTileType(x, y+h/2, TILE_DOOR)
	} else if h == 1 && w >= 3 {
		//up-down rooms
		if w > 3 {
			r.RoomMap.ChangeTileType(x+w/2-1, y, TILE_DOOR)
		}
		r.RoomMap.ChangeTileType(x+w/2, y, TILE_DOOR)
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
	roomstatus := r.Name + ": Status "
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
