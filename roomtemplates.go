package main

import (
	"github.com/bennicholls/burl-E/burl"
)

//This file is just going to hold room data until I can be bothered to
//write a system that imports all of this from raws.
//CONSIDER: Changing the word "room" to module?? That's what we'll be referring
//to these as in the game.
//THINK: If we change the order of the roomtype list or add ones in between, saved ships break. How to fix???

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
	ROOM_CORRIDOR
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

	roomTemplates[ROOM_BRIDGE] = RoomTemplate{
		name:        "Bridge",
		description: "A large module that acts as the command centre of the ship. Home of The Captain's Chair.",
		roomType:    ROOM_BRIDGE,
		width:       7,
		height:      9,
		stats: []RoomStat{
			RoomStat{
				Stat:     STAT_POWER_DRAW,
				Modifier: 10,
			},
			RoomStat{
				Stat:     STAT_CO2_SCRUBRATE,
				Modifier: 20, //volume of gas (L) that can be purged of CO2 per second
			},
		},
	}

	roomTemplates[ROOM_COCKPIT] = RoomTemplate{
		name:        "Cockpit",
		description: "A small module with a pilot's chair and the ship's controls.",
		roomType:    ROOM_COCKPIT,
		width:       4,
		height:      3,
		stats: []RoomStat{
			RoomStat{
				Stat:     STAT_POWER_DRAW,
				Modifier: 7,
			},
		},
	}

	roomTemplates[ROOM_CORRIDOR] = RoomTemplate{
		name:        "Corridor",
		description: "A place to walk and have Sorkin-esque conversations.",
		roomType:    ROOM_CORRIDOR,
		width:       12,
		height:      4,
		resizable:   true,
		min_width:   6,
		max_width:   15,
		min_height:  3,
		max_height:  5,
	}

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
			RoomStat{
				Stat:     STAT_POWER_GEN,
				Modifier: 10,
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
			RoomStat{
				Stat:     STAT_POWER_GEN,
				Modifier: 20,
			},
		},
	}

	roomTemplates[ROOM_ENGINE_LARGE] = RoomTemplate{
		name:        "Engine Room - Large",
		description: "A large engine room, with both sub-light and FTL engines",
		roomType:    ROOM_ENGINE_LARGE,
		width:       7,
		height:      9,
		stats: []RoomStat{
			RoomStat{
				Stat:     STAT_POWER_GEN,
				Modifier: 25,
			},
			RoomStat{
				Stat:     STAT_SUBLIGHT_THRUST,
				Modifier: 40,
			},
			RoomStat{
				Stat:     STAT_SUBLIGHT_FUELUSE,
				Modifier: 8,
			},
			RoomStat{
				Stat:     STAT_SUBLIGHT_POWER,
				Modifier: 60,
			},
			RoomStat{
				Stat:     STAT_FTL_THRUST,
				Modifier: 500,
			},
			RoomStat{
				Stat:     STAT_FTL_FUELUSE,
				Modifier: 10,
			},
			RoomStat{
				Stat:     STAT_FTL_POWER,
				Modifier: 150,
			},
		},
	}

	roomTemplates[ROOM_CARGOBAY] = RoomTemplate{
		name:        "Cargo Bay",
		description: "A module for storing large things, or many medium things, or lots and lots of small things.",
		roomType:    ROOM_CARGOBAY,
		width:       8,
		height:      8,
		stats: []RoomStat{
			RoomStat{
				Stat:     STAT_POWER_DRAW,
				Modifier: 2,
			},
		},
	}

	roomTemplates[ROOM_QUARTERS] = RoomTemplate{
		name:        "Crew Quarters",
		description: "A place for the crew to sleep and do other personal things ;)",
		roomType:    ROOM_QUARTERS,
		width:       6,
		height:      6,
		stats: []RoomStat{
			RoomStat{
				Stat:     STAT_POWER_DRAW,
				Modifier: 4,
			},
		},
	}

	roomTemplates[ROOM_COMMONAREA] = RoomTemplate{
		name:        "Common Area",
		description: "A central area for crew to congregate, work, and relax",
		roomType:    ROOM_COMMONAREA,
		width:       9,
		height:      7,
		stats: []RoomStat{
			RoomStat{
				Stat:     STAT_POWER_DRAW,
				Modifier: 5,
			},
		},
	}

	roomTemplates[ROOM_LABORATORY] = RoomTemplate{
		name:        "Laboratory",
		description: "A module filled to the brim with beakers and titration tubes and microscopes and... small probes?",
		roomType:    ROOM_LABORATORY,
		width:       7,
		height:      7,
		stats: []RoomStat{
			RoomStat{
				Stat:     STAT_POWER_DRAW,
				Modifier: 8,
			},
		},
	}

	roomTemplates[ROOM_MEDBAY] = RoomTemplate{
		name:        "Medbay",
		description: "A room for doctors and their professional kin to heal the wounded.",
		roomType:    ROOM_MEDBAY,
		width:       5,
		height:      5,
		stats: []RoomStat{
			RoomStat{
				Stat:     STAT_POWER_DRAW,
				Modifier: 8,
			},
		},
	}

	roomTemplates[ROOM_LASERTURRET] = RoomTemplate{
		name:        "Laser Turret",
		description: "A bona-fide laser gun, for blowing up things with a laser.",
		roomType:    ROOM_LASERTURRET,
		width:       3,
		height:      3,
		stats: []RoomStat{
			RoomStat{
				Stat:     STAT_POWER_DRAW,
				Modifier: 1,
			},
		},
	}
}

//Instantiates a room from a template. "dims" (width, height) is optional, allowing
//the caller to provide custom dimensions. If either dimension is zero, uses the default.
func CreateRoomFromTemplate(room RoomType, rotate bool, dims ...int) (r *Room) {
	r = new(Room)

	if temp, ok := roomTemplates[room]; ok {
		r = NewRoomFromTemplate(temp)
		if temp.resizable && len(dims) == 2 {
			if dims[0] == 0 {
				r.Width = temp.width
			} else {
				r.Width = burl.Clamp(dims[0], temp.min_width, temp.max_width)
			}

			if dims[1] == 0 {
				r.Height = temp.height
			} else {
				r.Height = burl.Clamp(dims[1], temp.min_height, temp.max_height)
			}
			r.CreateRoomMap()
		}

		if rotate {
			r.Rotate()
		}
	}

	return
}
