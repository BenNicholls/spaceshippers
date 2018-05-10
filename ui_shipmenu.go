package main

import (
	"github.com/bennicholls/burl-E/burl"
	"github.com/veandco/go-sdl2/sdl"
)

type ShipMenu struct {
	burl.PagedContainer
}

func NewShipMenu() (sm *ShipMenu) {
	sm = new(ShipMenu)
	sm.PagedContainer = *burl.NewPagedContainer(40, 28, 39, 3, 10, true)

	sm.SetVisibility(false)

	return
}

func (sg *SpaceshipGame) HandleKeypressShipMenu(key sdl.Keycode) {
	// switch key {
	// case sdl.K_PAGEUP:
	// 	sg.shipMenu.pages[sg.shipMenu.pageList.GetSelection()].ToggleVisible()
	// 	sg.shipMenu.pageList.Prev()
	// 	sg.shipMenu.pages[sg.shipMenu.pageList.GetSelection()].ToggleVisible()
	// case sdl.K_PAGEDOWN:
	// 	sg.shipMenu.pages[sg.shipMenu.pageList.GetSelection()].ToggleVisible()
	// 	sg.shipMenu.pageList.Next()
	// 	sg.shipMenu.pages[sg.shipMenu.pageList.GetSelection()].ToggleVisible()
	// }
}

// type ShipMenu struct {
// 	burl.Container

// 	pageList *burl.List

// 	pages []*burl.Container

// 	shipPage       *burl.Container
// 	propulsionPage *burl.Container
// 	powerPage      *burl.Container
// }

// func NewShipMenu() (sm *ShipMenu) {
// 	sm = new(ShipMenu)

// 	sm.Container = *burl.NewContainer(40, 26, 39, 4, 5, true)
// 	sm.SetTitle("Ship")
// 	sm.SetVisibility(false)

// 	sm.pageList = burl.NewList(7, 26, 0, 0, 0, true, "No pages??")
// 	sm.pageList.Append("Ship", "Propulsion", "Power")

// 	sm.Add(sm.pageList)

// 	sm.pages = make([]*burl.Container, 0)

// 	sm.shipPage = burl.NewContainer(32, 26, 8, 0, 0, true)
// 	sm.shipPage.Add(burl.NewTextbox(10, 1, 10, 10, 0, true, true, "Ship Stuff"))

// 	sm.propulsionPage = burl.NewContainer(32, 26, 8, 0, 0, true)
// 	sm.propulsionPage.SetVisibility(false)
// 	sm.propulsionPage.Add(burl.NewTextbox(10, 1, 10, 10, 0, true, true, "propulsion Stuff"))

// 	sm.powerPage = burl.NewContainer(32, 26, 8, 0, 0, true)
// 	sm.powerPage.SetVisibility(false)
// 	sm.powerPage.Add(burl.NewTextbox(10, 1, 10, 10, 0, true, true, "power Stuff"))

// 	sm.pages = append(sm.pages, sm.shipPage, sm.propulsionPage, sm.powerPage)

// 	for p := range sm.pages {
// 		sm.Add(sm.pages[p])
// 	}

// 	return
// }
