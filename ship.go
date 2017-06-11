package main

import "github.com/bennicholls/burl/core"
import "github.com/bennicholls/burl/util"
import "math/rand"

type Ship struct {
	Location //ship is technically a location, but you can't go there... you *are* there!

	Crew  []*Crewman
	Rooms []*Room

	Engine *PropulsionSystem

	//status numbers.
	Hull core.Stat
	Fuel core.Stat

	Heading util.Vec2Polar
	Course  util.Vec2Polar

	ShipMap *core.TileMap

	X, Y, Width, Height int //bounding box of the ship on the shipMap
	Volume              int //total floorspace volume of the ship

	CurrentLocation Locatable //where we at
	Destination     Locatable //where we're going
}

//Inits a new Ship. For now, starts with a bridge and 6 crew.
func NewShip(n string) *Ship {
	s := new(Ship)
	s.ShipMap = core.NewMap(100, 100)
	s.Crew = make([]*Crewman, 6)
	s.Rooms = make([]*Room, 0, 10)
	s.Engine = NewPropulsionSystem(s)

	s.Fuel = core.NewStat(1000000)

	s.locationType = loc_SHIP
	s.name = n
	s.coords = NewCoordinate(coord_LOCAL)
	s.SetExplored()
	s.SetKnown()

	s.AddRoom(NewRoom("Bridge", 20, 6, 6, 12, 500, 1000))

	for i, _ := range s.Crew {
		s.Crew[i] = NewCrewman()
	}

	s.PlaceCrew()

	return s
}

func (s *Ship) SetLocation(l Locatable) {
	s.CurrentLocation = l
	s.coords = l.GetCoords()
	s.coords.resolution = coord_LOCAL
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

func (s *Ship) SetCourse(l Locatable) {
	s.Destination = l

	if l.GetCoords().resolution == coord_LOCAL {
		g := s.coords.CalcVector(l.GetCoords())
		s.Course = g.local.ToVector().ToPolar()
		s.Engine.Firing = true
		s.Engine.Braking = false
	}
}

func (s Ship) GetSpeed() int {
	return int(s.Heading.R)
}

func (s *Ship) Update(spaceTime int) {

	//update course since we don't calculate it properly from the start since calculus is hard
	if s.Engine.Firing && spaceTime%1 == 0 {
		g := s.coords.CalcVector(s.Destination.GetCoords())
		s.Course = g.local.ToVector().ToPolar()

		//calculate decel curve, start decelerating if we're close enough
		if !s.Engine.Braking {
			t := (s.Heading.R - float64(s.Destination.GetVisitSpeed())) / float64(s.Engine.Thrust)
			decelDistance := s.Heading.R*t - float64(s.Engine.Thrust)*t*t/2
			if g.local.Mag()-s.Destination.GetVisitDistance() < int(decelDistance) {
				s.Engine.Braking = true
			}
		}
	}

	s.Engine.Update()

	x, y := s.Heading.ToRect().GetInt()

	s.coords.Move(x, y, coord_LOCAL)

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
