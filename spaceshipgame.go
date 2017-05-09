package main

import "github.com/bennicholls/burl/ui"
import "github.com/bennicholls/burl/util"
import "github.com/bennicholls/burl/core"
import "github.com/veandco/go-sdl2/sdl"
import "fmt"

//time values in DIGITAL SECONDS. One digital day = 100000 seconds, which is 14% longer than a regular day.
const (
	MINUTE int = 100
	HOUR   int = 10000
	DAY    int = 100000
)

//load some tile data
var TILE_FLOOR = core.LoadTileData("Floor", true, true, 0xB0, 0xFF444444)
var TILE_WALL = core.LoadTileData("Wall", false, false, 0x23, 0xFF888888)
var TILE_DOOR = core.LoadTileData("Door", true, false, 0xF0, 0xFF888888)

type SpaceshipGame struct {

	//ui stuff
	window *ui.Container
	input *ui.Inputbox
	output *ui.List
	crew *ui.Container
	shipstatus *ui.Container
	missiontime *ui.Textbox
	speeddisplay *ui.TileView
	menubar *ui.Container
	shipdisplay *ui.TileView
	crewUI []*ui.Container //one container for each crew member.

	//Time Globals.
	SpaceTime int //measured in Standard Galactic Seconds
	SimSpeed int  //4 speeds, plus pause (0)

	starField []int
	starFrequency int

	viewX, viewY int

	PlayerShip *Ship

}

func NewSpaceshipGame() *SpaceshipGame {
	sg := new(SpaceshipGame)

	sg.SpaceTime = 0
	sg.SimSpeed = 1
	sg.PlayerShip = NewShip("The Undestructable")
	sg.PlayerShip.AddRoom(NewRoom("Engineering", 5, 8, 5, 8, 700, 1000))
	sg.PlayerShip.AddRoom(NewRoom("Messhall", 15, 5, 6, 6, 1000, 500))
	sg.PlayerShip.AddRoom(NewRoom("Medbay", 9, 5, 6, 6, 1000, 700))
	sg.PlayerShip.AddRoom(NewRoom("Quarters 1", 15, 13, 6, 6, 900, 500))
	sg.PlayerShip.AddRoom(NewRoom("Quarters 2", 9, 13, 6, 6, 900, 500))
	sg.PlayerShip.AddRoom(NewRoom("Hallway", 9, 10, 12, 4, 0, 500))
	sg.starFrequency = 20

	sg.SetupUI()
	sg.UpdateSpeedUI()
	sg.initStarField()
	sg.CenterShip()

	return sg
}

//Centers the map of the ship in the main view.
func (sg *SpaceshipGame) CenterShip() {
	sg.viewX = sg.shipdisplay.Width/2 - sg.PlayerShip.Width/2 - sg.PlayerShip.X
	sg.viewY = sg.shipdisplay.Height/2 - sg.PlayerShip.Height/2 - sg.PlayerShip.Y
}

func (sg *SpaceshipGame) UpdateSpeedUI() {
	sg.speeddisplay.Clear()
	for i := 0; i < 4; i++ {
		if i < sg.SimSpeed {
			sg.speeddisplay.Draw(i, 0, 0x10, 0xFFFFFFFF, 0x00000000)
		} else {
			sg.speeddisplay.Draw(i, 0, 0x5F, 0xFFFFFFFF, 0x00000000)
		}
	}
}

func (sg *SpaceshipGame) SetupUI() {
	sg.window = ui.NewContainer(80, 45, 0, 0, 0, false)

	sg.missiontime = ui.NewTextbox(8, 1, 1, 1, 1, true, true, "")
	sg.speeddisplay = ui.NewTileView(8, 1, 1, 3, 1, true)
	sg.menubar = ui.NewContainer(69, 1, 10, 1, 1, true)

	sg.shipdisplay = ui.NewTileView(80, 28, 0, 3, 0, false)

	sg.shipstatus = ui.NewContainer(26, 12, 1, 32, 1, true)
	sg.shipstatus.Add(ui.NewTextbox(26, 1, 0, 11, 0, false, true, "The USS Prototype"))

	sg.input = ui.NewInputbox(51, 1, 28, 43, 2, true)
	sg.input.ToggleFocus()
	sg.input.SetTitle("SCIPPIE V6.18")
	sg.output = ui.NewList(51, 10, 28, 32, 1, true, "The Ship Computer Interactive Parameter Parser/Interface Entity, or SCIPPIE, is your computerized second in command. Ask questions, give commands and observe your ship through the high-tech text-tacular wonders of 38th century UI technology! Ask SCIPPIE a question, or give him a command!")
	sg.output.ToggleHighlight()

	sg.crew = ui.NewContainer(26, 18, 27, 13, 3, true)
	sg.crew.SetTitle("Crew Roster")
	sg.crew.SetVisibility(false)
	sg.crewUI = make([]*ui.Container, 6)

	sg.window.Add(sg.input, sg.output, sg.shipstatus, sg.shipdisplay, sg.speeddisplay, sg.missiontime, sg.menubar)
}

