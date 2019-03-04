package main

import (
	"math/rand"

	"github.com/bennicholls/burl-E/burl"
	"github.com/veandco/go-sdl2/sdl"
)

type CreateGalaxyMenu struct {
	burl.StatePrototype

	nameInput      *burl.Inputbox
	densityChoice  *burl.ChoiceBox //choice between some pre-defined densities
	shapeChoice    *burl.ChoiceBox //choice between blob, spiral, maybe more esoteric shapes??
	sizeChoice     *burl.ChoiceBox
	explainText    *burl.Textbox
	randomButton   *burl.Button
	generateButton *burl.Button
	cancelButton   *burl.Button

	focusedField burl.UIElem

	galaxyMap *GalaxyMapView

	galaxy *Galaxy
}

func NewCreateGalaxyMenu() (cgm *CreateGalaxyMenu) {
	cgm = new(CreateGalaxyMenu)
	cgm.InitWindow(true)
	cgm.Window.SetTitle("CREATE A WHOLE GALAXY WHY NOT")

	cgm.Window.Add(burl.NewTextbox(5, 1, 2, 2, 1, false, false, "Name:"))
	cgm.nameInput = burl.NewInputbox(20, 1, 9, 2, 1, true)
	cgm.Window.Add(burl.NewTextbox(5, 1, 2, 5, 1, false, false, "Density:"))
	cgm.densityChoice = burl.NewChoiceBox(20, 1, 9, 5, 2, true, burl.HORIZONTAL, "Sparse", "Normal", "Dense")
	cgm.Window.Add(burl.NewTextbox(5, 1, 2, 8, 1, false, false, "Shape:"))
	cgm.shapeChoice = burl.NewChoiceBox(20, 1, 9, 8, 1, true, burl.HORIZONTAL, "Disk", "Spiral")
	cgm.Window.Add(burl.NewTextbox(5, 1, 2, 11, 1, false, false, "Size:"))
	cgm.sizeChoice = burl.NewChoiceBox(20, 1, 9, 11, 2, true, burl.HORIZONTAL, "Small", "Medium", "Large")

	cgm.explainText = burl.NewTextbox(40, 15, 2, 35, 1, true, false, "explanations")

	cgm.randomButton = burl.NewButton(15, 1, 74, 30, 2, true, true, "Randomize Galaxy")
	cgm.generateButton = burl.NewButton(15, 1, 74, 34, 1, true, true, "Generate the Galaxy as Shown!")
	cgm.cancelButton = burl.NewButton(15, 1, 74, 38, 2, true, true, "Return to Main Menu")

	cgm.galaxyMap = NewGalaxyMapView(25, 69, 0, 0, true, cgm.galaxy)

	cgm.Window.Add(cgm.nameInput, cgm.densityChoice, cgm.shapeChoice, cgm.generateButton, cgm.explainText, cgm.cancelButton, cgm.galaxyMap, cgm.randomButton, cgm.sizeChoice)

	cgm.nameInput.SetTabID(1)
	cgm.densityChoice.SetTabID(2)
	cgm.shapeChoice.SetTabID(3)
	cgm.sizeChoice.SetTabID(4)
	cgm.randomButton.SetTabID(5)
	cgm.generateButton.SetTabID(6)
	cgm.cancelButton.SetTabID(7)

	cgm.focusedField = cgm.nameInput
	cgm.focusedField.ToggleFocus()
	cgm.UpdateExplanation()

	cgm.Generate()

	return
}

