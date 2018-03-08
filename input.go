package main

import (
	"github.com/bennicholls/burl-E/burl"
	"github.com/veandco/go-sdl2/sdl"
)

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
	case sdl.K_F5:
		sg.commsMenuButton.Press()
		sg.ActivateMenu(sg.commsMenu)
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
		case sg.commsMenu:
			sg.HandleKeypressCommsMenu(key)
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
	sg.input.HandleKeypress(key)

	switch key {
	case sdl.K_RETURN:
		sg.Execute()
		sg.input.Reset()
		sg.DeactivateMenu()
	case sdl.K_ESCAPE:
		sg.DeactivateMenu()
		sg.input.Reset()
	}
}

func (sg *SpaceshipGame) HandleKeypressCrewMenu(key sdl.Keycode) {
	sg.crewMenu.crewList.HandleKeypress(key)

	if key == sdl.K_RETURN {
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
			if l != sg.playerShip && l != sg.playerShip.currentLocation {
				sg.dialog = NewSetCourseDialog(sg.playerShip, l, sg.GetTime())
			}
		}
	}
}

func (sg *SpaceshipGame) HandleKeypressMissionMenu(key sdl.Keycode) {
	sg.missionMenu.missionList.HandleKeypress(key)

	switch key {
	case sdl.K_UP, sdl.K_DOWN:
		burl.PushEvent(burl.NewEvent(burl.EV_UPDATE_UI, "missions"))
	}
}

func (sg *SpaceshipGame) HandleKeypressCommsMenu(key sdl.Keycode) {
	sg.commsMenu.HandleKeypress(key)

	switch sg.commsMenu.CurrentIndex() {
	case 0: //Inbox
		sg.commsMenu.inboxList.HandleKeypress(key)
		if key == sdl.K_RETURN && len(sg.commsMenu.comms.Inbox) > 0 {
			s := sg.commsMenu.inboxList.GetSelection()
			msg := sg.commsMenu.comms.Inbox[s]
			sg.dialog = NewCommDialog(msg.sender.Name, "You", msg.sender.Pic, msg.message)
		}
	case 2: //Transmissions
		sg.commsMenu.transmissionsList.HandleKeypress(key)
		if key == sdl.K_RETURN && len(sg.commsMenu.comms.Transmissions) > 0 {
			s := sg.commsMenu.transmissionsList.GetSelection()
			msg := sg.commsMenu.comms.Transmissions[s]
			sg.dialog = NewCommDialog(msg.sender.Name, "You", msg.sender.Pic, msg.message)
		}
	}
}
