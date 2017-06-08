package main

import "math/rand"
import "math"

type StarSystem struct {
	Location

	Star    Star
	Planets []Planet //Starsystems have max 10 planets. Right? Yeah sounds right.
}

func NewStarSystem(c Coordinates) (s *StarSystem) {
	s = new(StarSystem)

	s.name = "The Sol System"
	s.locationType = loc_STARSYSTEM
	s.coords = c
	s.coords.resolution = coord_STARSYSTEM
	s.Star = NewStar(s.coords, "The Sun", 695700e3)
	s.Planets = make([]Planet, 0, 10)

	s.Planets = append(s.Planets, NewPlanet(s.coords, 57.3e9, 2493e3, "Mercury"))
	s.Planets = append(s.Planets, NewPlanet(s.coords, 108.2e9, 6051e3, "Venus"))
	s.Planets = append(s.Planets, NewPlanet(s.coords, 149.6e9, 6378e3, "Earth"))
	s.Planets = append(s.Planets, NewPlanet(s.coords, 227.9e9, 3396e3, "Mars"))
	s.Planets = append(s.Planets, NewPlanet(s.coords, 778.3e9, 71492e3, "Jupiter"))
	s.Planets = append(s.Planets, NewPlanet(s.coords, 1427e9, 60268e3, "Saturn"))
	s.Planets = append(s.Planets, NewPlanet(s.coords, 2871e9, 25559e3, "Uranus"))
	s.Planets = append(s.Planets, NewPlanet(s.coords, 4497e9, 24764e3, "Neptune"))

	return
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
	radius int
	orbitRange int
}

//NewStar creates a star. c is the coordinates of the starsystem. Defaults to center of system.
func NewStar(c Coordinates, name string, radius int) (s Star) {
	s.name = name
	s.coords = c
	s.coords.resolution = coord_LOCAL
	s.coords.xLocal = coord_LOCAL_MAX / 2
	s.coords.yLocal = coord_LOCAL_MAX / 2
	s.radius = radius
	s.orbitRange = radius + 15000e3 //nice "safe" 15000km sub orbit radius.

	return
}

func (s Star) GetOrbitDistance() int {
	return s.orbitRange
}

type Planet struct {
	Location
	oDistance  int     //orbital distance
	oPosition  float64 //orbital position in degrees
	radius     int     //radius of planet in meters
	orbitRange int     //range for standard orbit around planet.
}

//NewPlanet creates a planet. c is the coords of the starsystem. orbit is the distance in meters from the star
func NewPlanet(c Coordinates, orbit, radius int, name string) (p Planet) {
	p.name = name
	p.locationType = loc_PLANET
	p.coords = c

	p.oDistance = orbit
	p.oPosition = rand.Float64() * 2 * math.Pi
	p.orbitRange = radius + 1000e3 //right now: radius + 1000km
	p.coords.xLocal = (coord_LOCAL_MAX / 2) + int(float64(p.oDistance)*math.Cos(p.oPosition))
	p.coords.yLocal = (coord_LOCAL_MAX / 2) + int(float64(p.oDistance)*math.Sin(p.oPosition))

	return
}

//NOTE: When moons and orbitting stuff gets implemented, be sure to add those here.
func (p Planet) GetLocations() []Locatable {
	l := make([]Locatable, 0)
	l = append(l, p)
	return l
}

func (p Planet) GetOrbitDistance() int {
	return p.orbitRange
}

type Orbitable interface {
	GetOrbitDistance() int //orbit distance in meters.
	//AddOrbitObject() //someday this will allow us to attach objects to the things they're orbiting
}
