package main

import (
	//"strconv"

	"github.com/bennicholls/burl-E/burl"
	"github.com/veandco/go-sdl2/sdl"
)

type GalaxyMenu struct {
	burl.PagedContainer

	galaxy     *Galaxy
	playerShip *Ship

	starchartPage           *burl.Container
	//starchartMapView        *GalaxyMapView
	starchartTitleText      *burl.Textbox
	localInfo               *burl.Container
	localLocationList       *burl.List
	galaxyInfo              *burl.Container
	selectedInfo            *burl.Container
	selectedSetCourseButton *burl.Button

	coursePage  *burl.Container
	scannerPage *burl.Container
}

func NewGalaxyMenu(g *Galaxy, s *Ship) (gm *GalaxyMenu) {
	gm = new(GalaxyMenu)
	gm.galaxy = g
	gm.playerShip = s

	gm.PagedContainer = *burl.NewPagedContainer(56, 45, 39, 4, 10, true)
	gm.SetVisibility(false)

	pw, ph := gm.GetPageDims()

	gm.starchartPage = gm.AddPage("Star Chart")
	// gm.starchartMapView = NewGalaxyMapView(25, 0, 0, 0, true, gm.galaxy)
	// gm.starchartMapView.cursor = gm.playerShip.Coords.Sector
	// gm.starchartMapView.ToggleHighlight()
	// gm.starchartMapView.systemFocus = g.GetStarSystem(gm.playerShip.GetCoords()) //NOTE: eventually be able to choose this!
	// gm.starchartMapView.localFocus = gm.starchartMapView.systemFocus.Star
	gm.starchartTitleText = burl.NewTextbox(25, 1, 0, 26, 0, true, true, "")
	gm.localInfo = burl.NewContainer(25, ph-28, 0, 28, 0, false)
	gm.localInfo.SetVisibility(false)
	gm.localLocationList = burl.NewList(24, ph-29, 1, 1, 0, false, "No Locations in this system.")
	gm.localInfo.Add(burl.NewTextbox(25, 1, 0, 0, 0, false, false, "Locations:"))
	gm.localInfo.Add(gm.localLocationList)
	gm.galaxyInfo = burl.NewContainer(25, ph-28, 0, 28, 0, false)
	gm.selectedInfo = burl.NewContainer(pw-26, ph-2, 26, 0, 0, true)
	gm.selectedSetCourseButton = burl.NewButton(pw-26, 1, 26, ph-1, 0, true, true, "[S]et Course for this location!")

	gm.Update()

	// gm.starchartPage.Add(gm.starchartMapView, gm.localInfo, gm.galaxyInfo, gm.selectedInfo, gm.starchartTitleText, gm.selectedSetCourseButton)

	gm.coursePage = gm.AddPage("Course")
	gm.scannerPage = gm.AddPage("Scanners")

	return
}

func (gm *GalaxyMenu) HandleKeypress(key sdl.Keycode) {
	gm.PagedContainer.HandleKeypress(key)

	switch gm.PagedContainer.CurrentIndex() {
	case 0: //map view
		// gm.starchartMapView.HandleKeypress(key)
		// if gm.starchartMapView.zoom == zoom_LOCAL {
		// 	switch key {
		// 	case sdl.K_DOWN, sdl.K_UP:
		// 		if len(gm.localLocationList.Elements) > 0 {
		// 			gm.localLocationList.HandleKeypress(key)
		// 			locations := gm.starchartMapView.systemFocus.GetLocations()
		// 			gm.starchartMapView.localFocus = locations[gm.localLocationList.GetSelection()]
		// 			gm.UpdateSelectedInfo()
		// 			gm.DrawMap()
		// 		}
		// 	case sdl.K_s:
		// 		gm.selectedSetCourseButton.Press()
		// 		locations := gm.starchartMapView.systemFocus.GetLocations()
		// 		l := locations[gm.localLocationList.GetSelection()]
		// 		if l != gm.playerShip && l != gm.playerShip.currentLocation {
		// 			burl.OpenDialog(NewSetCourseDialog(gm.playerShip, l, gm.galaxy.spaceTime))
		// 		}
		// 	}
		// }
		gm.Update()
	}
}

func (gm *GalaxyMenu) DrawMap() {
	// gm.starchartMapView.DrawMap()

	// switch gm.starchartMapView.zoom {
	// case zoom_GALAXY:
	// 	//TODO: make it so this is only drawn if you've scanned/discovered the location of earth
	// 	gm.starchartMapView.DrawMapMarker(gm.galaxy.GetEarth().GetCoords().Sector, burl.Visuals{
	// 		Glyph:      burl.GLYPH_DONUT,
	// 		ForeColour: burl.COL_BLUE,
	// 		BackColour: burl.COL_BLACK})
	// case zoom_LOCAL:
	// }

	// gm.starchartMapView.DrawMapObject(gm.playerShip)
}

