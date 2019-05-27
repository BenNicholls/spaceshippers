package main

import (
	"github.com/bennicholls/burl-E/burl"
	"github.com/veandco/go-sdl2/sdl"

	"fmt"
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

	sm.storesPage.Add(sm.storesStats, sm.storesInventoryList, sm.storesItemDetails)

	sm.UpdateStoreMenu()

	return
}

func (sm *ShipMenu) UpdateStoreMenu() {
	var capacities string
	capacities += fmt.Sprint("General Storage: ", sm.ship.Storage.GetStorageVolume(STAT_GENERAL_STORAGE), "/", sm.ship.Storage.GetStorageCapacity(STAT_GENERAL_STORAGE), "/n")
	capacities += fmt.Sprint("Volatile Storage: ", sm.ship.Storage.GetStorageVolume(STAT_VOLATILE_STORAGE), "/", sm.ship.Storage.GetStorageCapacity(STAT_VOLATILE_STORAGE), "/n")
	capacities += fmt.Sprint("Cold Storage: ", sm.ship.Storage.GetStorageVolume(STAT_COLD_STORAGE), "/", sm.ship.Storage.GetStorageCapacity(STAT_COLD_STORAGE), "/n")
	capacities += fmt.Sprint("Fuel Storage: ", sm.ship.Storage.GetStorageVolume(STAT_FUEL_STORAGE), "/", sm.ship.Storage.GetStorageCapacity(STAT_FUEL_STORAGE), "/n")
	sm.storesStatsCapacities.ChangeText(capacities)

	i := sm.storesInventoryList.GetSelection()
	sm.storesInventoryList.ClearElements()
	for _, item := range sm.ship.Storage.items {
		sm.storesInventoryList.Append(fmt.Sprint(item.GetVolume()) + " - " + item.GetName())
	}
	sm.storesInventoryList.Select(i)

	sm.UpdateStoreItemDescription()
}

func (sm *ShipMenu) UpdateStoreItemDescription() {
	item := sm.ship.Storage.items[sm.storesInventoryList.GetSelection()]

	sm.storesItemNameText.ChangeText(item.GetName())
	sm.storesItemDescriptionText.ChangeText(item.GetDescription())
	sm.storesItemVolumeText.ChangeText("Volume: " + fmt.Sprint(item.GetVolume()))
	switch item.GetStorageType() {
	case STAT_GENERAL_STORAGE:
		sm.storesItemStorageTypeText.ChangeText("Stored in: General Storage")
	case STAT_VOLATILE_STORAGE:
		sm.storesItemStorageTypeText.ChangeText("Stored in: Volatile Storage")
	case STAT_COLD_STORAGE:
		sm.storesItemStorageTypeText.ChangeText("Stored in: Cold Storage")
	case STAT_FUEL_STORAGE:
		sm.storesItemStorageTypeText.ChangeText("Stored in: Fuel Storage")
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
