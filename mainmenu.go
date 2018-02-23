package main

import (
	"github.com/bennicholls/burl-E/burl"
	"github.com/veandco/go-sdl2/sdl"
)

type MainMenu struct {
	menu *burl.List
}

func NewMainMenu() (mm *MainMenu) {
	mm = new(MainMenu)
	mm.menu = burl.NewList(10, 5, 10, 10, 1, true, "")
	mm.menu.CenterInConsole()

	mm.menu.Append("New Game", "Load Game", "Ship Designer", "Options", "Quit")
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
			//Galaxy Generator Dialog/state (need to decide)
		case 1: //Load Game
			//Load Game dialog
		case 2: //Ship Designer
			//Not sure if this one stays in.
		case 3: //Options
			//Options Dialog
		case 4: //Quit
			burl.PushEvent(burl.NewEvent(burl.QUIT_EVENT, ""))
		}
	}
}

func (mm *MainMenu) Update() {

}

func (mm *MainMenu) Render() {
	mm.menu.Render()
}

func (mm MainMenu) GetTick() int {
	return 0
}

func (mm MainMenu) Shutdown() {

}
