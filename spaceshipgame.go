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

func init() {
	//need to register types that might be hidden by an interface, in order for them to be serializable
	gob.Register(&SleepJob{})
}

type SpaceshipGame struct {
	burl.StatePrototype

	//ui stuff
	output      *burl.List
	quickstats  *QuickStatsWindow
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

	//Time Globals.
	startTime int //time since launch, measured in Standard Galactic Seconds
	simSpeed  int //4 speeds, plus pause (0)
	paused    bool

	Stars StarField

	viewX, viewY int
	viewMode     int

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

	//fuel up the ship!!! Take up 50% of liquid storage.
	sg.playerShip.Storage.Store(&Item{
		Name:        "Fuel",
		Volume:      s.Storage.GetCapacity(STORE_LIQUID)/2 - sg.playerShip.Storage.GetItemVolume("Fuel"),
		StorageType: STORE_LIQUID,
	})

	sg.playerShip.Storage.Store(&Item{
		Name:        "Candy",
		Volume:      50,
		StorageType: STORE_GENERAL,
	})

	sg.playerShip.Storage.Store(&Gas{
		gasType: GAS_O2,
		molar:   s.Storage.GetCapacity(STORE_GAS),
	})

	sg.SetupUI() //must be done after ship setup

	sg.LoadSpaceEvents()

	sg.OpenDialog(NewSpaceEventDialog(spaceEvents[1]))

	burl.RegisterDebugCommand("fuel", func() {
		sg.playerShip.Storage.Store(&Item{
			Name:        "Fuel",
			Volume:      s.Storage.GetCapacity(STORE_LIQUID)/2 - sg.playerShip.Storage.GetItemVolume("Fuel"),
			StorageType: STORE_LIQUID,
		})
	})

	// burl.RegisterDebugCommand("earth", func() {
	// 	sg.playerShip.SetLocation(sg.galaxy.GetEarth())
	// 	sg.galaxyMenu.starchartMapView.systemFocus = sg.galaxy.GetStarSystem(sg.playerShip.GetCoords()) //this is messy.
	// })

	// for _, r := range sg.playerShip.Rooms {
	// 	burl.RegisterWatch(r.Name+" o2", &r.atmo.gasses)
	// 	burl.RegisterWatch(r.Name+" co2", &r.atmo.CO2)
	// 	burl.RegisterWatch(r.Name+" total", &r.atmo.pressure)
	// }

	return sg
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
	sg.InitWindow(false)

	sg.shipdisplay = burl.NewTileView(96, 46, 0, 3, 1, false)
	//sg.Stars = NewStarField(20, sg.shipdisplay)

	sg.output = burl.NewList(37, 8, 1, 45, 10, true, "Nothing to report, Captain!")
	sg.output.SetHint("PgUp/PgDown to scroll")
	sg.quickstats = NewQuickStatsWindow(39, 50, sg.playerShip)

	sg.timeDisplay = NewTimeDisplay(79, 50, sg.galaxy)
	sg.timeDisplay.UpdateSpeed(sg.simSpeed)

	sg.Window.Add(sg.output, sg.shipdisplay, sg.timeDisplay, sg.quickstats)

	sg.menus = make([]burl.UIElem, 0, MAX_MENUS)

	sg.gameMenu = NewGameMenu(sg.player)
	sg.shipMenu = NewShipMenu(sg.playerShip)
	sg.galaxyMenu = NewGalaxyMenu(sg.galaxy, sg.player.SpaceShip)
	sg.crewMenu = NewCrewMenu(sg.playerShip)
	sg.commMenu = NewCommsMenu(sg.playerShip.Comms)
	sg.viewMenu = NewViewMenu()
	sg.mainMenu = NewMainMenu()

	sg.menus = append(sg.menus, sg.gameMenu, sg.shipMenu, sg.galaxyMenu, sg.crewMenu, sg.commMenu, sg.viewMenu, sg.mainMenu)
	sg.Window.Add(sg.menus...)

	sg.menuButtons = make([]*burl.Button, 0, MAX_MENUS)

	sg.gameMenuButton = burl.NewButton(10, 1, 4, 1, 12, true, true, "Game")
	sg.gameMenuButton.SetHint("F1")
	sg.shipMenuButton = burl.NewButton(10, 1, 17, 1, 12, true, true, "Ship")
	sg.shipMenuButton.SetHint("F2")
	sg.galaxyMenuButton = burl.NewButton(10, 1, 30, 1, 12, true, true, "Galaxy")
	sg.galaxyMenuButton.SetHint("F3")
	sg.crewMenuButton = burl.NewButton(10, 1, 43, 1, 12, true, true, "Crew")
	sg.crewMenuButton.SetHint("F4")
	sg.commMenuButton = burl.NewButton(10, 1, 56, 1, 12, true, true, "Communications")
	sg.commMenuButton.SetHint("F5")
	sg.viewMenuButton = burl.NewButton(10, 1, 69, 1, 12, true, true, "View  Mode")
	sg.viewMenuButton.SetHint("F6")
	sg.mainMenuButton = burl.NewButton(10, 1, 82, 1, 12, true, true, "Main  Menu")
	sg.mainMenuButton.SetHint("ESC")

	sg.menuButtons = append(sg.menuButtons, sg.gameMenuButton, sg.shipMenuButton, sg.galaxyMenuButton, sg.crewMenuButton, sg.commMenuButton, sg.viewMenuButton, sg.mainMenuButton)

	for i := range sg.menuButtons {
		sg.Window.Add(sg.menuButtons[i])
	}

	//setup view mode colour palettes and other runtime jazz
	viewModeData[VIEW_ATMO_PRESSURE].SetTarget(sg.playerShip.LifeSupport.targetPressure)
	viewModeData[VIEW_ATMO_O2].SetTarget(sg.playerShip.LifeSupport.targetO2)
	viewModeData[VIEW_ATMO_TEMP].SetTarget(sg.playerShip.LifeSupport.targetTemp)
	viewModeData[VIEW_ATMO_CO2].SetTarget(sg.playerShip.LifeSupport.targetCO2)

	sg.CenterShip()
}

