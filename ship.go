package main

import "strconv"

type Ship struct {
    name string

    Crew []*Crewman
    Rooms []*Room
    
    //status numbers.
    Hull Stat
}

func NewShip(n string) *Ship {
    s := new(Ship)
    s.Crew = make([]*Crewman, 6)
    s.Rooms = make([]*Room, 5)

    s.Rooms[BRIDGE] = &Room{"Bridge", 100, 1, 1000}
    s.Rooms[ENGINEERING] = &Room{"Engineering", 100, 1, 1000}
    s.Rooms[MESSHALL] = &Room{"Messhall", 100, 1, 500}
    s.Rooms[MEDBAY] = &Room{"Medbay", 100 , 1, 700}
    s.Rooms[QUARTERS] = &Room{"Quarters", 100, 1, 500}

    for i, _ := range s.Crew {
        s.Crew[i] = NewCrewman()
    }

    s.name = n

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

//room indexes
const (
    BRIDGE int = iota
    ENGINEERING
    MESSHALL
    MEDBAY
    QUARTERS
    MAX_ROOMS
)

type Room struct {
    name string

    state int //state of repair. 
    upkeep int //periodic decay of repair state.
    repairDifficulty int //default time to repair by 1 unit.
}

func (r Room) PrintStatus() {
    roomstatus := r.name + ": Status "
    if r.state > 80 {
        roomstatus += "NOMINAL."
    } else if r.state > 50 {
        roomstatus += "FINE."
    } else if r.state > 20 {
        roomstatus += "NEEDS REPAIR."
    } else if r.state > 0 {
        roomstatus += "CRITICAL."
    } else {
        roomstatus += "DESTROYED."
    }

    roomstatus += " (" + strconv.Itoa(r.state) + "/100)"

    output.Append(roomstatus)
}

func (r *Room) ApplyUpkeep() {
    r.state -= r.upkeep
}

func (r *Room) Update() {

}