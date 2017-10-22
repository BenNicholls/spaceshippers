package main

import "github.com/bennicholls/burl-E/burl"
import "fmt"

//time values in DIGITAL SECONDS. One digital day = 100000 seconds, which is 14% longer than a regular day.
const (
	MINUTE int = 100
	HOUR   int = 10000
	DAY    int = 100000
)

//load some tile data
var TILE_FLOOR = burl.LoadTileData("Floor", true, true, burl.GLYPH_FILL_SPARSE, burl.COL_DARKGREY)
var TILE_WALL = burl.LoadTileData("Wall", false, false, burl.GLYPH_HASH, burl.COL_GREY)
var TILE_DOOR = burl.LoadTileData("Door", true, false, burl.GLYPH_IDENTICAL, burl.COL_GREY)

type SpaceshipGame struct {

	//ui stuff
	window       *burl.Container
	input        *burl.Inputbox
	output       *burl.List
	shipstatus   *ShipStatsWindow
	missiontime  *burl.Textbox
	speeddisplay *burl.TileView
	shipdisplay  *burl.TileView

	//top menu. contains buttons for submenus
	menubar             *burl.Container
	crewMenuButton      *burl.Button
	shipMenuButton      *burl.Button
	missionMenuButton   *burl.Button
	mainMenuButton      *burl.Button
	starchartMenuButton *burl.Button
	scippieMenuButton   *burl.Button

	//submenus. these are stored always for fast switching.
	//crew menu (F1)
	crewMenu    *burl.Container
	crewList    *burl.List
	crewDetails *burl.Container

	shipMenu      *ShipMenu      //shipmenu (F2)
	missionMenu   *MissionMenu   //missionmenu (F3)
	starchartMenu *StarchartMenu //starchart (F4)

	activeMenu burl.UIElem
	dialog     Dialog //dialog presented to the player. higher priority than everything else!

	//Time Globals.
	spaceTime int //measured in Standard Galactic Seconds
	simSpeed  int //4 speeds, plus pause (0)
	paused    bool

	Stars StarField

	viewX, viewY int

	galaxy     *Galaxy
	playerShip *Ship

	missionLog []Mission
}

func NewSpaceshipGame() *SpaceshipGame {
	sg := new(SpaceshipGame)

	sg.spaceTime = 0
	sg.simSpeed = 1

	sg.galaxy = NewGalaxy()

	sg.playerShip = NewShip("The Undestructable", sg.galaxy)
	sg.playerShip.description = "This is your ship! Look at it's heroic hull valiantly floating amongst the stars. One could almost weep."
	sg.playerShip.AddRoom(NewRoom("Engineering", 5, 8, 5, 8, 700, 1000))
	sg.playerShip.AddRoom(NewRoom("Messhall", 15, 5, 6, 6, 1000, 500))
	sg.playerShip.AddRoom(NewRoom("Medbay", 9, 5, 6, 6, 1000, 700))
	sg.playerShip.AddRoom(NewRoom("Quarters 1", 15, 13, 6, 6, 900, 500))
	sg.playerShip.AddRoom(NewRoom("Quarters 2", 9, 13, 6, 6, 900, 500))
	sg.playerShip.AddRoom(NewRoom("Hallway", 9, 10, 12, 4, 0, 500))

	ss := sg.galaxy.GetSector(8, 8).GenerateSubSector(250, 171)
	ss.starSystem = NewStarSystem(ss.GetCoords())
	sg.playerShip.SetLocation(ss.starSystem.Planets[2]) //Earth!!

	sg.missionLog = make([]Mission, 0)

	sg.SetupUI() //must be done after ship setup
	sg.UpdateSpeedUI()

	sg.CenterShip()
	sg.AddMission(GenerateGoToMission(sg.playerShip, ss.starSystem.Planets[4], ss.starSystem.Star))

	return sg
}

func (sg *SpaceshipGame) AddMission(m *Mission) {
	sg.missionLog = append(sg.missionLog, *m)
	sg.missionMenu.UpdateMissionList()
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

	sg.ResetShipView()
}

func (sg *SpaceshipGame) UpdateSpeedUI() {
	sg.speeddisplay.Reset()
	for i := 0; i < 4; i++ {
		if i < sg.simSpeed {
			sg.speeddisplay.Draw(i, 0, burl.GLYPH_TRIANGLE_RIGHT, burl.COL_WHITE, burl.COL_BLACK)
		} else {
			sg.speeddisplay.Draw(i, 0, burl.GLYPH_UNDERSCORE, burl.COL_WHITE, burl.COL_BLACK)
		}
	}
}