func (sg *SpaceshipGame) Update() {
	//simulation!
	for i := 0; i < sg.GetIncrement(); i++ {
		sg.galaxy.spaceTime++
		sg.playerShip.Update(sg.GetTime())

		//need starfield shift speed controlled here (currently hardcoded to shift every 100 seconds as long as the ship is moving)
		if sg.playerShip.GetSpeed() != 0 && sg.GetTick()%100 == 0 {
			//sg.Stars.Shift()
		}

		for i := range sg.player.MissionLog {
			sg.player.MissionLog[i].Update()
		}
	}

	sg.timeDisplay.UpdateTime()
}

func (sg *SpaceshipGame) HandleEvent(event *burl.Event) {
	switch event.ID {
	case burl.EV_UPDATE_UI:
		switch event.Message {
		case "inbox":
			sg.commMenu.UpdateInbox()
		case "transmissions":
			sg.commMenu.UpdateTransmissions()
		case "missions":
			sg.gameMenu.UpdateMissions()
		case "crew":
			if sg.activeMenu == sg.crewMenu {
				sg.crewMenu.UpdateCrewDetails()
			}
		case "stores":
			sg.shipMenu.storesMenu.Update()
		case "ship status":
			sg.quickstats.Update()
		case "ship move":
			if sg.activeMenu == sg.galaxyMenu {
				sg.galaxyMenu.Update()
			}
			sg.quickstats.Update()
		}
	case burl.EV_LIST_CYCLE:
		switch event.Caller {
		case sg.crewMenu.crewList:
			sg.crewMenu.UpdateCrewDetails()
		}
	case LOG_EVENT:
		sg.AddMessage(event.Message)
	}
}

func (sg *SpaceshipGame) Render() {
	//sg.Stars.Draw()
	sg.playerShip.DrawToTileView(sg.shipdisplay, sg.viewMenu.GetViewMode(), sg.viewX, sg.viewY)
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
	sg.CenterShip()
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
	//sg.Stars.dirty = true
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
