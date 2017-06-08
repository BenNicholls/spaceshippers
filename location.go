package main

import "strconv"
import "github.com/bennicholls/burl/util"
import "math"

const (
	METERS_PER_LY int = 9.461e15
	LY_PER_SECTOR int = 1000
)

//defines any place where your ship can travel to
type Locatable interface {
	GetName() string
	GetLocationType() LocationType
	IsExplored() bool
	IsKnown() bool //it is known
	GetCoords() Coordinates
	GetLocations() []Locatable
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
	name         string
	locationType LocationType //see location type defs above
	explored     bool         //have we been there
	known        bool         //do we know about this place
	coords       Coordinates  //where it at
}

func (l Location) GetName() string {
	return l.name
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
	xSector    int
	xSubSector int
	xStarCoord int
	xLocal     int

	ySector    int
	ySubSector int
	yStarCoord int
	yLocal     int

	resolution CoordResolution //how deep into the rabbit hole this coordinate goes. see above
}

//NewCoordinate makes a new Coordinate object, defaulted to the center of the galaxy.
func NewCoordinate(res CoordResolution) (c Coordinates) {
	c.resolution = res
	c.xSector = coord_SECTOR_MAX / 2
	c.ySector = coord_SECTOR_MAX / 2
	c.xSubSector = coord_SUBSECTOR_MAX / 2
	c.ySubSector = coord_SUBSECTOR_MAX / 2
	c.xStarCoord = coord_STARSYSTEM_MAX / 2
	c.yStarCoord = coord_STARSYSTEM_MAX / 2
	c.xLocal = coord_LOCAL_MAX / 2
	c.yLocal = coord_LOCAL_MAX / 2

	return
}

//NewCoordinate makes a new Coordinate object, defaulted to the center of a sector.
func NewSectorCoordinate(x, y int) (c Coordinates) {
	c = NewCoordinate(coord_SECTOR)
	c.xSector, c.ySector = x, y

	return
}

//Returns a string for each coordinate in the form
//SECTOR:SUBSECTOR:STARCOORD:LOCAL, subject to resolution limits.
func (c Coordinates) GetCoordStrings() (xString string, yString string) {
	xString, yString = strconv.Itoa(c.xSector), strconv.Itoa(c.ySector)
	if c.resolution == coord_SECTOR {
		return
	}

	xString += ":" + strconv.Itoa(c.xSubSector)
	yString += ":" + strconv.Itoa(c.ySubSector)
	if c.resolution == coord_SUBSECTOR {
		return
	}

	xString += ":" + strconv.FormatInt(int64(c.xStarCoord), 36)
	yString += ":" + strconv.FormatInt(int64(c.yStarCoord), 36)
	if c.resolution == coord_STARSYSTEM {
		return
	}

	xString += ":" + strconv.FormatInt(int64(c.xLocal), 36)
	yString += ":" + strconv.FormatInt(int64(c.yLocal), 36)

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

	if c1.xSector != c2.xSector || c1.ySector != c2.ySector {
		return false
	} else if c1.resolution == coord_SECTOR {
		return true
	}

	if c1.xSubSector != c2.xSubSector || c1.ySubSector != c2.ySubSector {
		return false
	} else if c1.resolution == coord_SUBSECTOR {
		return true
	}

	if c1.xStarCoord != c2.xStarCoord || c1.yStarCoord != c2.yStarCoord {
		return false
	} else if c1.resolution == coord_STARSYSTEM {
		return true
	}

	//if this point is reached, then we know c1 and c2 are both local points in the same starsystem,
	//so we want to see if c1 is orbitting/docking with l
	dist := int(c1.CalcVector(c2).Distance * float64(METERS_PER_LY))

	switch loc := l.(type) {
	case Planet:
		if dist < loc.orbitRange {
			return true
		}
	case Star:
		if dist < loc.orbitRange {
			return true
		}
	}

	return false
}

func (c Coordinates) Sector() (int, int) {
	return c.xSector, c.ySector
}

//returns the subsector portion of the coord. REMEMBER: not all coords handle these!
func (c Coordinates) SubSector() (int, int) {
	return c.xSubSector, c.ySubSector
}

