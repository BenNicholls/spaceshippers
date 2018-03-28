package main

import (
	"github.com/bennicholls/burl-E/burl"
	"github.com/veandco/go-sdl2/sdl"
)

type MainMenu struct {
	burl.BaseState

	menu  *burl.List
	title *burl.Textbox

	background *burl.TileView
	stars      StarField
}

func NewMainMenu() (mm *MainMenu) {
	mm = new(MainMenu)

	mm.title = burl.NewTextbox(20, 1, 0, 10, 1, true, true, "SPACE SHIPPERS: The Ones Who Space Ship!")
	mm.title.CenterX(80, 0)

	mm.menu = burl.NewList(10, 5, 10, 10, 1, true, "")
	mm.menu.CenterInConsole()
	mm.menu.Append("New Game", "Load Game", "Ship Designer", "Options", "Quit")

	mm.background = burl.NewTileView(80, 45, 0, 0, 0, false)
	mm.stars = NewStarField(25, mm.background)
	return
}

func (mm *MainMenu) HandleKeypress(key sdl.Keycode) {
	switch key {
	case sdl.K_UP:
		mm.menu.Prev()
	case sdl.K_DOWN:
		mm.menu.Next()
	case sdl.K_RETURN:
		switch mm.menu.GetSelection() {
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
		s.SetupFromTemplate(defaultShipTemplates[SHIPTYPE_TRANSPORT])
		burl.ChangeState(NewSpaceshipGame(g, s))
	}
}

func (mm *MainMenu) Update() {
	mm.Tick++

	if mm.Tick%10 == 0 {
		mm.stars.Shift()
	}

}

func (mm *MainMenu) Render() {
	mm.stars.Draw()

	mm.title.Render()
	mm.menu.Render()
	mm.background.Render()
}