func (cgm *CreateGalaxyMenu) UpdateExplanation() {
	switch cgm.focusedField {
	case cgm.nameInput:
		cgm.explainText.ChangeText("GALAXY NAME:/n/nIt is believed that one of the main ways in which all sentient races of the galaxy are similar is a common desire to name and label the universe. No Galaxy is complete without a name!")
	case cgm.densityChoice:
		cgm.explainText.ChangeText("GALAXY DENSITY:/n/nGalaxies come in all shapes, sizes and consistencies. Some are small and dense, with stars but a stone's throw away from each. Others have stars so spread out that many sentient species decide to never even attempt inter-system travel, instead deciding to focus efforts on art and philosophy and creating better and better tofu-based meat substitutes.")
	case cgm.shapeChoice:
		cgm.explainText.ChangeText("GALAXY SHAPE:/n/nGalaxies, like cookies, come in many different shapes. Some are globular, some are spirals, some are simple disks, and during certain times of year some are shaped like Christmas trees. (Note: currently only disk galaxies are created).")
	case cgm.sizeChoice:
		cgm.explainText.ChangeText("GALAXY SIZE:/n/nAll people need to live in a galaxy, even the very tall. Choose the largest galaxy you can afford to.")
	case cgm.randomButton:
		cgm.explainText.ChangeText("RANDOMIZE:/n/nIndecisive? Stunned by the marvelous array of choices before you? Let me do the work!")
	case cgm.generateButton:
		cgm.explainText.ChangeText("GENERATE:/n/n If this galaxy looks good, we can then generate the galaxy and move on to Ship Selection.")
	case cgm.cancelButton:
		cgm.explainText.ChangeText("CANCEL:/n/n Return to the main menu, discarding everything here.")
	}
}

func (cgm *CreateGalaxyMenu) Randomize() {
	names := []string{"The Biggest Galaxy", "The Galaxy of Terror", "The Lactose Blob", "The Thing Fulla Stars", "Andromeda 2", "Home"}

	cgm.nameInput.ChangeText(names[rand.Intn(len(names))])
	cgm.densityChoice.RandomizeChoice()
	cgm.shapeChoice.RandomizeChoice()
	cgm.sizeChoice.RandomizeChoice()

	cgm.Generate()
}

func (cgm *CreateGalaxyMenu) Generate() {

	var densityFactor int
	var size int

	switch cgm.densityChoice.GetChoice() {
	case 0: //Sparse
		densityFactor = GAL_SPARSE
	case 1: //Normal
		densityFactor = GAL_NORMAL
	case 2: //Dense
		densityFactor = GAL_DENSE
	}

	switch cgm.sizeChoice.GetChoice() {
	case 0: //Small
		size = GAL_MIN_RADIUS
	case 1: //Normal
		size = GAL_MIN_RADIUS + (GAL_MAX_RADIUS-GAL_MIN_RADIUS)/2
	case 2: //Dense
		size = GAL_MAX_RADIUS
	}

	cgm.galaxy = NewGalaxy(cgm.nameInput.GetText(), size, densityFactor)
	cgm.galaxyMap.galaxy = cgm.galaxy
	cgm.galaxyMap.DrawGalaxyMap()
}

func (cgm *CreateGalaxyMenu) HandleKeypress(key sdl.Keycode) {
	cgm.focusedField.HandleKeypress(key)

	switch key {
	case sdl.K_UP:
		burl.PushEvent(burl.NewEvent(burl.EV_TAB_FIELD, "-"))
	case sdl.K_DOWN, sdl.K_TAB:
		burl.PushEvent(burl.NewEvent(burl.EV_TAB_FIELD, "+"))
	default:

	}
}

func (cgm *CreateGalaxyMenu) HandleEvent(e *burl.Event) {
	switch e.ID {
	case burl.EV_TAB_FIELD:
		cgm.focusedField.ToggleFocus()
		if e.Message == "+" {
			cgm.focusedField = cgm.Window.FindNextTab(cgm.focusedField)
		} else {
			cgm.focusedField = cgm.Window.FindPrevTab(cgm.focusedField)
		}
		cgm.focusedField.ToggleFocus()
		cgm.UpdateExplanation()
	case burl.EV_BUTTON_PRESS:
		if e.Caller == cgm.randomButton {
			cgm.Randomize()
		}
	case burl.EV_LIST_CYCLE:
		switch e.Caller {
		case cgm.shapeChoice, cgm.sizeChoice, cgm.densityChoice:
			cgm.Generate()
		}
	case burl.EV_ANIMATION_DONE:
		if e.Caller == cgm.generateButton {
			if cgm.nameInput.GetText() == "" {
				burl.OpenDialog(NewCommDialog("", "", "", "You must give your galaxy a name before you can continue!"))
			} else {
				burl.ChangeState(NewShipCreateMenu(cgm.galaxy))
			}
		} else if e.Caller == cgm.cancelButton {
			burl.ChangeState(NewStartMenu())
		}
	}
}
