package main

import (
	"github.com/bennicholls/tyumi"
	"github.com/bennicholls/tyumi/event"
	"github.com/bennicholls/tyumi/gfx/ui"
	"github.com/bennicholls/tyumi/input"
	"github.com/bennicholls/tyumi/vec"
)

type StartMenu struct {
	tyumi.State

	menu  ui.List
	stars StarField
}

func (sm *StartMenu) Init() {
	sm.State.Init()

	title := ui.Textbox{}
	title.Init(vec.Dims{ui.FIT_TEXT, 1}, vec.Coord{0, 10}, 1, "SPACE SHIPPERS: The Ones Who Space Ship!", ui.JUSTIFY_CENTER)
	title.EnableBorder()
	sm.Window().AddChild(&title)
	title.CenterHorizontal()

	sm.menu.Init(vec.Dims{16, 5}, vec.Coord{0, 0}, 1)
	sm.menu.EnableBorder()
	sm.Window().AddChild(&sm.menu)
	sm.menu.Center()
	sm.menu.AddTextItems(ui.JUSTIFY_CENTER, "New Game", "Load Game", "Ship Designer", "Options", "Quit")
	sm.menu.EnableHighlight()
	sm.menu.Focus()

	sm.stars.Init(vec.Dims{96, 54}, vec.ZERO_COORD, 0, 25, 10)
	sm.Window().AddChildren(&sm.stars)

	sm.SetInputHandler(sm.HandleInput)
}

func (sm *StartMenu) HandleInput(input_event event.Event) (event_handled bool) {
	if input_event.Handled() {
		return
	}

	if input_event.ID() == input.EV_KEYBOARD {
		key_event := input_event.(*input.KeyboardEvent)
		if key_event.PressType == input.KEY_RELEASED {
			return
		}

		switch key_event.Key {
		case input.K_RETURN:
			switch sm.menu.GetSelectionIndex() {
			case 0: //New Game
				tyumi.ChangeState(NewCreateGalaxyMenu())
			case 1: //Load Game
				//Load Game dialog
			case 2: //Ship Designer
				//burl.ChangeState(NewShipDesignMenu())
			case 3: //Options
				//Options Dialog
			case 4: //Quit
				event.FireSimple(tyumi.EV_QUIT)
			}
		case input.K_SPACE: //FOR TESTING PURPOSES ONLY DAMMIT
			// g := NewGalaxy("Test Galaxy", GAL_MAX_RADIUS, GAL_DENSE)
			// s := NewShip("The Greatest Spaceship There Is", g)
			// temp, _ := LoadShipTemplate("raws/ship/Transport Ship.shp")
			// s.SetupFromTemplate(temp)
			// burl.ChangeState(NewSpaceshipGame(g, s))
		}
	}

	return
}
