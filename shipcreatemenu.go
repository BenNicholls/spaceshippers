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

	dialog Dialog

	galaxy *Galaxy //ship exists in a galaxy, for real

	ship *Ship

	stars StarField
}

func NewShipCreateMenu(g *Galaxy) (scm *ShipCreateMenu) {
	scm = new(ShipCreateMenu)

	scm.window = burl.NewContainer(78, 43, 1, 1, 0, true)
	scm.window.SetTitle("YOUR SHIP IS YOUR WHOLE WORLD")

	scm.shipView = burl.NewTileView(78, 21, 0, 0, 0, true)

	scm.stars = NewStarField(20, scm.shipView)

	scm.shipNameInput = burl.NewInputbox(20, 1, 0, 22, 0, true)
	scm.shipNameInput.SetTitle("Ship Name")
	scm.shipNameInput.CenterX(78, 1)

	scm.shipTypeList = burl.NewList(15, 19, 2, 23, 1, true, "No Ships! How'd this happen?")
	scm.shipTypeList.SetTitle("Ships!")

	scm.shipDescriptionText = burl.NewTextbox(38, 17, 20, 25, 1, true, false, "DESCRIPTIOS")

	scm.generateButton = burl.NewButton(15, 1, 61, 30, 1, true, true, "Confirm Ship Selection!")
	scm.cancelButton = burl.NewButton(15, 1, 61, 35, 1, true, true, "Return to Galaxy Creation")

	scm.window.Add(scm.shipView, scm.shipNameInput, scm.shipTypeList, scm.shipDescriptionText, scm.generateButton, scm.cancelButton)

	scm.shipTypeList.Append("Civilian Craft", "Transport", "Mining Ship", "Fighter", "Explorer")
	scm.UpdateShipDescription()

	scm.focusedField = scm.shipNameInput
	scm.focusedField.ToggleFocus()

	scm.shipNameInput.SetTabID(1)
	scm.shipTypeList.SetTabID(2)
	scm.generateButton.SetTabID(3)
	scm.cancelButton.SetTabID(4)

	scm.galaxy = g

	scm.CreateShip()

	return
}

func (scm *ShipCreateMenu) UpdateShipDescription() {
	switch scm.shipTypeList.GetSelection() {
	case 0: //Civilian
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

func (scm *ShipCreateMenu) CreateShip() {
	scm.ship = NewShip(scm.shipNameInput.GetText(), scm.galaxy)
	switch scm.shipTypeList.GetSelection() {
	case 0: //Civilian
		scm.ship.SetupFromTemplate(SHIPTYPE_CIVILIAN)
	case 1: //Transport
		scm.ship.SetupFromTemplate(SHIPTYPE_TRANSPORT)
	case 2: //Mining Ship
		scm.ship.SetupFromTemplate(SHIPTYPE_CIVILIAN)
	case 3: //Fighter
		scm.ship.SetupFromTemplate(SHIPTYPE_CIVILIAN)
	case 4: //Explorer
		scm.ship.SetupFromTemplate(SHIPTYPE_CIVILIAN)
	}

	scm.ship.SetupShip(scm.galaxy)
}

func (scm *ShipCreateMenu) HandleKeypress(key sdl.Keycode) {
	if scm.dialog != nil {
		scm.dialog.HandleKeypress(key)
		return
	}

	//non-standard ui behaviour
	if key == sdl.K_TAB || (key == sdl.K_RETURN && scm.focusedField == scm.shipNameInput) {
		burl.PushEvent(burl.NewEvent(burl.EV_TAB_FIELD, "+"))
	}

	scm.focusedField.HandleKeypress(key)
}

func (scm *ShipCreateMenu) HandleEvent(e *burl.Event) {
	switch e.ID {
	case burl.EV_LIST_CYCLE:
		if e.Caller == scm.shipTypeList {
			scm.UpdateShipDescription()
			scm.CreateShip()
		}
	case burl.EV_TAB_FIELD:
		scm.focusedField.ToggleFocus()
		scm.focusedField = scm.window.FindNextTab(scm.focusedField)
		scm.focusedField.ToggleFocus()
	case burl.EV_ANIMATION_DONE:
		if e.Caller == scm.generateButton {
			if scm.shipNameInput.GetText() == "" {
				scm.dialog = NewCommDialog("", "", "", "You must give your ship a name before you can continue!")
			} else {
				scm.ship.Name = scm.shipNameInput.GetText()
				burl.ChangeState(NewSpaceshipGame(scm.galaxy, scm.ship))
			}
		} else if e.Caller == scm.cancelButton {
			burl.ChangeState(NewCreateGalaxyMenu())
		}
	}
}

func (scm *ShipCreateMenu) Update() {
	if scm.dialog != nil && scm.dialog.Done() {
		scm.dialog = nil
	}

	scm.Tick++

	if scm.Tick%10 == 0 {
		scm.stars.Shift()
	}

	//move around the crew, for fun!
	for i, _ := range scm.ship.Crew {
		scm.ship.Crew[i].Update()
		if scm.Tick%20 == 0 {
			dx, dy := burl.RandomDirection()
			if scm.ship.shipMap.GetTile(scm.ship.Crew[i].X+dx, scm.ship.Crew[i].Y+dy).Empty() {
				scm.ship.shipMap.MoveEntity(scm.ship.Crew[i].X, scm.ship.Crew[i].Y, dx, dy)
				scm.ship.Crew[i].Move(dx, dy)
			}
		}
	}
}

func (scm *ShipCreateMenu) Render() {
	scm.stars.Draw()

	//calculate offset for ship based on ship size, so ship is centered
	displayWidth, displayHeight := scm.shipView.Dims()
	offX := scm.ship.width/2 + scm.ship.x - displayWidth/2
	offY := scm.ship.height/2 + scm.ship.y - displayHeight/2

	scm.ship.DrawToTileView(scm.shipView, offX, offY)

	scm.window.Render()

	if scm.dialog != nil {
		scm.dialog.Render()
	}
}
