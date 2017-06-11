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
	window       *ui.Container
	input        *ui.Inputbox
	output       *ui.List
	shipstatus   *ShipStatsWindow
	missiontime  *ui.Textbox
	speeddisplay *ui.TileView
	shipdisplay  *ui.TileView

	//top menu. contains buttons for submenus
	menubar             *ui.Container
	crewMenuButton      *ui.Button
	shipMenuButton      *ui.Button
	roomMenuButton      *ui.Button
	mainMenuButton      *ui.Button
	starchartMenuButton *ui.Button
	scippieMenuButton   *ui.Button

	//submenus. these are stored always for fast switching.
	//crew menu (F1)
	crewMenu    *ui.Container
	crewList    *ui.List
	crewDetails *ui.Container

	shipMenu      *ShipMenu      //shipmenu (F3)
	starchartMenu *StarchartMenu //starchart (F4)

	activeMenu ui.UIElem

	//Time Globals.
	spaceTime int //measured in Standard Galactic Seconds
	simSpeed  int //4 speeds, plus pause (0)
	paused    bool

	Stars StarField

	viewX, viewY int

	galaxy     *Galaxy
	playerShip *Ship
}

func NewSpaceshipGame() *SpaceshipGame {
	sg := new(SpaceshipGame)

	sg.spaceTime = 0
	sg.simSpeed = 1

	sg.galaxy = NewGalaxy()

	sg.playerShip = NewShip("The Undestructable")
	sg.playerShip.AddRoom(NewRoom("Engineering", 5, 8, 5, 8, 700, 1000))
	sg.playerShip.AddRoom(NewRoom("Messhall", 15, 5, 6, 6, 1000, 500))
	sg.playerShip.AddRoom(NewRoom("Medbay", 9, 5, 6, 6, 1000, 700))
	sg.playerShip.AddRoom(NewRoom("Quarters 1", 15, 13, 6, 6, 900, 500))
	sg.playerShip.AddRoom(NewRoom("Quarters 2", 9, 13, 6, 6, 900, 500))
	sg.playerShip.AddRoom(NewRoom("Hallway", 9, 10, 12, 4, 0, 500))

	ss := sg.galaxy.GetSector(8, 8).GenerateSubSector(250, 171)
	ss.star = NewStarSystem(ss.GetCoords())
	sg.playerShip.SetLocation(ss.star.Planets[2]) //Earth!!

	sg.SetupUI() //must be done after ship setup
	sg.UpdateSpeedUI()

	sg.CenterShip()

	return sg
}

//Centers the map of the ship in the main view.
func (sg *SpaceshipGame) CenterShip() {
	displayWidth, displayHeight := sg.shipdisplay.Dims()
	sg.viewX = displayWidth/2 - sg.playerShip.Width/2 - sg.playerShip.X
	sg.viewY = displayHeight/2 - sg.playerShip.Height/2 - sg.playerShip.Y
	if sg.activeMenu != nil {
		w, _ := sg.activeMenu.Dims()
		sg.viewX -= w / 2
	}
}

