package main

import (
	"github.com/bennicholls/burl-E/burl"
	"github.com/veandco/go-sdl2/sdl"

	"fmt"
	"sort"
)

type ShipMenu struct {
	burl.PagedContainer

	statusPage      *burl.Container
	powerPage       *burl.Container
	enginePage      *burl.Container
	combatPage      *burl.Container
	lifeSupportPage *burl.Container
	storesMenu      StorageSubmenu
	modulesPage     *burl.Container

	ship *Ship
}

func NewShipMenu(s *Ship) (sm *ShipMenu) {
	sm = new(ShipMenu)
	sm.PagedContainer = *burl.NewPagedContainer(56, 45, 39, 4, 10, true)
	sm.SetVisibility(false)
	sm.SetHint("TAB to switch submenus")
	sm.ship = s

	sm.statusPage = sm.AddPage("Status")
	sm.powerPage = sm.AddPage("Energy")
	sm.enginePage = sm.AddPage("Propulsion")
	sm.combatPage = sm.AddPage("Combat/Shields")
	sm.lifeSupportPage = sm.AddPage("Life Support")
	sm.storesMenu.window = sm.AddPage("Stores")
	sm.modulesPage = sm.AddPage("Module")

	sm.storesMenu.Init(sm.ship)

	return
}

func (sm *ShipMenu) HandleKeypress(key sdl.Keycode) {
	sm.PagedContainer.HandleKeypress(key)

	switch sm.PagedContainer.CurrentIndex() {
	case 5: //stores page
		sm.storesMenu.HandleKeypress(key)
	}
}

type StorageSubmenu struct {
	window *burl.Container

	stats               *burl.Container
	capacitiesText      *burl.Textbox
	inventoryList       *burl.List
	itemDetails         *burl.Container
	itemNameText        *burl.Textbox
	itemDescriptionText *burl.Textbox
	itemStorageTypeText *burl.Textbox
	itemVolumeText      *burl.Textbox

	inventory []string
	ship      *Ship
}

func (ss *StorageSubmenu) Init(s *Ship) {
	ss.ship = s

	w, h := ss.window.Dims()
	ss.stats = burl.NewContainer(w, 6, 0, 0, 0, false)
	ss.capacitiesText = burl.NewTextbox(20, 6, 0, 0, 0, false, false, "No Capacities????")
	ss.stats.Add(ss.capacitiesText)
	ss.inventoryList = burl.NewList((w-4)/2, h-8, 1, 7, 1, true, "No Items in Stores!")
	ss.inventoryList.SetHint("PgUp/PgDown to Scroll")
	ss.itemDetails = burl.NewContainer((w-4)/2, h-8, (w-4)/2+3, 7, 2, true)
	itemw, _ := ss.itemDetails.Dims()
	ss.itemNameText = burl.NewTextbox(itemw, 1, 0, 0, 0, true, true, "")
	ss.itemDescriptionText = burl.NewTextbox(itemw, 3, 0, 3, 0, false, false, "")
	ss.itemStorageTypeText = burl.NewTextbox(itemw, 1, 0, 7, 0, false, false, "")
	ss.itemVolumeText = burl.NewTextbox(itemw, 1, 0, 8, 0, false, false, "")
	ss.itemDetails.Add(ss.itemNameText, ss.itemDescriptionText, ss.itemVolumeText, ss.itemStorageTypeText)

	ss.CompileInventoryList()

	ss.window.Add(ss.stats, ss.inventoryList, ss.itemDetails)

	ss.Update()
}

func (ss *StorageSubmenu) CompileInventoryList() {
	ss.inventory = make([]string, 0)

	for item := range ss.ship.Storage.items {
		ss.inventory = append(ss.inventory, item)
	}

	sort.Strings(ss.inventory)
}

func (ss *StorageSubmenu) Update() {
	var capacities string
	capacities += fmt.Sprint("General Storage: ", ss.ship.Storage.GetFilledVolume(STORE_GENERAL), "/", ss.ship.Storage.GetCapacity(STORE_GENERAL), "/n")
	capacities += fmt.Sprint("Liquid Storage: ", ss.ship.Storage.GetFilledVolume(STORE_LIQUID), "/", ss.ship.Storage.GetCapacity(STORE_LIQUID), "/n")
	capacities += fmt.Sprint("Gas Storage: ", ss.ship.Storage.GetFillPct(STORE_GAS), "% full/n")
	ss.capacitiesText.ChangeText(capacities)

	selectedItem := ss.inventory[ss.inventoryList.GetSelection()]
	ss.inventoryList.ClearElements()
	for i, item := range ss.inventory {
		ss.inventoryList.Append(fmt.Sprint(ss.ship.Storage.GetItemVolume(item)) + " - " + item)
		if item == selectedItem {
			ss.inventoryList.Select(i)
		}
	}

	ss.UpdateItemDescription()
}

func (ss *StorageSubmenu) UpdateItemDescription() {
	itemName := ss.inventory[ss.inventoryList.GetSelection()]
	item := ss.ship.Storage.items[itemName]

	ss.itemNameText.ChangeText(item.GetName())
	ss.itemDescriptionText.ChangeText(item.GetDescription())
	ss.itemVolumeText.ChangeText("Amount: " + fmt.Sprint(item.GetAmount()))
	switch item.GetStorageType() {
	case STORE_GENERAL:
		ss.itemStorageTypeText.ChangeText("Stored in: General Storage")
	case STORE_LIQUID:
		ss.itemStorageTypeText.ChangeText("Stored in: Liquid Storage")
	case STORE_GAS:
		ss.itemStorageTypeText.ChangeText("Stored in: Gas Storage")
	}
}

func (ss *StorageSubmenu) HandleKeypress(key sdl.Keycode) {
	selection := ss.inventoryList.GetSelection()
	ss.inventoryList.HandleKeypress(key)
	if selection != ss.inventoryList.GetSelection() {
		ss.UpdateItemDescription() //update if selected item has changed.
	}
}
