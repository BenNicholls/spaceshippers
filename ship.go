package main

import (
	"slices"

	"github.com/bennicholls/burl-E/burl"
	"github.com/bennicholls/tyumi/gfx"
	"github.com/bennicholls/tyumi/gfx/col"
	"github.com/bennicholls/tyumi/log"
	"github.com/bennicholls/tyumi/rl"
	"github.com/bennicholls/tyumi/vec"
)

type Ship struct {
	Location    //ship is technically a location, but you can't go there... you *are* there!
	gfx.Visuals //for drawing on the galaxy map

	Crew  []*Crewman
	Rooms []*Room

	Engine      *PropulsionSystem
	Navigation  *NavigationSystem
	Comms       *CommSystem
	LifeSupport *LifeSupportSystem
	Storage     *StorageSystem
	//Power *PowerSystem
	//Computer *ComputerSystem
	//Weapons *WeaponSystem
	//Shields *ShieldSystem

	Systems map[SystemType]ShipSystem

	//status numbers.
	Hull burl.Stat

	Velocity vec.Vec2Polar

	shipMap rl.TileMap

	x, y, width, height int //bounding box of the ship on the shipMap
	volume              int //total floorspace volume of the ship

	currentLocation Locatable //where we at
	destination     Locatable //where we're going
}

// Inits a new Ship. Ship systems but no modules or crew. Add those yourself you lazy bum.
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

	s.Systems = make(map[SystemType]ShipSystem)

	//systems
	s.Engine = NewPropulsionSystem(s)
	s.Systems[SYS_PROPULSION] = s.Engine
	s.Navigation = NewNavigationSystem(s, g)
	s.Systems[SYS_NAVIGATION] = s.Navigation
	s.Comms = NewCommSystem()
	s.Systems[SYS_COMMS] = s.Comms
	s.LifeSupport = NewLifeSupportSystem(s)
	s.Systems[SYS_LIFESUPPORT] = s.LifeSupport
	s.Storage = NewStorageSystem(s)
	s.Systems[SYS_STORAGE] = s.Storage

	s.Hull = burl.NewStat(100)

	s.shipMap.Init(vec.Dims{100, 100}, rl.TILE_NONE)

	s.Visuals = gfx.Visuals{
		Glyph:   gfx.GLYPH_FACE2,
		Colours: col.Pair{col.WHITE, col.BLACK},
	}

	return s
}

// SetupShip performs post-init calculations. Do after loading.
func (s *Ship) SetupShip(g *Galaxy) {
	s.shipMap.Init(vec.Dims{100, 100}, rl.TILE_NONE)

	//rooms -- need to process room connections and add them to shipmap
	for i := range s.Rooms {
		for c := i + 1; c < len(s.Rooms); c++ {
			s.ConnectRooms(s.Rooms[i], s.Rooms[c])
		}
		s.DrawRoom(s.Rooms[i])
	}
	s.CalcShipBounds()

	//crew -- need to set the crew's jobs to point back at the crew, add crew to shipmap
	// for i := range s.Crew {
	// 	if s.Crew[i].CurrentTask != nil {
	// 		s.Crew[i].CurrentTask.SetWorker(s.Crew[i])
	// 	}
	// 	rx, ry := s.Crew[i].X, s.Crew[i].Y
	// 	s.shipMap.AddEntity(rx, ry, s.Crew[i])
	// 	s.Crew[i].MoveTo(rx, ry)
	// }

	s.Engine.ship = s
	s.Navigation.ship = s
	s.Navigation.galaxy = g
}

func (s Ship) Bounds() vec.Rect {
	return vec.Rect{vec.Coord{s.x, s.y}, vec.Dims{s.width, s.height}}
}

func (s *Ship) CompileStats() {
	//remove stats from ship systems, we're starting fresh
	for i := range s.Systems {
		s.Systems[i].InitStats()
	}

	for _, room := range s.Rooms {
		for _, stat := range room.Stats {
			if _, ok := s.Systems[stat.GetSystem()]; ok {
				s.Systems[stat.GetSystem()].AddStat(stat)
			}
		}
	}

	for i := range s.Systems {
		s.Systems[i].OnStatUpdate()
	}
}

func (s *Ship) SetLocation(l Locatable) {
	s.currentLocation = l
	s.Coords = l.GetCoords()
	s.Coords.Resolution = coord_LOCAL
}

// Adds a room to the ship and connects it. If room is an invalid add
// (ex. overlaps too much with an existing room), does nothing.
func (s *Ship) AddRoom(pos vec.Coord, r *Room) {
	r.pos = pos

	if !s.CheckRoomValidAdd(r) {
		log.Error("Invalid room add attempt: " + r.Name)
		return
	}

	s.Rooms = append(s.Rooms, r)

	//attempt to connect to each current room
	for _, room := range s.Rooms {
		s.ConnectRooms(room, r)
	}

	s.DrawRoom(r)
	s.CalcShipBounds()
	s.CompileStats()

}

