package main

import "github.com/bennicholls/tyumi/input"

func (sg *SpaceshipGame) HandleKeypress(key_event *input.KeyboardEvent) (event_handled bool) {
	if key_event.PressType == input.KEY_RELEASED {
		return
	}

	switch key_event.Key {
	case input.K_F1:
		sg.gameMenuButton.Press()
		return true
	case input.K_F2:
		sg.shipMenuButton.Press()
		return true
	case input.K_F3:
		sg.galaxyMenuButton.Press()
		return true
	case input.K_F4:
		sg.crewMenuButton.Press()
		return true
	case input.K_F5:
		sg.commMenuButton.Press()
		return true
	case input.K_F6:
		sg.viewMenuButton.Press()
		return true
	case input.K_ESCAPE:
		sg.mainMenuButton.Press()
		return true
	case input.K_SPACE:
		sg.paused = !sg.paused
		if sg.paused {
			sg.AddLogMessage("Game Paused")
		} else {
			sg.AddLogMessage("Game Unpaused")
		}
		return true
	case input.K_KP_PLUS:
		if sg.simSpeed < 4 {
			sg.simSpeed++
			sg.timeDisplay.UpdateSpeed(sg.simSpeed)
			return true
		}
	case input.K_KP_MINUS:
		if sg.simSpeed > 0 {
			sg.simSpeed--
			sg.timeDisplay.UpdateSpeed(sg.simSpeed)
			return true
		}
	}

	if sg.activeMenu == nil {
		switch key_event.Key {
		case input.K_PAGEUP:
			sg.logOutput.Prev()
			return true
		case input.K_PAGEDOWN:
			sg.logOutput.Next() // NOTE: this is supposed to just scroll up, probably not cycle
			return true
		case input.K_HOME:
			sg.CenterShip()
			return true
		case input.K_UP:
			sg.shipView.MoveCamera(0, 1)
			return true
		case input.K_DOWN:
			sg.shipView.MoveCamera(0, -1)
			return true
		case input.K_LEFT:
			sg.shipView.MoveCamera(1, 0)
			return true
		case input.K_RIGHT:
			sg.shipView.MoveCamera(-1, 0)
			return true
		}
	}

	return
}
