package main

//import "strconv"
import (
	"github.com/bennicholls/burl-E/burl"
)

type Room struct {
	Name          string
	Description   string
	Roomtype      RoomType
	Width, Height int

	X, Y    int
	RoomMap *burl.TileMap

	connected []*Room

	Stats []RoomStat
}

func NewRoom(name string, t RoomType, w, h int) *Room {
	r := new(Room)
	r.Width, r.Height = w, h
	r.Name = name
	r.Stats = make([]RoomStat, 0)
	r.Description = "a room"

	r.CreateRoomMap()
	r.connected = make([]*Room, 0, 10)

	return r
}

func NewRoomFromTemplate(temp RoomTemplate) (r *Room) {
	r = new(Room)
	r.Name = temp.name
	r.Description = temp.description
	r.Width = temp.width
	r.Height = temp.height
	r.Roomtype = temp.roomType
	r.Stats = temp.stats

	r.CreateRoomMap()
	r.connected = make([]*Room, 0, 10)

	return
}

func (r *Room) CreateRoomMap() {
	r.RoomMap = burl.NewMap(r.Width, r.Height)
	for i := 0; i < r.Width*r.Height; i++ {
		if i < r.Width || i%r.Width == 0 || i%r.Width == r.Width-1 || i/r.Width == r.Height-1 {
			r.RoomMap.ChangeTileType(i%r.Width, i/r.Width, TILE_WALL)
		} else {
			r.RoomMap.ChangeTileType(i%r.Width, i/r.Width, TILE_FLOOR)
		}
	}
}

func (r *Room) Rotate() {
	r.Width, r.Height = r.Height, r.Width
	r.CreateRoomMap()
}

//Tries to connect room to another. Finds the intersection of the two rooms and puts doors there!
//If rooms not properly lined up, does nothing.
func (r *Room) AddConnection(c *Room) {
	//ensure room isn't trying to connect with itself
	if r == c {
		return
	}

	//check if room is already connected
	for _, room := range r.connected {
		if room == c {
			return
		}
	}

	//check if rooms intersect properly
	x, y, w, h := burl.FindIntersectionRect(c, r)
	if w != 1 && h != 1 {
		return
	}

	//translate coords from shipspace to roomspace
	x, y = x-r.X, y-r.Y

	if w == 1 && h >= 3 {
		if h%2 == 0 {
			r.RoomMap.ChangeTileType(x, y+h/2-1, TILE_DOOR)
		}
		r.RoomMap.ChangeTileType(x, y+h/2, TILE_DOOR)
	} else if h == 1 && w >= 3 {
		//up-down rooms
		if w%2 == 0 {
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

			//redraw walls over now-nonexistent doors
			x, y, w, h := burl.FindIntersectionRect(r, c)
			for i := 0; i < w*h; i++ {
				r.RoomMap.ChangeTileType(x+i%w-r.X, y+i/w-r.Y, TILE_WALL)
			}
		}
	}
}

func (r Room) GetStatus() string {
	roomstatus := r.Name + ": Status OKAY FOR NOW"
	// if r.RepairState.Get() > 80 {
	// 	roomstatus += "NOMINAL."
	// } else if r.RepairState.Get() > 50 {
	// 	roomstatus += "FINE."
	// } else if r.RepairState.Get() > 20 {
	// 	roomstatus += "NEEDS REPAIR."
	// } else if r.RepairState.Get() > 0 {
	// 	roomstatus += "CRITICAL."
	// } else {
	// 	roomstatus += "DESTROYED."
	// }

	// roomstatus += " (" + strconv.Itoa(r.RepairState.Get()) + "/100)"

	return roomstatus
}

func (r *Room) ApplyUpkeep(spaceTime int) {
	// if r.Upkeep == 0 {
	// 	return
	// } else if spaceTime%r.Upkeep == 0 {
	// 	r.RepairState.Mod(-1)
	// }
}

func (r *Room) Update(spaceTime int) {
	//r.ApplyUpkeep(spaceTime)
}

func (r Room) Rect() (int, int, int, int) {
	return r.X, r.Y, r.Width, r.Height
}
