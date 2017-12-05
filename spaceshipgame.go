package main

import "github.com/bennicholls/burl-E/burl"
import "os"
import "encoding/gob"

//load some tile data
var TILE_FLOOR = burl.LoadTileData("Floor", true, true, burl.GLYPH_FILL_SPARSE, burl.COL_DARKGREY)
var TILE_WALL = burl.LoadTileData("Wall", false, false, burl.GLYPH_HASH, burl.COL_GREY)
var TILE_DOOR = burl.LoadTileData("Door", true, false, burl.GLYPH_IDENTICAL, burl.COL_GREY)

func init() {
	//need to register types that might be hidden by an interface, in order for them to be serializable
	gob.Register(&SleepJob{})
}

type SpaceshipGame struct {
	//ui stuff
	window       *burl.Container
	input        *burl.Inputbox
	output       *burl.List
	shipstatus   *ShipStatsWindow
	timeDisplay  *TimeDisplay
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

	crewMenu      *CrewMenu      //crew menu (F1)
	shipMenu      *ShipMenu      //shipmenu (F2)
	missionMenu   *MissionMenu   //missionmenu (F3)
	starchartMenu *StarchartMenu //starchart (F4)

	activeMenu burl.UIElem
	dialog     Dialog //dialog presented to the player. higher priority than everything else!

	//Time Globals.
	startTime int //time since launch, measured in Standard Galactic Seconds
	simSpeed  int //4 speeds, plus pause (0)
	paused    bool

	Stars StarField

	viewX, viewY int

	galaxy     *Galaxy
	player     *Player
	playerShip *Ship //THINK: do we need this? it's just a pointer to player.Spaceship
}

func NewSpaceshipGame() *SpaceshipGame {
	sg := new(SpaceshipGame)

	sg.simSpeed = 1

	sg.galaxy = NewGalaxy()
	sg.startTime = sg.galaxy.spaceTime

	sg.player = NewPlayer("Ol Cappy")
	sg.player.SpaceShip = NewShip("The Undestructable", sg.galaxy)
	sg.playerShip = sg.player.SpaceShip

	ss := sg.galaxy.GetSector(8, 8).GenerateSubSector(250, 171)
	ss.starSystem = NewStarSystem(ss.GetCoords())
	sg.playerShip.SetLocation(ss.starSystem.Planets[2]) //Earth!!

	sg.SetupUI() //must be done after ship setup
	sg.CenterShip()

	sg.AddMission(GenerateGoToMission(sg.playerShip, ss.starSystem.Planets[4], ss.starSystem.Star))
	sg.AddMission(GenerateGoToMission(sg.playerShip, ss.starSystem.Planets[5], ss.starSystem.Planets[2]))

	welcomeMessage := "Hi Captain! Welcome to " + sg.playerShip.GetName() + "! I am the Ship Computer Interactive Parameter-Parsing Intelligence Entity, but you can call me SCIPPIE! "
	sg.dialog = NewCommDialog("SCIPPIE", sg.player.Name + ", Captain of "+sg.playerShip.GetName(), "res/art/scippie.csv", welcomeMessage)

	return sg
}

//Adds a mission to the player's list.
//THINK ABOUT: this could be a method for the player object??? Then how does UI get updated? mm.
func (sg *SpaceshipGame) AddMission(m *Mission) {
	sg.player.MissionLog = append(sg.player.MissionLog, *m)
	sg.missionMenu.UpdateMissionList()
}

