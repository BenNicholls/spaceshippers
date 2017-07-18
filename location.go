package main

import "strconv"
import "github.com/bennicholls/burl/util"
import "math"

const (
	METERS_PER_LY float64 = 9.461e15
	LY_PER_SECTOR int     = 1000
)

//defines any place where your ship can travel to
type Locatable interface {
	GetName() string
	GetDescription() string
	GetLocationType() LocationType
	IsExplored() bool
	IsKnown() bool //it is known
	GetCoords() Coordinates
	GetLocations() []Locatable
	GetVisitDistance() float64
	GetVisitSpeed() float64
}

type LocationType int

//different types of locations
const (
	loc_NONE LocationType = iota
	loc_SECTOR
	loc_STARSYSTEM
	loc_PLANET
	loc_MOON
	loc_ANOMALY
	loc_SHIP
)

type Location struct {
	name          string
	description   string
	locationType  LocationType //see location type defs above
	explored      bool         //have we been there
	known         bool         //do we know about this place
	coords        Coordinates  //where it at
	visitDistance float64      //how close do you have to be to visit? ie orbit distance.
	visitSpeed    float64      //how fast should we be going when we get there? ie orbital velocity or relative docking speed
}

func (l Location) GetName() string {
	return l.name
}

func (l Location) GetDescription() string {
	return l.description
}

func (l Location) GetLocationType() LocationType {
	return l.locationType
}

func (l Location) IsExplored() bool {
	return l.explored
}

func (l *Location) SetExplored() {
	l.explored = true
}

func (l *Location) SetKnown() {
	l.known = true
}

func (l Location) IsKnown() bool {
	return l.known
}

func (l Location) GetCoords() Coordinates {
	return l.coords
}

//Returns orbital distance or docking distance. Only works for local locations, returns 0 otherwise.
func (l Location) GetVisitDistance() float64 {
	return l.visitDistance
}

//Returns orbital speed or docking speed. Only works for local locations, returns 0 otherwise.
func (l Location) GetVisitSpeed() float64 {
	return l.visitSpeed
}

//This default method actually doesn't work, we don't want the Location object. Hmm.
func (l Location) GetLocations() []Locatable {
	a := make([]Locatable, 1)
	a[0] = l
	return a
}

type CoordResolution int

//coordinate resolution levels
const (
	coord_SECTOR CoordResolution = iota
	coord_SUBSECTOR
	coord_STARSYSTEM
	coord_LOCAL
)

const (
	coord_SECTOR_MAX     = 25                   //25x25 sectors in a galaxy
	coord_SUBSECTOR_MAX  = 1000                 //1000x1000 subsectors in a sector
	coord_STARSYSTEM_MAX = 1000                 //1000x1000 locations for stars in sector
	coord_LOCAL_MAX      = METERS_PER_LY / 1000 //9.416e12 meters to a side for a starsystem
)

type Coordinates struct {
	sector    util.Coord
	subSector util.Coord
	starCoord util.Coord
	local     util.Vec2

	resolution CoordResolution //how deep into the rabbit hole this coordinate goes. see above
}

//NewCoordinate makes a new Coordinate object, defaulted to the center of the galaxy.
func NewCoordinate(res CoordResolution) (c Coordinates) {
	c.resolution = res
	c.sector.MoveTo(coord_SECTOR_MAX/2, coord_SECTOR_MAX/2)
	c.subSector.MoveTo(coord_SUBSECTOR_MAX/2, coord_SUBSECTOR_MAX/2)
	c.starCoord.MoveTo(coord_STARSYSTEM_MAX/2, coord_STARSYSTEM_MAX/2)
	c.local.Set(coord_LOCAL_MAX/2, coord_LOCAL_MAX/2)

	return
}

//NewCoordinate makes a new Coordinate object, defaulted to the center of a sector.
func NewSectorCoordinate(x, y int) (c Coordinates) {
	c = NewCoordinate(coord_SECTOR)
	c.sector.MoveTo(x, y)

	return
}

//Returns a string for each coordinate in the form
//SECTOR:SUBSECTOR:STARCOORD:LOCAL, subject to resolution limits.
func (c Coordinates) GetCoordStrings() (xString string, yString string) {
	xString, yString = strconv.Itoa(c.sector.X), strconv.Itoa(c.sector.Y)
	if c.resolution == coord_SECTOR {
		return
	}

	xString += ":" + strconv.Itoa(c.subSector.X)
	yString += ":" + strconv.Itoa(c.subSector.Y)
	if c.resolution == coord_SUBSECTOR {
		return
	}

	xString += ":" + strconv.FormatInt(int64(c.starCoord.X), 36)
	yString += ":" + strconv.FormatInt(int64(c.starCoord.Y), 36)
	if c.resolution == coord_STARSYSTEM {
		return
	}

	xString += ":" + strconv.FormatInt(int64(c.local.X), 36)
	yString += ":" + strconv.FormatInt(int64(c.local.Y), 36)

	return
}