func (sg *SpaceshipGame) UpdateCrewUI() {
	w, _ := sg.crew.Dims()
	sg.crew.ClearElements()
	for i := range sg.PlayerShip.Crew {
		sg.crewUI[i] = ui.NewContainer(w, 3, 0, i*3, 1, false)
		name := ui.NewProgressBar(w, 1, 0, 0, 0, false, false, sg.PlayerShip.Crew[i].Name, 0xFFFF0000)
		name.SetProgress(sg.PlayerShip.Crew[i].Awakeness.GetPct())
		sg.crewUI[i].Add(name)
		sg.crewUI[i].Add(ui.NewTextbox(w-2, 1, 2, 1, 0, false, false, "is "+sg.PlayerShip.Crew[i].GetStatus()))
		jobstring := ""
		if sg.PlayerShip.Crew[i].CurrentTask != nil {
			jobstring = "is " + sg.PlayerShip.Crew[i].CurrentTask.GetDescription()
		}
		sg.crewUI[i].Add(ui.NewTextbox(w-2, 1, 2, 2, 0, false, false, jobstring))
		sg.crew.Add(sg.crewUI[i])
	}
}

func (sg *SpaceshipGame) HandleKeypress(key sdl.Keycode) {
	if util.ValidText(rune(key)) {
		sg.input.InsertText(rune(key))
	} else {
		switch key {
		case sdl.K_RETURN:
			sg.Execute()
			sg.input.Reset()
		case sdl.K_BACKSPACE:
			sg.input.Delete()
		case sdl.K_SPACE:
			sg.input.Insert(" ")
		case sdl.K_PAGEUP:
			sg.output.ScrollUp()
		case sdl.K_PAGEDOWN:
			sg.output.ScrollDown()
		case sdl.K_HOME:
			sg.CenterShip()
		case sdl.K_KP_PLUS:
			if sg.SimSpeed < 4 {
				sg.SimSpeed++
				sg.UpdateSpeedUI()
			}
		case sdl.K_KP_MINUS:
			if sg.SimSpeed > 0 {
				sg.SimSpeed--
				sg.UpdateSpeedUI()
			}
		case sdl.K_UP:
			sg.viewY -= 1
		case sdl.K_DOWN:
			sg.viewY += 1
		case sdl.K_LEFT:
			sg.viewX -= 1
		case sdl.K_RIGHT:
			sg.viewX += 1
		case sdl.K_F1:
			sg.crew.ToggleVisible()
		}
	}
}

func (sg *SpaceshipGame) Update() {

	for i := 0; i < sg.GetIncrement(); i++ {
		sg.SpaceTime++

		sg.PlayerShip.Update(sg.SpaceTime)

		for i := range sg.PlayerShip.Crew {
			sg.PlayerShip.Crew[i].Update()
		}

		if sg.SpaceTime%100 == 0 {
			sg.shiftStarField()
		}
	}

	if sg.crew.IsVisible() {
		sg.UpdateCrewUI()
	}
	sg.missiontime.ChangeText(fmt.Sprintf("%.4d", sg.SpaceTime/100000) + "d:" + fmt.Sprintf("%.1d", (sg.SpaceTime/10000)%10) + "h:" + fmt.Sprintf("%.2d", (sg.SpaceTime/100)%100) + "m:" + fmt.Sprintf("%.2d", sg.SpaceTime%100) + "s")

}

func (sg *SpaceshipGame) Render() {
	sg.DrawStarfield()

	w, h := sg.PlayerShip.ShipMap.Dims()
	x, y := 0, 0

	for i := 0; i < w*h; i++ {
		//shipdisplay-space coords
		x = i%w + sg.viewX
		y = i/w + sg.viewY

		if util.CheckBounds(x, y, sg.shipdisplay.Width, sg.shipdisplay.Height) {
			t := sg.PlayerShip.ShipMap.GetTile(i%w, i/w)
			if t.TileType() != 0 {
				tv := t.GetVisuals()
				sg.shipdisplay.Draw(x, y, tv.Glyph, tv.ForeColour, 0xFF000000)
			}

			if sg.PlayerShip.ShipMap.GetEntity(i%w, i/w) != nil {
				e := sg.PlayerShip.ShipMap.GetEntity(i%w, i/w)
				sg.shipdisplay.Draw(x, y, e.GetVisuals().Glyph, e.GetVisuals().ForeColour, 0xFF000000)
			}
		}
	}

	sg.window.Render()
	sg.crew.Render()
}

func (sg *SpaceshipGame) GetIncrement() int {
	switch sg.SimSpeed {
	case 0:
		return 0
	case 1:
		return 1
	case 2:
		return 10
	case 3:
		return 100
	case 4:
		return 1000
	}

	return 0
}