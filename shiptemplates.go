package main

//this file is going to keep ship template data so we can set up a few kinds of ship
//eventually this will all be read from raw text files, but for now i want to fly a
//spaceship damn it so that will have to wait

const (
	SHIPTYPE_CIVILIAN ShipType = iota
	SHIPTYPE_TRANSPORT
	SHIPTYPE_MINING
	SHIPTYPE_FIGHTER
	SHIPTYPE_EXPLORER
	SHIPTYPE_CUSTOM //unused
)

//type of ship
type ShipType int

//we are assuming the ship's defaults are already set, the templates add Rooms,
//Systems and a default number of non-descript Crew that can be modified afterwards
func (s *Ship) SetupFromTemplate(shiptype ShipType) {

	crewNum := 0

	switch shiptype {
	case SHIPTYPE_CIVILIAN:
		crewNum = 4
		s.AddRoom(CreateRoomFromTemplate(ROOM_COCKPIT, false), 34, 29)
		s.AddRoom(CreateRoomFromTemplate(ROOM_COMMONAREA, false), 26, 27)
		s.AddRoom(CreateRoomFromTemplate(ROOM_ENGINE_SMALL, false), 24, 28)
	case SHIPTYPE_TRANSPORT:
		crewNum = 6
		s.AddRoom(CreateRoomFromTemplate(ROOM_BRIDGE, false), 45, 21)
		s.AddRoom(CreateRoomFromTemplate(ROOM_CORRIDOR, false, 16, 4), 34, 29)
		s.AddRoom(CreateRoomFromTemplate(ROOM_CORRIDOR, true), 36, 18)
		s.AddRoom(CreateRoomFromTemplate(ROOM_QUARTERS, false), 39, 18)
		s.AddRoom(CreateRoomFromTemplate(ROOM_QUARTERS, false), 39, 24)
		s.AddRoom(CreateRoomFromTemplate(ROOM_CARGOBAY, false), 29, 20)
		s.AddRoom(CreateRoomFromTemplate(ROOM_ENGINE_LARGE, false), 28, 27)
	case SHIPTYPE_MINING:

	case SHIPTYPE_FIGHTER:

	case SHIPTYPE_EXPLORER:

	}

	s.ShipType = shiptype

	for i := 0; i < crewNum; i++ {
		s.AddCrewman(NewCrewman())
	}
}
