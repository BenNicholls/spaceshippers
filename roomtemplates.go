package main

//Similarly to the shiptemplates.go file, this file is just going to hold
//room data until I can be bothered to write a ststem that imports all of this
//from raws.

var roomTemplates map[RoomType]RoomTemplate

type RoomTemplate struct {
	name                  string
	description           string
	roomType              RoomType
	width, height         int  //default dimensions
	resizable             bool //if a room is resizable, this will be set and
	min_width, min_height int  //these paramaeters will be set. Otherwise,
	max_width, max_height int  //these will be zero.

	stats []RoomStat

	//decorations []Decor //eventually will hold things like how many tables/computers/chairs/crap for the room to contain
}

type RoomType int //specific variety of room

const (
	ROOM_BRIDGE RoomType = iota
	ROOM_COCKPIT
	ROOM_CORRIDOR_H
	ROOM_CORRIDOR_V
	ROOM_ENGINE_SMALL
	ROOM_ENGINE_MEDIUM
	ROOM_ENGINE_LARGE
	ROOM_CARGOBAY
	ROOM_QUARTERS
	ROOM_COMMONAREA
	ROOM_LABORATORY
	ROOM_MEDBAY
	ROOM_LASERTURRET
)

func init() {
	roomTemplates = make(map[RoomType]RoomTemplate)

	roomTemplates[ROOM_ENGINE_SMALL] = RoomTemplate{
		name:        "Engine Room - Small",
		description: "A small engine room, with small sub-light engines.",
		roomType:    ROOM_ENGINE_SMALL,
		width:       3,
		height:      5,
		stats: []RoomStat{
			RoomStat{
				Stat:     STAT_SUBLIGHT_THRUST,
				Modifier: 10,
			},
			RoomStat{
				Stat:     STAT_SUBLIGHT_FUELUSE,
				Modifier: 2,
			},
			RoomStat{
				Stat:     STAT_SUBLIGHT_POWER,
				Modifier: 30,
			},
		},
	}

	roomTemplates[ROOM_ENGINE_MEDIUM] = RoomTemplate{
		name:        "Engine Room - Medium",
		description: "A medium sized engine room, with decent sub-light engines.",
		roomType:    ROOM_ENGINE_MEDIUM,
		width:       5,
		height:      8,
		stats: []RoomStat{
			RoomStat{
				Stat:     STAT_SUBLIGHT_THRUST,
				Modifier: 25,
			},
			RoomStat{
				Stat:     STAT_SUBLIGHT_FUELUSE,
				Modifier: 3,
			},
			RoomStat{
				Stat:     STAT_SUBLIGHT_POWER,
				Modifier: 50,
			},
		},
	}
}

func CreateRoomFromTemplate(room RoomType) (r *Room) {
	r = new(Room)

	if temp, ok := roomTemplates[room]; ok {
		r = NewRoomFromTemplate(temp)
	} else {
		switch room {
		case ROOM_BRIDGE:
			r = NewRoom("Bridge", ROOM_BRIDGE, 6, 12)
		case ROOM_COCKPIT:
			r = NewRoom("Cockpit", ROOM_COCKPIT, 4, 3)
		case ROOM_CORRIDOR_H:
			r = NewRoom("Corridor", ROOM_CORRIDOR_H, 12, 4)
		case ROOM_CORRIDOR_V:
			r = NewRoom("Corridor", ROOM_CORRIDOR_V, 4, 12)
		case ROOM_ENGINE_LARGE:
			r = NewRoom("Engineering", ROOM_ENGINE_LARGE, 7, 9)
		case ROOM_CARGOBAY:
			r = NewRoom("Cargo Bay", ROOM_CARGOBAY, 8, 8)
		case ROOM_QUARTERS:
			r = NewRoom("Quarters", ROOM_QUARTERS, 6, 6)
		case ROOM_COMMONAREA:
			r = NewRoom("Common Area", ROOM_COMMONAREA, 9, 7)
		case ROOM_LABORATORY:
			r = NewRoom("Laboratory", ROOM_LABORATORY, 7, 7)
		case ROOM_MEDBAY:
			r = NewRoom("Medical Bay", ROOM_MEDBAY, 5, 5)
		case ROOM_LASERTURRET:
			r = NewRoom("Laser Turret", ROOM_LASERTURRET, 3, 3)
		}
	}

	return
}
