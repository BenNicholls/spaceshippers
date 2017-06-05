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

	X, Y, Width, Height int //bounding box of the ship on the shipMap
	Volume              int //total floorspace volume of the ship

	ShipCoords  Coordinates //actual coordinates on the galactic map
	Location    Locatable   //Current location (planet, star system, sector, whatever)
	Destination Locatable   //where we're going
}

//Inits a new Ship. For now, starts with a bridge and 6 crew.
func NewShip(n string) *Ship {
	s := new(Ship)
	s.ShipMap = core.NewMap(100, 100)
	s.Crew = make([]*Crewman, 6)
	s.Rooms = make([]*Room, 0, 10)
	s.name = n

	s.AddRoom(NewRoom("Bridge", 20, 6, 6, 12, 500, 1000))

	for i, _ := range s.Crew {
		s.Crew[i] = NewCrewman()
	}

	s.ShipCoords = NewCoordinate(coord_LOCAL)

	s.PlaceCrew()

	return s
}

func (s *Ship) SetLocation(l Locatable) {
	s.Location = l
	c := s.Location.GetCoords()
	s.ShipCoords.xSector = c.xSector
	s.ShipCoords.ySector = c.ySector
}

//Adds a room to the ship and connects it.
//TODO: Check if room is a valid add.
func (s *Ship) AddRoom(r *Room) {
	s.Rooms = append(s.Rooms, r)

	//attempt to connect to each current room
	for _, room := range s.Rooms {
		s.ConnectRooms(room, r)
	}

	s.DrawRoom(r)
	s.CalcShipDims()
}

func (s *Ship) ConnectRooms(r1, r2 *Room) {
	r1.AddConnection(r2)
	r2.AddConnection(r1)
}

//Draws a room onto the shipmap
func (s *Ship) DrawRoom(r *Room) {
	for i := 0; i < r.Width*r.Height; i++ {
		s.ShipMap.ChangeTileType(r.X+i%r.Width, r.Y+i/r.Width, r.RoomMap.GetTileType(i%r.Width, i/r.Width))
	}
}

//Calculates the bounding box for the current ship configuration, as well as the volume.
func (s *Ship) CalcShipDims() {
	s.X, s.Y = s.ShipMap.Dims()
	x2, y2 := 0, 0

	for _, r := range s.Rooms {
		s.X = util.Min(s.X, r.X)
		x2 = util.Max(x2, r.X+r.Width)
		s.Y = util.Min(s.Y, r.Y)
		y2 = util.Max(y2, r.Y+r.Height)
		s.Volume += (r.Width - 2) * (r.Height - 2)
	}

	s.Width = x2 - s.X
	s.Height = y2 - s.Y
}

//Inits crew. For now just randomizes their positions.
func (s *Ship) PlaceCrew() {
	for i, _ := range s.Crew {
		start := s.Rooms[rand.Intn(len(s.Rooms))]
		for {
			rx, ry := util.GenerateCoord(start.X, start.Y, start.Width, start.Height)
			if s.ShipMap.GetTile(rx, ry).Empty() {
				s.ShipMap.AddEntity(rx, ry, s.Crew[i])
				s.Crew[i].MoveTo(rx, ry)
				break
			}
		}
	}
}

func (s *Ship) Update(spaceTime int) {

	s.ShipCoords.Move(-1000, 1000, coord_LOCAL)

	for i, _ := range s.Rooms {
		s.Rooms[i].Update(spaceTime)
	}

	for i, _ := range s.Crew {
		s.Crew[i].Update()
		if spaceTime%20 == 0 && s.Crew[i].IsAwake() {
			dx, dy := util.RandomDirection()
			if s.ShipMap.GetTile(s.Crew[i].X+dx, s.Crew[i].Y+dy).Empty() {
				s.ShipMap.MoveEntity(s.Crew[i].X, s.Crew[i].Y, dx, dy)
				s.Crew[i].Move(dx, dy)
			}
		}
	}
}
