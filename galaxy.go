package main

import "math"
import "math/rand"
import "github.com/bennicholls/burl-E/burl"

const (
	//density parameters
	GAL_DENSE  int = 100
	GAL_NORMAL int = 70
	GAL_SPARSE int = 50

	GAL_MIN_RADIUS int = 5
	GAL_MAX_RADIUS int = 12
)

//oh the places you'll go...
type Galaxy struct {
	name          string //all good galaxies have names
	width, height int
	spaceTime     int //time since beginning of the Digital Era
	radius        int //radius in sectors

	sectors []*Sector //galactic map data goes in here! potential for crazy hugeness in here

	earth Coordinates //coordinates of EARTH!!
}

func NewGalaxy(name string, radius, densityFactor int) (g *Galaxy) {
	g = new(Galaxy)
	g.name = name
	g.width, g.height = coord_SECTOR_MAX, coord_SECTOR_MAX
	g.radius = radius

	g.sectors = make([]*Sector, 0, g.width*g.height)
	g.spaceTime = 50*CYCLE + 80*DAY + 8*HOUR //start time for the game. super arbitrary.

	for i := 0; i < cap(g.sectors); i++ {
		x, y := i%g.width, i/g.width
		dist := math.Sqrt(float64(burl.Distance(12, 12, x, y))) + rand.Float64()*2
		density := burl.Clamp(densityFactor-int(float64(densityFactor)*dist/float64(g.radius)), 0, densityFactor)
		g.sectors = append(g.sectors, NewSector(x, y, density))
	}

	//generate the sol system and save the coordinates of Earth
	ss := g.GetSector(8, 8).GenerateSubSector(250, 171)
	ss.starSystem = NewStarSystem(ss.GetCoords())
	g.earth = ss.starSystem.Planets[2].GetCoords()

	return
}

func (g Galaxy) Dims() (int, int) {
	return g.width, g.height
}

//Retreives sector at (x, y). Returns nil if x,y out of bounds (bad).
func (g Galaxy) GetSector(x, y int) *Sector {
	if !burl.CheckBounds(x, y, coord_SECTOR_MAX, coord_SECTOR_MAX) {
		return nil
	}
	return g.sectors[y*g.width+x]
}

func (g Galaxy) GetLocation(c Coordinates) Locatable {
	sector := g.GetSector(c.Sector.Get())
	if c.Resolution == coord_SECTOR {
		return sector
	}

	subsector := sector.GetSubSector(c.SubSector.Get())
	if c.Resolution == coord_SUBSECTOR {
		return subsector
	}

	if !subsector.HasStar() {
		return subsector
	} else {
		if star := subsector.starSystem; star.Coords.StarCoord == c.StarCoord {
			for _, l := range star.GetLocations() {
				if c.IsIn(l) {
					return l
				}
			}
			return star
		} else {
			return subsector
		}
	}
}

func (g Galaxy) GetEarth() Locatable {
	return g.GetLocation(g.earth)
}

type Sector struct {
	Location
	Explored bool
	Density  int //0-100, determines propensity of stars

	subSectors map[int]*SubSector
}

func NewSector(x, y, density int) (s *Sector) {
	s = new(Sector)
	name := ""
	if density == 0 {
		name = "Non-Galactic Space"
	} else if density < 10 {
		name = "The Void Zone"
	} else if density < 30 {
		name = "Outworlder Space"
	} else if density < 75 {
		name = "Main Space Zone Area"
	} else {
		name = "Galactic Core Space"
	}
	s.Location = Location{name, "Sectors are 1000x1000 lightyears! Wow!", loc_SECTOR, false, true, NewSectorCoordinate(x, y), 0, 0}
	s.Density = burl.Max(density, 0) //ensures density is at least 0

	s.subSectors = make(map[int]*SubSector)

	return
}

//generates the name of the sector based on its (x, y).
func (s Sector) ProperName() string {
	x, y := s.Coords.GetCoordStrings()
	return x + "-" + y
}

//GetSubSector attempts to retreive a subsector. If none exists, returns nil.
func (s Sector) GetSubSector(x, y int) *SubSector {
	if s, ok := s.subSectors[x+y*coord_SUBSECTOR_MAX]; ok {
		return s
	} else {
		return nil
	}
}

//generates a subsector and adds it to the subsector map. if (x, y) already exists, just returns the old one
func (s *Sector) GenerateSubSector(x, y int) *SubSector {
	if s, ok := s.subSectors[x+y*coord_SUBSECTOR_MAX]; ok {
		return s
	}

	ss := new(SubSector)
	ss.Coords = s.Coords
	ss.Coords.Resolution = coord_SUBSECTOR
	ss.Coords.SubSector.MoveTo(x, y)

	//PUT STAR GENERATION CODE HERE WHY DON'T YOU.

	s.subSectors[x+y*coord_SUBSECTOR_MAX] = ss
	return ss
}

type SubSector struct {
	Location
	starSystem *StarSystem
}

func (s SubSector) HasStar() bool {
	if s.starSystem != nil {
		return true
	}
	return false
}
