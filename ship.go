package main

type Ship struct {
	name string

	Crew  []*Crewman
	Rooms []*Room

	//status numbers.
	Hull  Stat
	Pilot *Crewman
}

func NewShip(n string) *Ship {
	s := new(Ship)
	s.Crew = make([]*Crewman, 6)
	s.Rooms = make([]*Room, MAX_ROOMS)

	s.Rooms[BRIDGE] = &Room{"Bridge", 38, 16, 3, 3, Stat{100, 100}, 500, 1000}
	s.Rooms[ENGINEERING] = &Room{"Engineering", 30, 15, 2, 5, Stat{100, 100}, 700, 1000}
	s.Rooms[MESSHALL] = &Room{"Messhall", 32, 14, 3, 3, Stat{100, 100}, 1000, 500}
	s.Rooms[MEDBAY] = &Room{"Medbay", 35, 14, 3, 3, Stat{100, 100}, 1000, 700}
	s.Rooms[QUARTERS] = &Room{"Quarters", 32, 18, 6, 3, Stat{100, 100}, 900, 500}
	s.Rooms[HALLWAY] = &Room{"Hallway", 32, 17, 6, 1, Stat{100, 100}, 0, 500}

	for i, _ := range s.Crew {
		s.Crew[i] = NewCrewman()
	}

	s.name = n
	s.Pilot = nil

	return s
}

func (s *Ship) Update() {
	for i, _ := range s.Rooms {
		s.Rooms[i].Update()
	}

	for i, _ := range s.Crew {
		s.Crew[i].Update()
	}
}

func (s Ship) Draw() {
	for _, r := range s.Rooms {
		r.Draw()
	}
}
