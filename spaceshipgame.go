package main

import (
	"encoding/gob"
	"os"

	"github.com/bennicholls/burl-E/burl"
)

const (
	MENU_GAME int = iota
	MENU_SHIP
	MENU_GALAXY
	MENU_CREW
	MENU_COMM
	MENU_VIEW
	MENU_MAIN
	MAX_MENUS
)

//Event types for spaceshippers!
var LOG_EVENT = burl.RegisterCustomEvent()

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
	window      *burl.Container
	output      *burl.List
	shipstatus  *ShipStatsWindow
	timeDisplay *TimeDisplay
	shipdisplay *burl.TileView

	//top menu. contains buttons for submenus
	gameMenuButton   *burl.Button
	shipMenuButton   *burl.Button
	galaxyMenuButton *burl.Button
	crewMenuButton   *burl.Button
	commMenuButton   *burl.Button
	viewMenuButton   *burl.Button
	mainMenuButton   *burl.Button

	gameMenu   *GameMenu   //(F1)
	shipMenu   *ShipMenu   //(F2)
	galaxyMenu *GalaxyMenu //(F3)
	crewMenu   *CrewMenu   //(F4)
	commMenu   *CommMenu   //(F5)
	viewMenu   *ViewMenu   //(F6)
	mainMenu   *MainMenu   //(ESC)

	menus       []burl.UIElem
	menuButtons []*burl.Button

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

func NewSpaceshipGame(g *Galaxy, s *Ship) *SpaceshipGame {
	sg := new(SpaceshipGame)

	sg.simSpeed = 1

	sg.galaxy = g
	sg.startTime = sg.galaxy.spaceTime

	sg.player = NewPlayer("Ol Cappy")
	sg.player.SpaceShip = s
	sg.playerShip = sg.player.SpaceShip

	sg.playerShip.SetLocation(sg.galaxy.GenerateStart())
	//sg.playerShip.SetLocation(sg.galaxy.GetEarth())

	sg.SetupUI() //must be done after ship setup

	sg.LoadSpaceEvents()

	sg.dialog = NewSpaceEventDialog(spaceEvents[1])

	return sg
}

//Adds a mission to the player's list.
//THINK ABOUT: this could be a method for the player object???
func (sg *SpaceshipGame) AddMission(m *Mission) {
	sg.player.MissionLog = append(sg.player.MissionLog, *m)
	burl.PushEvent(burl.NewEvent(burl.EV_UPDATE_UI, "missions"))
}

//Centers the map of the ship in the main view.
func (sg *SpaceshipGame) CenterShip() {
	displayWidth, displayHeight := sg.shipdisplay.Dims()
	sg.viewX = sg.playerShip.x + sg.playerShip.width/2 - displayWidth/2
	sg.viewY = sg.playerShip.y + sg.playerShip.height/2 - displayHeight/2

	if sg.activeMenu != nil {
		w, _ := sg.activeMenu.Dims()
		sg.viewX += w / 2
	}

	sg.ResetShipView()
}

func (sg *SpaceshipGame) SetupUI() {
	sg.window = burl.NewContainer(80, 45, 0, 0, 0, false)

	sg.shipdisplay = burl.NewTileView(80, 28, 0, 3, 1, false)
	sg.Stars = NewStarField(20, sg.shipdisplay)

	sg.shipstatus = NewShipStatsWindow(sg.playerShip)
	sg.shipstatus.Update()
	sg.timeDisplay = NewTimeDisplay(sg.galaxy)
	sg.timeDisplay.UpdateSpeed(sg.simSpeed)

	sg.output = burl.NewList(51, 12, 28, 32, 10, true, "Nothing to report, Captain!")
	sg.output.ToggleHighlight()
	sg.output.SetHint("PgUp/PgDown to scroll")

	sg.window.Add(sg.output, sg.shipstatus, sg.shipdisplay, sg.timeDisplay)

	sg.menus = make([]burl.UIElem, 0, MAX_MENUS)

	sg.gameMenu = NewGameMenu()
	sg.shipMenu = NewShipMenu()
	sg.galaxyMenu = NewGalaxyMenu()
	sg.crewMenu = NewCrewMenu()
	sg.commMenu = NewCommsMenu(sg.playerShip.Comms)
	sg.viewMenu = NewViewMenu()
	sg.mainMenu = NewMainMenu()

	sg.menus = append(sg.menus, sg.gameMenu, sg.shipMenu, sg.galaxyMenu, sg.crewMenu, sg.commMenu, sg.viewMenu, sg.mainMenu)
	sg.window.Add(sg.menus...)

	sg.menuButtons = make([]*burl.Button, 0, MAX_MENUS)

	sg.gameMenuButton = burl.NewButton(8, 1, 3, 1, 10, true, true, "Game")
	sg.gameMenuButton.SetHint("F1")
	sg.shipMenuButton = burl.NewButton(8, 1, 14, 1, 10, true, true, "Ship")
	sg.shipMenuButton.SetHint("F2")
	sg.galaxyMenuButton = burl.NewButton(8, 1, 25, 1, 10, true, true, "Galaxy")
	sg.galaxyMenuButton.SetHint("F3")
	sg.crewMenuButton = burl.NewButton(8, 1, 36, 1, 10, true, true, "Crew")
	sg.crewMenuButton.SetHint("F4")
	sg.commMenuButton = burl.NewButton(8, 1, 47, 1, 10, true, true, "Communications")
	sg.commMenuButton.SetHint("F5")
	sg.viewMenuButton = burl.NewButton(8, 1, 58, 1, 10, true, true, "View  Mode")
	sg.viewMenuButton.SetHint("F6")
	sg.mainMenuButton = burl.NewButton(8, 1, 69, 1, 10, true, true, "Main  Menu")
	sg.mainMenuButton.SetHint("ESC")

	sg.menuButtons = append(sg.menuButtons, sg.gameMenuButton, sg.shipMenuButton, sg.galaxyMenuButton, sg.crewMenuButton, sg.commMenuButton, sg.viewMenuButton, sg.mainMenuButton)

	for i := range sg.menuButtons {
		sg.window.Add(sg.menuButtons[i])
	}

	sg.CenterShip()
}

