package main

import "github.com/bennicholls/burl-E/burl"

type MainMenu struct {
	burl.PagedContainer
}

func NewMainMenu() (mm *MainMenu) {
	mm = new(MainMenu)
	mm.PagedContainer = *burl.NewPagedContainer(40, 28, 39, 3, 10, true)

	mm.SetVisibility(false)

	return
}
