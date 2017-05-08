package main

import "github.com/bennicholls/burl/core"
import "github.com/bennicholls/burl/util"
import "math/rand"

type Ship struct {
	name string

	Crew  []*Crewman
	Rooms []*Room

	//status numbers.
	Hull  core.Stat
	Pilot *Crewman

	ShipMap *core.TileMap
}

func NewShip(n string) *Ship {
	s := new(Ship)
	s.ShipMap = core.NewMap(100, 100)
	s.Crew = make([]*Crewman, 6)
	s.Rooms = make([]*Room, 7)

	s.Rooms[0] = NewRoom("Bridge", 20, 6, 6, 12, 500, 1000)
	s.Rooms[1] = NewRoom("Engineering", 5, 8, 5, 8, 700, 1000)
	s.Rooms[2] = NewRoom("Messhall", 15, 5, 6, 6, 1000, 500)
	s.Rooms[3] = NewRoom("Medbay", 9, 5, 6, 6, 1000, 700)
	s.Rooms[4] = NewRoom("Quarters 1", 15, 13, 6, 6, 900, 500)
	s.Rooms[5] = NewRoom("Quarters 2", 9, 13, 6, 6, 900, 500)
	s.Rooms[6] = NewRoom("Hallway", 9, 10, 12, 4, 0, 500)

	for _, r := range s.Rooms {
		for i := 0; i < r.W*r.H; i ++ {
			if i < r.W || i%r.W == 0 || i%r.W == r.W - 1 || i/r.W == r.H - 1 {
				s.ShipMap.ChangeTileType(r.X + i%r.W, r.Y + i/r.W, TILE_WALL)	
			} else {
				s.ShipMap.ChangeTileType(r.X + i%r.W, r.Y + i/r.W, TILE_FLOOR)
			}
		}
	}

	s.ConnectRooms(s.Rooms[6], s.Rooms[0])
	s.ConnectRooms(s.Rooms[6], s.Rooms[1])
	s.ConnectRooms(s.Rooms[6], s.Rooms[2])
	s.ConnectRooms(s.Rooms[6], s.Rooms[3])
	s.ConnectRooms(s.Rooms[6], s.Rooms[4])
	s.ConnectRooms(s.Rooms[6], s.Rooms[5])

	for i, _ := range s.Crew {
		s.Crew[i] = NewCrewman()
		start := s.Rooms[rand.Intn(len(s.Rooms))]
		for {
			rx, ry := util.GenerateCoord(start.X, start.Y, start.W, start.H)
			if s.ShipMap.GetTile(rx, ry).Empty() {
				s.ShipMap.AddEntity(rx, ry, s.Crew[i])
				s.Crew[i].MoveTo(rx, ry)
				break
			}
		}
	}

	s.name = n
	s.Pilot = nil

	return s
}

//Finds the intersection of the two rooms and puts doors there!
//TODO: smarter. currently only works if they are set up right in the first place
//TODO: this is a monstrosity.
func (s *Ship) ConnectRooms(r1, r2 *Room) {
	r1.AddConnection(r2)
	r2.AddConnection(r1)

	x, y, l := 0,0,0

	//left/right connection
	if r1.X + r1.W - 1 == r2.X || r2.X + r2.W - 1 == r1.X  {
		if r1.X + r1.W - 1 == r2.X {
			x = r2.X
		} else {
			x = r1.X
		}
		
		y = util.Max(r1.Y, r2.Y)
		l = util.Min(r1.Y + r1.H, r2.Y + r2.H) - y - 1
		
		s.ShipMap.ChangeTileType(x, y+l/2, TILE_DOOR)
		s.ShipMap.ChangeTileType(x, y+l/2 + 1, TILE_DOOR)

	} else if r1.Y + r1.H - 1 == r2.Y || r2.Y + r2.H - 1 == r1.Y  {
		if r1.Y + r1.H - 1 == r2.Y {
			y = r2.Y
		} else {
			y = r1.Y
		}
		
		x = util.Max(r1.X, r2.X)
		l = util.Min(r1.X + r1.W, r2.X + r2.W) - x - 1
		
		s.ShipMap.ChangeTileType(x + l/2, y, TILE_DOOR)
		s.ShipMap.ChangeTileType(x + l/2 + 1, y, TILE_DOOR)

	}
}

func (s *Ship) Update(spaceTime int) {
	for i, _ := range s.Rooms {
		s.Rooms[i].Update(spaceTime)
	}

	for i, _ := range s.Crew {
		s.Crew[i].Update()
		if spaceTime%20 == 0 && s.Crew[i].IsAwake() {
			dx, dy := util.GenerateDirection()
			if s.ShipMap.GetTile(s.Crew[i].X + dx, s.Crew[i].Y + dy).Empty() {
				s.ShipMap.MoveEntity(s.Crew[i].X, s.Crew[i].Y, dx, dy)
				s.Crew[i].Move(dx, dy)
			}
		}
	}
}
