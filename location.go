package main

import "strconv"
import "github.com/bennicholls/burl/util"

//defines any place where your ship can travel to
type Locatable interface {
	GetName() string
	GetLocationType() int
	IsExplored() bool
	IsKnown() bool //it is known
	GetCoords() Coordinates
}

//different types of locations
const (
	loc_NONE int = iota
	loc_SECTOR
	loc_STARSYSTEM
	loc_PLANET
	loc_ANOMALY
)

type Location struct {
	name             string
	locationType     int  //see location type defs above
	explored         bool //have we been there
	known            bool //do we know about this place
	coords           Coordinates //where it at
}

func (l Location) GetName() string {
	return l.name
}

func (l Location) GetLocationType() int {
	return l.locationType
}

func (l Location) IsExplored() bool {
	return l.explored
}

func (l *Location) SetExplored() {
	l.explored = true
}

func (l Location) IsKnown() bool {
	return l.known
}

func (l Location) GetCoords() Coordinates {
	return l.coords
}

//coordinate resolution levels
const (
	coord_SECTOR int = iota
	coord_SUBSECTOR
	coord_STARSYSTEM
	coord_LOCAL
)

const (
	coord_SECTOR_MAX = 25 //25x25 sectors in a galaxy
	coord_SUBSECTOR_MAX = 1000 //1000x1000 subsectors in a sector
	coord_STARSYSTEM_MAX = 1000 //1000x1000 locations for stars in a subsector
	coord_LOCAL_MAX int = 9461000000000 //9.416e12 meters in a starsystem segment
)

type Coordinates struct {
	xSector int
	xSubSector int
	xStarCoord int
	xLocal int

	ySector int
	ySubSector int
	yStarCoord int
	yLocal int

	resolution int //how deep into the rabbit hole this coordinate goes. see above
}

func NewCoordinate(res int) Coordinates {
	c := Coordinates{}
	c.resolution = res
	return c
}

func NewSectorCoordinate(x, y int) Coordinates {
	c := Coordinates{}
	c.xSector, c.ySector = x, y
	c.resolution = coord_SECTOR

	return c
}

//Returns a string for each coordinate in the form
//SECTOR:SUBSECTOR:STARCOORD:LOCAL, subject to resolution limits.
func (c Coordinates) GetCoordStrings() (xString string, yString string) {
	xString, yString = strconv.Itoa(c.xSector), strconv.Itoa(c.ySector)
	if c.resolution == coord_SECTOR { return }

	xString += ":" + strconv.Itoa(c.xSubSector)
	yString += ":" + strconv.Itoa(c.ySubSector)
	if c.resolution == coord_SUBSECTOR { return }

	xString += ":" + strconv.Itoa(c.xStarCoord)
	yString += ":" + strconv.Itoa(c.yStarCoord)
	if c.resolution == coord_STARSYSTEM { return }

	xString += ":" + strconv.Itoa(c.xLocal)
	yString += ":" + strconv.Itoa(c.yLocal)

	return
}

func (c Coordinates) Sector() (int, int) {
	return c.xSector, c.ySector
}

//Holy damn this is not right.
func (c *Coordinates) Move(dx, dy, res int) {
	//check to ensure this coord handles the proper resolution
	if res > c.resolution {
		return
	}

	switch res {
	case coord_LOCAL:
		c.xLocal = util.ModularClamp(c.xLocal + dx, 0, coord_LOCAL_MAX)
		c.yLocal = util.ModularClamp(c.yLocal + dy, 0, coord_LOCAL_MAX)
	case coord_STARSYSTEM:
		c.xStarCoord = util.ModularClamp(c.xStarCoord + dx, 0, coord_STARSYSTEM_MAX)
		c.yStarCoord = util.ModularClamp(c.yStarCoord + dy, 0, coord_STARSYSTEM_MAX)
	case coord_SUBSECTOR:
		c.xSubSector = util.ModularClamp(c.xSubSector + dx, 0, coord_SUBSECTOR_MAX)
		c.ySubSector = util.ModularClamp(c.ySubSector + dy, 0, coord_SUBSECTOR_MAX)
	case coord_SECTOR:
		c.xSector = util.Clamp(c.xSector + dx, 0, coord_SECTOR_MAX - 1)
		c.ySector = util.Clamp(c.ySector + dy, 0, coord_SECTOR_MAX - 1)
	}
}