package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"strings"
	"time"

	"github.com/bennicholls/delveengine/console"
	"github.com/bennicholls/delveengine/ui"
	"github.com/bennicholls/delvetown/util"
	"github.com/veandco/go-sdl2/sdl"
)

var window *ui.Container
var input *ui.Inputbox
var output *ui.List
var crew *ui.Container
var shipstatus *ui.Container
var missiontime *ui.Textbox
var speeddisplay *ui.TileView
var menubar *ui.Container
var shipdisplay *ui.TileView
var crewUI []*ui.Container //one container for each crew member.

//Time Globals.
var SpaceTime int //measured in Standard Galactic Seconds
var SimSpeed int  //4 speeds, plus pause (0)

//time values in DIGITAL SECONDS. One digital day = 100000 seconds, which is 14% longer than a regular day.
const (
	MINUTE int = 100
	HOUR   int = 10000
	DAY    int = 100000
)

var PlayerShip *Ship

func main() {

	runtime.LockOSThread()
	rand.Seed(time.Now().UTC().UnixNano())

	var event sdl.Event

	err := console.Setup(80, 45, "res/curses24x24.bmp", "res/DelveFont12x24.bmp", "Spaceshippers")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer console.Cleanup()

	console.SetFullscreen()

	SetupUI()

	
	initStarField()

	SpaceTime = 0
	SimSpeed = 1
	UpdateSpeedUI()

	PlayerShip = NewShip("The Undestructable")

	running := true

	for running {
		for event = sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				running = false
			// case *sdl.MouseMotionEvent:
			// 	fmt.Printf("[%d ms] MouseMotion\ttype:%d\tid:%d\tx:%d\ty:%d\txrel:%d\tyrel:%d\n",
			// 		t.Timestamp, t.Type, t.Which, t.X, t.Y, t.XRel, t.YRel)
			// case *sdl.MouseButtonEvent:
			// 	fmt.Printf("[%d ms] MouseButton\ttype:%d\tid:%d\tx:%d\ty:%d\tbutton:%d\tstate:%d\n",
			// 		t.Timestamp, t.Type, t.Which, t.X, t.Y, t.Button, t.State)
			// case *sdl.MouseWheelEvent:
			// 	fmt.Printf("[%d ms] MouseWheel\ttype:%d\tid:%d\tx:%d\ty:%d\n",
			// 		t.Timestamp, t.Type, t.Which, t.X, t.Y)
			case *sdl.KeyUpEvent:
				//fmt.Printf("[%d ms] Keyboard\ttype:%d\tsym:%c\tmodifiers:%d\tstate:%d\trepeat:%d\n",
				//	t.Timestamp, t.Type, t.Keysym.Sym, t.Keysym.Mod, t.State, t.Repeat)
				HandleKeypress(t.Keysym.Sym)
			}
		}

		Update()

		Render()
		console.Render()
	}
}

func SetupUI() {
	window = ui.NewContainer(80, 45, 0, 0, 0, false)

	missiontime = ui.NewTextbox(8, 1, 1, 1, 1, true, true, "")
	speeddisplay = ui.NewTileView(8, 1, 1, 3, 1, true)
	menubar = ui.NewContainer(69, 1, 10, 1, 1, true)

	shipdisplay = ui.NewTileView(80, 28, 0, 3, 0, false)

	shipstatus = ui.NewContainer(26, 12, 1, 32, 1, true)
	shipstatus.Add(ui.NewTextbox(26, 1, 0, 11, 0, false, true, "The Unsinkable"))

	input = ui.NewInputbox(51, 1, 28, 43, 2, true)
	input.ToggleFocus()
	input.SetTitle("SCIPPIE V6.18")
	output = ui.NewList(51, 10, 28, 32, 1, true, "The Ship Computer Interactive Parameter Parser/Interface Entity, or SCIPPIE, is your computerized second in command. Ask questions, give commands and observe your ship through the high-tech text-tacular wonders of 38th century UI technology! Ask SCIPPIE a question, or give him a command!")
	output.ToggleHighlight()

	crew = ui.NewContainer(26, 18, 27, 13, 3, true)
	crew.SetTitle("Crew Roster")
	crew.SetVisibility(false)
	crewUI = make([]*ui.Container, 6)

	window.Add(input, output, shipstatus, shipdisplay, speeddisplay, missiontime, menubar)
}