//Centers the map of the ship in the main view.
func (sg *SpaceshipGame) CenterShip() {
	displayWidth, displayHeight := sg.shipdisplay.Dims()
	sg.viewX = displayWidth/2 - sg.playerShip.width/2 - sg.playerShip.x
	sg.viewY = displayHeight/2 - sg.playerShip.height/2 - sg.playerShip.y
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

	sg.timeDisplay = NewTimeDisplay(sg.galaxy)
	sg.speeddisplay = burl.NewTileView(4, 1, 1, 4, 10, true)

	sg.menubar = burl.NewContainer(69, 1, 12, 1, 10, false)
	sg.crewMenuButton = burl.NewButton(9, 1, 1, 0, 1, true, true, "Crew")
	sg.shipMenuButton = burl.NewButton(9, 1, 12, 0, 2, true, true, "Ship")
	sg.missionMenuButton = burl.NewButton(9, 1, 23, 0, 1, true, true, "Missions")
	sg.starchartMenuButton = burl.NewButton(9, 1, 34, 0, 2, true, true, "Star Chart")
	sg.scippieMenuButton = burl.NewButton(9, 1, 45, 0, 1, true, true, "S.C.I.P.P.I.E.")
	sg.mainMenuButton = burl.NewButton(9, 1, 56, 0, 2, true, true, "Main  Menu")
	sg.menubar.Add(sg.crewMenuButton, sg.shipMenuButton, sg.missionMenuButton, sg.starchartMenuButton, sg.scippieMenuButton, sg.mainMenuButton)

	sg.shipdisplay = burl.NewTileView(80, 28, 0, 3, 1, false)
	w, h := sg.shipdisplay.Dims()
	sg.Stars = NewStarField(w, h, 20, sg.shipdisplay)

	sg.shipstatus = NewShipStatsWindow(sg.playerShip)

	sg.input = burl.NewInputbox(50, 1, 15, 27, 100, true)
	sg.input.ToggleFocus()
	sg.input.SetVisibility(false)
	sg.input.SetTitle("SCIPPIE V6.18")

	sg.output = burl.NewList(51, 12, 28, 32, 10, true, "Nothing to report, Captain!")
	sg.output.ToggleHighlight()

	sg.crewMenu = NewCrewMenu(sg.playerShip)
	sg.starchartMenu = NewStarchartMenu(sg.galaxy, sg.playerShip)
	sg.shipMenu = NewShipMenu()
	sg.missionMenu = NewMissionMenu(&sg.player.MissionLog)

	sg.window.Add(sg.input, sg.output, sg.shipstatus, sg.shipdisplay, sg.speeddisplay, sg.timeDisplay, sg.menubar, sg.shipMenu, sg.starchartMenu)

	sg.UpdateSpeedUI()
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

	startCoords := sg.playerShip.Coords

	//simulation!
	for i := 0; i < sg.GetIncrement(); i++ {
		sg.galaxy.spaceTime++
		sg.playerShip.Update(sg.GetTime())

		for i := range sg.playerShip.Crew {
			sg.playerShip.Crew[i].Update()
		}

		//need starfield shift speed controlled here (currently hardcoded to shift every 100 seconds as long as the ship is moving)
		if sg.playerShip.GetSpeed() != 0 && sg.GetTick()%100 == 0 {
			sg.Stars.Shift()
		}

		for i := range sg.player.MissionLog {
			sg.player.MissionLog[i].Update()
		}
	}

	//update starchart if ship has moved
	if sg.activeMenu == sg.starchartMenu && sg.playerShip.GetSpeed() != 0 {
		sg.starchartMenu.Update()
		delta := startCoords.CalcVector(sg.playerShip.Coords)
		if sec := delta.Sector; sg.starchartMenu.mapMode == coord_SECTOR && (sec.X != 0 || sec.Y != 0) {
			sg.starchartMenu.DrawMap()
		} else if sg.starchartMenu.mapMode == coord_LOCAL {
			sg.starchartMenu.DrawSystem() //should we really do this every update tick??? ugh.
		}
	}

	if sg.activeMenu == sg.crewMenu && sg.crewMenu.crewDetails.IsVisible() {
		sg.crewMenu.UpdateCrewDetails()
	}

	if sg.activeMenu == sg.missionMenu {
		sg.missionMenu.Update()
	}

	sg.shipstatus.Update()
	sg.timeDisplay.Update()
}

func (sg *SpaceshipGame) Render() {
	sg.Stars.Draw()

	w, h := sg.playerShip.shipMap.Dims()
	x, y := 0, 0
	displayWidth, displayHeight := sg.shipdisplay.Dims()

	for i := 0; i < w*h; i++ {
		//shipdisplay-space coords
		x = i%w + sg.viewX
		y = i/w + sg.viewY

		if burl.CheckBounds(x, y, displayWidth, displayHeight) {
			t := sg.playerShip.shipMap.GetTile(i%w, i/w)
			if t.TileType != 0 {
				tv := t.GetVisuals()
				sg.shipdisplay.Draw(x, y, tv.Glyph, tv.ForeColour, burl.COL_BLACK)
			}

			if e := sg.playerShip.shipMap.GetEntity(i%w, i/w); e != nil {
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

//returns the number of simulated seconds since launch
func (sg SpaceshipGame) GetTick() int {
	return sg.galaxy.spaceTime - sg.startTime
}

//gets the time from the Galaxy
func (sg SpaceshipGame) GetTime() int {
	return sg.galaxy.spaceTime
}

func (sg *SpaceshipGame) SaveShip() {
	f, err := os.Create("savefile")
	if err != nil {
		burl.LogError("Could not open file for saving: " + err.Error())
	}
	defer f.Close()

	enc := gob.NewEncoder(f)
	err = enc.Encode(sg.playerShip)
	if err != nil {
		burl.LogError("Could not save ship: " + err.Error())
	}
}

func (sg *SpaceshipGame) LoadShip() {
	f, err := os.Open("savefile")
	if err != nil {
		burl.LogError("Could not open file for loading: " + err.Error())
	}
	defer f.Close()

	s := new(Ship)

	dec := gob.NewDecoder(f)
	err = dec.Decode(s)
	if err != nil {
		burl.LogError("Could not load ship: " + err.Error())
	}

	//data loaded, now to re-init everything
	s.SetupShip(sg.galaxy)

	//load complete, make the switch!!
	sg.player.SpaceShip = s
	sg.playerShip = s

	sg.SetupUI()
}
