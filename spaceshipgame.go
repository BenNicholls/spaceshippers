package main

import "github.com/bennicholls/burl/ui"
import "github.com/bennicholls/burl/util"
import "github.com/veandco/go-sdl2/sdl"
import "fmt"
import "strings"

//time values in DIGITAL SECONDS. One digital day = 100000 seconds, which is 14% longer than a regular day.
const (
	MINUTE int = 100
	HOUR   int = 10000
	DAY    int = 100000
)

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

	PlayerShip *Ship

}

func NewSpaceshipGame() *SpaceshipGame {
	sg := new(SpaceshipGame)

	sg.SpaceTime = 0
	sg.SimSpeed = 1
	sg.PlayerShip = NewShip("The Undestructable")
	sg.starFrequency = 20


	sg.SetupUI()
	sg.UpdateSpeedUI()
	sg.initStarField()

	return sg
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
	sg.shipstatus.Add(ui.NewTextbox(26, 1, 0, 11, 0, false, true, "The Unsinkable"))

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

//n = crew index
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
		case sdl.K_UP:
			sg.output.ScrollUp()
		case sdl.K_DOWN:
			sg.output.ScrollDown()
		case sdl.K_PAGEUP:
			if sg.SimSpeed < 4 {
				sg.SimSpeed++
				sg.UpdateSpeedUI()
			}
		case sdl.K_PAGEDOWN:
			if sg.SimSpeed > 0 {
				sg.SimSpeed--
				sg.UpdateSpeedUI()
			}
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
	sg.window.Render()
	sg.crew.Render()
}

func (sg *SpaceshipGame) Execute() {
	sg.output.Append("")
	sg.output.Append(">>> " + sg.input.GetText())
	sg.output.Append("")
	switch strings.ToLower(sg.input.GetText()) {
	case "status":
		for _, r := range sg.PlayerShip.Rooms {
			sg.output.Append(r.GetStatus())
		}
	case "help":
		sg.output.Append("S.C.I.P.P.I.E. is your AI helper. Give him one of the following commands, and he'll get 'r done!")
		sg.output.Append("   status     prints ship room status")
		sg.output.Append("   help       prints a mysterious menu")
	default:
		sg.output.Append("I do not understand that command, you dummo. Try \"help\"")
	}
	sg.output.ScrollToBottom()
}

func (sg *SpaceshipGame) AddMessage(s string) {
	sg.output.Append(s)
	sg.output.ScrollToBottom()
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