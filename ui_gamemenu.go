package main

import (
	"github.com/bennicholls/burl-E/burl"
	"github.com/veandco/go-sdl2/sdl"
)

type GameMenu struct {
	burl.PagedContainer
}

func NewGameMenu() (gm *GameMenu) {
	gm = new(GameMenu)
	gm.PagedContainer = *burl.NewPagedContainer(40, 28, 39, 3, 10, true)

	gm.SetVisibility(false)

	return
}

func (sg *SpaceshipGame) HandleKeypressGameMenu(key sdl.Keycode) {
	// sg.missionMenu.missionList.HandleKeypress(key)

	// switch key {
	// case sdl.K_UP, sdl.K_DOWN:
	// 	burl.PushEvent(burl.NewEvent(burl.EV_UPDATE_UI, "missions"))
	// }
}
