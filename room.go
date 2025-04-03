package main

import (
	"slices"

	"github.com/bennicholls/tyumi/gfx"
	"github.com/bennicholls/tyumi/rl"
	"github.com/bennicholls/tyumi/vec"
)

type Room struct {
	Name          string
	Description   string
	Roomtype      RoomType
	Rotated       bool
	Width, Height int

	pos     vec.Coord
	RoomMap rl.TileMap
	atmo    GasMixture

	connected []*Room

	Stats []RoomStat
}

func NewRoom(name string, t RoomType, w, h int) *Room {
	r := new(Room)
	r.Width, r.Height = w, h
	r.Name = name
	r.Roomtype = t
	r.Stats = make([]RoomStat, 0)
	r.Description = "a room"

	r.CreateRoomMap()
	r.connected = make([]*Room, 0, 10)
	r.atmo.InitStandardAtmosphere(float64(r.Volume() * 1000)) //conversion to litres

	return r
}

func NewRoomFromTemplate(temp RoomTemplate) (r *Room) {
	r = NewRoom(temp.name, temp.roomType, temp.width, temp.height)
	r.Description = temp.description
	r.Stats = temp.stats

	return
}

func (r Room) Size() vec.Dims {
	return vec.Dims{r.Width, r.Height}
}

func (r *Room) Bounds() vec.Rect {
	return vec.Rect{r.pos, r.Size()}
}

// Creates the map for the room.
func (r *Room) CreateRoomMap() {
	r.RoomMap.Init(r.Size(), rl.TILE_NONE)
	for cursor := range vec.EachCoordInArea(r.RoomMap) {
		if cursor.IsInPerimeter(r.RoomMap) {
			r.RoomMap.SetTileType(cursor, TILE_WALL)
		} else {
			r.RoomMap.SetTileType(cursor, TILE_FLOOR)
		}
	}
}

func (r *Room) Rotate() {
	r.Width, r.Height = r.Height, r.Width
	r.Rotated = !r.Rotated
	r.CreateRoomMap()
}

// Tries to connect room to another. Finds the intersection of the two rooms and puts doors there!
// If rooms not properly lined up, does nothing.
func (r *Room) AddConnection(c *Room) {
	//ensure room isn't trying to connect with itself
	if r == c {
		return
	}

	//check if room is already connected
	if slices.Contains(r.connected, c) {
		return
	}

	//check if rooms intersect properly
	i := vec.FindIntersectionRect(c, r)
	if i.W != 1 && i.H != 1 {
		return
	}

	i.Coord = i.Coord.Subtract(r.pos) // translate from shipspace to roomspace
	if i.W == 1 && i.H >= 3 {         // left-right rooms
		if i.H%2 == 0 {
			r.RoomMap.SetTileType(i.Coord.Add(vec.Coord{0, i.H/2 - 1}), TILE_DOOR)
		}
		r.RoomMap.SetTileType(i.Coord.Add(vec.Coord{0, i.H / 2}), TILE_DOOR)
	} else if i.H == 1 && i.W >= 3 { //up-down rooms
		if i.W%2 == 0 {
			r.RoomMap.SetTileType(i.Coord.Add(vec.Coord{i.W/2 - 1, 0}), TILE_DOOR)
		}
		r.RoomMap.SetTileType(i.Coord.Add(vec.Coord{i.W / 2, 0}), TILE_DOOR)
	}

	r.connected = append(r.connected, c)
}

func (r *Room) RemoveConnection(c *Room) {
	c_i := slices.Index(r.connected, c)
	r.connected = slices.Delete(r.connected, c_i, c_i+1)

	//redraw walls over now-nonexistent doors
	sect := vec.FindIntersectionRect(r, c)
	for cursor := range vec.EachCoordInArea(sect) {
		r.RoomMap.SetTileType(cursor.Subtract(r.pos), TILE_WALL)
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

// Volume returns the interior volume of the room in m^3
func (r Room) Volume() int {
	return (r.Width - 2) * (r.Height - 2) * 3 //rooms are 3 meters high... right?
}

func (r Room) Draw(dst_canvas *gfx.Canvas, offset vec.Coord, depth int) {
	r.RoomMap.Draw(dst_canvas, offset, depth)
}
