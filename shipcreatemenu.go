package main

import (
	"github.com/bennicholls/burl-E/burl"
	"github.com/veandco/go-sdl2/sdl"
)

type ShipCreateMenu struct {
	burl.BaseState

	window              *burl.Container
	shipView            *burl.TileView
	shipNameInput       *burl.Inputbox
	shipTypeList        *burl.List
	shipDescriptionText *burl.Textbox
	generateButton      *burl.Button
	cancelButton        *burl.Button
	focusedField        burl.UIElem

	galaxy *Galaxy //ship exists in a galaxy, for real

	ship *Ship
}

func NewShipCreateMenu(g *Galaxy) (scm *ShipCreateMenu) {
	scm = new(ShipCreateMenu)

	scm.window = burl.NewContainer(78, 43, 1, 1, 0, true)
	scm.window.SetTitle("YOUR SHIP IS YOUR WHOLE WORLD")
	scm.shipView = burl.NewTileView(78, 21, 0, 0, 0, true)

	scm.shipNameInput = burl.NewInputbox(20, 1, 0, 22, 0, true)
	scm.shipNameInput.SetTitle("Ship Name")
	scm.shipNameInput.CenterX(78, 1)

	scm.shipTypeList = burl.NewList(15, 19, 2, 23, 1, true, "No Ships! How'd this happen?")
	scm.shipTypeList.SetTitle("Ships!")

	scm.shipDescriptionText = burl.NewTextbox(38, 17, 20, 25, 1, true, false, "DESCRIPTIOS")

	scm.generateButton = burl.NewButton(15, 1, 61, 30, 1, true, true, "Confirm Ship Selection!")
	scm.cancelButton = burl.NewButton(15, 1, 61, 35, 1, true, true, "Return to Galaxy Creation")

	scm.window.Add(scm.shipView, scm.shipNameInput, scm.shipTypeList, scm.shipDescriptionText, scm.generateButton, scm.cancelButton)

	scm.galaxy = g

	scm.shipTypeList.Append("Civilian Craft", "Transport", "Mining Ship", "Fighter", "Explorer")
	scm.UpdateShipDescription()

	scm.focusedField = scm.shipNameInput
	scm.focusedField.ToggleFocus()

	scm.shipNameInput.SetTabID(1)
	scm.shipTypeList.SetTabID(2)
	scm.generateButton.SetTabID(3)
	scm.cancelButton.SetTabID(4)

	return
}

func (scm *ShipCreateMenu) UpdateShipDescription() {
	switch scm.shipTypeList.GetSelection() {
	case 0: //Scout
		scm.shipDescriptionText.ChangeText("CIVILIAN CRAFT/n/nThe Toyota Camry of spaceships. The Civilian Craft has everything the casual spacegoer might need: door, engine, steering wheel of some variety, snack table, all the cool space things. While it may not look like much, it's cheap to repair, very moddable, and will last forever if you take care of it right!")
	case 1: //Transport
		scm.shipDescriptionText.ChangeText("TRANSPORT SHIP/n/nUsed for transporting passengers and cargo. The Transport Ship begins with a larger cargo bay, additional dormitories, and a bulkier engine. Guzzles like an Irishman though.")
	case 2: //Mining Ship
		scm.shipDescriptionText.ChangeText("MINING SHIP/n/nAn industrial craft used in the dangerous but lucrative asteroid/comet mining field. The Mining Ship comes with an expanded cargo bay, medical facility, and high durability plating.")
	case 3: //Fighter
		scm.shipDescriptionText.ChangeText("FIGHER/n/nA small ship employed as an assault craft to be launched from a larger Carrier vessel. The Fighter begins with a powerful, efficient engine, punchy laser weapons, but with only room for 2 crew!")
	case 4: //Explorer
		scm.shipDescriptionText.ChangeText("EXPLORER/n/nA deep-space science vessel used to map new systems. The Explorer has efficient engines, extra fuel capacity, and a laboratory.")
	}
}

func (scm *ShipCreateMenu) HandleKeypress(key sdl.Keycode) {
	switch key {
	// case sdl.K_UP:
	// 	scm.focusedField.ToggleFocus()
	// 	scm.focusedField = scm.window.FindPrevTab(scm.focusedField)
	// 	scm.focusedField.ToggleFocus()
	// 	scm.UpdateShipDescription()
	case sdl.K_TAB:
		scm.focusedField.ToggleFocus()
		scm.focusedField = scm.window.FindNextTab(scm.focusedField)
		scm.focusedField.ToggleFocus()
	}

	switch scm.focusedField {
	case scm.shipNameInput:
		switch key {
		case sdl.K_BACKSPACE:
			scm.shipNameInput.Delete()
		case sdl.K_SPACE:
			scm.shipNameInput.Insert(" ")
		default:
			scm.shipNameInput.InsertText(rune(key))
		}
	case scm.shipTypeList:
		switch key {
		case sdl.K_UP:
			scm.shipTypeList.Prev()
			scm.UpdateShipDescription()
		case sdl.K_DOWN:
			scm.shipTypeList.Next()
			scm.UpdateShipDescription()
		}
	case scm.generateButton:
		switch key {
		case sdl.K_RETURN:
			scm.generateButton.Press()
		}
	case scm.cancelButton:
		switch key {
		case sdl.K_RETURN:
			scm.cancelButton.Press()
		}
	}
}

func (scm *ShipCreateMenu) Render() {
	scm.window.Render()
}
