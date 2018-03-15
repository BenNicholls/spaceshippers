package main

import (
	"math"
	"math/rand"
	"sort"
	"strconv"

	"github.com/bennicholls/burl-E/burl"
)

const GRAVCONST float64 = 6.674e-11

type StarSystem struct {
	Location

	Star    Star
	Planets []Planet
}

func NewStarSystem(c Coordinates) (s *StarSystem) {
	s = new(StarSystem)
	s.LocationType = loc_STARSYSTEM
	s.Coords = c
	s.Coords.Resolution = coord_STARSYSTEM
	s.Coords.StarCoord = burl.Coord{X: rand.Intn(coord_STARSYSTEM_MAX), Y: rand.Intn(coord_STARSYSTEM_MAX)}
	s.Planets = make([]Planet, 0, 10) //Starsystems have max 10 planets, right? Yeah sounds about right.

	return
}

func (s *StarSystem) GenerateRandom() {
	s.Name = "The Glurmglormp System"                  //TODO: random system name generator
	s.Description = "This system is extremely random." //TODO ditto

	//insert star type randomization code here. gotta re-learn my star types!
	//for now, all stars are effectively just The Sun
	s.Star = NewStar(s.Coords, "Glurmglormp", 695700e3, 1.988435e30)

	numPlanets := 1 + rand.Intn(9) //between 1-10 planets
	s.Planets = make([]Planet, 0, numPlanets)
	orbits := make([]float64, numPlanets)
	for i := range orbits {
		orbits[i] = 50e9 + rand.Float64()*4950e9
	}
	sort.Float64s(orbits)

	for i := 0; i < numPlanets; i++ {

		var pType PlanetType
		switch rand.Intn(5) {
		case 0, 1:
			pType = PLANET_ROCKY
		case 2, 3:
			pType = PLANET_GASGIANT
		case 4:
			pType = PLANET_DWARF
		}

		p := NewPlanet(s.Coords, orbits[i], s.Star.Name+" "+strconv.Itoa(i), pType)
		s.Planets = append(s.Planets, p)
	}

}

func (s *StarSystem) GenerateSolSystem() {
	s.Name = "The Sol System"
	s.Description = "Almost certainly the best starsystem in the galaxy. No other system contains more Burger Kings or Weird Al CDs."

	s.Star = NewStar(s.Coords, "The Sun", 695700e3, 1.988435e30)

	s.Planets = make([]Planet, 0, 8)
	s.Planets = append(s.Planets, NewUniquePlanet(s.Coords, 57.3e9, 2493e3, 3.301e23, "Mercury", PLANET_ROCKY))
	s.Planets = append(s.Planets, NewUniquePlanet(s.Coords, 108.2e9, 6051e3, 4.867e24, "Venus", PLANET_ROCKY))
	s.Planets = append(s.Planets, NewUniquePlanet(s.Coords, 149.6e9, 6378e3, 5.972e24, "Earth", PLANET_ROCKY))
	s.Planets = append(s.Planets, NewUniquePlanet(s.Coords, 227.9e9, 3390e3, 6.417e23, "Mars", PLANET_ROCKY))
	s.Planets = append(s.Planets, NewUniquePlanet(s.Coords, 778.3e9, 71492e3, 1.898e27, "Jupiter", PLANET_GASGIANT))
	s.Planets = append(s.Planets, NewUniquePlanet(s.Coords, 1427e9, 60268e3, 5.68319e26, "Saturn", PLANET_GASGIANT))
	s.Planets = append(s.Planets, NewUniquePlanet(s.Coords, 2871e9, 25559e3, 8.681e25, "Uranus", PLANET_GASGIANT))
	s.Planets = append(s.Planets, NewUniquePlanet(s.Coords, 4497e9, 24764e3, 1.024e26, "Neptune", PLANET_GASGIANT))
}

//Returns a list of all locations in the system. Right now that means the star and planets,
//but it's set up to automatically grab locations in a hierarchy, so if planets get moons
//at some point it will pick them up.
func (s StarSystem) GetLocations() []Locatable {
	l := make([]Locatable, 0)
	l = append(l, s.Star)
	for _, p := range s.Planets {
		l = append(l, p.GetLocations()...)
	}
	return l
}

type Star struct {
	Location
	radius float64
	mass   float64
}

//NewStar creates a star. c is the coordinates of the starsystem. Defaults to center of system.
func NewStar(c Coordinates, name string, radius, mass float64) (s Star) {
	s.Name = name
	s.Description = "This is a star. Stars are big hot balls of lava that float in space like magic."
	s.Coords = c
	s.Coords.Resolution = coord_LOCAL
	s.Coords.Local.Set(coord_LOCAL_MAX/2, coord_LOCAL_MAX/2)
	s.radius = radius
	s.mass = mass
	s.VisitDistance = radius * 1.2
	s.VisitSpeed = math.Sqrt(GRAVCONST * s.mass / s.VisitDistance)

	return
}

type PlanetType int

const (
	PLANET_ROCKY PlanetType = iota
	PLANET_GASGIANT
	PLANET_DWARF
)

type Planet struct {
	Location
	oDistance float64 //orbital distance
	oPosition float64 //orbital position in degrees
	radius    float64 //radius of planet in meters
	mass      float64
	ptype     PlanetType
}

//NewPlanet creates a planet. c is the coords of the starsystem. orbit is the distance in meters from the star
func NewPlanet(c Coordinates, orbit float64, name string, pType PlanetType) (p Planet) {
	p.Name = name
	p.Description = "This is a planet. Planets are rocks that are big enough to be important. Some planets have life on them, but most of them are super boring."
	p.LocationType = loc_PLANET
	p.Coords = c
	p.Coords.Resolution = coord_LOCAL

	p.oDistance = orbit
	p.oPosition = rand.Float64() * 2 * math.Pi
	p.Coords.Local.X = (coord_LOCAL_MAX / 2) + p.oDistance*math.Cos(p.oPosition)
	p.Coords.Local.Y = (coord_LOCAL_MAX / 2) + p.oDistance*math.Sin(p.oPosition)

	p.Generate(pType)

	return
}

func NewUniquePlanet(c Coordinates, orbit, radius, mass float64, name string, pType PlanetType) (p Planet) {
	p = NewPlanet(c, orbit, name, pType)
	p.radius = radius
	p.mass = mass

	p.VisitDistance = p.radius * 1.2
	p.VisitSpeed = math.Sqrt(GRAVCONST * p.mass / p.VisitDistance)

	return
}

func (p *Planet) Generate(t PlanetType) {
	var min_r, max_r float64
	var min_d, max_d float64

	switch t {
	case PLANET_ROCKY:
		min_r, max_r = 2000e3, 10000e3
		min_d, max_d = 3000e3, 6000e3
	case PLANET_GASGIANT:
		min_r, max_r = 20000e3, 100000e3
		min_d, max_d = 600e3, 2000e3
	case PLANET_DWARF:
		min_r, max_r = 1000e3, 1800e3
		min_d, max_d = 2000e3, 4000e3
	}

	p.radius = min_r + rand.Float64()*(max_r-min_r)
	p.mass = (min_d + rand.Float64()*(max_d-min_d)) * (math.Pi * p.radius * p.radius * 4 / 3)
	p.VisitDistance = p.radius * 1.5
	p.VisitSpeed = math.Sqrt(GRAVCONST * p.mass / p.VisitDistance)
	p.ptype = t

}

//NOTE: When moons and orbitting stuff gets implemented, be sure to add those here.
func (p Planet) GetLocations() []Locatable {
	l := make([]Locatable, 0)
	l = append(l, p)
	return l
}
