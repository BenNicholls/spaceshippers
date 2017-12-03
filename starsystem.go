package main

import "github.com/bennicholls/burl-E/burl"
import "math/rand"
import "math"

const GRAVCONST float64 = 6.674e-11

type StarSystem struct {
	Location

	Star    Star
	Planets []Planet
}

func NewStarSystem(c Coordinates) (s *StarSystem) {
	s = new(StarSystem)

	s.Name = "The Sol System"
	s.Description = "The only star system in the game, cosmological principle be damned! Leave it at your peril."
	s.LocationType = loc_STARSYSTEM
	s.Coords = c
	s.Coords.Resolution = coord_STARSYSTEM
	s.Coords.StarCoord = burl.Coord{X: 500, Y: 500}
	s.Star = NewStar(s.Coords, "The Sun", 695700e3, 1.988435e30)
	s.Planets = make([]Planet, 0, 10) //Starsystems have max 10 planets, right? Yeah sounds about right.

	s.Planets = append(s.Planets, NewPlanet(s.Coords, 57.3e9, 2493e3, 3.301e23, "Mercury"))
	s.Planets = append(s.Planets, NewPlanet(s.Coords, 108.2e9, 6051e3, 4.867e24, "Venus"))
	s.Planets = append(s.Planets, NewPlanet(s.Coords, 149.6e9, 6378e3, 5.972e24, "Earth"))
	s.Planets = append(s.Planets, NewPlanet(s.Coords, 227.9e9, 3390e3, 6.417e23, "Mars"))
	s.Planets = append(s.Planets, NewPlanet(s.Coords, 778.3e9, 71492e3, 1.898e27, "Jupiter"))
	s.Planets = append(s.Planets, NewPlanet(s.Coords, 1427e9, 60268e3, 5.68319e26, "Saturn"))
	s.Planets = append(s.Planets, NewPlanet(s.Coords, 2871e9, 25559e3, 8.681e25, "Uranus"))
	s.Planets = append(s.Planets, NewPlanet(s.Coords, 4497e9, 24764e3, 1.024e26, "Neptune"))

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

type Planet struct {
	Location
	oDistance float64 //orbital distance
	oPosition float64 //orbital position in degrees
	radius    float64 //radius of planet in meters
	mass      float64
}

//NewPlanet creates a planet. c is the coords of the starsystem. orbit is the distance in meters from the star
func NewPlanet(c Coordinates, orbit, radius, mass float64, name string) (p Planet) {
	p.Name = name
	p.Description = "This is a planet. Planets are rocks that are big enough to be important. Some planets have life on them, but most of them are super boring."
	p.LocationType = loc_PLANET
	p.Coords = c
	p.Coords.Resolution = coord_LOCAL

	p.oDistance = orbit
	p.oPosition = rand.Float64() * 2 * math.Pi
	p.radius = radius
	p.mass = mass
	p.VisitDistance = radius * 1.2
	p.VisitSpeed = math.Sqrt(GRAVCONST * p.mass / p.VisitDistance)
	p.Coords.Local.X = (coord_LOCAL_MAX / 2) + p.oDistance*math.Cos(p.oPosition)
	p.Coords.Local.Y = (coord_LOCAL_MAX / 2) + p.oDistance*math.Sin(p.oPosition)

	return
}

//NOTE: When moons and orbitting stuff gets implemented, be sure to add those here.
func (p Planet) GetLocations() []Locatable {
	l := make([]Locatable, 0)
	l = append(l, p)
	return l
}
