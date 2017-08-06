package main

import "github.com/bennicholls/burl-E/burl"

type ShipMenu struct {
	burl.Container
}

func NewShipMenu() (sm *ShipMenu) {
	sm = new(ShipMenu)

	sm.Container = *burl.NewContainer(20, 27, 59, 4, 3, true)
	sm.SetTitle("Ship")
	sm.SetVisibility(false)
	sm.ToggleFocus()

	sm.Add(burl.NewTextbox(13, 1, 2, 2, 1, true, true, "ship stuff?"))

	return
}