//returns the starcoord portion of the coord. REMEMBER: not all coords handle these!
func (c Coordinates) StarCoord() (int, int) {
	return c.xStarCoord, c.yStarCoord
}

//returns the local portion of the coord. REMEMBER: not all coords handle these!
func (c Coordinates) LocalCoord() (int, int) {
	return c.xLocal, c.yLocal
}

func (c *Coordinates) Move(dx, dy int, res CoordResolution) {
	//check to ensure this coord handles the proper resolution
	if res > c.resolution {
		return
	}

	switch res {
	case coord_LOCAL:
		c.moveLocal(dx, dy)
	case coord_STARSYSTEM:
		c.moveStarSystem(dx, dy)
	case coord_SUBSECTOR:
		c.moveSubSector(dx, dy)
	case coord_SECTOR:
		c.moveSector(dx, dy)
	}
}

func (c *Coordinates) moveLocal(dx, dy int) {
	var odx, ody int

	c.xLocal, odx = util.ModularClamp(c.xLocal+dx, 0, coord_LOCAL_MAX-1)
	c.yLocal, ody = util.ModularClamp(c.yLocal+dy, 0, coord_LOCAL_MAX-1)

	if odx != 0 || ody != 0 {
		c.moveStarSystem(odx, ody)
	}
}

func (c *Coordinates) moveStarSystem(dx, dy int) {
	var odx, ody int

	c.xStarCoord, odx = util.ModularClamp(c.xStarCoord+dx, 0, coord_STARSYSTEM_MAX-1)
	c.yStarCoord, ody = util.ModularClamp(c.yStarCoord+dy, 0, coord_STARSYSTEM_MAX-1)

	if odx != 0 || ody != 0 {
		c.moveSubSector(odx, ody)
	}
}

func (c *Coordinates) moveSubSector(dx, dy int) {
	var odx, ody int

	c.xSubSector, odx = util.ModularClamp(c.xSubSector+dx, 0, coord_SUBSECTOR_MAX-1)
	c.ySubSector, ody = util.ModularClamp(c.ySubSector+dy, 0, coord_SUBSECTOR_MAX-1)

	if odx != 0 || ody != 0 {
		c.moveSector(odx, ody)
	}
}

func (c *Coordinates) moveSector(dx, dy int) {
	c.xSector = util.Clamp(c.xSector+dx, 0, coord_SECTOR_MAX-1)
	c.ySector = util.Clamp(c.ySector+dy, 0, coord_SECTOR_MAX-1)
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

	g.xSector = c2.xSector - c1.xSector
	xDistance := float64(g.xSector) * float64(LY_PER_SECTOR)
	g.ySector = c2.ySector - c1.ySector
	yDistance := float64(g.ySector) * float64(LY_PER_SECTOR)

	g.xSubSector = c2.xSubSector - c1.xSubSector
	g.ySubSector = c2.ySubSector - c1.ySubSector
	xDistance += float64(g.xSubSector) * float64(LY_PER_SECTOR) / float64(coord_SUBSECTOR_MAX)
	yDistance += float64(g.ySubSector) * float64(LY_PER_SECTOR) / float64(coord_SUBSECTOR_MAX)

	g.xStarCoord = c2.xStarCoord - c1.xStarCoord
	g.yStarCoord = c2.yStarCoord - c1.yStarCoord
	xDistance += float64(g.xStarCoord) * float64(LY_PER_SECTOR) / float64(coord_SUBSECTOR_MAX) / float64(coord_STARSYSTEM_MAX)
	yDistance += float64(g.yStarCoord) * float64(LY_PER_SECTOR) / float64(coord_SUBSECTOR_MAX) / float64(coord_STARSYSTEM_MAX)

	g.xLocal = c2.xLocal - c1.xLocal
	g.yLocal = c2.yLocal - c1.yLocal
	xDistance += float64(g.xLocal) / float64(METERS_PER_LY)
	yDistance += float64(g.yLocal) / float64(METERS_PER_LY)

	g.Distance = math.Sqrt(xDistance*xDistance + yDistance*yDistance)

	return
}
