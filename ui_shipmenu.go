package main

import "github.com/bennicholls/burl/ui"

type ShipMenu struct {
	ui.Container

}

func InitShipMenu() (sm *ShipMenu) {
	sm = new(ShipMenu)

	sm.Container = *ui.NewContainer(20, 27, 59, 4, 3, true)
	sm.SetTitle("Ship")
	sm.SetVisibility(false)
	sm.ToggleFocus()

	sm.Add(ui.NewTextbox(13, 1, 2, 2, 1, true, true, "ship stuff?"))

	return
}