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
	storesPage      *burl.Container
	modulesPage     *burl.Container

	//stores page
	storesStats               *burl.Container
	storesStatsCapacities     *burl.Textbox
	storesInventoryList       *burl.List
	storesItemDetails         *burl.Container
	storesItemNameText        *burl.Textbox
	storesItemDescriptionText *burl.Textbox
	storesItemStorageTypeText *burl.Textbox
	storesItemVolumeText      *burl.Textbox
	inventory                 []string

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
	sm.storesPage = sm.AddPage("Stores")
	sm.modulesPage = sm.AddPage("Module")

	w, h := sm.GetPageDims()
	sm.storesStats = burl.NewContainer(w, 6, 0, 0, 0, false)
	sm.storesStatsCapacities = burl.NewTextbox(20, 6, 0, 0, 0, false, false, "No Capacities????")
	sm.storesStats.Add(sm.storesStatsCapacities)
	sm.storesInventoryList = burl.NewList((w-4)/2, h-8, 1, 7, 1, true, "No Items in Stores!")
	sm.storesInventoryList.SetHint("PgUp/PgDown to Scroll")
	sm.storesItemDetails = burl.NewContainer((w-4)/2, h-8, (w-4)/2+3, 7, 2, true)
	itemw, _ := sm.storesItemDetails.Dims()
	sm.storesItemNameText = burl.NewTextbox(itemw, 1, 0, 0, 0, true, true, "")
	sm.storesItemDescriptionText = burl.NewTextbox(itemw, 3, 0, 3, 0, false, false, "")
	sm.storesItemStorageTypeText = burl.NewTextbox(itemw, 1, 0, 7, 0, false, false, "")
	sm.storesItemVolumeText = burl.NewTextbox(itemw, 1, 0, 8, 0, false, false, "")
	sm.storesItemDetails.Add(sm.storesItemNameText, sm.storesItemDescriptionText, sm.storesItemVolumeText, sm.storesItemStorageTypeText)

	sm.CompileInventoryList()

	sm.storesPage.Add(sm.storesStats, sm.storesInventoryList, sm.storesItemDetails)

	sm.UpdateStoreMenu()

	return
}

func (sm *ShipMenu) CompileInventoryList() {
	sm.inventory = make([]string, 0)

	for item := range sm.ship.Storage.items {
		sm.inventory = append(sm.inventory, item)
	}

	sort.Strings(sm.inventory)
}

func (sm *ShipMenu) UpdateStoreMenu() {
	var capacities string
	capacities += fmt.Sprint("General Storage: ", sm.ship.Storage.GetFilledVolume(STORE_GENERAL), "/", sm.ship.Storage.GetCapacity(STORE_GENERAL), "/n")
	capacities += fmt.Sprint("Liquid Storage: ", sm.ship.Storage.GetFilledVolume(STORE_LIQUID), "/", sm.ship.Storage.GetCapacity(STORE_LIQUID), "/n")
	capacities += fmt.Sprint("Gas Storage: ", sm.ship.Storage.GetFillPct(STORE_GAS), "% full/n")
	sm.storesStatsCapacities.ChangeText(capacities)

	selectedItem := sm.inventory[sm.storesInventoryList.GetSelection()]
	sm.storesInventoryList.ClearElements()
	for i, item := range sm.inventory {
		sm.storesInventoryList.Append(fmt.Sprint(sm.ship.Storage.GetItemVolume(item)) + " - " + item)
		if item == selectedItem {
			sm.storesInventoryList.Select(i)
		}

	}

	sm.UpdateStoreItemDescription()
}

func (sm *ShipMenu) UpdateStoreItemDescription() {
	itemName := sm.inventory[sm.storesInventoryList.GetSelection()]
	item := sm.ship.Storage.items[itemName]

	sm.storesItemNameText.ChangeText(item.GetName())
	sm.storesItemDescriptionText.ChangeText(item.GetDescription())
	sm.storesItemVolumeText.ChangeText("Amount: " + fmt.Sprint(item.GetAmount()))
	switch item.GetStorageType() {
	case STORE_GENERAL:
		sm.storesItemStorageTypeText.ChangeText("Stored in: General Storage")
	case STORE_LIQUID:
		sm.storesItemStorageTypeText.ChangeText("Stored in: Liquid Storage")
	case STORE_GAS:
		sm.storesItemStorageTypeText.ChangeText("Stored in: Gas Storage")
	}
}

func (sm *ShipMenu) HandleKeypress(key sdl.Keycode) {
	sm.PagedContainer.HandleKeypress(key)

	switch sm.PagedContainer.CurrentIndex() {
	case 5: //stores page
		sm.storesInventoryList.HandleKeypress(key)
		sm.UpdateStoreItemDescription()
	}
}
