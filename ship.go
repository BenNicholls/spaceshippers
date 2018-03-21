package main

import (
	"math/rand"

	"github.com/bennicholls/burl-E/burl"
)

type Ship struct {
	Location //ship is technically a location, but you can't go there... you *are* there!
	ShipType ShipType

	Crew  []*Crewman
	Rooms []*Room

	Engine     *PropulsionSystem
	Navigation *NavigationSystem
	Comms      *CommSystem

	//status numbers.
	Hull burl.Stat
	Fuel burl.Stat

	Velocity burl.Vec2Polar

	shipMap *burl.TileMap

	x, y, width, height int //bounding box of the ship on the shipMap
	volume              int //total floorspace volume of the ship

	currentLocation Locatable //where we at
	destination     Locatable //where we're going
}

//Inits a new Ship. For now, starts with a bridge and 6 crew.
func NewShip(n string, g *Galaxy) *Ship {
	s := new(Ship)
	s.Name = n
	s.Description = "This is your ship! Look at it's heroic hull valiantly floating amongst the stars. One could almost weep."
	s.LocationType = loc_SHIP
	s.Coords = NewCoordinate(coord_LOCAL)
	s.Explored = true
	s.Known = true

	s.Crew = make([]*Crewman, 0, 50)
	s.Rooms = make([]*Room, 0, 50)

	//systems
	s.Engine = NewPropulsionSystem(s)
	s.Navigation = NewNavigationSystem(s, g)
	s.Comms = NewCommSystem()
	s.Fuel = burl.NewStat(1000000)

	s.shipMap = burl.NewMap(100, 100)

	return s
}

//SetupShip performs post-init calculations. Do after loading.
func (s *Ship) SetupShip(g *Galaxy) {
	s.shipMap = burl.NewMap(100, 100)

	//rooms -- need to process room connections and add them to shipmap
	for i := range s.Rooms {
		for c := i + 1; c < len(s.Rooms); c++ {
			s.ConnectRooms(s.Rooms[i], s.Rooms[c])
		}
		s.DrawRoom(s.Rooms[i])
	}
	s.CalcShipDims()

	//crew -- need to set the crew's jobs to point back at the crew, add crew to shipmap
	for i := range s.Crew {
		if s.Crew[i].CurrentTask != nil {
			s.Crew[i].CurrentTask.SetWorker(s.Crew[i])
		}
		rx, ry := s.Crew[i].X, s.Crew[i].Y
		s.shipMap.AddEntity(rx, ry, s.Crew[i])
		s.Crew[i].MoveTo(rx, ry)
	}

	s.Engine.ship = s
	s.Navigation.ship = s
	s.Navigation.galaxy = g
}

func (s *Ship) SetLocation(l Locatable) {
	s.currentLocation = l
	s.Coords = l.GetCoords()
	s.Coords.Resolution = coord_LOCAL
}

//Adds a room to the ship and connects it. If room is an invalid add
//(ex. overlaps too much with an existing room), does nothing.
func (s *Ship) AddRoom(r *Room, x, y int) {
	r.X = x
	r.Y = y

	if s.CheckRoomValidAdd(r, x, y) {
		s.Rooms = append(s.Rooms, r)

		//attempt to connect to each current room
		for _, room := range s.Rooms {
			s.ConnectRooms(room, r)
		}

		s.DrawRoom(r)
		s.CalcShipDims()
	} else {
		burl.LogError("Invalid room add attempt: " + r.Name)
	}
}

//Checks to see if the provided room collides illegally with another
//in the ship. If there is no collision at all, still reports true
func (s *Ship) CheckRoomValidAdd(r *Room, x, y int) bool {
	for _, room := range s.Rooms {
		_, _, w, h := burl.FindIntersectionRect(r, room)
		if w >= 2 && h >= 2 {
			return false
		}
	}

	return true
}

func (s *Ship) AddCrewman(c *Crewman) {
	s.Crew = append(s.Crew, c)

	//place randomly in ship
	start := s.Rooms[rand.Intn(len(s.Rooms))]
	for {
		rx, ry := burl.GenerateCoord(start.X, start.Y, start.Width, start.Height)
		if s.shipMap.GetTile(rx, ry).Empty() {
			s.shipMap.AddEntity(rx, ry, c)
			c.MoveTo(rx, ry)
			break
		}
	}
}

func (s *Ship) ConnectRooms(r1, r2 *Room) {
	r1.AddConnection(r2)
	r2.AddConnection(r1)
}

//Draws a room onto the shipmap
func (s *Ship) DrawRoom(r *Room) {
	for i := 0; i < r.Width*r.Height; i++ {
		s.shipMap.ChangeTileType(r.X+i%r.Width, r.Y+i/r.Width, r.RoomMap.GetTileType(i%r.Width, i/r.Width))
	}
}

//Calculates the bounding box for the current ship configuration, as well as the volume.
func (s *Ship) CalcShipDims() {
	s.x, s.y = s.shipMap.Dims()
	x2, y2 := 0, 0

	for _, r := range s.Rooms {
		s.x = burl.Min(s.x, r.X)
		x2 = burl.Max(x2, r.X+r.Width)
		s.y = burl.Min(s.y, r.Y)
		y2 = burl.Max(y2, r.Y+r.Height)
		s.volume += (r.Width - 2) * (r.Height - 2)
	}

	s.width = x2 - s.x
	s.height = y2 - s.y
}

func (s *Ship) SetCourse(l Locatable, c Course) {
	s.destination = l
	if l.GetCoords().Resolution == coord_LOCAL {
		s.Navigation.CurrentCourse = c
		s.Engine.Firing = true
	}
}

func (s Ship) GetSpeed() int {
	return int(s.Velocity.R)
}

func (s *Ship) Update(spaceTime int) {
	s.Navigation.Update(spaceTime)
	s.Engine.Update()
	s.Comms.Update(spaceTime)

	x, y := s.Velocity.ToRect().Get()
	s.Coords.moveLocal(x, y)

	for i, _ := range s.Rooms {
		s.Rooms[i].Update(spaceTime)
	}

	for i, _ := range s.Crew {
		s.Crew[i].Update()
		if spaceTime%20 == 0 && s.Crew[i].IsAwake() {
			dx, dy := burl.RandomDirection()
			if s.shipMap.GetTile(s.Crew[i].X+dx, s.Crew[i].Y+dy).Empty() {
				s.shipMap.MoveEntity(s.Crew[i].X, s.Crew[i].Y, dx, dy)
				s.Crew[i].Move(dx, dy)
			}
		}
	}
}

//draws the ship to the provided TileView UI object, offset by (offX, offY)
func (s *Ship) DrawToTileView(view *burl.TileView, offX, offY int) {
	x, y := 0, 0
	displayWidth, displayHeight := view.Dims()

	for i := 0; i < s.width*s.height; i++ {
		//tileView-space coords
		x = i%s.width + s.x - offX
		y = i/s.width + s.y - offY

		if burl.CheckBounds(x, y, displayWidth, displayHeight) {
			if t := s.shipMap.GetTile(i%s.width+s.x, i/s.width+s.y); t.TileType != 0 {
				tv := t.GetVisuals()
				view.Draw(x, y, tv.Glyph, tv.ForeColour, burl.COL_BLACK)
			}

			if e := s.shipMap.GetEntity(i%s.width+s.x, i/s.width+s.y); e != nil {
				view.Draw(x, y, e.GetVisuals().Glyph, e.GetVisuals().ForeColour, burl.COL_BLACK)
			}
		}
	}
}
