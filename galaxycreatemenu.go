package main

import (
	"github.com/bennicholls/tyumi"
	"github.com/bennicholls/tyumi/event"
	"github.com/bennicholls/tyumi/gfx/ui"
	"github.com/bennicholls/tyumi/util"
	"github.com/bennicholls/tyumi/vec"
)

type CreateGalaxyMenu struct {
	tyumi.Scene

	nameInput      ui.InputBox
	densityChoice  ui.ChoiceBox //choice between some pre-defined densities
	shapeChoice    ui.ChoiceBox //choice between blob, spiral, maybe more esoteric shapes??
	sizeChoice     ui.ChoiceBox
	explainText    ui.Textbox
	randomButton   ui.Button
	generateButton ui.Button
	cancelButton   ui.Button

	galaxy    *Galaxy
	galaxyMap GalaxyMapView
}

func NewCreateGalaxyMenu() (cgm *CreateGalaxyMenu) {
	cgm = new(CreateGalaxyMenu)
	cgm.Scene.InitBordered()
	cgm.Events().Listen(ui.EV_FOCUS_CHANGED, ui.EV_CHOICE_CHANGED)
	cgm.SetEventHandler(cgm.HandleEvent)

	cgm.Window().SetupBorder("CREATE A WHOLE GALAXY WHY NOT", "")
	windowStyle := ui.DefaultBorderStyle
	windowStyle.TitleJustification = ui.JUSTIFY_CENTER
	cgm.Window().Border.SetStyle(ui.BORDER_STYLE_CUSTOM, windowStyle)

	cgm.Window().AddChild(ui.NewTextbox(vec.Dims{5, 1}, vec.Coord{2, 2}, 1, "Name:", ui.JUSTIFY_LEFT))
	cgm.nameInput.Init(vec.Dims{20, 1}, vec.Coord{9, 2}, 1, 0)
	cgm.nameInput.EnableBorder()

	cgm.Window().AddChild(ui.NewTextbox(vec.Dims{5, 1}, vec.Coord{2, 5}, 1, "Density:", ui.JUSTIFY_LEFT))
	cgm.densityChoice.Init(vec.Dims{20, 1}, vec.Coord{9, 5}, 1, "Sparse", "Normal", "Dense")
	cgm.densityChoice.EnableBorder()

	cgm.Window().AddChild(ui.NewTextbox(vec.Dims{5, 1}, vec.Coord{2, 8}, 1, "Shape:", ui.JUSTIFY_LEFT))
	cgm.shapeChoice.Init(vec.Dims{20, 1}, vec.Coord{9, 8}, 1, "Disk", "Spiral")
	cgm.shapeChoice.EnableBorder()

	cgm.Window().AddChild(ui.NewTextbox(vec.Dims{5, 1}, vec.Coord{2, 11}, 1, "Size:", ui.JUSTIFY_LEFT))
	cgm.sizeChoice.Init(vec.Dims{20, 1}, vec.Coord{9, 11}, 1, "Small", "Medium", "Large")
	cgm.sizeChoice.EnableBorder()

	cgm.Window().AddChildren(&cgm.nameInput, &cgm.densityChoice, &cgm.shapeChoice, &cgm.sizeChoice)

	cgm.explainText.Init(vec.Dims{40, 15}, vec.Coord{2, 35}, 1, "explanations", ui.JUSTIFY_LEFT)
	cgm.explainText.EnableBorder()
	cgm.Window().AddChild(&cgm.explainText)

	cgm.randomButton.Init(vec.Dims{15, 1}, vec.Coord{74, 30}, 1, "Randomize Galaxy", cgm.Randomize)
	cgm.randomButton.EnableBorder()
	cgm.generateButton.Init(vec.Dims{15, 1}, vec.Coord{74, 34}, 1, "Generate the Galaxy as Shown!", cgm.Generate)
	cgm.generateButton.EnableBorder()
	cgm.cancelButton.Init(vec.Dims{15, 1}, vec.Coord{74, 38}, 1, "Return to Main Menu", func() {
		sm := StartMenu{}
		sm.Init()
		tyumi.ChangeScene(&sm)
	})
	cgm.cancelButton.EnableBorder()
	cgm.Window().AddChildren(&cgm.randomButton, &cgm.generateButton, &cgm.cancelButton)

	cgm.galaxyMap.Init(vec.Dims{25, 25}, vec.Coord{69, 0}, ui.BorderDepth, cgm.galaxy)
	cgm.galaxyMap.EnableBorder()
	cgm.galaxyMap.ToggleHighlight()
	cgm.Window().AddChild(&cgm.galaxyMap)

	cgm.nameInput.Focus()
	cgm.Window().SetTabbingOrder(&cgm.nameInput, &cgm.densityChoice, &cgm.shapeChoice, &cgm.sizeChoice, &cgm.randomButton, &cgm.generateButton, &cgm.cancelButton)

	cgm.UpdateExplanation()
	cgm.GeneratePreview()

	return
}