func (sg *SpaceshipGame) Update() {
	//check if we should be handling a dialog
	if sg.dialog != nil {
		if sg.dialog.Done() {
			sg.dialog = nil
		} else {
			sg.dialog.Update()
			return
		}
	}

	//startCoords := sg.playerShip.Coords

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
	// if sg.activeMenu == sg.starchartMenu && sg.playerShip.GetSpeed() != 0 {
	// 	sg.starchartMenu.Update()
	// 	delta := startCoords.CalcVector(sg.playerShip.Coords)
	// 	if sec := delta.Sector; sg.starchartMenu.mapMode == coord_SECTOR && (sec.X != 0 || sec.Y != 0) {
	// 		sg.starchartMenu.DrawMap()
	// 	} else if sg.starchartMenu.mapMode == coord_LOCAL {
	// 		sg.starchartMenu.DrawSystem() //should we really do this every update tick??? ugh.
	// 	}
	// }

	sg.timeDisplay.UpdateTime()
}

func (sg *SpaceshipGame) HandleEvent(event *burl.Event) {
	if sg.dialog != nil {
		sg.dialog.HandleEvent(event)
	}

	switch event.ID {
	case burl.EV_UPDATE_UI:
		switch event.Message {
		case "inbox":
			sg.commMenu.UpdateInbox()
		case "transmissions":
			sg.commMenu.UpdateTransmissions()
		// case "missions":
		// 	sg.missionMenu.Update()
		// case "crew":
		// 	if sg.activeMenu == sg.crewMenu {
		// 		sg.crewMenu.UpdateCrewDetails()
		// 	}
		case "ship status":
			sg.shipstatus.Update()
		}

	case LOG_EVENT:
		sg.AddMessage(event.Message)
	}
}

func (sg *SpaceshipGame) Render() {
	sg.Stars.Draw()
	sg.playerShip.DrawToTileView(sg.shipdisplay, sg.viewX, sg.viewY)

	sg.window.Render()
	if sg.activeMenu != nil {
		sg.activeMenu.Render()
	}

	if sg.dialog != nil {
		sg.dialog.Render()
	}
}

//Activates a menu (crew, rooms, systems, etc). Deactivates menu if menu already active.
func (sg *SpaceshipGame) ActivateMenu(menu int) {
	sg.menuButtons[menu].Press()
	m := sg.menus[menu]

	if sg.activeMenu == m {
		sg.DeactivateMenu()
		return
	}

	m.SetVisibility(true)
	if sg.activeMenu != nil {
		sg.activeMenu.SetVisibility(false)
	}
	sg.activeMenu = m
}

//deactivates the open menu (if there is one)
func (sg *SpaceshipGame) DeactivateMenu() {
	if sg.activeMenu == nil {
		return
	}
	sg.activeMenu.SetVisibility(false)
	sg.activeMenu = nil

	for i := range sg.menuButtons {
		sg.menuButtons[i].SetVisibility(true)
	}
	sg.CenterShip()
}

func (sg *SpaceshipGame) MoveShipCamera(dx, dy int) {
	sg.ResetShipView()

	sg.viewX -= dx
	sg.viewY -= dy
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

func (sg SpaceshipGame) Shutdown() {
	sg.SaveShip()
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
