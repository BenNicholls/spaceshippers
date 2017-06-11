package main

import "math/rand"
import "math"

const GRAVCONST float64 = 6.674e-11

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
	s.Star = NewStar(s.coords, "The Sun", 695700e3, 1.988435e30)
	s.Planets = make([]Planet, 0, 10)

	s.Planets = append(s.Planets, NewPlanet(s.coords, 57.3e9, 2493e3, 3.301e23, "Mercury"))
	s.Planets = append(s.Planets, NewPlanet(s.coords, 108.2e9, 6051e3, 4.867e24, "Venus"))
	s.Planets = append(s.Planets, NewPlanet(s.coords, 149.6e9, 6378e3, 5.972e24, "Earth"))
	s.Planets = append(s.Planets, NewPlanet(s.coords, 227.9e9, 3396e3, 6.417e23, "Mars"))
	s.Planets = append(s.Planets, NewPlanet(s.coords, 778.3e9, 71492e3, 1.898e27, "Jupiter"))
	s.Planets = append(s.Planets, NewPlanet(s.coords, 1427e9, 60268e3, 5.68319e26, "Saturn"))
	s.Planets = append(s.Planets, NewPlanet(s.coords, 2871e9, 25559e3, 8.681e25, "Uranus"))
	s.Planets = append(s.Planets, NewPlanet(s.coords, 4497e9, 24764e3, 1.024e26, "Neptune"))

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
	radius float64
	mass   float64
}

//NewStar creates a star. c is the coordinates of the starsystem. Defaults to center of system.
func NewStar(c Coordinates, name string, radius, mass float64) (s Star) {
	s.name = name
	s.coords = c
	s.coords.resolution = coord_LOCAL
	s.coords.local.MoveTo(coord_LOCAL_MAX/2, coord_LOCAL_MAX/2)
	s.radius = radius
	s.mass = mass
	s.visitDistance = int(radius * 1.2)
	s.visitSpeed = int(math.Sqrt(GRAVCONST * float64(s.mass) / float64(s.visitDistance)))

	return
}

type Planet struct {
	Location
	oDistance float64 //orbital distance
	oPosition float64 //orbital position in degrees
	radius    float64 //radius of planet in meters
	mass      float64
}

//NewPlanet creates a planet. c is the coords of the starsystem. orbit is the distance in meters from the star
func NewPlanet(c Coordinates, orbit, radius, mass float64, name string) (p Planet) {
	p.name = name
	p.locationType = loc_PLANET
	p.coords = c
	p.coords.resolution = coord_LOCAL

	p.oDistance = orbit
	p.oPosition = rand.Float64() * 2 * math.Pi
	p.radius = radius
	p.mass = mass
	p.visitDistance = int(radius * 1.2)
	p.visitSpeed = int(math.Sqrt(GRAVCONST * p.mass / float64(p.visitDistance)))
	p.coords.local.X = (coord_LOCAL_MAX / 2) + int(p.oDistance*math.Cos(p.oPosition))
	p.coords.local.Y = (coord_LOCAL_MAX / 2) + int(p.oDistance*math.Sin(p.oPosition))

	return
}

//NOTE: When moons and orbitting stuff gets implemented, be sure to add those here.
func (p Planet) GetLocations() []Locatable {
	l := make([]Locatable, 0)
	l = append(l, p)
	return l
}
