package main

import "math"
import "math/rand"
import "github.com/bennicholls/burl/util"

//oh the places you'll go...
type Galaxy struct {
	name string //all good galaxies have names

	width, height int //dimensions of galaxy in sectors.
	sectorSize    int //dimensions of a sector (always square). measured in ly

	sectors []*Sector
}

func NewGalaxy() (g *Galaxy) {
	g = new(Galaxy)
	g.name = "The Galaxy of Terror"
	g.width, g.height = 25, 25
	g.sectorSize = 1000

	g.sectors = make([]*Sector, 0, g.width*g.height)

	for i := 0; i < cap(g.sectors); i++ {
		x, y := i%g.width, i/g.width
		dist := math.Sqrt(float64(util.Distance(12, 12, x, y))) + rand.Float64()*2
		density := util.Clamp(100-int(100.0*dist/11.0), 0, 100)
		g.sectors = append(g.sectors, NewSector(x, y, g.sectorSize, density))
	}

	return
}

func (g Galaxy) Dims() (int, int) {
	return g.width, g.height
}

//Retreives sector at (x, y). Returns nil if x,y out of bounds (bad).
func (g Galaxy) GetSector(x, y int) *Sector {
	if y*g.width+x >= len(g.sectors) {
		return nil
	}
	return g.sectors[y*g.width+x]
}

type Sector struct {
	Location
	size int //size to a side for a sector

	Explored bool
	Density  int //0-100, determines propensity of stars

	Stars []*StarSystem
}

func NewSector(x, y, size, density int) (s *Sector) {
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
	s.size = size
	s.Density = util.Max(density, 0) //ensures density is at least 0

	s.Stars = make([]*StarSystem, 0, 50)

	return
}

//generates the name of the sector based on its (x, y).
func (s Sector) ProperName() string {
	x, y := s.coords.GetCoordStrings()
	return x + "-" + y
}

type StarSystem struct {
	Location
}
