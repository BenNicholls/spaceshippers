package main

import "strconv"
import "github.com/bennicholls/burl-E/burl"

//TODO: incorpoate GalaxyMapView object here, maybe once i put together a systemMapView

//Menu for viewing star charts, getting location data, setting courses, etc.
type StarchartMenu struct {
	burl.Container

	//general map stuff
	mapTitleText *burl.Textbox
	mapView      *burl.TileView
	mapHighlight *burl.PulseAnimation

	xCursor, yCursor int //map cursor coordinates

	//sector details for galaxy map mode
	sectorDetails      *burl.Container
	sectorCoordsText   *burl.Textbox
	sectorNameText     *burl.Textbox
	sectorDensityText  *burl.Textbox
	sectorExploredText *burl.Textbox
	sectorKnownText    *burl.Textbox
	sectorLocationText *burl.Textbox

	//StarSystem details for system mode
	systemDetails         *burl.Container
	systemLocTitleText    *burl.Textbox
	systemLocNameText     *burl.Textbox
	systemLocDescText     *burl.Textbox
	systemLocDistText     *burl.Textbox
	systemSetCourseButton *burl.Button
	systemLocationsList   *burl.List

	systemLocations []Locatable

	galaxy     *Galaxy         //to know what we're drawing
	playerShip *Ship           //to know where we are, current course, etc
	mapMode    CoordResolution //defines the zoom level.
	localZoom  int             //zoom level for local mode.
}

func NewStarchartMenu(gal *Galaxy, ship *Ship) (sm *StarchartMenu) {
	sm = new(StarchartMenu)
	sm.galaxy = gal
	sm.playerShip = ship
	sm.xCursor, sm.yCursor = ship.Coords.Sector.Get() //start sector picker on player ship
	sm.mapMode = coord_SECTOR

	//ui setup
	sm.Container = *burl.NewContainer(40, 26, 39, 4, 5, true)
	sm.SetTitle("Star Charts")
	sm.SetVisibility(false)

	sm.mapTitleText = burl.NewTextbox(25, 1, 0, 25, 1, false, true, sm.galaxy.name)
	sm.mapView = burl.NewTileView(25, 25, 0, 0, 1, false)

	sm.sectorDetails = burl.NewContainer(15, 26, 25, 0, 1, false)
	sm.sectorCoordsText = burl.NewTextbox(15, 1, 0, 0, 1, false, true, "")
	sm.sectorNameText = burl.NewTextbox(15, 2, 0, 1, 1, false, true, "")
	sm.sectorDensityText = burl.NewTextbox(15, 1, 0, 4, 1, false, false, "")
	sm.sectorExploredText = burl.NewTextbox(15, 1, 0, 5, 1, false, false, "")
	sm.sectorLocationText = burl.NewTextbox(15, 2, 0, 6, 1, false, false, "")
	sm.sectorKnownText = burl.NewTextbox(15, 2, 0, 9, 1, false, true, "We know nothing about this sector.")
	sm.sectorDetails.Add(sm.sectorCoordsText, sm.sectorNameText, sm.sectorDensityText, sm.sectorExploredText, sm.sectorLocationText, sm.sectorKnownText)

	sm.systemDetails = burl.NewContainer(15, 26, 25, 0, 1, false)
	sm.systemDetails.SetVisibility(false)
	sm.systemLocTitleText = burl.NewTextbox(15, 1, 0, 0, 1, false, true, "Currently viewing:")
	sm.systemLocNameText = burl.NewTextbox(15, 1, 0, 1, 1, false, true, "")
	sm.systemLocDescText = burl.NewTextbox(15, 5, 0, 3, 1, false, true, "")
	sm.systemLocDistText = burl.NewTextbox(15, 1, 0, 10, 1, false, false, "")
	sm.systemSetCourseButton = burl.NewButton(13, 1, 1, 15, 1, true, true, "Press Enter to Set Course!")
	sm.systemSetCourseButton.ToggleFocus()

	sm.systemLocationsList = burl.NewList(13, 5, 1, 20, 1, true, "NO LOCATIONS")
	sm.LoadLocalInfo()
	sm.systemDetails.Add(sm.systemLocTitleText, sm.systemLocNameText, sm.systemLocDescText, sm.systemLocDistText, sm.systemSetCourseButton, sm.systemLocationsList)

	sm.mapHighlight = burl.NewPulseAnimation(sm.xCursor, sm.yCursor, 1, 1, 50, 10, true)
	sm.mapHighlight.Activate()
	sm.mapView.AddAnimation(sm.mapHighlight)
	sm.Add(sm.mapView, sm.mapTitleText, sm.sectorDetails, sm.systemDetails)

	sm.ZoomIn() //start us off in local mode.

	return
}

