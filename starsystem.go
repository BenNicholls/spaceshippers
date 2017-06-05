package main

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
	s.Star = NewStar(s.coords, "The Sun")
	s.Planets = make([]Planet, 0, 10)

	s.Planets = append(s.Planets, NewPlanet(s.coords, 57.3e9, "Mercury"))
	s.Planets = append(s.Planets, NewPlanet(s.coords, 108.2e9, "Venus"))
	s.Planets = append(s.Planets, NewPlanet(s.coords, 149.6e9, "Earth"))
	s.Planets = append(s.Planets, NewPlanet(s.coords, 227.9e9, "Mars"))
	s.Planets = append(s.Planets, NewPlanet(s.coords, 778.3e9, "Jupiter"))
	s.Planets = append(s.Planets, NewPlanet(s.coords, 1427e9, "Saturn"))
	s.Planets = append(s.Planets, NewPlanet(s.coords, 2871e9, "Uranus"))
	s.Planets = append(s.Planets, NewPlanet(s.coords, 4497e9, "Neptune"))

	return
}

type Star struct {
	Location
}

//NewStar creates a star. c is the coordinates of the starsystem. Defaults to center of system.
func NewStar(c Coordinates, name string) (s Star) {
	s.name = name
	s.coords = c
	s.coords.resolution = coord_LOCAL
	s.coords.xLocal = coord_LOCAL_MAX / 2
	s.coords.yLocal = coord_LOCAL_MAX / 2

	return
}

type Planet struct {
	Location
	oDistance int //orbital distance
}

//NewPlanet creates a planet. c is the coords of the starsystem. orbit is the distance in meters from the star
func NewPlanet(c Coordinates, orbit int, name string) (p Planet) {
	p.name = name
	p.oDistance = orbit
	p.coords = c
	p.coords.xLocal = (coord_LOCAL_MAX / 2) + p.oDistance
	p.coords.yLocal = coord_LOCAL_MAX / 2

	return
}
