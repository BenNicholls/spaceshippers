package main

import (
	"github.com/bennicholls/burl-E/burl"
	"github.com/veandco/go-sdl2/sdl"
)

type ShipMenu struct {
	burl.PagedContainer

	statusPage      *burl.Container
	powerPage       *burl.Container
	enginePage      *burl.Container
	combatPage      *burl.Container
	lifeSupportPage *burl.Container
	inventoryPage   *burl.Container
	modulesPage     *burl.Container
}

func NewShipMenu() (sm *ShipMenu) {
	sm = new(ShipMenu)
	sm.PagedContainer = *burl.NewPagedContainer(40, 36, 39, 4, 10, true)
	sm.SetVisibility(false)
	sm.SetHint("TAB to switch submenus")

	sm.statusPage = sm.AddPage("Status")
	sm.powerPage = sm.AddPage("Energy")
	sm.enginePage = sm.AddPage("Propulsion")
	sm.combatPage = sm.AddPage("Combat/Shields")
	sm.lifeSupportPage = sm.AddPage("Life Support")
	sm.inventoryPage = sm.AddPage("Stores")
	sm.modulesPage = sm.AddPage("Module")

	return
}

func (sm *ShipMenu) HandleKeypress(key sdl.Keycode) {
	sm.PagedContainer.HandleKeypress(key)
}
