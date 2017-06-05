package main

import "math"
import "math/rand"
import "github.com/bennicholls/burl/util"

//galaxy creation parameters
const (
	galaxy_TOTALSTARS = 10000000 //ten million stars for now
	galaxy_RADIUS     = 12       //in sectors
)

//oh the places you'll go...
type Galaxy struct {
	name          string //all good galaxies have names
	width, height int
	sectors       []*Sector

	starFactor int //number of stars per density leve for a sector
}

func NewGalaxy() (g *Galaxy) {
	g = new(Galaxy)
	g.name = "The Galaxy of Terror"
	g.width, g.height = coord_SECTOR_MAX, coord_SECTOR_MAX

	g.sectors = make([]*Sector, 0, g.width*g.height)

	cumDens := 0
	nonEmpty := 0
	for i := 0; i < cap(g.sectors); i++ {
		x, y := i%g.width, i/g.width
		dist := math.Sqrt(float64(util.Distance(12, 12, x, y))) + rand.Float64()*2
		density := util.Clamp(100-int(100.0*dist/galaxy_RADIUS), 0, 100)
		g.sectors = append(g.sectors, NewSector(x, y, density))
		cumDens += density
		if density > 0 {
			nonEmpty++
		}
	}

	g.starFactor = galaxy_TOTALSTARS / cumDens

	return
}

func (g Galaxy) Dims() (int, int) {
	return g.width, g.height
}

//Retreives sector at (x, y). Returns nil if x,y out of bounds (bad).
func (g Galaxy) GetSector(x, y int) *Sector {
	if !util.CheckBounds(x, y, coord_SECTOR_MAX, coord_SECTOR_MAX) {
		return nil
	}
	return g.sectors[y*g.width+x]
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
	s.Location = Location{name, loc_SECTOR, false, true, NewSectorCoordinate(x, y)}
	s.Density = util.Max(density, 0) //ensures density is at least 0

	s.subSectors = make(map[int]*SubSector)

	return
}

//generates the name of the sector based on its (x, y).
func (s Sector) ProperName() string {
	x, y := s.coords.GetCoordStrings()
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
	ss.coords = s.coords
	ss.coords.xStarCoord = x
	ss.coords.yStarCoord = y

	//PUT STAR GENERATION CODE HERE WHY DON'T YOU.

	s.subSectors[x+y*coord_SUBSECTOR_MAX] = ss
	return ss
}

type SubSector struct {
	Location
	star *StarSystem
}

func (s SubSector) HasStar() bool {
	if s.star != nil {
		return true
	}
	return false
}
