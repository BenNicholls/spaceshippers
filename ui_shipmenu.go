package main

import (
	"github.com/bennicholls/tyumi/event"
	"github.com/bennicholls/tyumi/gfx/ui"
	"github.com/bennicholls/tyumi/vec"

	"fmt"
	"sort"
)

type ShipMenu struct {
	ui.PageContainer

	statusPage      *ui.Page
	powerPage       *ui.Page
	enginePage      *ui.Page
	combatPage      *ui.Page
	lifeSupportPage *ui.Page
	storesMenu      StorageSubmenu
	modulesPage     *ui.Page

	ship *Ship
}

func (sm *ShipMenu) Init(s *Ship) {
	sm.PageContainer.Init(menuSize, menuPos, menuDepth)
	sm.EnableBorder()
	sm.Hide()
	sm.AcceptInput = true
	sm.ship = s

	sm.statusPage = sm.CreatePage("Status")
	sm.powerPage = sm.CreatePage("Energy")
	sm.enginePage = sm.CreatePage("Propulsion")
	sm.combatPage = sm.CreatePage("Combat/Shields")
	sm.lifeSupportPage = sm.CreatePage("Life Support")

	sm.AddPage("Stores", &sm.storesMenu)
	sm.storesMenu.Setup(s)

	sm.modulesPage = sm.CreatePage("Module")

	return
}

type StorageSubmenu struct {
	ui.Page

	capacitiesText      ui.Textbox
	inventoryList       ui.List
	itemNameText        ui.Textbox
	itemDescriptionText ui.Textbox
	itemStorageTypeText ui.Textbox
	itemVolumeText      ui.Textbox

	inventory []string
	ship      *Ship
}

func (ss *StorageSubmenu) Setup(s *Ship) {
	ss.ship = s
	ss.OnActivate = ss.UpdateStorage

	ss.Listen(EV_STORAGEITEMCHANGED, EV_STORAGECAPACITYCHANGED)
	ss.SuppressDuplicateEvents(event.KeepFirst)
	ss.SetEventHandler(ss.handleEvent)

	size := ss.Size()
	stats := ui.Element{}
	stats.Init(vec.Dims{size.W, 6}, vec.ZERO_COORD, 0)
	ss.capacitiesText.Init(vec.Dims{20, 6}, vec.ZERO_COORD, 0, "No Capacities???", ui.JUSTIFY_LEFT)
	stats.AddChild(&ss.capacitiesText)
	ss.inventoryList.Init(vec.Dims{(size.W - 4) / 2, size.H - 8}, vec.Coord{1, 7}, 1)
	ss.inventoryList.SetupBorder("", "Up/Down to Scroll")
	ss.inventoryList.SetEmptyText("No Items in Stores!")
	ss.inventoryList.OnChangeSelection = ss.UpdateItemDescription
	ss.inventoryList.AcceptInput = true
	ss.inventoryList.ToggleHighlight()

	itemDetails := ui.Element{}
	itemDetails.Init(vec.Dims{(size.W - 4) / 2, size.H - 8}, vec.Coord{(size.W-4)/2 + 3, 7}, 2)
	itemDetails.EnableBorder()

	itemw := itemDetails.Size().W
	ss.itemNameText.Init(vec.Dims{itemw, 1}, vec.ZERO_COORD, ui.BorderDepth, "", ui.JUSTIFY_CENTER)
	ss.itemNameText.EnableBorder()
	ss.itemDescriptionText.Init(vec.Dims{itemw, 3}, vec.Coord{0, 3}, 0, "", ui.JUSTIFY_CENTER)
	ss.itemStorageTypeText.Init(vec.Dims{itemw, 1}, vec.Coord{0, 7}, 0, "", ui.JUSTIFY_LEFT)
	ss.itemVolumeText.Init(vec.Dims{itemw, 1}, vec.Coord{0, 8}, 0, "", ui.JUSTIFY_LEFT)

	itemDetails.AddChildren(&ss.itemNameText, &ss.itemDescriptionText, &ss.itemVolumeText, &ss.itemStorageTypeText)

	ss.AddChildren(&stats, &ss.inventoryList, &itemDetails)

	ss.UpdateStorage()
	ss.UpdateItemDescription()
}

func (ss *StorageSubmenu) handleEvent(e event.Event) (event_handled bool) {
	switch e.ID() {
	case EV_STORAGECAPACITYCHANGED:
		ss.UpdateStorageOverview()
		event_handled = true
	case EV_STORAGEITEMCHANGED:
		ss.UpdateStorage()
		event_handled = true
	}

	return
}

func (ss *StorageSubmenu) UpdateStorage() {
	ss.UpdateStorageOverview()
	ss.UpdateInventoryList()
}

func (ss *StorageSubmenu) UpdateStorageOverview() {
	var capacities string
	capacities += fmt.Sprint("General Storage: ", ss.ship.Storage.GetFilledVolume(STORE_GENERAL), "/", ss.ship.Storage.GetCapacity(STORE_GENERAL), "/n")
	capacities += fmt.Sprint("Liquid Storage: ", ss.ship.Storage.GetFilledVolume(STORE_LIQUID), "/", ss.ship.Storage.GetCapacity(STORE_LIQUID), "/n")
	capacities += fmt.Sprint("Gas Storage: ", ss.ship.Storage.GetFillPct(STORE_GAS), "% full/n")
	ss.capacitiesText.ChangeText(capacities)
}

func (ss *StorageSubmenu) CompileInventoryList() {
	ss.inventory = make([]string, 0)

	for item := range ss.ship.Storage.items {
		ss.inventory = append(ss.inventory, item)
	}

	sort.Strings(ss.inventory)
}

func (ss *StorageSubmenu) UpdateInventoryList() {
	ss.CompileInventoryList()
	var selectedItem string
	if ss.inventoryList.GetSelectionIndex() != -1 {
		selectedItem = ss.inventory[ss.inventoryList.GetSelectionIndex()]
	}

	ss.inventoryList.RemoveAll()
	for i, item := range ss.inventory {
		ss.inventoryList.InsertText(ui.JUSTIFY_LEFT, fmt.Sprintf("%.0f - %s", ss.ship.Storage.GetItemVolume(item), item))
		if item == selectedItem {
			ss.inventoryList.Select(i)
		}
	}
}

func (ss *StorageSubmenu) UpdateItemDescription() {
	if index := ss.inventoryList.GetSelectionIndex(); index == -1 {
		ss.itemNameText.ChangeText("No Item!")
		ss.itemDescriptionText.ChangeText("The nothingness here makes you wonder why we even *have* a storage menu.")
		ss.itemStorageTypeText.ChangeText("Stored in: Nowhere. Or everywhere, I guess.")
		ss.itemVolumeText.ChangeText("Amount: No.")
	} else {
		itemName := ss.inventory[index]
		item := ss.ship.Storage.items[itemName]

		ss.itemNameText.ChangeText(item.GetName())
		ss.itemDescriptionText.ChangeText(item.GetDescription())
		ss.itemVolumeText.ChangeText(fmt.Sprintf("Amount: %.0f", item.GetAmount()))

		switch item.GetStorageType() {
		case STORE_GENERAL:
			ss.itemStorageTypeText.ChangeText("Stored in: General Storage")
		case STORE_LIQUID:
			ss.itemStorageTypeText.ChangeText("Stored in: Liquid Storage")
		case STORE_GAS:
			ss.itemStorageTypeText.ChangeText("Stored in: Gas Storage")
		}
	}
}
