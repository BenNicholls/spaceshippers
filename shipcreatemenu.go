package main

import (
	"github.com/bennicholls/tyumi"
	"github.com/bennicholls/tyumi/gfx"
	"github.com/bennicholls/tyumi/gfx/ui"
	"github.com/bennicholls/tyumi/log"
	"github.com/bennicholls/tyumi/rl"
	"github.com/bennicholls/tyumi/util"
	"github.com/bennicholls/tyumi/vec"
)

type ShipCreateMenu struct {
	tyumi.State

	shipView            rl.TileMapView
	shipNameInput       ui.InputBox
	shipTypeList        ui.List
	shipDescriptionText ui.Textbox
	generateButton      ui.Button
	cancelButton        ui.Button
	stars               StarField

	galaxy *Galaxy //ship exists in a galaxy, for real
	ship   *Ship

	templates []ShipTemplate
}

func NewShipCreateMenu(galaxy *Galaxy) (scm *ShipCreateMenu) {
	scm = new(ShipCreateMenu)

	scm.InitBordered()
	scm.Window().SetupBorder("YOUR SHIP IS YOUR WHOLE WORLD", "[TAB]/[ENTER]")
	windowStyle := ui.DefaultBorderStyle
	windowStyle.TitleJustification = ui.JUSTIFY_CENTER
	scm.Window().Border.SetStyle(ui.BORDER_STYLE_CUSTOM, windowStyle)

	scm.shipView.Init(vec.Dims{94, 30}, vec.ZERO_COORD, ui.BorderDepth, nil)
	scm.shipView.SetDefaultVisuals(gfx.Visuals{Mode: gfx.DRAW_NONE})
	scm.Window().AddChild(&scm.shipView)

	scm.shipNameInput.Init(vec.Dims{20, 1}, vec.Coord{0, 31}, ui.BorderDepth, 0)
	scm.shipNameInput.SetupBorder("Ship Name", "")
	scm.shipNameInput.Border.SetStyle(ui.BORDER_STYLE_CUSTOM, windowStyle)
	scm.Window().AddChild(&scm.shipNameInput)
	scm.shipNameInput.CenterHorizontal()

	scm.shipTypeList.Init(vec.Dims{15, 19}, vec.Coord{2, 32}, 1)
	scm.shipTypeList.SetupBorder("Ships!", "")
	scm.shipTypeList.OnChangeSelection = func() {
		scm.UpdateShipDescription()
		scm.CreateShip()
	}
	scm.shipTypeList.ToggleHighlight()
	scm.Window().AddChild(&scm.shipTypeList)
	//TODO: default text for empty list

	scm.shipDescriptionText.Init(vec.Dims{52, 17}, vec.Coord{20, 34}, 1, "Descriptios", ui.JUSTIFY_LEFT)
	scm.shipDescriptionText.EnableBorder()
	scm.Window().AddChild(&scm.shipDescriptionText)

	scm.generateButton.Init(vec.Dims{15, 1}, vec.Coord{76, 39}, 1, "Confirm Ship Selection!", scm.onGeneratePress)
	scm.generateButton.EnableBorder()
	scm.cancelButton.Init(vec.Dims{15, 1}, vec.Coord{76, 44}, 1, "Return to Galaxy Creation", func() {
		tyumi.ChangeState(NewCreateGalaxyMenu())
	})
	scm.cancelButton.EnableBorder()
	scm.Window().AddChildren(&scm.generateButton, &scm.cancelButton)

	scm.stars.Init(scm.Window().Size(), vec.ZERO_COORD, 0, 25, 10)
	scm.Window().AddChild(&scm.stars)

	scm.templates = make([]ShipTemplate, 0)
	shipFiles, _ := util.GetFileList("raws/ship/", ".shp")

	for _, file := range shipFiles {
		temp, err := LoadShipTemplate("raws/ship/" + file)
		if err != nil {
			log.Error(err.Error())
		} else {
			scm.templates = append(scm.templates, temp)
			scm.shipTypeList.AddTextItems(ui.JUSTIFY_LEFT, temp.Name)
		}
	}

	scm.UpdateShipDescription()

	scm.shipNameInput.Focus()
	scm.Window().SetTabbingOrder(&scm.shipNameInput, &scm.shipTypeList, &scm.generateButton, &scm.cancelButton)

	scm.galaxy = galaxy
	scm.CreateShip()

	return
}

func (scm *ShipCreateMenu) UpdateShipDescription() {
	temp := scm.templates[scm.shipTypeList.GetSelectionIndex()]
	scm.shipDescriptionText.ChangeText(temp.Name + "/n/n" + temp.Description)
}

func (scm *ShipCreateMenu) CreateShip() {
	temp := scm.templates[scm.shipTypeList.GetSelectionIndex()]
	scm.ship = NewShip(scm.shipNameInput.InputtedText(), scm.galaxy)
	scm.ship.SetupFromTemplate(temp)
	for range temp.CrewNum {
		scm.ship.AddCrewman(NewCrewman())
	}
	scm.ship.SetupShip(scm.galaxy)
	scm.shipView.SetTileMap(&scm.ship.shipMap)

	//calculate offset for ship based on ship size, so ship is centered
	//TODO: the tilemap view should do this automatically maybe
	offX := scm.ship.width/2 + scm.ship.x - scm.shipView.Size().W/2
	offY := scm.ship.height/2 + scm.ship.y - scm.shipView.Size().H/2
	scm.shipView.SetCameraOffset(vec.Coord{-offX, -offY})
	scm.Window().ForceRedraw()

}

func (scm *ShipCreateMenu) onGeneratePress() {
	if scm.shipNameInput.InputtedText() == "" {
		scm.OpenDialog(NewSimpleCommDialog("You must give your ship a name before you can continue!"))
	} else {
		scm.ship.Name = scm.shipNameInput.InputtedText()
		//burl.ChangeState(NewSpaceshipGame(scm.galaxy, scm.ship))
	}
}

func (scm *ShipCreateMenu) Update() {
	//move around the crew, for fun!
	for i := range scm.ship.Crew {
		scm.ship.Crew[i].Update(tyumi.GetTick())
		if tyumi.GetTick()%20 == 0 {
			// dx, dy := util.RandomDirection()
			// if scm.ship.shipMap.GetTile(scm.ship.Crew[i].X+dx, scm.ship.Crew[i].Y+dy).Empty() {
			// 	scm.ship.shipMap.MoveEntity(scm.ship.Crew[i].X, scm.ship.Crew[i].Y, dx, dy)
			// 	scm.ship.Crew[i].Move(dx, dy)
			// }
		}
	}
}