package main

import (
	"github.com/bennicholls/burl-E/burl"
	"github.com/veandco/go-sdl2/sdl"
)

type StartMenu struct {
	burl.BaseState

	menu  *burl.List
	title *burl.Textbox

	background *burl.TileView
	stars      StarField
}

func NewStartMenu() (sm *StartMenu) {
	sm = new(StartMenu)

	sm.title = burl.NewTextbox(20, 1, 0, 10, 1, true, true, "SPACE SHIPPERS: The Ones Who Space Ship!")
	sm.title.CenterX(80, 0)

	sm.menu = burl.NewList(10, 5, 10, 10, 1, true, "")
	sm.menu.CenterInConsole()
	sm.menu.Append("New Game", "Load Game", "Ship Designer", "Options", "Quit")

	sm.background = burl.NewTileView(80, 45, 0, 0, 0, false)
	sm.stars = NewStarField(25, sm.background)
	return
}

func (sm *StartMenu) HandleKeypress(key sdl.Keycode) {
	switch key {
	case sdl.K_UP:
		sm.menu.Prev()
	case sdl.K_DOWN:
		sm.menu.Next()
	case sdl.K_RETURN:
		switch sm.menu.GetSelection() {
		case 0: //New Game
			burl.ChangeState(NewCreateGalaxyMenu())
		case 1: //Load Game
			//Load Game dialog
		case 2: //Ship Designer
			burl.ChangeState(NewShipDesignMenu())
		case 3: //Options
			//Options Dialog
		case 4: //Quit
			burl.PushEvent(burl.NewEvent(burl.EV_QUIT, ""))
		}
	case sdl.K_SPACE: //FOR TESTING PURPOSES ONLY DAMMIT
		g := NewGalaxy("Test Galaxy", GAL_MAX_RADIUS, GAL_DENSE)
		s := NewShip("The Greatest Spaceship There Is", g)
		temp, _ := LoadShipTemplate("raws/ship/Transport Ship.shp")
		s.SetupFromTemplate(temp)
		burl.ChangeState(NewSpaceshipGame(g, s))
	}
}

func (sm *StartMenu) Update() {
	sm.Tick++

	if sm.Tick%10 == 0 {
		sm.stars.Shift()
	}

}

func (sm *StartMenu) Render() {
	sm.stars.Draw()

	sm.title.Render()
	sm.menu.Render()
	sm.background.Render()
}
