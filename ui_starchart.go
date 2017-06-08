package main

import "strconv"
import "github.com/bennicholls/burl/ui"
import "github.com/bennicholls/burl/console"
import "github.com/bennicholls/burl/util"

//Menu for viewing star charts, getting location data, setting courses, etc.
type StarchartMenu struct {
	ui.Container

	//general map stuff
	mapTitleText *ui.Textbox
	mapView      *ui.TileView
	mapHighlight *ui.PulseAnimation

	xCursor, yCursor int //map cursor coordinates

	//sector details for galaxy map mode
	sectorDetails      *ui.Container
	sectorCoordsText   *ui.Textbox
	sectorNameText     *ui.Textbox
	sectorDensityText  *ui.Textbox
	sectorExploredText *ui.Textbox
	sectorKnownText    *ui.Textbox
	sectorLocationText *ui.Textbox

	//StarSystem details for system mode
	systemDetails         *ui.Container
	systemLocTitleText    *ui.Textbox
	systemLocNameText     *ui.Textbox
	systemLocDescText     *ui.Textbox
	systemLocDistText     *ui.Textbox
	systemSetCourseButton *ui.Button
	systemLocationsList   *ui.List

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
	sm.xCursor, sm.yCursor = ship.coords.Sector() //start sector picker on player ship
	sm.mapMode = coord_SECTOR

	//ui setup
	sm.Container = *ui.NewContainer(40, 26, 39, 4, 1, true)
	sm.SetTitle("Star Charts")
	sm.SetVisibility(false)

	sm.mapTitleText = ui.NewTextbox(25, 1, 0, 25, 1, false, true, sm.galaxy.name)
	sm.mapView = ui.NewTileView(25, 25, 0, 0, 1, false)

	sm.sectorDetails = ui.NewContainer(15, 26, 25, 0, 1, false)
	sm.sectorCoordsText = ui.NewTextbox(15, 1, 0, 0, 1, false, true, "")
	sm.sectorNameText = ui.NewTextbox(15, 2, 0, 1, 1, false, true, "")
	sm.sectorDensityText = ui.NewTextbox(15, 1, 0, 4, 1, false, false, "")
	sm.sectorExploredText = ui.NewTextbox(15, 1, 0, 5, 1, false, false, "")
	sm.sectorLocationText = ui.NewTextbox(15, 2, 0, 6, 1, false, false, "")
	sm.sectorKnownText = ui.NewTextbox(15, 2, 0, 9, 1, false, true, "We know nothing about this sector.")
	sm.sectorDetails.Add(sm.sectorCoordsText, sm.sectorNameText, sm.sectorDensityText, sm.sectorExploredText, sm.sectorLocationText, sm.sectorKnownText)

	sm.systemDetails = ui.NewContainer(15, 26, 25, 0, 1, false)
	sm.systemDetails.SetVisibility(false)
	sm.systemLocTitleText = ui.NewTextbox(15, 1, 0, 0, 0, false, true, "Currently viewing:")
	sm.systemLocNameText = ui.NewTextbox(15, 1, 0, 1, 0, false, true, "")
	sm.systemLocDescText = ui.NewTextbox(15, 5, 0, 3, 0, false, true, "")
	sm.systemLocDistText = ui.NewTextbox(15, 1, 0, 10, 0, false, false, "")
	sm.systemSetCourseButton = ui.NewButton(13, 1, 1, 15, 0, true, true, "Press Enter to Go There!")

	sm.systemLocationsList = ui.NewList(13, 5, 1, 20, 0, true, "NO LOCATIONS")
	sm.systemDetails.Add(sm.systemLocTitleText, sm.systemLocNameText, sm.systemLocDescText, sm.systemLocDistText, sm.systemSetCourseButton, sm.systemLocationsList)

	sm.mapHighlight = ui.NewPulseAnimation(sm.xCursor, sm.yCursor, 1, 1, 50, 10, true)
	sm.mapHighlight.Activate()
	sm.mapView.AddAnimation(sm.mapHighlight)
	sm.Add(sm.mapView, sm.mapTitleText, sm.sectorDetails, sm.systemDetails)

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

	if x, y := sm.playerShip.coords.Sector(); sm.xCursor == x && sm.yCursor == y {
		sm.sectorLocationText.ChangeText("We're currently here!")
	} else if sm.playerShip.Destination != nil {
		if x, y := sm.playerShip.Destination.GetCoords().Sector(); sm.xCursor == x && sm.yCursor == y {
			sm.sectorLocationText.ChangeText("We're currently going here!")
		}
	} else {
		d := sm.playerShip.coords.CalcVector(sector.GetCoords()).Distance
		sm.sectorLocationText.ChangeText("Distance to sector center: " + strconv.FormatFloat(d, 'f', 2, 64) + "Ly.")
	}
}

//Loads the location list for the system being viewed.
func (sm *StarchartMenu) LoadLocalInfo() {
	c := sm.playerShip.coords
	system := sm.galaxy.GetSector(c.xSector, c.ySector).GetSubSector(c.xSubSector, c.ySubSector).star
	sm.systemLocations = make([]Locatable, 1)
	sm.systemLocations[0] = sm.playerShip
	sm.systemLocations = append(sm.systemLocations, system.GetLocations()...)
	for _, l := range sm.systemLocations {
		sm.systemLocationsList.Append(l.GetName())
	}
	sm.UpdateLocalInfo()
}