func (sm *StarchartMenu) Update() {
	switch sm.mapMode {
	case coord_SECTOR:
		sm.UpdateSectorInfo()
	case coord_LOCAL:
		sm.UpdateLocalInfo()
	}
}

func (sm *StarchartMenu) UpdateSectorInfo() {
	sector := sm.galaxy.GetSector(sm.xCursor, sm.yCursor)
	sm.sectorCoordsText.ChangeText("Sector (" + sector.ProperName() + ")")
	sm.sectorNameText.ChangeText("\"" + sector.GetName() + "\"")
	sm.sectorDensityText.ChangeText("Star Density: " + strconv.Itoa(sector.Density) + "%")
	if sector.IsExplored() {
		sm.sectorExploredText.ChangeText("SECTOR EXPLORED!")
	} else {
		sm.sectorExploredText.ChangeText("SECTOR UNEXPLORED")
	}

	if sec := sm.playerShip.Coords.Sector; sm.xCursor == sec.X && sm.yCursor == sec.Y {
		sm.sectorLocationText.ChangeText("We're currently here!")
	} else if sm.playerShip.destination != nil {
		if sec := sm.playerShip.destination.GetCoords().Sector; sm.xCursor == sec.X && sm.yCursor == sec.Y {
			sm.sectorLocationText.ChangeText("We're currently going here!")
		}
	} else {
		d := sm.playerShip.Coords.CalcVector(sector.GetCoords()).Distance
		sm.sectorLocationText.ChangeText("Distance to sector center: " + strconv.FormatFloat(d, 'f', 2, 64) + "Ly.")
	}
}

//Loads the location list for the system being viewed.
func (sm *StarchartMenu) LoadLocalInfo() {
	c := sm.playerShip.Coords
	system := sm.galaxy.GetSector(c.Sector.Get()).GetSubSector(c.SubSector.Get()).starSystem
	sm.systemLocations = make([]Locatable, 1)
	sm.systemLocations[0] = sm.playerShip
	sm.systemLocations = append(sm.systemLocations, system.GetLocations()...)
	for _, l := range sm.systemLocations {
		sm.systemLocationsList.Append(l.GetName())
	}
	sm.UpdateLocalInfo()
}

func (sm *StarchartMenu) UpdateLocalInfo() {
	c := sm.playerShip.Coords
	system := sm.galaxy.GetSector(c.Sector.Get()).GetSubSector(c.SubSector.Get()).starSystem
	sm.mapTitleText.ChangeText(system.GetName())

	loc := sm.systemLocations[sm.systemLocationsList.GetSelection()]
	sm.systemLocNameText.ChangeText(loc.GetName())

	if loc.GetDescription() == "" {
		sm.systemLocDescText.ChangeText("What is this? How did you select this?")
	} else {
		sm.systemLocDescText.ChangeText(loc.GetDescription())
	}

	d := int(c.CalcVector(loc.GetCoords()).Distance * METERS_PER_LY / 1000)
	sm.systemLocDistText.ChangeText("Distance:" + strconv.Itoa(d) + "km.")
	if loc == sm.playerShip {
		sm.systemSetCourseButton.ChangeText("This is us!")
	} else if sm.playerShip.Coords.IsIn(loc) {
		sm.systemSetCourseButton.ChangeText("We are currently here!")
	} else {
		sm.systemSetCourseButton.ChangeText("Press Enter to Set Course!")
	}
}

//draws the required map. galaxy map, sector map, star system map
func (sm *StarchartMenu) DrawMap() {
	switch sm.mapMode {
	case coord_SECTOR:
		sm.DrawGalaxy()
	case coord_LOCAL:
		sm.DrawSystem()
	}
}

func (sm *StarchartMenu) DrawGalaxy() {

	/////NOTE: Replace this with GalaxyMapView.DrawGalaxy
	//--------------------------
	w, h := sm.galaxy.Dims()
	for i := 0; i < w*h; i++ {
		x, y := i%w, i/w
		s := sm.galaxy.GetSector(x, y)
		bright := burl.Lerp(0, 255, s.Density, 100)
		g := burl.GLYPH_FILL_SPARSE
		if bright == 0 {
			g = burl.GLYPH_NONE
		}
		sm.mapView.Draw(x, y, g, burl.MakeOpaqueColour(bright, bright, bright), burl.COL_BLACK)
	}
	//----------------------------

	x, y := sm.playerShip.Coords.Sector.Get()
	sm.mapView.Draw(x, y, burl.GLYPH_FACE2, burl.COL_WHITE, burl.COL_BLACK)

	ex, ey := sm.galaxy.earth.Sector.Get()
	sm.mapView.Draw(ex, ey, burl.GLYPH_DONUT, burl.COL_BLUE, burl.COL_BLACK)
}