func (gm *GalaxyMenu) Update() {
	// switch gm.CurrentIndex() {
	// case 0: //starcharts
	// 	switch gm.starchartMapView.zoom {
	// 	case zoom_GALAXY:
	// 		gm.UpdateGalaxyInfo()
	// 	case zoom_LOCAL:
	// 		gm.UpdateLocalInfo()
	// 	}

	// 	gm.DrawMap()
	// }
}

func (gm *GalaxyMenu) UpdateGalaxyInfo() {
	gm.localInfo.SetVisibility(false)
	gm.galaxyInfo.SetVisibility(true)

	gm.starchartTitleText.ChangeText(gm.galaxy.name)

	gm.UpdateSelectedInfo()
}

func (gm *GalaxyMenu) UpdateLocalInfo() {
	// gm.galaxyInfo.SetVisibility(false)
	// gm.localInfo.SetVisibility(true)

	// gm.starchartTitleText.ChangeText(gm.starchartMapView.systemFocus.GetName())
	// gm.localLocationList.ClearElements()
	// for _, l := range gm.starchartMapView.systemFocus.GetLocations() {
	// 	gm.localLocationList.Append(l.GetName())
	// }

	// gm.UpdateSelectedInfo()
}

func (gm *GalaxyMenu) UpdateSelectedInfo() {
	// gm.selectedInfo.ClearElements()

	// gm.selectedInfo.Add(burl.NewTextbox(30, 1, 0, 0, 0, true, true, "Selection Info"))

	// switch gm.starchartMapView.zoom {
	// case zoom_GALAXY:
	// 	s := gm.galaxy.GetSector(gm.starchartMapView.cursor.Get())
	// 	gm.selectedInfo.Add(burl.NewTextbox(30, 1, 0, 2, 0, false, true, "Sector ("+s.ProperName()+")"))
	// 	gm.selectedInfo.Add(burl.NewTextbox(30, 1, 0, 3, 0, false, true, s.Name))
	// 	gm.selectedInfo.Add(burl.NewTextbox(30, 2, 0, 5, 0, false, true, s.Description))
	// 	gm.selectedInfo.Add(burl.NewTextbox(30, 2, 0, 8, 0, false, false, "Star Density: "+strconv.Itoa(s.Density)+"%"))
	// 	if s.Explored {
	// 		gm.selectedInfo.Add(burl.NewTextbox(30, 2, 0, 9, 0, false, false, "SECTOR EXPLORED"))
	// 	} else {
	// 		gm.selectedInfo.Add(burl.NewTextbox(30, 2, 0, 9, 0, false, false, "SECTOR UNEXPLORED"))
	// 	}

	// 	if gm.playerShip.GetCoords().IsIn(s) {
	// 		gm.selectedSetCourseButton.ChangeText("We are currently here!")
	// 	} else {
	// 		gm.selectedSetCourseButton.ChangeText("Cannot travel to different sectors. :(")
	// 	}

	// case zoom_LOCAL:
	// 	s := gm.starchartMapView.localFocus
	// 	gm.selectedInfo.Add(burl.NewTextbox(30, 1, 0, 2, 0, false, false, "Name: "+s.GetName()))
	// 	gm.selectedInfo.Add(burl.NewTextbox(30, 1, 0, 3, 0, false, false, "Type: "+s.GetLocationType().String()))
	// 	dist := int(gm.playerShip.GetCoords().CalcVector(s.GetCoords()).Distance * METERS_PER_LY / 1000)
	// 	gm.selectedInfo.Add(burl.NewTextbox(30, 1, 0, 4, 0, false, false, "Distance: "+strconv.Itoa(dist)+"km"))
	// 	gm.selectedInfo.Add(burl.NewTextbox(30, 4, 0, 6, 0, false, false, s.GetDescription()))

	// 	var explored string
	// 	if s.IsExplored() {
	// 		explored = "We have explored this location."
	// 	} else {
	// 		explored = "We have not explored this location."
	// 	}
	// 	gm.selectedInfo.Add(burl.NewTextbox(30, 2, 0, 11, 0, false, false, explored))

	// 	gm.selectedSetCourseButton.SetVisibility(true)
	// 	if !gm.playerShip.GetCoords().IsIn(s) {
	// 		if gm.playerShip.destination == s {
	// 			gm.selectedSetCourseButton.ChangeText("Currently on course. [S]et new course?")
	// 		} else {
	// 			gm.selectedSetCourseButton.ChangeText("[S]et course for this location!")
	// 		}
	// 	} else {
	// 		gm.selectedSetCourseButton.ChangeText("We are currently here.")
	// 	}
	// }
}


// DEAR FUTURE BEN: THIS WAS ALREADY COMMENTED OUT BEFORE FOR SOME REASON, DON'T FEEL LIKE YOU HAVE TO REIMPLEMENT THIS UNTIL YOU LOOK INTO IT
// func (sm *StarchartMenu) OnActivate() {
// 	sm.xCursor, sm.yCursor = sm.playerShip.Coords.Sector.Get()
// 	sm.mapHighlight.MoveTo(sm.xCursor, sm.yCursor)
// 	sm.Update()
// 	sm.DrawMap()
// }
