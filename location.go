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
	name         string
	locationType int         //see location type defs above
	explored     bool        //have we been there
	known        bool        //do we know about this place
	coords       Coordinates //where it at
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
	coord_SECTOR_MAX         = 25            //25x25 sectors in a galaxy
	coord_SUBSECTOR_MAX      = 1000          //1000x1000 subsectors in a sector
	coord_STARSYSTEM_MAX     = 1000          //1000x1000 locations for stars in a subsector
	coord_LOCAL_MAX          = 9461000000000 //9.416e12 meters in a starsystem segment
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
	if c.resolution == coord_SECTOR {
		return
	}

	xString += ":" + strconv.Itoa(c.xSubSector)
	yString += ":" + strconv.Itoa(c.ySubSector)
	if c.resolution == coord_SUBSECTOR {
		return
	}

	xString += ":" + strconv.Itoa(c.xStarCoord)
	yString += ":" + strconv.Itoa(c.yStarCoord)
	if c.resolution == coord_STARSYSTEM {
		return
	}

	xString += ":" + strconv.Itoa(c.xLocal)
	yString += ":" + strconv.Itoa(c.yLocal)

	return
}

func (c Coordinates) Sector() (int, int) {
	return c.xSector, c.ySector
}

func (c *Coordinates) Move(dx, dy, res int) {
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
