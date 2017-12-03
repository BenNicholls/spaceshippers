package main

import "strconv"
import "math"
import "github.com/bennicholls/burl-E/burl"

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
	Name          string
	Description   string
	LocationType  LocationType //see location type defs above
	Explored      bool         //have we been there
	Known         bool         //do we know about this place
	Coords        Coordinates  //where it at
	VisitDistance float64      //how close do you have to be to visit? ie orbit distance.
	VisitSpeed    float64      //how fast should we be going when we get there? ie orbital velocity or relative docking speed
}

func (l Location) GetName() string {
	return l.Name
}

func (l Location) GetDescription() string {
	return l.Description
}

func (l Location) GetLocationType() LocationType {
	return l.LocationType
}

func (l Location) IsExplored() bool {
	return l.Explored
}

func (l *Location) SetExplored() {
	l.Explored = true
}

func (l *Location) SetKnown() {
	l.Known = true
}

func (l Location) IsKnown() bool {
	return l.Known
}

func (l Location) GetCoords() Coordinates {
	return l.Coords
}

//Returns orbital distance or docking distance. Only works for local locations, returns 0 otherwise.
func (l Location) GetVisitDistance() float64 {
	return l.VisitDistance
}

//Returns orbital speed or docking speed. Only works for local locations, returns 0 otherwise.
func (l Location) GetVisitSpeed() float64 {
	return l.VisitSpeed
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
	Sector    burl.Coord
	SubSector burl.Coord
	StarCoord burl.Coord
	Local     burl.Vec2

	Resolution CoordResolution //how deep into the rabbit hole this coordinate goes. see above
}

//NewCoordinate makes a new Coordinate object, defaulted to the center of the galaxy.
func NewCoordinate(res CoordResolution) (c Coordinates) {
	c.Resolution = res
	c.Sector.MoveTo(coord_SECTOR_MAX/2, coord_SECTOR_MAX/2)
	c.SubSector.MoveTo(coord_SUBSECTOR_MAX/2, coord_SUBSECTOR_MAX/2)
	c.StarCoord.MoveTo(coord_STARSYSTEM_MAX/2, coord_STARSYSTEM_MAX/2)
	c.Local.Set(coord_LOCAL_MAX/2, coord_LOCAL_MAX/2)

	return
}

//NewCoordinate makes a new Coordinate object, defaulted to the center of a sector.
func NewSectorCoordinate(x, y int) (c Coordinates) {
	c = NewCoordinate(coord_SECTOR)
	c.Sector.MoveTo(x, y)

	return
}

//Returns a string for each coordinate in the form
//SECTOR:SUBSECTOR:STARCOORD:LOCAL, subject to resolution limits.
func (c Coordinates) GetCoordStrings() (xString string, yString string) {
	xString, yString = strconv.Itoa(c.Sector.X), strconv.Itoa(c.Sector.Y)
	if c.Resolution == coord_SECTOR {
		return
	}

	xString += ":" + strconv.Itoa(c.SubSector.X)
	yString += ":" + strconv.Itoa(c.SubSector.Y)
	if c.Resolution == coord_SUBSECTOR {
		return
	}

	xString += ":" + strconv.FormatInt(int64(c.StarCoord.X), 36)
	yString += ":" + strconv.FormatInt(int64(c.StarCoord.Y), 36)
	if c.Resolution == coord_STARSYSTEM {
		return
	}

	xString += ":" + strconv.FormatInt(int64(c.Local.X), 36)
	yString += ":" + strconv.FormatInt(int64(c.Local.Y), 36)

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
	if c2.Resolution > c1.Resolution {
		return false
	}

	if c1.Sector != c2.Sector {
		return false
	} else if c2.Resolution == coord_SECTOR {
		return true
	}

	if c1.SubSector != c2.SubSector {
		return false
	} else if c2.Resolution == coord_SUBSECTOR {
		return true
	}

	if c1.StarCoord != c2.StarCoord {
		return false
	} else if c2.Resolution == coord_STARSYSTEM {
		return true
	}

	//if this point is reached, then we know c1 and c2 are both local points in the same starsystem,
	//so we want to see if c1 is orbitting/docking with l
	if dist := c1.CalcVector(c2).Distance * METERS_PER_LY; dist <= l.GetVisitDistance() {
		return true
	}

	return false
}

func (c *Coordinates) Move(dx, dy int, res CoordResolution) {
	//check to ensure this coord handles the proper resolution
	if res > c.Resolution {
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
	xDecimals := (c.Local.X - math.Trunc(c.Local.X)) + (dx - math.Trunc(dx))
	yDecimals := (c.Local.Y - math.Trunc(c.Local.Y)) + (dy - math.Trunc(dy))

	x, odx := burl.ModularClamp(int(c.Local.X)+int(dx), 0, int(coord_LOCAL_MAX)-1)
	y, ody := burl.ModularClamp(int(c.Local.Y)+int(dy), 0, int(coord_LOCAL_MAX)-1)

	c.Local.X = float64(x) + xDecimals
	c.Local.Y = float64(y) + yDecimals

	if odx != 0 || ody != 0 {
		c.moveStarSystem(odx, ody)
	}
}

func (c *Coordinates) moveStarSystem(dx, dy int) {
	var odx, ody int

	c.StarCoord.X, odx = burl.ModularClamp(c.StarCoord.X+dx, 0, coord_STARSYSTEM_MAX-1)
	c.StarCoord.Y, ody = burl.ModularClamp(c.StarCoord.Y+dy, 0, coord_STARSYSTEM_MAX-1)

	if odx != 0 || ody != 0 {
		c.moveSubSector(odx, ody)
	}
}

func (c *Coordinates) moveSubSector(dx, dy int) {
	var odx, ody int

	c.SubSector.X, odx = burl.ModularClamp(c.SubSector.X+dx, 0, coord_SUBSECTOR_MAX-1)
	c.SubSector.Y, ody = burl.ModularClamp(c.SubSector.Y+dy, 0, coord_SUBSECTOR_MAX-1)

	if odx != 0 || ody != 0 {
		c.moveSector(odx, ody)
	}
}

func (c *Coordinates) moveSector(dx, dy int) {
	c.Sector.X = burl.Clamp(c.Sector.X+dx, 0, coord_SECTOR_MAX-1)
	c.Sector.Y = burl.Clamp(c.Sector.Y+dy, 0, coord_SECTOR_MAX-1)
}

//Galactic Vector. represents the vector between two points in the galaxy
type GalVec struct {
	Coordinates         //0-centered vector
	Distance    float64 //magnitude in Ly
}

func (c1 Coordinates) CalcVector(c2 Coordinates) (g GalVec) {
	g.Sector = c2.Sector.Sub(c1.Sector)
	g.SubSector = c2.SubSector.Sub(c1.SubSector)
	g.StarCoord = c2.StarCoord.Sub(c1.StarCoord)
	g.Local = c2.Local.Sub(c1.Local)

	xDistance := g.Local.X/METERS_PER_LY + float64(LY_PER_SECTOR)*(float64(g.Sector.X)+(float64(g.SubSector.X)+float64(g.StarCoord.X)/float64(coord_STARSYSTEM_MAX))/float64(coord_SUBSECTOR_MAX))
	yDistance := g.Local.Y/METERS_PER_LY + float64(LY_PER_SECTOR)*(float64(g.Sector.Y)+(float64(g.SubSector.Y)+float64(g.StarCoord.Y)/float64(coord_STARSYSTEM_MAX))/float64(coord_SUBSECTOR_MAX))
	g.Distance = math.Sqrt(xDistance*xDistance + yDistance*yDistance)

	return
}