//n = crew index
func UpdateCrewUI() {
	w, _ := crew.Dims()
	crew.ClearElements()
	for i := range PlayerShip.Crew {
		crewUI[i] = ui.NewContainer(w, 3, 0, i*3, 1, false)
		name := ui.NewProgressBar(w, 1, 0, 0, 0, false, false, PlayerShip.Crew[i].Name, 0xFFFF0000)
		name.SetProgress(PlayerShip.Crew[i].Awakeness.GetPct())
		crewUI[i].Add(name)
		crewUI[i].Add(ui.NewTextbox(w-2, 1, 2, 1, 0, false, false, "is "+PlayerShip.Crew[i].GetStatus()))
		jobstring := ""
		if PlayerShip.Crew[i].CurrentTask != nil {
			jobstring = "is " + PlayerShip.Crew[i].CurrentTask.GetDescription()
		}
		crewUI[i].Add(ui.NewTextbox(w-2, 1, 2, 2, 0, false, false, jobstring))
		crew.Add(crewUI[i])
	}
}

func UpdateSpeedUI() {
	speeddisplay.Clear()
	for i := 0; i < 4; i++ {
		if i < SimSpeed {
			speeddisplay.Draw(i, 0, 0x10, 0xFFFFFFFF, 0x00000000)
		} else {
			speeddisplay.Draw(i, 0, 0x5F, 0xFFFFFFFF, 0x00000000)
		}
	}
}

func HandleKeypress(key sdl.Keycode) {
	if util.ValidText(rune(key)) {
		input.InsertText(rune(key))
	} else {
		switch key {
		case sdl.K_RETURN:
			Execute()
			input.Reset()
		case sdl.K_BACKSPACE:
			input.Delete()
		case sdl.K_SPACE:
			input.Insert(" ")
		case sdl.K_UP:
			output.ScrollUp()
		case sdl.K_DOWN:
			output.ScrollDown()
		case sdl.K_PAGEUP:
			if SimSpeed < 4 {
				SimSpeed++
				UpdateSpeedUI()
			}
		case sdl.K_PAGEDOWN:
			if SimSpeed > 0 {
				SimSpeed--
				UpdateSpeedUI()
			}
		case sdl.K_F1:
			crew.ToggleVisible()
		}
	}
}

func Update() {

	for i := 0; i < GetIncrement(); i++ {
		SpaceTime++

		PlayerShip.Update()

		for i := range PlayerShip.Crew {
			PlayerShip.Crew[i].Update()
		}

		if SpaceTime%100 == 0 {
			shiftStarField()
		}
	}

	if crew.IsVisible() {
		UpdateCrewUI()
	}
	missiontime.ChangeText(fmt.Sprintf("%.4d", SpaceTime/100000) + "d:" + fmt.Sprintf("%.1d", (SpaceTime/10000)%10) + "h:" + fmt.Sprintf("%.2d", (SpaceTime/100)%100) + "m:" + fmt.Sprintf("%.2d", SpaceTime%100) + "s")

}

func Render() {
	DrawStarfield()
	PlayerShip.Draw()
	window.Render()
	crew.Render()
}

func Execute() {
	output.Append("")
	output.Append(">>> " + input.GetText())
	output.Append("")
	switch strings.ToLower(input.GetText()) {
	case "status":
		for _, r := range PlayerShip.Rooms {
			r.PrintStatus()
		}
	case "help":
		output.Append("S.C.I.P.P.I.E. is your AI helper. Give him one of the following commands, and he'll get 'r done!")
		output.Append("   status     prints ship room status")
		output.Append("   help       prints a mysterious menu")
	default:
		output.Append("I do not understand that command, you dummo. Try \"help\"")
	}
	output.ScrollToBottom()
}

func AddMessage(s string) {
	output.Append(s)
	output.ScrollToBottom()
}

func GetIncrement() int {
	switch SimSpeed {
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