func (cgm *CreateGalaxyMenu) UpdateExplanation() {
	switch cgm.Window().GetFocusedElementID() {
	case cgm.nameInput.ID():
		cgm.explainText.ChangeText("GALAXY NAME:/n/nIt is believed that one of the main ways in which all sentient races of the galaxy are similar is a common desire to name and label the universe. No Galaxy is complete without a name!")
	case cgm.densityChoice.ID():
		cgm.explainText.ChangeText("GALAXY DENSITY:/n/nGalaxies come in all shapes, sizes and consistencies. Some are small and dense, with stars but a stone's throw away from each. Others have stars so spread out that many sentient species decide to never even attempt inter-system travel, instead deciding to focus efforts on art and philosophy and creating better and better tofu-based meat substitutes.")
	case cgm.shapeChoice.ID():
		cgm.explainText.ChangeText("GALAXY SHAPE:/n/nGalaxies, like cookies, come in many different shapes. Some are globular, some are spirals, some are simple disks, and during certain times of year some are shaped like Christmas trees. (Note: currently only disk galaxies are created).")
	case cgm.sizeChoice.ID():
		cgm.explainText.ChangeText("GALAXY SIZE:/n/nAll people need to live in a galaxy, even the very tall. Choose the largest galaxy you can afford to.")
	case cgm.randomButton.ID():
		cgm.explainText.ChangeText("RANDOMIZE:/n/nIndecisive? Stunned by the marvelous array of choices before you? Let me do the work!")
	case cgm.generateButton.ID():
		cgm.explainText.ChangeText("GENERATE:/n/n If this galaxy looks good, we can then generate the galaxy and move on to Ship Selection.")
	case cgm.cancelButton.ID():
		cgm.explainText.ChangeText("CANCEL:/n/n Return to the main menu, discarding everything here.")
	}
}

func (cgm *CreateGalaxyMenu) Randomize() {
	names := []string{"The Biggest Galaxy", "The Galaxy of Terror", "The Lactose Blob", "The Thing Fulla Stars", "Andromeda 2", "Home"}

	cgm.nameInput.ChangeText(util.PickOne(names))
	cgm.densityChoice.RandomizeChoice()
	cgm.shapeChoice.RandomizeChoice()
	cgm.sizeChoice.RandomizeChoice()

	cgm.GeneratePreview()
}

func (cgm *CreateGalaxyMenu) GeneratePreview() {
	var densityFactor int
	var size int

	switch cgm.densityChoice.GetChoiceIndex() {
	case 0: //Sparse
		densityFactor = GAL_SPARSE
	case 1: //Normal
		densityFactor = GAL_NORMAL
	case 2: //Dense
		densityFactor = GAL_DENSE
	}

	switch cgm.sizeChoice.GetChoiceIndex() {
	case 0: //Small
		size = GAL_MIN_RADIUS
	case 1: //Normal
		size = GAL_MIN_RADIUS + (GAL_MAX_RADIUS-GAL_MIN_RADIUS)/2
	case 2: //Dense
		size = GAL_MAX_RADIUS
	}

	cgm.galaxy = NewGalaxy(cgm.nameInput.InputtedText(), size, densityFactor)
	cgm.galaxyMap.galaxy = cgm.galaxy
	cgm.galaxyMap.Updated = true
}

func (cgm *CreateGalaxyMenu) Generate() {
	if cgm.nameInput.InputtedText() == "" {
		cgm.OpenDialog(NewSimpleCommDialog("You must give your galaxy a name before you can continue!"))
	} else {
		tyumi.ChangeScene(NewShipCreateMenu(cgm.galaxy))
	}
}

func (cgm *CreateGalaxyMenu) HandleEvent(game_event event.Event) (event_handled bool) {
	switch game_event.ID() {
	case ui.EV_FOCUS_CHANGED:
		cgm.UpdateExplanation()
		event_handled = true
	case ui.EV_CHOICE_CHANGED:
		cgm.GeneratePreview()
		event_handled = true
	}

	return
}
