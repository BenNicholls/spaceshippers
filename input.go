package main

import "github.com/veandco/go-sdl2/sdl"
import "github.com/bennicholls/burl-E/burl"

func (sg *SpaceshipGame) HandleKeypress(key sdl.Keycode) {

	//dialogs have the highest priority and they handle their own input
	if sg.dialog != nil {
		sg.dialog.HandleInput(key)
		return
	}

	//general keys -- works in all menus, modes, etc. Mainly menu switching stuff
	switch key {
	case sdl.K_F1:
		sg.crewMenuButton.Press()
		sg.ActivateMenu(sg.crewMenu)
	case sdl.K_F2:
		sg.shipMenuButton.Press()
		sg.ActivateMenu(sg.shipMenu)
	case sdl.K_F3:
		sg.missionMenuButton.Press()
		sg.ActivateMenu(sg.missionMenu)
	case sdl.K_F4:
		sg.starchartMenuButton.Press()
		if sg.activeMenu != sg.starchartMenu {
			sg.starchartMenu.OnActivate()
		}
		sg.ActivateMenu(sg.starchartMenu)
	case sdl.K_KP_PLUS:
		if sg.simSpeed < 4 {
			sg.simSpeed++
			sg.UpdateSpeedUI()
		}
	case sdl.K_KP_MINUS:
		if sg.simSpeed > 0 {
			sg.simSpeed--
			sg.UpdateSpeedUI()
		}
	case sdl.K_SPACE:
		sg.paused = !sg.paused
		if sg.paused {
			sg.AddMessage("Game Paused")
		} else {
			sg.AddMessage("Game Unpaused")
		}
	default:
		//Check for active menus. If nothing, apply to base game.
		switch sg.activeMenu {
		case sg.input:
			sg.HandleKeypressInput(key)
		case sg.crewMenu:
			sg.HandleKeypressCrewMenu(key)
		case sg.starchartMenu:
			sg.HandleKeypressStarchartMenu(key)
		case sg.missionMenu:
			sg.HandleKeypressMissionMenu(key)
		default:
			switch key {
			case sdl.K_PAGEUP:
				sg.output.ScrollUp()
			case sdl.K_PAGEDOWN:
				sg.output.ScrollDown()
			case sdl.K_HOME:
				sg.CenterShip()
			case sdl.K_UP:
				sg.MoveShipCamera(0, -1)
			case sdl.K_DOWN:
				sg.MoveShipCamera(0, 1)
			case sdl.K_LEFT:
				sg.MoveShipCamera(-1, 0)
			case sdl.K_RIGHT:
				sg.MoveShipCamera(1, 0)
			case sdl.K_ESCAPE:
				sg.ActivateMenu(sg.input)
			}
		}
	}
}

func (sg *SpaceshipGame) HandleKeypressInput(key sdl.Keycode) {
	if burl.ValidText(rune(key)) {
		sg.input.InsertText(rune(key))
	} else {
		switch key {
		case sdl.K_RETURN:
			sg.Execute()
			sg.input.Reset()
			sg.DeactivateMenu()
		case sdl.K_BACKSPACE:
			sg.input.Delete()
		case sdl.K_SPACE:
			sg.input.Insert(" ")
		case sdl.K_ESCAPE:
			sg.DeactivateMenu()
			sg.input.Reset()
		}
	}
}

func (sg *SpaceshipGame) HandleKeypressCrewMenu(key sdl.Keycode) {
	switch key {
	case sdl.K_UP:
		sg.crewMenu.crewList.Prev()
	case sdl.K_DOWN:
		sg.crewMenu.crewList.Next()
	case sdl.K_RETURN:
		sg.crewMenu.ToggleCrewDetails()
	}
}

func (sg *SpaceshipGame) HandleKeypressStarchartMenu(key sdl.Keycode) {
	switch key {
	case sdl.K_UP:
		sg.starchartMenu.MoveMapCursor(0, -1)
	case sdl.K_DOWN:
		sg.starchartMenu.MoveMapCursor(0, 1)
	case sdl.K_LEFT:
		sg.starchartMenu.MoveMapCursor(-1, 0)
	case sdl.K_RIGHT:
		sg.starchartMenu.MoveMapCursor(1, 0)
	case sdl.K_PAGEUP:
		sg.starchartMenu.ZoomIn()
	case sdl.K_PAGEDOWN:
		sg.starchartMenu.ZoomOut()
	case sdl.K_RETURN:
		if sg.starchartMenu.mapMode == coord_LOCAL {
			sg.starchartMenu.systemSetCourseButton.Press()
			l := sg.starchartMenu.systemLocations[sg.starchartMenu.systemLocationsList.GetSelection()]
			if l != sg.playerShip && l != sg.playerShip.CurrentLocation {
				sg.dialog = NewSetCourseDialog(sg.playerShip, l, sg.spaceTime)
			}
		}
	}
}

func (sg *SpaceshipGame) HandleKeypressMissionMenu(key sdl.Keycode) {
	switch key {
	case sdl.K_UP:
		sg.missionMenu.missionList.Prev()
		sg.missionMenu.Update()
	case sdl.K_DOWN:
		sg.missionMenu.missionList.Next()
		sg.missionMenu.Update()
	}
}
