package main

import (
	"fmt"

	"github.com/bennicholls/tyumi"
	"github.com/bennicholls/tyumi/gfx/ui"
	"github.com/bennicholls/tyumi/input"
	"github.com/bennicholls/tyumi/vec"
)

type GalaxyMenu struct {
	ui.PageContainer

	galaxy     *Galaxy
	playerShip *Ship

	starchartPage *ui.Page

	starchartMapView        GalaxyMapView
	starchartTitleText      ui.Textbox
	localInfo               ui.Element
	localLocationList       ui.List
	galaxyInfo              ui.Element
	selectedInfo            ui.Element
	selectedSetCourseButton ui.Button

	coursePage  *ui.Page
	scannerPage *ui.Page
}

func (gm *GalaxyMenu) Init(g *Galaxy, s *Ship) {
	gm.galaxy = g
	gm.playerShip = s

	gm.PageContainer.Init(menuSize, menuPos, menuDepth)
	gm.EnableBorder()
	gm.Hide()
	gm.AcceptInput = true

	gm.starchartPage = gm.CreatePage("Star Chart")
	pw, ph := gm.starchartPage.Size().W, gm.starchartPage.Size().H
	gm.starchartMapView.Init(vec.Dims{25, 25}, vec.ZERO_COORD, ui.BorderDepth, gm.galaxy)
	gm.starchartMapView.EnableBorder()
	gm.starchartMapView.SetCursor(gm.playerShip.Coords.Sector)
	gm.starchartMapView.systemFocus = g.GetStarSystem(gm.playerShip.GetCoords()) //NOTE: eventually be able to choose this!
	gm.starchartMapView.localFocus = gm.starchartMapView.systemFocus.Star
	gm.starchartMapView.AcceptInput = true
	gm.starchartMapView.OnCursorMove = func(new_pos vec.Coord) { gm.UpdateGalaxyInfo() }
	gm.starchartMapView.OnZoomTypeChange = gm.OnZoomChange
	gm.starchartMapView.AddMapObjectMarker(gm.playerShip)

	gm.starchartTitleText.Init(vec.Dims{25, 1}, vec.Coord{0, 26}, ui.BorderDepth, "", ui.JUSTIFY_CENTER)
	gm.starchartTitleText.EnableBorder()

	gm.localInfo.Init(vec.Dims{25, ph - 28}, vec.Coord{0, 28}, 0)
	gm.localInfo.Hide()
	gm.localLocationList.Init(vec.Dims{24, ph - 29}, vec.Coord{1, 1}, 0)
	gm.localLocationList.SetEmptyText("No Locations in this system.")
	gm.localLocationList.ToggleHighlight()
	gm.localLocationList.AcceptInput = true
	gm.localLocationList.OnChangeSelection = gm.OnLocationChange
	gm.localInfo.AddChild(ui.NewTextbox(vec.Dims{25, 1}, vec.ZERO_COORD, 0, "Locations:", ui.JUSTIFY_LEFT))
	gm.localInfo.AddChild(&gm.localLocationList)

	gm.galaxyInfo.Init(vec.Dims{25, ph - 28}, vec.Coord{0, 28}, 0)

	gm.selectedInfo.Init(vec.Dims{pw - 26, ph - 2}, vec.Coord{26, 0}, ui.BorderDepth)
	gm.selectedInfo.EnableBorder()
	gm.selectedSetCourseButton.Init(vec.Dims{pw - 26, 1}, vec.Coord{26, ph - 1}, ui.BorderDepth, "[S]et Course for this Location!", nil)
	gm.selectedSetCourseButton.EnableBorder()
	gm.selectedSetCourseButton.OnPressCallback = gm.OpenSetCourseDialog

	gm.UpdateGalaxyInfo()

	gm.starchartPage.AddChildren(&gm.starchartMapView, &gm.localInfo, &gm.galaxyInfo, &gm.selectedInfo, &gm.starchartTitleText, &gm.selectedSetCourseButton)

	gm.coursePage = gm.CreatePage("Course")
	gm.scannerPage = gm.CreatePage("Scanners")

	return
}

func (gm *GalaxyMenu) HandleKeypress(key_event *input.KeyboardEvent) (event_handled bool) {
	if key_event.PressType != input.KEY_PRESSED {
		return
	}

	switch gm.GetPageIndex() {
	case 0: //map view
		switch key_event.Key {
		case input.K_s:
			gm.selectedSetCourseButton.Press()
		}
	}

	return
}

func (gm *GalaxyMenu) OpenSetCourseDialog() {
	locations := gm.starchartMapView.systemFocus.GetLocations()
	l := locations[gm.localLocationList.GetSelectionIndex()]
	if l != gm.playerShip && l != gm.playerShip.currentLocation {
		tyumi.OpenDialog(NewSetCourseDialog(gm.playerShip, l, gm.galaxy.spaceTime))
	}

}

func (gm *GalaxyMenu) OnZoomChange(level zoomLevel) {
	switch level {
	case zoom_GALAXY:
		gm.UpdateGalaxyInfo()
		gm.localInfo.Hide()
		gm.galaxyInfo.Show()
	case zoom_LOCAL:
		gm.galaxyInfo.Hide()
		gm.localInfo.Show()
		gm.UpdateLocalInfo()
	}
}