func (sg *SpaceshipGame) SetupUI() {
	sg.window = burl.NewContainer(80, 45, 0, 0, 0, false)

	sg.missiontime = burl.NewTextbox(8, 1, 1, 1, 2, true, true, "")
	sg.speeddisplay = burl.NewTileView(8, 1, 1, 3, 3, true)

	sg.menubar = burl.NewContainer(69, 1, 12, 1, 1, false)
	sg.crewMenuButton = burl.NewButton(9, 1, 1, 0, 1, true, true, "Crew")
	sg.shipMenuButton = burl.NewButton(9, 1, 12, 0, 1, true, true, "Ship")
	sg.missionMenuButton = burl.NewButton(9, 1, 23, 0, 1, true, true, "Missions")
	sg.starchartMenuButton = burl.NewButton(9, 1, 34, 0, 1, true, true, "Star Chart")
	sg.scippieMenuButton = burl.NewButton(9, 1, 45, 0, 1, true, true, "S.C.I.P.P.I.E.")
	sg.mainMenuButton = burl.NewButton(9, 1, 56, 0, 1, true, true, "Main  Menu")
	sg.menubar.Add(sg.crewMenuButton, sg.shipMenuButton, sg.missionMenuButton, sg.starchartMenuButton, sg.scippieMenuButton, sg.mainMenuButton)

	sg.shipdisplay = burl.NewTileView(80, 28, 0, 3, 1, false)
	w, h := sg.shipdisplay.Dims()
	sg.Stars = NewStarField(w, h, 20, sg.shipdisplay)

	sg.shipstatus = NewShipStatsWindow(sg.playerShip)

	sg.input = burl.NewInputbox(50, 1, 15, 27, 2, true)
	sg.input.ToggleFocus()
	sg.input.SetVisibility(false)
	sg.input.SetTitle("SCIPPIE V6.18")

	sg.output = burl.NewList(51, 12, 28, 32, 1, true, "The Ship Computer Interactive Parameter Parser/Interface Entity, or SCIPPIE, is your computerized second in command. Ask questions, give commands and observe your ship through the high-tech text-tacular wonders of 38th century UI technology! Ask SCIPPIE a question, or give him a command!")
	sg.output.ToggleHighlight()

	sg.SetupCrewMenu()
	sg.starchartMenu = NewStarchartMenu(sg.galaxy, sg.playerShip)
	sg.shipMenu = NewShipMenu()
	sg.missionMenu = NewMissionMenu(&sg.missionLog)

	sg.window.Add(sg.input, sg.output, sg.shipstatus, sg.shipdisplay, sg.speeddisplay, sg.missiontime, sg.menubar, sg.shipMenu, sg.starchartMenu)
}

func (sg *SpaceshipGame) Update() {

	//check if we should be handling a dialog
	if sg.dialog != nil {
		if sg.dialog.Done() {
			sg.dialog.ToggleVisible()
			sg.dialog = nil
		} else {
			return
		}
	}

	startCoords := sg.playerShip.coords

	//simulation!
	for i := 0; i < sg.GetIncrement(); i++ {
		sg.spaceTime++
		sg.playerShip.Update(sg.spaceTime)

		for i := range sg.playerShip.Crew {
			sg.playerShip.Crew[i].Update()
		}

		//need starfield shift speed controlled here (currently hardcoded to shift every 100 seconds as long as the ship is moving)
		if sg.playerShip.GetSpeed() != 0 && sg.spaceTime%100 == 0 {
			sg.Stars.Shift()
		}

		for i := range sg.missionLog {
			sg.missionLog[i].Update()
		}
	}

	//update starchart if ship has moved
	if sg.activeMenu == sg.starchartMenu && sg.playerShip.GetSpeed() != 0 {
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

	if sg.activeMenu == sg.missionMenu {
		sg.missionMenu.Update()
	}

	sg.shipstatus.Update()
	sg.missiontime.ChangeText(GetTimeString(sg.spaceTime))
}

func GetTimeString(t int) string {
	return fmt.Sprintf("%.4d", t/100000) + "d:" + fmt.Sprintf("%.1d", (t/10000)%10) + "h:" + fmt.Sprintf("%.2d", (t/100)%100) + "m:" + fmt.Sprintf("%.2d", t%100) + "s"
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

		if burl.CheckBounds(x, y, displayWidth, displayHeight) {
			t := sg.playerShip.ShipMap.GetTile(i%w, i/w)
			if t.TileType() != 0 {
				tv := t.GetVisuals()
				sg.shipdisplay.Draw(x, y, tv.Glyph, tv.ForeColour, burl.COL_BLACK)
			}

			if e := sg.playerShip.ShipMap.GetEntity(i%w, i/w); e != nil {
				sg.shipdisplay.Draw(x, y, e.GetVisuals().Glyph, e.GetVisuals().ForeColour, burl.COL_BLACK)
			}
		}
	}

	sg.window.Render()
	if sg.activeMenu != nil {
		sg.activeMenu.Render()
	}

	if sg.dialog != nil {
		sg.dialog.Render()
	}
}

//Activates a menu (crew, rooms, systems, etc). Deactivates menu if menu already active.
func (sg *SpaceshipGame) ActivateMenu(m burl.UIElem) {
	if sg.activeMenu == m {
		sg.DeactivateMenu()
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

func (sg *SpaceshipGame) MoveShipCamera(dx, dy int) {
	sg.ResetShipView()

	sg.viewX += dx
	sg.viewY += dy
}

func (sg *SpaceshipGame) ResetShipView() {
	sg.shipdisplay.Reset()
	sg.Stars.dirty = true
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

func (sg SpaceshipGame) GetTick() int {
	return sg.spaceTime
}
