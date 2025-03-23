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

	title := ui.NewTitleTextbox(vec.Dims{ui.FIT_TEXT, 1}, vec.Coord{0, 10}, 1, "SPACE SHIPPERS: The Ones Who Space Ship!")
	sm.Window().AddChild(title)
	title.CenterHorizontal()

	sm.menu.Init(vec.Dims{16, 5}, vec.Coord{0, 0}, 1)
	sm.menu.EnableBorder()
	sm.Window().AddChild(&sm.menu)
	sm.menu.Center()
	sm.menu.InsertText(ui.JUSTIFY_CENTER, "New Game", "Load Game", "Ship Designer", "Options", "Quit")
	sm.menu.EnableHighlight()
	sm.menu.Focus()

	sm.stars.Init(vec.Dims{96, 54}, vec.ZERO_COORD, 0, 25, 10)
	sm.Window().AddChildren(&sm.stars)

	sm.SetKeypressHandler(sm.HandleKeypress)
}

func (sm *StartMenu) HandleKeypress(key_event *input.KeyboardEvent) (event_handled bool) {
	if key_event.Handled() {
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
			tyumi.ChangeState(NewShipDesignMenu())
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

	return
}
