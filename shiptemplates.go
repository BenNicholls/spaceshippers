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
		s.AddRoom(CreateRoomFromTemplate(ROOM_COCKPIT), 34, 29)
		s.AddRoom(CreateRoomFromTemplate(ROOM_COMMONAREA), 26, 27)
		s.AddRoom(CreateRoomFromTemplate(ROOM_ENGINE_SMALL), 24, 28)
	case SHIPTYPE_TRANSPORT:
		crewNum = 6
		s.AddRoom(CreateRoomFromTemplate(ROOM_BRIDGE), 45, 21)
		s.AddRoom(CreateRoomFromTemplate(ROOM_CORRIDOR_H), 34, 29)
		s.AddRoom(CreateRoomFromTemplate(ROOM_CORRIDOR_V), 36, 18)
		s.AddRoom(CreateRoomFromTemplate(ROOM_QUARTERS), 39, 18)
		s.AddRoom(CreateRoomFromTemplate(ROOM_QUARTERS), 39, 24)
		s.AddRoom(CreateRoomFromTemplate(ROOM_CARGOBAY), 29, 20)
		s.AddRoom(CreateRoomFromTemplate(ROOM_ENGINE_LARGE), 28, 27)
	case SHIPTYPE_MINING:

	case SHIPTYPE_FIGHTER:

	case SHIPTYPE_EXPLORER:

	}

	s.ShipType = shiptype

	for i := 0; i < crewNum; i++ {
		s.AddCrewman(NewCrewman())
	}
}