func (gm *GalaxyMenu) OnLocationChange() {
	if gm.localLocationList.Count() > 0 {
		locations := gm.starchartMapView.systemFocus.GetLocations()
		gm.starchartMapView.localFocus = locations[gm.localLocationList.GetSelectionIndex()]
		gm.starchartMapView.Updated = true
		gm.starchartMapView.markerLayer.Updated = true
		gm.UpdateSelectedInfo()
		gm.Updated = true
	}
}

func (gm *GalaxyMenu) UpdateGalaxyInfo() {
	gm.starchartTitleText.ChangeText(gm.galaxy.name)
	gm.UpdateSelectedInfo()
}

func (gm *GalaxyMenu) UpdateLocalInfo() {
	gm.starchartTitleText.ChangeText(gm.starchartMapView.systemFocus.GetName())
	gm.localLocationList.RemoveAll()
	for _, l := range gm.starchartMapView.systemFocus.GetLocations() {
		gm.localLocationList.InsertText(ui.JUSTIFY_LEFT, l.GetName())
	}

	gm.UpdateSelectedInfo()
}

func (gm *GalaxyMenu) UpdateSelectedInfo() {
	gm.selectedInfo.RemoveAllChildren()

	gm.selectedInfo.AddChild(ui.NewTitleTextbox(vec.Dims{30, 1}, vec.ZERO_COORD, ui.BorderDepth, "Selection Info"))

	switch gm.starchartMapView.zoom {
	case zoom_GALAXY:
		s := gm.galaxy.GetSector(gm.starchartMapView.cursor)
		gm.selectedInfo.AddChild(ui.NewTextbox(vec.Dims{30, 1}, vec.Coord{0, 2}, 0, "Sector ("+s.ProperName()+")", ui.JUSTIFY_CENTER))
		gm.selectedInfo.AddChild(ui.NewTextbox(vec.Dims{30, 1}, vec.Coord{0, 3}, 0, s.Name, ui.JUSTIFY_CENTER))
		gm.selectedInfo.AddChild(ui.NewTextbox(vec.Dims{30, 2}, vec.Coord{0, 5}, 0, s.Description, ui.JUSTIFY_CENTER))
		gm.selectedInfo.AddChild(ui.NewTextbox(vec.Dims{30, 2}, vec.Coord{0, 8}, 0, fmt.Sprintf("Star Density: %d%%", s.Density), ui.JUSTIFY_LEFT))
		if s.Explored {
			gm.selectedInfo.AddChild(ui.NewTextbox(vec.Dims{30, 2}, vec.Coord{0, 9}, 0, "SECTOR EXPLORED", ui.JUSTIFY_LEFT))
		} else {
			gm.selectedInfo.AddChild(ui.NewTextbox(vec.Dims{30, 2}, vec.Coord{0, 9}, 0, "SECTOR UNEXPLORED", ui.JUSTIFY_LEFT))
		}

		gm.selectedSetCourseButton.DisablePress = true
		if gm.playerShip.GetCoords().IsIn(s) {
			gm.selectedSetCourseButton.ChangeText("We are currently here!")
		} else {
			gm.selectedSetCourseButton.ChangeText("Cannot travel to different sectors. :(")
		}

	case zoom_LOCAL:
		s := gm.starchartMapView.localFocus
		gm.selectedInfo.AddChild(ui.NewTextbox(vec.Dims{30, 1}, vec.Coord{0, 2}, 0, "Name: "+s.GetName(), ui.JUSTIFY_LEFT))
		gm.selectedInfo.AddChild(ui.NewTextbox(vec.Dims{30, 1}, vec.Coord{0, 3}, 0, "Type: "+s.GetLocationType().String(), ui.JUSTIFY_LEFT))
		dist := int(gm.playerShip.GetCoords().CalcVector(s.GetCoords()).Distance * METERS_PER_LY / 1000)
		gm.selectedInfo.AddChild(ui.NewTextbox(vec.Dims{30, 1}, vec.Coord{0, 4}, 0, fmt.Sprintf("Distance: %dkm", dist), ui.JUSTIFY_LEFT))
		gm.selectedInfo.AddChild(ui.NewTextbox(vec.Dims{30, 4}, vec.Coord{0, 6}, 0, s.GetDescription(), ui.JUSTIFY_LEFT))

		var explored string
		if s.IsExplored() {
			explored = "We have explored this location."
		} else {
			explored = "We have not explored this location."
		}
		gm.selectedInfo.AddChild(ui.NewTextbox(vec.Dims{30, 2}, vec.Coord{0, 11}, 0, explored, ui.JUSTIFY_LEFT))

		if gm.playerShip.GetCoords().IsIn(s) {
			gm.selectedSetCourseButton.ChangeText("We are currently here!")
			gm.selectedSetCourseButton.DisablePress = true
		} else {
			gm.selectedSetCourseButton.DisablePress = false
			if gm.playerShip.destination == s {
				gm.selectedSetCourseButton.ChangeText("Currently on course. [S]et new course?")
			} else {
				gm.selectedSetCourseButton.ChangeText("[S]et course for this location!")
			}
		}
	}
}