func (sg *SpaceshipGame) UpdateSpeedUI() {
	sg.speeddisplay.Clear()
	for i := 0; i < 4; i++ {
		if i < sg.simSpeed {
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
	w, h := sg.shipdisplay.Dims()
	sg.Stars = NewStarField(w, h, 20, sg.shipdisplay)

	sg.shipstatus = NewShipStatsWindow(sg.playerShip)

	sg.input = ui.NewInputbox(50, 1, 15, 27, 2, true)
	sg.input.ToggleFocus()
	sg.input.SetVisibility(false)
	sg.input.SetTitle("SCIPPIE V6.18")

	sg.output = ui.NewList(51, 12, 28, 32, 1, true, "The Ship Computer Interactive Parameter Parser/Interface Entity, or SCIPPIE, is your computerized second in command. Ask questions, give commands and observe your ship through the high-tech text-tacular wonders of 38th century UI technology! Ask SCIPPIE a question, or give him a command!")
	sg.output.ToggleHighlight()

	sg.SetupCrewMenu()
	sg.starchartMenu = NewStarchartMenu(sg.galaxy, sg.playerShip)
	sg.starchartMenu.LoadLocalInfo()
	sg.shipMenu = NewShipMenu()

	sg.window.Add(sg.input, sg.output, sg.shipstatus, sg.shipdisplay, sg.speeddisplay, sg.missiontime, sg.menubar, sg.shipMenu, sg.starchartMenu)
}

func (sg *SpaceshipGame) Update() {

	startCoords := sg.playerShip.coords

	//simulation!
	for i := 0; i < sg.GetIncrement(); i++ {
		sg.spaceTime++
		sg.playerShip.Update(sg.spaceTime)

		//change location if we move away,
		if !sg.playerShip.coords.IsIn(sg.playerShip.CurrentLocation) {
			c := sg.playerShip.coords
			sg.playerShip.CurrentLocation = sg.galaxy.GetSector(c.sector.Get()).GetSubSector(c.subSector.Get()).star
		}

		//change destination/location when we arrive!
		if sg.playerShip.coords.IsIn(sg.playerShip.Destination) {
			sg.playerShip.CurrentLocation = sg.playerShip.Destination
			sg.playerShip.Destination = nil
			sg.playerShip.Speed = 0
			sg.playerShip.Engine.Firing = false
		}

		for i := range sg.playerShip.Crew {
			sg.playerShip.Crew[i].Update()
		}

		//need starfield shift speed controlled here (currently hardcoded to shift every 100 seconds as long as the ship is moving)
		if sg.playerShip.Speed != 0 && sg.spaceTime%100 == 0 {
			sg.Stars.Shift()
		}
	}

	//update starchart if ship has moved sectors
	if sg.activeMenu == sg.starchartMenu {
		sg.starchartMenu.Update()
		delta := startCoords.CalcVector(sg.playerShip.coords)
		if sec := delta.Sector(); sg.starchartMenu.mapMode == coord_SECTOR && (sec.X != 0 || sec.Y != 0) {
			sg.starchartMenu.DrawMap()
		} else if sg.starchartMenu.mapMode == coord_LOCAL {
			sg.starchartMenu.DrawSystem() //should we really do this every update tick??? ugh.
		}
	}

	if sg.activeMenu == sg.crewMenu && sg.crewDetails.IsVisible() {
		sg.UpdateCrewDetails()
	}
	sg.shipstatus.Update()
	sg.missiontime.ChangeText(fmt.Sprintf("%.4d", sg.spaceTime/100000) + "d:" + fmt.Sprintf("%.1d", (sg.spaceTime/10000)%10) + "h:" + fmt.Sprintf("%.2d", (sg.spaceTime/100)%100) + "m:" + fmt.Sprintf("%.2d", sg.spaceTime%100) + "s")
}

func (sg *SpaceshipGame) Render() {
	sg.Stars.Draw()

	w, h := sg.playerShip.ShipMap.Dims()
	x, y := 0, 0
	displayWidth, displayHeight := sg.shipdisplay.Dims()

	for i := 0; i < w*h; i++ {
		//shipdisplay-space coords
		x = i%w + sg.viewX
		y = i/w + sg.viewY

		if util.CheckBounds(x, y, displayWidth, displayHeight) {
			t := sg.playerShip.ShipMap.GetTile(i%w, i/w)
			if t.TileType() != 0 {
				tv := t.GetVisuals()
				sg.shipdisplay.Draw(x, y, tv.Glyph, tv.ForeColour, 0xFF000000)
			}

			if e := sg.playerShip.ShipMap.GetEntity(i%w, i/w); e != nil {
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

func (sg *SpaceshipGame) MoveShipCamera(dx, dy int) {
	sg.Stars.dirty = true //clears the tileview

	sg.viewX += dx
	sg.viewY += dy
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

	switch sg.simSpeed {
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
