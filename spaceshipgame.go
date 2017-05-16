package main

import "github.com/bennicholls/burl/ui"
import "github.com/bennicholls/burl/util"
import "github.com/bennicholls/burl/core"
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
	shipstatus *ui.Container
	missiontime *ui.Textbox
	speeddisplay *ui.TileView
	shipdisplay *ui.TileView

	//top menu. contains buttons for submenus
	menubar *ui.Container
	crewMenuButton *ui.Button
	shipMenuButton *ui.Button
	roomMenuButton *ui.Button
	mainMenuButton *ui.Button
	starchartMenuButton *ui.Button
	scippieMenuButton *ui.Button

	//submenus. these are stored always for fast switching.
	//crew menu (F1)
	crewMenu *ui.Container
	crewList *ui.List
	crewDetails *ui.Container

	activeMenu ui.UIElem

	//Time Globals.
	SpaceTime int //measured in Standard Galactic Seconds
	SimSpeed int  //4 speeds, plus pause (0)
	paused bool

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

	sg.SetupUI() //must be done after ship setup
	sg.UpdateSpeedUI()
	sg.initStarField()
	sg.CenterShip()

	return sg
}

//Centers the map of the ship in the main view.
func (sg *SpaceshipGame) CenterShip() {
	displayWidth, displayHeight := sg.shipdisplay.Dims()
	sg.viewX = displayWidth/2 - sg.PlayerShip.Width/2 - sg.PlayerShip.X
	sg.viewY = displayHeight/2 - sg.PlayerShip.Height/2 - sg.PlayerShip.Y
	if sg.activeMenu != nil {
		w, _ := sg.activeMenu.Dims()
		sg.viewX -= w/2
	}
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

	sg.menubar = ui.NewContainer(69, 1, 12, 1, 1, false)
	sg.crewMenuButton = ui.NewButton(9, 1, 1, 0, 1, true, true, "Crew")
	sg.shipMenuButton = ui.NewButton(9, 1, 12, 0, 1, true, true, "Ship")
	sg.roomMenuButton = ui.NewButton(9, 1, 23, 0, 1, true, true, "Room")
	sg.starchartMenuButton = ui.NewButton(9, 1, 34, 0, 1, true, true, "Star Chart")
	sg.scippieMenuButton = ui.NewButton(9, 1, 45, 0, 1, true, true, "S.C.I.P.P.I.E.")
	sg.mainMenuButton = ui.NewButton(9, 1, 56, 0, 1, true, true, "Main  Menu")
	sg.menubar.Add(sg.crewMenuButton, sg.shipMenuButton, sg.roomMenuButton, sg.roomMenuButton, sg.starchartMenuButton, sg.scippieMenuButton, sg.mainMenuButton)

	sg.shipdisplay = ui.NewTileView(80, 28, 0, 3, 0, false)

	sg.shipstatus = ui.NewContainer(26, 12, 1, 32, 1, true)
	sg.shipstatus.Add(ui.NewTextbox(26, 1, 0, 11, 0, false, true, "The USS Prototype"))

	sg.input = ui.NewInputbox(50, 1, 15, 27, 2, true)
	sg.input.ToggleFocus()
	sg.input.SetVisibility(false)
	sg.input.SetTitle("SCIPPIE V6.18")

	sg.output = ui.NewList(51, 12, 28, 32, 1, true, "The Ship Computer Interactive Parameter Parser/Interface Entity, or SCIPPIE, is your computerized second in command. Ask questions, give commands and observe your ship through the high-tech text-tacular wonders of 38th century UI technology! Ask SCIPPIE a question, or give him a command!")
	sg.output.ToggleHighlight()

	sg.SetupCrewMenu()

	sg.window.Add(sg.input, sg.output, sg.shipstatus, sg.shipdisplay, sg.speeddisplay, sg.missiontime, sg.menubar)
}

func (sg *SpaceshipGame) Update() {

	//simulation!
	for i := 0; i < sg.GetIncrement(); i++ {
		sg.SpaceTime++
		sg.PlayerShip.Update(sg.SpaceTime)

		for i := range sg.PlayerShip.Crew {
			sg.PlayerShip.Crew[i].Update()
		}

		//need starfield shift speed controlled here (currently hardcoded to shift every 100 seconds)
		if sg.SpaceTime%100 == 0 {
			sg.shiftStarField()
		}
	}

	if sg.activeMenu == sg.crewMenu && sg.crewDetails.IsVisible() {
		sg.UpdateCrewDetails()
	}

	sg.missiontime.ChangeText(fmt.Sprintf("%.4d", sg.SpaceTime/100000) + "d:" + fmt.Sprintf("%.1d", (sg.SpaceTime/10000)%10) + "h:" + fmt.Sprintf("%.2d", (sg.SpaceTime/100)%100) + "m:" + fmt.Sprintf("%.2d", sg.SpaceTime%100) + "s")

}

func (sg *SpaceshipGame) Render() {
	sg.DrawStarfield()

	w, h := sg.PlayerShip.ShipMap.Dims()
	x, y := 0, 0
	displayWidth, displayHeight := sg.shipdisplay.Dims()

	for i := 0; i < w*h; i++ {
		//shipdisplay-space coords
		x = i%w + sg.viewX
		y = i/w + sg.viewY
		

		if util.CheckBounds(x, y, displayWidth, displayHeight) {
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
	if sg.activeMenu != nil {
		sg.activeMenu.Render()
	}
}

//Activates a menu (crew, rooms, systems, etc). Does nothing if menu already active.
func (sg *SpaceshipGame) ActivateMenu(m ui.UIElem) {
	if sg.activeMenu == m {
		return
	}

	m.SetVisibility(true)
	if sg.activeMenu != nil {
		sg.activeMenu.SetVisibility(false)
	}
	sg.activeMenu = m

	if m != sg.input {
		sg.CenterShip()
	}
}

//deactivates the open menu (if there is one)
func (sg *SpaceshipGame) DeactivateMenu() {
	if sg.activeMenu == nil {
		return
	}
	sg.activeMenu.SetVisibility(false)
	sg.activeMenu = nil
	sg.CenterShip()
}

func (sg SpaceshipGame) GetIncrement() int {
	if sg.paused {
		return 0
	}

	switch sg.SimSpeed {
	case 1:
		return 1
	case 2:
		return 10
	case 3:
		return 100
	case 4:
		return 1000
	default:
		return 0
	}
}