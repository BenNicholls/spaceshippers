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

type zoomLevel int

const (
	zoom_GALAXY zoomLevel = iota
	zoom_LOCAL
)

type GalaxyMapView struct {
	burl.TileView

	galaxy *Galaxy

	cursor    burl.Coord
	highlight *burl.PulseAnimation

	zoom zoomLevel

	//for local map drawing
	localZoom   int
	localFocus  Locatable
	systemFocus *StarSystem
}

func NewGalaxyMapView(w, x, y, z int, bord bool, g *Galaxy) (gmv *GalaxyMapView) {
	gmv = new(GalaxyMapView)
	gmv.TileView = *burl.NewTileView(w, w, x, y, z, bord)
	gmv.galaxy = g

	gmv.highlight = burl.NewPulseAnimation(1, 1, 0, 1, 1, 100, 10, true)
	gmv.AddAnimation(gmv.highlight)

	return
}

func (gmv *GalaxyMapView) ZoomIn() {
	if gmv.zoom == zoom_GALAXY {
		gmv.zoom = zoom_LOCAL
		gmv.DrawLocalMap()
		gmv.ToggleHighlight()
	} else if gmv.zoom == zoom_LOCAL {
		if gmv.localZoom < 7 {
			gmv.localZoom += 1
			gmv.DrawLocalMap()
		}
	}
}

func (gmv *GalaxyMapView) ZoomOut() {
	if gmv.zoom == zoom_LOCAL {
		if gmv.localZoom == 0 {
			gmv.zoom = zoom_GALAXY
			gmv.DrawGalaxyMap()
			gmv.ToggleHighlight()
		} else {
			gmv.localZoom -= 1
			gmv.DrawLocalMap()
		}
	}
}

func (gmv *GalaxyMapView) DrawMap() {
	switch gmv.zoom {
	case zoom_GALAXY:
		gmv.DrawGalaxyMap()
	case zoom_LOCAL:
		gmv.DrawLocalMap()
	default:
		burl.LogError("No draw function for selected zoom level.")
	}
}

func (gmv *GalaxyMapView) DrawGalaxyMap() {
	if gmv.galaxy != nil {
		w, h := gmv.galaxy.Dims()
		for i := 0; i < w*h; i++ {
			x, y := i%w, i/w
			s := gmv.galaxy.GetSector(x, y)
			bright := burl.Lerp(0, 255, s.Density, 100)
			g := burl.GLYPH_FILL_SPARSE
			if bright == 0 {
				g = burl.GLYPH_NONE
			}
			gmv.Draw(x, y, g, burl.MakeOpaqueColour(bright, bright, bright), burl.COL_BLACK)
		}
	}
}

func (gmv *GalaxyMapView) calcLocalMapCoord(c Coordinates) (mc burl.Coord) {
	w, h := gmv.Dims()
	gFactor := gmv.localZoomFactor()

	//camera computation
	cX := gmv.localFocus.GetCoords().Local.X - (gFactor * float64(w) / 2)
	cY := gmv.localFocus.GetCoords().Local.Y - (gFactor * float64(h) / 2)

	mc.X = int((c.Local.X - cX) / gFactor)
	mc.Y = int((c.Local.Y - cY) / gFactor)

	return
}

func (gmv *GalaxyMapView) localZoomFactor() float64 {
	w, _ := gmv.Dims()
	return coord_LOCAL_MAX / float64(w) / float64(burl.Pow(2, gmv.localZoom))
}

func (gmv *GalaxyMapView) DrawLocalMap() {
	gmv.Reset()
	smc := gmv.calcLocalMapCoord(gmv.systemFocus.Star.GetCoords())

	//draw system things!
	for _, p := range gmv.systemFocus.Planets {
		//draw orbits. TODO: some way of culling orbit drawing. currently drawing all of them
		gmv.DrawCircle(smc.X, smc.Y, burl.RoundFloatToInt(p.oDistance/gmv.localZoomFactor()), burl.GLYPH_PERIOD, 0xFF114411, burl.COL_BLACK)
		gmv.DrawObject(p, burl.Visuals{int(rune(p.Name[0])), 0xFF825814, burl.COL_BLACK})
	}

	gmv.DrawObject(gmv.systemFocus.Star, burl.Visuals{burl.GLYPH_STAR, burl.COL_YELLOW, burl.COL_BLACK})
}

func (gmv *GalaxyMapView) DrawObject(l Locatable, v burl.Visuals) {
	var c burl.Coord

	switch gmv.zoom {
	case zoom_GALAXY:
		c = l.GetCoords().Sector
	case zoom_LOCAL:
		c = gmv.calcLocalMapCoord(l.GetCoords())
	}

	w, h := gmv.Dims()

	if burl.CheckBounds(c.X, c.Y, w, h) {
		gmv.Draw(c.X, c.Y, v.Glyph, v.ForeColour, v.BackColour)
	}
}

func (gmv *GalaxyMapView) ToggleHighlight() {
	gmv.highlight.MoveTo(gmv.cursor.X, gmv.cursor.Y)
	gmv.highlight.Toggle()
}

func (gmv *GalaxyMapView) HandleKeypress(key sdl.Keycode) {
	//generic
	switch key {
	case sdl.K_PAGEUP:
		gmv.ZoomIn()
	case sdl.K_PAGEDOWN:
		gmv.ZoomOut()
	}

	//zoom-level specific
	switch gmv.zoom {
	case zoom_GALAXY:
		switch key {
		case sdl.K_UP:
			gmv.MoveCursor(0, -1)
		case sdl.K_DOWN:
			gmv.MoveCursor(0, 1)
		case sdl.K_LEFT:
			gmv.MoveCursor(-1, 0)
		case sdl.K_RIGHT:
			gmv.MoveCursor(1, 0)
		}
	}
}

func (gmv *GalaxyMapView) MoveCursor(dx, dy int) {
	w, h := gmv.Dims()
	if burl.CheckBounds(gmv.cursor.X+dx, gmv.cursor.Y+dy, w, h) {
		gmv.cursor.Move(dx, dy)
		gmv.highlight.Move(dx, dy, 0)
	}
}