//Reports whether coordinate c1 is "inside" a location l.
//NOTES: if c1 and l are the same, this reports true. if c1 and l are both Local
//objects, it tests instead to see if c1 is in orbit/docking range.
func (c1 Coordinates) IsIn(l Locatable) bool {
	if l == nil {
		return false
	}
	c2 := l.GetCoords()
	if c2.resolution > c1.resolution {
		return false
	}

	if c1.sector != c2.sector {
		return false
	} else if c2.resolution == coord_SECTOR {
		return true
	}

	if c1.subSector != c2.subSector {
		return false
	} else if c2.resolution == coord_SUBSECTOR {
		return true
	}

	if c1.starCoord != c2.starCoord {
		return false
	} else if c2.resolution == coord_STARSYSTEM {
		return true
	}

	//if this point is reached, then we know c1 and c2 are both local points in the same starsystem,
	//so we want to see if c1 is orbitting/docking with l
	if dist := c1.CalcVector(c2).Distance * METERS_PER_LY; dist <= l.GetVisitDistance() {
		return true
	}

	return false
}

func (c Coordinates) Sector() util.Coord {
	return c.sector
}

//returns the subsector portion of the coord. REMEMBER: not all coords handle these!
func (c Coordinates) SubSector() util.Coord {
	return c.subSector
}

//returns the starcoord portion of the coord. REMEMBER: not all coords handle these!
func (c Coordinates) StarCoord() util.Coord {
	return c.starCoord
}

//returns the local portion of the coord. REMEMBER: not all coords handle these!
func (c Coordinates) LocalCoord() util.Vec2 {
	return c.local
}

func (c *Coordinates) Move(dx, dy int, res CoordResolution) {
	//check to ensure this coord handles the proper resolution
	if res > c.resolution {
		return
	}

	switch res {
	case coord_LOCAL:
		c.moveLocal(float64(dx), float64(dy))
	case coord_STARSYSTEM:
		c.moveStarSystem(dx, dy)
	case coord_SUBSECTOR:
		c.moveSubSector(dx, dy)
	case coord_SECTOR:
		c.moveSector(dx, dy)
	}
}

func (c *Coordinates) moveLocal(dx, dy float64) {
	xDecimals := (c.local.X - math.Trunc(c.local.X)) + (dx - math.Trunc(dx))
	yDecimals := (c.local.Y - math.Trunc(c.local.Y)) + (dy - math.Trunc(dy))

	x, odx := util.ModularClamp(int(c.local.X)+int(dx), 0, int(coord_LOCAL_MAX)-1)
	y, ody := util.ModularClamp(int(c.local.Y)+int(dy), 0, int(coord_LOCAL_MAX)-1)

	c.local.X = float64(x) + xDecimals
	c.local.Y = float64(y) + yDecimals

	if odx != 0 || ody != 0 {
		c.moveStarSystem(odx, ody)
	}
}

func (c *Coordinates) moveStarSystem(dx, dy int) {
	var odx, ody int

	c.starCoord.X, odx = util.ModularClamp(c.starCoord.X+dx, 0, coord_STARSYSTEM_MAX-1)
	c.starCoord.Y, ody = util.ModularClamp(c.starCoord.Y+dy, 0, coord_STARSYSTEM_MAX-1)

	if odx != 0 || ody != 0 {
		c.moveSubSector(odx, ody)
	}
}

func (c *Coordinates) moveSubSector(dx, dy int) {
	var odx, ody int

	c.subSector.X, odx = util.ModularClamp(c.subSector.X+dx, 0, coord_SUBSECTOR_MAX-1)
	c.subSector.Y, ody = util.ModularClamp(c.subSector.Y+dy, 0, coord_SUBSECTOR_MAX-1)

	if odx != 0 || ody != 0 {
		c.moveSector(odx, ody)
	}
}

func (c *Coordinates) moveSector(dx, dy int) {
	c.sector.X = util.Clamp(c.sector.X+dx, 0, coord_SECTOR_MAX-1)
	c.sector.Y = util.Clamp(c.sector.Y+dy, 0, coord_SECTOR_MAX-1)
}

//Galactic Vector. represents the vector between two points in the galaxy
type GalVec struct {
	Coordinates             //0-centered vector
	c1, c2      Coordinates //endpoints in galactic-space. stored for... some reason.
	Distance    float64     //magnitude in Ly
}

func (c1 Coordinates) CalcVector(c2 Coordinates) (g GalVec) {
	g.c1 = c1
	g.c2 = c2

	g.sector = c2.sector.Sub(c1.sector)
	g.subSector = c2.subSector.Sub(c1.subSector)
	g.starCoord = c2.starCoord.Sub(c1.starCoord)
	g.local = c2.local.Sub(c1.local)

	xDistance := g.local.X/METERS_PER_LY + float64(LY_PER_SECTOR)*(float64(g.sector.X)+(float64(g.subSector.X)+float64(g.starCoord.X)/float64(coord_STARSYSTEM_MAX))/float64(coord_SUBSECTOR_MAX))
	yDistance := g.local.Y/METERS_PER_LY + float64(LY_PER_SECTOR)*(float64(g.sector.Y)+(float64(g.subSector.Y)+float64(g.starCoord.Y)/float64(coord_STARSYSTEM_MAX))/float64(coord_SUBSECTOR_MAX))
	g.Distance = math.Sqrt(xDistance*xDistance + yDistance*yDistance)

	return
}
