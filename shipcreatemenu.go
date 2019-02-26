package main

import (
	"github.com/bennicholls/burl-E/burl"
	"github.com/veandco/go-sdl2/sdl"
)

type ShipCreateMenu struct {
	burl.StatePrototype

	shipView            *burl.TileView
	shipNameInput       *burl.Inputbox
	shipTypeList        *burl.List
	shipDescriptionText *burl.Textbox
	generateButton      *burl.Button
	cancelButton        *burl.Button
	focusedField        burl.UIElem

	galaxy *Galaxy //ship exists in a galaxy, for real
	ship   *Ship

	stars StarField

	templates []ShipTemplate
}

func NewShipCreateMenu(g *Galaxy) (scm *ShipCreateMenu) {
	scm = new(ShipCreateMenu)
	scm.InitWindow(true)
	scm.Window.SetTitle("YOUR SHIP IS YOUR WHOLE WORLD")

	scm.shipView = burl.NewTileView(94, 30, 0, 0, 0, true)

	scm.stars = NewStarField(20, scm.shipView)

	scm.shipNameInput = burl.NewInputbox(20, 1, 0, 31, 0, true)
	scm.shipNameInput.SetTitle("Ship Name")
	scm.shipNameInput.CenterX(94, 1)

	scm.shipTypeList = burl.NewList(15, 19, 2, 32, 1, true, "No Ships! How'd this happen?")
	scm.shipTypeList.SetTitle("Ships!")

	scm.shipDescriptionText = burl.NewTextbox(52, 17, 20, 34, 1, true, false, "DESCRIPTIOS")

	scm.generateButton = burl.NewButton(15, 1, 76, 39, 1, true, true, "Confirm Ship Selection!")
	scm.cancelButton = burl.NewButton(15, 1, 76, 44, 1, true, true, "Return to Galaxy Creation")

	scm.Window.Add(scm.shipView, scm.shipNameInput, scm.shipTypeList, scm.shipDescriptionText, scm.generateButton, scm.cancelButton)

	scm.templates = make([]ShipTemplate, 0)
	shipFiles, _ := burl.GetFileList("raws/ship/", ".shp")

	for _, file := range shipFiles {
		temp, err := LoadShipTemplate("raws/ship/" + file)
		if err != nil {
			burl.LogError(err.Error())
		} else {
			scm.templates = append(scm.templates, temp)
			scm.shipTypeList.Append(temp.Name)
		}
	}

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
	temp := scm.templates[scm.shipTypeList.GetSelection()]
	scm.shipDescriptionText.ChangeText(temp.Name + "/n/n" + temp.Description)
}

func (scm *ShipCreateMenu) CreateShip() {
	temp := scm.templates[scm.shipTypeList.GetSelection()]
	scm.ship = NewShip(scm.shipNameInput.GetText(), scm.galaxy)
	scm.ship.SetupFromTemplate(temp)
	for i := 0; i < temp.CrewNum; i++ {
		scm.ship.AddCrewman(NewCrewman())
	}
	scm.ship.SetupShip(scm.galaxy)
}

func (scm *ShipCreateMenu) HandleKeypress(key sdl.Keycode) {
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
		scm.focusedField = scm.Window.FindNextTab(scm.focusedField)
		scm.focusedField.ToggleFocus()
	case burl.EV_ANIMATION_DONE:
		if e.Caller == scm.generateButton {
			if scm.shipNameInput.GetText() == "" {
				scm.OpenDialog(NewCommDialog("", "", "", "You must give your ship a name before you can continue!"))
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
	scm.Tick++

	if scm.Tick%10 == 0 {
		scm.stars.Shift()
	}

	//move around the crew, for fun!
	for i, _ := range scm.ship.Crew {
		scm.ship.Crew[i].Update(scm.Tick)
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

	scm.ship.DrawToTileView(scm.shipView, VIEW_DEFAULT, offX, offY)
}
