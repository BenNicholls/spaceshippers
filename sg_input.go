package main

import "github.com/veandco/go-sdl2/sdl"
import "github.com/bennicholls/burl/util"

func (sg *SpaceshipGame) HandleKeypress(key sdl.Keycode) {
	//if we're inputting commands to scippie
	if sg.activeMenu == sg.input {
		if util.ValidText(rune(key)) {
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
	} else {
		switch key {
		case sdl.K_PAGEUP:
			sg.output.ScrollUp()
		case sdl.K_PAGEDOWN:
			sg.output.ScrollDown()
		case sdl.K_HOME:
			sg.CenterShip()
		case sdl.K_KP_PLUS:
			if sg.SimSpeed < 4 {
				sg.SimSpeed++
				sg.UpdateSpeedUI()
			}
		case sdl.K_KP_MINUS:
			if sg.SimSpeed > 0 {
				sg.SimSpeed--
				sg.UpdateSpeedUI()
			}
		case sdl.K_UP:
			sg.viewY -= 1
		case sdl.K_DOWN:
			sg.viewY += 1
		case sdl.K_LEFT:
			sg.viewX -= 1
		case sdl.K_RIGHT:
			sg.viewX += 1
		case sdl.K_F1:
			if sg.activeMenu == sg.crew {
				sg.DeactivateMenu()
			} else {
				sg.ActivateMenu(sg.crew)
			}
		case sdl.K_SPACE:
			sg.paused = !sg.paused
			if sg.paused {
				sg.AddMessage("Game Paused")
			} else {
				sg.AddMessage("Game Unpaused")
			}
		case sdl.K_ESCAPE:
			sg.ActivateMenu(sg.input)
		}
	}
}