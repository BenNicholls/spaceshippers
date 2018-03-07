package main

//Similarly to the shiptemplates.go file, this file is just going to hold
//room data until I can be bothered to write a ststem that imports all of this
//from raws.

type RoomClass int //general type of room. determines which system(s) it is keyed to
type RoomType int  //specific variety of room

const (
	ROOMCLASS_CORE        RoomClass = iota //heart of the ship, home of skippie, main controls, LIMIT OF 1
	ROOMCLASS_ENGINEERING                  //source of propulsion
	ROOMCLASS_PERSONAL                     //living quarters, barracks, etc.
	ROOMCLASS_STORAGE                      //cargo bays, fuel pods, etc.
	ROOMCLASS_WEAPON                       //laser guns, turrets, phasers, uhhh.. torpedos? whatever
	ROOMCLASS_TECH                         //Specialty rooms like labs, medical bay, scanners, etc.
	ROOMCLASS_STRUCTURAL                   //Non-functional rooms. Hallways, accessways, elevators maybe)
)

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

func CreateRoomFromTemplate(room RoomType) (r *Room) {
	r = new(Room)

	switch room {
	case ROOM_BRIDGE:
		r = NewRoom("Bridge", ROOM_BRIDGE, ROOMCLASS_CORE, 6, 12, 500, 1000)
	case ROOM_COCKPIT:
		r = NewRoom("Cockpit", ROOM_COCKPIT, ROOMCLASS_CORE, 4, 3, 500, 1000)
	case ROOM_CORRIDOR_H:
		r = NewRoom("Corridor", ROOM_CORRIDOR_H, ROOMCLASS_STRUCTURAL, 12, 4, 500, 1000)
	case ROOM_CORRIDOR_V:
		r = NewRoom("Corridor", ROOM_CORRIDOR_V, ROOMCLASS_STRUCTURAL, 4, 12, 500, 1000)
	case ROOM_ENGINE_SMALL:
		r = NewRoom("Engineering", ROOM_ENGINE_SMALL, ROOMCLASS_ENGINEERING, 3, 5, 500, 1000)
	case ROOM_ENGINE_MEDIUM:
		r = NewRoom("Engineering", ROOM_ENGINE_MEDIUM, ROOMCLASS_ENGINEERING, 5, 8, 500, 1000)
	case ROOM_ENGINE_LARGE:
		r = NewRoom("Engineering", ROOM_ENGINE_LARGE, ROOMCLASS_ENGINEERING, 7, 9, 500, 1000)
	case ROOM_CARGOBAY:
		r = NewRoom("Cargo Bay", ROOM_CARGOBAY, ROOMCLASS_STORAGE, 8, 8, 500, 1000)
	case ROOM_QUARTERS:
		r = NewRoom("Quarters", ROOM_QUARTERS, ROOMCLASS_PERSONAL, 6, 6, 500, 1000)
	case ROOM_COMMONAREA:
		r = NewRoom("Common Area", ROOM_COMMONAREA, ROOMCLASS_PERSONAL, 9, 7, 500, 1000)
	case ROOM_LABORATORY:
		r = NewRoom("Laboratory", ROOM_LABORATORY, ROOMCLASS_TECH, 7, 7, 500, 1000)
	case ROOM_MEDBAY:
		r = NewRoom("Medical Bay", ROOM_MEDBAY, ROOMCLASS_TECH, 5, 5, 500, 1000)
	case ROOM_LASERTURRET:
		r = NewRoom("Laser Turret", ROOM_LASERTURRET, ROOMCLASS_WEAPON, 3, 3, 500, 1000)
	}

	return
}
