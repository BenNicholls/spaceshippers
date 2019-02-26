package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

func (sg *SpaceshipGame) HandleKeypress(key sdl.Keycode) {
	//general keys -- works in all menus, modes, etc. Mainly menu switching stuff
	switch key {
	case sdl.K_F1:
		sg.ActivateMenu(MENU_GAME)
	case sdl.K_F2:
		sg.ActivateMenu(MENU_SHIP)
	case sdl.K_F3:
		// if sg.activeMenu != sg.starchartMenu {
		// 	sg.starchartMenu.OnActivate()
		// }
		sg.ActivateMenu(MENU_GALAXY)
	case sdl.K_F4:
		sg.ActivateMenu(MENU_CREW)
	case sdl.K_F5:
		sg.ActivateMenu(MENU_COMM)
	case sdl.K_F6:
		sg.ActivateMenu(MENU_VIEW)
	case sdl.K_ESCAPE:
		sg.ActivateMenu(MENU_MAIN)
	case sdl.K_KP_PLUS:
		if sg.simSpeed < 4 {
			sg.simSpeed++
			sg.timeDisplay.UpdateSpeed(sg.simSpeed)
		}
	case sdl.K_KP_MINUS:
		if sg.simSpeed > 0 {
			sg.simSpeed--
			sg.timeDisplay.UpdateSpeed(sg.simSpeed)
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
		case sg.crewMenu:
			sg.menus[MENU_CREW].HandleKeypress(key)
		case sg.shipMenu:
			sg.menus[MENU_SHIP].HandleKeypress(key)
		case sg.galaxyMenu:
			sg.menus[MENU_GALAXY].HandleKeypress(key)
		case sg.gameMenu:
			sg.menus[MENU_GAME].HandleKeypress(key)
		case sg.commMenu:
			sg.HandleKeypressCommMenu(key)
		case sg.viewMenu:
			sg.menus[MENU_VIEW].HandleKeypress(key)
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
			}
		}
	}
}