func (sm *StarchartMenu) DrawSystem() {
	sm.mapView.Reset()
	c := sm.systemLocations[sm.systemLocationsList.GetSelection()].GetCoords()
	w, h := sm.mapView.Dims()
	gFactor := int(coord_LOCAL_MAX) / w / burl.Pow(2, sm.localZoom)

	xCamera := int(c.Local.X) - (gFactor * w / 2)
	yCamera := int(c.Local.Y) - (gFactor * h / 2)

	system := sm.galaxy.GetSector(c.Sector.Get()).GetSubSector(c.SubSector.Get()).starSystem
	sx, sy := system.Star.GetCoords().Local.GetInt() // sun coords
	px, py := sm.playerShip.Coords.Local.GetInt()    //player coords

	//draw system things!
	for _, p := range system.Planets {
		x, y := p.GetCoords().Local.GetInt()
		//draw orbits
		//TODO: some way of culling orbt drawing. currently drawing all of them
		sm.mapView.DrawCircle((sx-xCamera)/gFactor, (sy-yCamera)/gFactor, burl.RoundFloatToInt(p.oDistance/float64(gFactor)), burl.GLYPH_PERIOD, 0xFF114411, burl.COL_BLACK)
		if burl.IsInside(x, y, xCamera, yCamera, gFactor*w, gFactor*h) {
			//draw planet on top of orbit (hopefully)
			sm.mapView.Draw((x-xCamera)/gFactor, (y-yCamera)/gFactor, int(rune(p.Name[0])), 0xFF825814, burl.COL_BLACK)
		}
	}

	if burl.IsInside(sx, sy, xCamera, yCamera, gFactor*w, gFactor*h) {
		sm.mapView.Draw((sx-xCamera)/gFactor, (sy-yCamera)/gFactor, burl.GLYPH_STAR, burl.COL_YELLOW, burl.COL_BLACK)
	}

	//draw player ship last to ensure nothing overwrites it
	if burl.IsInside(px, py, xCamera, yCamera, gFactor*w, gFactor*h) {
		sm.mapView.Draw((px-xCamera)/gFactor, (py-yCamera)/gFactor, burl.GLYPH_FACE2, burl.COL_WHITE, burl.COL_BLACK)
	}
}

//ugh, hack incoming:
func (sm *StarchartMenu) MoveMapCursor(dx, dy int) {
	if sm.mapMode == coord_SECTOR {
		w, h := sm.mapView.Dims()
		if burl.CheckBounds(sm.xCursor+dx, sm.yCursor+dy, w, h) {
			sm.xCursor += dx
			sm.yCursor += dy
			sm.mapHighlight.MoveTo(sm.xCursor, sm.yCursor)
			sm.Update()
		}
	} else if sm.mapMode == coord_LOCAL {
		if dy == 1 {
			sm.systemLocationsList.Next()
		} else if dy == -1 {
			sm.systemLocationsList.Prev()
		}

		sm.UpdateLocalInfo()
	}

	sm.DrawMap()
}

func (sm *StarchartMenu) ZoomIn() {
	if sm.mapMode == coord_SECTOR {
		sm.mapMode = coord_LOCAL
		sm.localZoom = 0
		sm.sectorDetails.ToggleVisible()
		sm.systemDetails.ToggleVisible()
		sm.UpdateLocalInfo()
		sm.mapHighlight.Toggle()
	} else {
		sm.localZoom = burl.Clamp(sm.localZoom+1, 0, 8)
		sm.mapView.Reset()
	}

	sm.DrawMap()
}

func (sm *StarchartMenu) ZoomOut() {
	if sm.mapMode == coord_LOCAL {
		if sm.localZoom == 0 {
			sm.mapMode = coord_SECTOR
			sm.sectorDetails.ToggleVisible()
			sm.systemDetails.ToggleVisible()
			sm.UpdateSectorInfo()
			sm.mapHighlight.Toggle()
		} else {
			sm.localZoom -= 1
		}
	}

	sm.DrawMap()
}

func (sm *StarchartMenu) OnActivate() {
	sm.xCursor, sm.yCursor = sm.playerShip.Coords.Sector.Get()
	sm.mapHighlight.MoveTo(sm.xCursor, sm.yCursor)
	sm.Update()
	sm.DrawMap()
}