func (sm *StarchartMenu) UpdateLocalInfo() {
	c := sm.playerShip.coords
	system := sm.galaxy.GetSector(c.xSector, c.ySector).GetSubSector(c.xSubSector, c.ySubSector).star
	sm.mapTitleText.ChangeText(system.GetName())

	loc := sm.systemLocations[sm.systemLocationsList.GetSelection()]
	sm.systemLocNameText.ChangeText(loc.GetName())

	switch loc.(type) {
	case Star:
		sm.systemLocDescText.ChangeText("This is a star. Stars are big hot balls of lava that float in space like magic.")
	case Planet:
		sm.systemLocDescText.ChangeText("This is a planet. Planets are rocks that are big enough to be important. Some planets have life on them, but most of them are super boring.")
	case *Ship:
		sm.systemLocDescText.ChangeText("This is your ship! Look at it's heroic hull valiantly floating amongst the stars. One could almost weep.")	
	default:
		sm.systemLocDescText.ChangeText("What is this? How did you select this?")
	}

	d := int(c.CalcVector(loc.GetCoords()).Distance * float64(METERS_PER_LY) / 1000)
	sm.systemLocDistText.ChangeText("Distance:" + strconv.Itoa(d) + "km.")
	if sm.playerShip.coords.IsIn(loc) {
		sm.systemSetCourseButton.ChangeText("We are currently here!")
	} else {
		sm.systemSetCourseButton.ChangeText("Lets go here!")
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
	w, h := sm.galaxy.Dims()
	for i := 0; i < w*h; i++ {
		x, y := i%w, i/w
		s := sm.galaxy.GetSector(x, y)
		bright := util.Lerp(0, 255, s.Density, 100)
		g := 0xB0
		if bright == 0 {
			g = 0
		}
		sm.mapView.Draw(x, y, g, console.MakeColour(bright, bright, bright), 0xFF000000)
	}

	x, y := sm.playerShip.coords.Sector()
	sm.mapView.Draw(x, y, 0x02, 0xFFFFFFFF, 0xFF000000)
}

func (sm *StarchartMenu) DrawSystem() {
	sm.mapView.Clear()
	c := sm.systemLocations[sm.systemLocationsList.GetSelection()].GetCoords()
	w, h := sm.mapView.Dims()
	gFactor := coord_LOCAL_MAX / w / util.Pow(2, sm.localZoom)

	xCamera := c.xLocal - (gFactor * w / 2)
	yCamera := c.yLocal - (gFactor * h / 2)

	system := sm.galaxy.GetSector(c.xSector, c.ySector).GetSubSector(c.xSubSector, c.ySubSector).star

	//draw system things!
	for _, p := range system.Planets {
		x, y := p.GetCoords().LocalCoord()
		if util.IsInside(x, y, xCamera, yCamera, gFactor*w, gFactor*h) {
			sm.mapView.Draw((x-xCamera)/gFactor, (y-yCamera)/gFactor, int(rune(p.name[0])), 0xFF825814, 0xFF000000)
		}
	}

	x, y := system.Star.GetCoords().LocalCoord()
	if util.IsInside(x, y, xCamera, yCamera, gFactor*w, gFactor*h) {
		sm.mapView.Draw((x-xCamera)/gFactor, (y-yCamera)/gFactor, 0x0F, 0xFFFFFF00, 0xFF000000)
	}

	//draw player ship last to ensure nothing overwrites it
	x, y = sm.playerShip.coords.LocalCoord()
	if util.IsInside(x, y, xCamera, yCamera, gFactor*w, gFactor*h) {
		sm.mapView.Draw((x-xCamera)/gFactor, (y-yCamera)/gFactor, 0x02, 0xFFFFFFFF, 0xFF000000)
	}
}

//ugh, hack incoming:
func (sm *StarchartMenu) MoveMapCursor(dx, dy int) {
	if sm.mapMode == coord_SECTOR {
		w, h := sm.mapView.Dims()
		if util.CheckBounds(sm.xCursor+dx, sm.yCursor+dy, w, h) {
			sm.xCursor += dx
			sm.yCursor += dy
			sm.mapHighlight.MoveTo(sm.xCursor, sm.yCursor)
			sm.Update()
		}
	} else if sm.mapMode == coord_LOCAL {
		if dy == 1 {
			sm.systemLocationsList.Next()
			sm.UpdateLocalInfo()
		} else if dy == -1 {
			sm.systemLocationsList.Prev()
			sm.UpdateLocalInfo()
		}
	}
}

func (sm *StarchartMenu) ZoomIn() {
	if sm.mapMode == coord_SECTOR {
		sm.mapMode = coord_LOCAL
		sm.localZoom = 0
		sm.sectorDetails.ToggleVisible()
		sm.systemDetails.ToggleVisible()
		sm.UpdateLocalInfo()
		sm.mapView.Clear()
		sm.mapHighlight.Toggle()
		sm.DrawMap()
	} else {
		sm.localZoom = util.Clamp(sm.localZoom+1, 0, 8)
		sm.mapView.Clear()
		sm.DrawMap()
	}
}

func (sm *StarchartMenu) ZoomOut() {
	if sm.mapMode == coord_LOCAL {
		if sm.localZoom == 0 {
			sm.mapMode = coord_SECTOR
			sm.sectorDetails.ToggleVisible()
			sm.systemDetails.ToggleVisible()
			sm.UpdateSectorInfo()
			sm.mapHighlight.Toggle()
			sm.DrawMap()
		} else {
			sm.localZoom -= 1
			sm.DrawMap()
		}
	}
}

func (sm *StarchartMenu) OnActivate() {
	sm.xCursor, sm.yCursor = sm.playerShip.coords.Sector()
	sm.mapHighlight.MoveTo(sm.xCursor, sm.yCursor)
	sm.Update()
	sm.DrawMap()
}