func (s *Ship) RemoveRoom(r *Room) {
	//find the room in the ship's roomlist
	roomIndex := slices.Index(s.Rooms, r)
	if roomIndex == -1 {
		return
	}

	s.Rooms = slices.Delete(s.Rooms, roomIndex, roomIndex+1)
	//erase room from shipmap
	for cursor := range vec.EachCoordInArea(r) {
		s.shipMap.SetTileType(cursor, rl.TILE_NONE)
	}

	//remove connections and re-draw
	for _, connected := range r.connected {
		connected.RemoveConnection(r)
		s.DrawRoom(connected)
	}

	s.CalcShipBounds()
	s.CompileStats()
}

// Checks to see if the provided room collides illegally with another
// in the ship. If there is no collision at all, still reports true
func (s *Ship) CheckRoomValidAdd(r *Room) bool {
	for _, room := range s.Rooms {
		i := vec.FindIntersectionRect(r, room)
		if i.W >= 2 && i.H >= 2 {
			return false
		}
	}

	return true
}

func (s *Ship) AddCrewman(c *Crewman) {
	// s.Crew = append(s.Crew, c)
	// c.ship = s

	// //place randomly in ship
	// for {
	// 	start := s.Rooms[rand.Intn(len(s.Rooms))]
	// 	rx, ry := burl.GenerateCoord(start.X, start.Y, start.Width, start.Height)
	// 	if s.shipMap.GetTile(rx, ry).Empty() {
	// 		s.shipMap.AddEntity(rx, ry, c)
	// 		c.MoveTo(rx, ry)
	// 		break
	// 	}
	// }
}

func (s *Ship) ConnectRooms(r1, r2 *Room) {
	r1.AddConnection(r2)
	r2.AddConnection(r1)
}

// Draws a room onto the shipmap
func (s *Ship) DrawRoom(r *Room) {
	r.RoomMap.CopyToTileMap(&s.shipMap, r.pos)
}

// Calculates the bounding box for the current ship configuration, as well as the volume.
func (s *Ship) CalcShipBounds() {
	if len(s.Rooms) == 0 {
		// if no rooms, just pretend the ship is a 0-area dot in the middle of the ship map.
		s.width, s.height = 0, 0
		s.x, s.y = s.width/2, s.height/2
		s.volume = 0
		return
	}

	s.x, s.y = s.shipMap.Size().W, s.shipMap.Size().H
	x2, y2 := 0, 0
	s.volume = 0

	for _, r := range s.Rooms {
		b := r.Bounds()
		s.x = min(s.x, b.X)
		x2 = max(x2, b.X+b.W)
		s.y = min(s.y, b.Y)
		y2 = max(y2, b.Y+b.H)
		s.volume += (b.W - 2) * (b.H - 2)
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
	for sys := range s.Systems {
		s.Systems[sys].Update(spaceTime)
	}

	for i := range s.Rooms {
		s.Rooms[i].Update(spaceTime)
	}

	for i := range s.Crew {
		s.Crew[i].Update(spaceTime)
	}
}

// Returns the room for a given coord on the shipmap. Returns nil if no room found.
// If multiple rooms occupy a space (ex. shared wall), returns the first one found.
// TODO: this behaviour could be better. Could return all the valid rooms? ugh.
func (s *Ship) GetRoom(c vec.Coord) *Room {
	for _, room := range s.Rooms {
		if c.IsInside(room) {
			return room
		}
	}

	return nil
}

// draws the ship to the provided TileView UI object, offset by (offX, offY). mode is the VIEWMODE
func (s *Ship) DrawToTileView(view *burl.TileView, mode, offX, offY int) {
	// x, y := 0, 0
	// displayWidth, displayHeight := view.Dims()

	// for i := 0; i < s.width*s.height; i++ {
	// 	//tileView-space coords
	// 	x = i%s.width + s.x - offX
	// 	y = i/s.width + s.y - offY

	// 	if burl.CheckBounds(x, y, displayWidth, displayHeight) {
	// 		if t := s.shipMap.GetTile(i%s.width+s.x, i/s.width+s.y); t.TileType != 0 {
	// 			tv := t.GetVisuals()

	// 			if t.TileType != TILE_WALL && t.TileType != TILE_DOOR {
	// 				r := s.GetRoom(vec.Coord{i%s.width + s.x, i/s.width + s.y})
	// 				switch mode {
	// 				case VIEW_ATMO_PRESSURE:
	// 					tv.BackColour = viewModeData[mode].GetColour(r.atmo.Pressure())
	// 				case VIEW_ATMO_O2:
	// 					tv.BackColour = viewModeData[mode].GetColour(r.atmo.PartialPressure(GAS_O2))
	// 				case VIEW_ATMO_TEMP:
	// 					tv.BackColour = viewModeData[mode].GetColour(r.atmo.Temp)
	// 				case VIEW_ATMO_CO2:
	// 					tv.BackColour = viewModeData[mode].GetColour(r.atmo.PartialPressure(GAS_CO2))
	// 				}
	// 			}

	// 			view.DrawObject(x, y, tv)
	// 		}

	// 		if e := s.shipMap.GetEntity(i%s.width+s.x, i/s.width+s.y); e != nil {
	// 			view.DrawObject(x, y, e)
	// 		}
	// 	}
	// }
}
