package main

import "strconv"
import "github.com/bennicholls/burl/ui"
import "github.com/bennicholls/burl/console"
import "github.com/bennicholls/burl/util"

//Menu for viewing star charts, getting location data, setting courses, etc.
type StarchartMenu struct {
	ui.Container

	galaxynameText *ui.Textbox
	mapView        *ui.TileView
	mapHighlight   *ui.PulseAnimation

	sectorDetails      *ui.Container
	sectorCoordsText   *ui.Textbox
	sectorNameText     *ui.Textbox
	sectorDensityText  *ui.Textbox
	sectorExploredText *ui.Textbox
	sectorKnownText    *ui.Textbox
	sectorLocationText *ui.Textbox

	selectedSector Coordinates

	galaxy     *Galaxy //to know what we're drawing
	playerShip *Ship   //to know where we are, current course, etc
}

func NewStarchartMenu(gal *Galaxy, ship *Ship) (sm *StarchartMenu) {
	sm = new(StarchartMenu)
	sm.galaxy = gal
	sm.playerShip = ship
	sm.selectedSector = ship.Location.GetCoords() //start sector picker on player ship

	//ui setup
	sm.Container = *ui.NewContainer(40, 26, 39, 4, 1, true)
	sm.SetTitle("Star Charts")
	sm.SetVisibility(false)

	sm.galaxynameText = ui.NewTextbox(25, 1, 0, 25, 1, false, true, sm.galaxy.name)
	sm.mapView = ui.NewTileView(25, 25, 0, 0, 1, false)

	sm.sectorDetails = ui.NewContainer(15, 26, 25, 0, 1, false)
	sm.sectorCoordsText = ui.NewTextbox(15, 1, 0, 0, 1, false, true, "")
	sm.sectorNameText = ui.NewTextbox(15, 2, 0, 1, 1, false, true, "")
	sm.sectorDensityText = ui.NewTextbox(15, 1, 0, 4, 1, false, false, "")
	sm.sectorExploredText = ui.NewTextbox(15, 1, 0, 5, 1, false, false, "")
	sm.sectorLocationText = ui.NewTextbox(15, 1, 0, 6, 1, false, false, "")
	sm.sectorKnownText = ui.NewTextbox(15, 2, 0, 8, 1, false, true, "We know nothing about this sector.")
	sm.sectorDetails.Add(sm.sectorCoordsText, sm.sectorNameText, sm.sectorDensityText, sm.sectorExploredText, sm.sectorLocationText, sm.sectorKnownText)

	x, y := sm.selectedSector.Sector()
	sm.mapHighlight = ui.NewPulseAnimation(x, y, 1, 1, 50, 10, true)
	sm.mapHighlight.Activate()
	sm.mapView.AddAnimation(sm.mapHighlight)
	sm.Add(sm.mapView, sm.galaxynameText, sm.sectorDetails)

	sm.DrawMap()
	sm.UpdateSectorInfo()

	return
}

func (sm *StarchartMenu) UpdateSectorInfo() {
	sx, sy := sm.selectedSector.Sector()
	sector := sm.galaxy.GetSector(sx, sy)
	sm.sectorCoordsText.ChangeText("Sector (" + sm.galaxy.GetSector(sx, sy).ProperName() + ")")
	sm.sectorNameText.ChangeText("\"" + sector.GetName() + "\"")
	sm.sectorDensityText.ChangeText("Star Density: " + strconv.Itoa(sector.Density) + "%")
	if sector.IsExplored() {
		sm.sectorExploredText.ChangeText("SECTOR EXPLORED!")
	} else {
		sm.sectorExploredText.ChangeText("SECTOR UNEXPLORED")
	}

	if x, y := sm.playerShip.Location.GetCoords().Sector(); sx == x && sy == y {
		sm.sectorLocationText.ChangeText("We're currently here!")
	} else if sm.playerShip.Destination != nil {
		if x, y := sm.playerShip.Destination.GetCoords().Sector(); sx == x && sy == y {
			sm.sectorLocationText.ChangeText("We're currently going here!")
		}
	} else {
		//TODO: could have code here saying distance to sector, estimated travel time, etc.
		sm.sectorLocationText.ChangeText("We could go here!")
	}

	sm.mapHighlight.MoveTo(sx, sy)
}

//draws the required map. galaxy map, sector map, star system map
func (sm *StarchartMenu) DrawMap() {
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

	x, y := sm.playerShip.Location.GetCoords().Sector()
	sm.mapView.Draw(x, y, 0x02, 0xFFFFFFFF, 0xFF000000)
}

func (sm *StarchartMenu) MoveSectorCursor(dx, dy int) {
	w, h := sm.mapView.Dims()
	sx, sy := sm.selectedSector.Sector()
	if util.CheckBounds(sx+dx, sy+dy, w, h) {
		sm.selectedSector.Move(dx, dy, coord_SECTOR)
		sm.UpdateSectorInfo()
	}
}
