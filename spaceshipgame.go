package main

import (
	"encoding/gob"
	"os"

	"github.com/bennicholls/tyumi"
	"github.com/bennicholls/tyumi/event"
	"github.com/bennicholls/tyumi/gfx"
	"github.com/bennicholls/tyumi/gfx/col"
	"github.com/bennicholls/tyumi/gfx/ui"
	"github.com/bennicholls/tyumi/log"
	"github.com/bennicholls/tyumi/rl"
	"github.com/bennicholls/tyumi/vec"
)

type Menu interface {
	Show()
	Hide()
}

func init() {
	//need to register types that might be hidden by an interface, in order for them to be serializable
	gob.Register(&SleepJob{})
}

type SpaceshipGame struct {
	tyumi.Scene

	//ui stuff
	logOutput   ui.List
	quickstats  QuickStatsWindow
	timeDisplay TimeDisplay
	shipView    rl.TileMapView
	stars       StarField

	//top menu. contains buttons for submenus
	gameMenuButton   ui.Button
	shipMenuButton   ui.Button
	galaxyMenuButton ui.Button
	crewMenuButton   ui.Button
	commMenuButton   ui.Button
	viewMenuButton   ui.Button
	mainMenuButton   ui.Button

	gameMenu   GameMenu    //(F1)
	shipMenu   *ShipMenu   //(F2)
	galaxyMenu *GalaxyMenu //(F3)
	crewMenu   CrewMenu    //(F4)
	commMenu   CommMenu    //(F5)
	viewMenu   ViewMenu    //(F6)
	mainMenu   MainMenu    //(ESC)

	activeMenu Menu

	//Time Globals.
	startTime int //time since launch, measured in Standard Galactic Seconds
	simSpeed  int //4 speeds, plus pause (0)
	paused    bool

	viewMode int

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

	//STUFF FOR TESTING
	sg.player.AddMission(GenerateGoToMission(sg.playerShip, sg.galaxy.GetEarth(), nil))
	sg.player.AddMission(GenerateGoToMission(sg.playerShip, sg.playerShip, sg.galaxy.GetEarth()))
	sg.playerShip.Comms.AddRandomTransmission(10)
	sg.gameMenu.UpdateMissions()
	// REMOVE THIS LATER

	//sg.OpenDialog(NewSpaceEventDialog(spaceEvents[1]))

	// burl.RegisterDebugCommand("fuel", func() {
	// 	sg.playerShip.Storage.Store(&Item{
	// 		Name:        "Fuel",
	// 		Volume:      s.Storage.GetCapacity(STORE_LIQUID)/2 - sg.playerShip.Storage.GetItemVolume("Fuel"),
	// 		StorageType: STORE_LIQUID,
	// 	})
	// })

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

// Centers the map of the ship in the main view.
func (sg *SpaceshipGame) CenterShip() {
	sg.shipView.CenterOnTileMapCoord(sg.playerShip.Bounds().Center())
	if sg.activeMenu != nil {
		sg.shipView.MoveCamera(28, 0)
	}
}

func (sg *SpaceshipGame) SetupUI() {
	sg.Init()
	sg.Window().SendEventsToUnfocused = true

	sg.shipView.Init(vec.Dims{96, 46}, vec.Coord{0, 3}, 1, &sg.playerShip.shipMap)
	sg.shipView.SetDefaultVisuals(gfx.Visuals{Mode: gfx.DRAW_NONE, Colours: col.Pair{col.WHITE, col.BLACK}})
	sg.shipView.CenterOnTileMapCoord(sg.playerShip.Bounds().Center())
	sg.stars.Init(sg.shipView.Size(), vec.Coord{0, 3}, 0, 20, 0)

	sg.logOutput.Init(vec.Dims{37, 8}, vec.Coord{1, 45}, menuDepth)
	sg.logOutput.SetupBorder("SPACE LOG", "PgUp/PgDown to scroll")
	sg.logOutput.SetEmptyText("Nothing to report, Captain!")

	sg.quickstats.Init(sg.playerShip)

	sg.timeDisplay.Init(vec.Coord{79, 50}, sg.galaxy)
	sg.timeDisplay.UpdateSpeed(sg.simSpeed)

	sg.Window().AddChildren(&sg.logOutput, &sg.shipView, &sg.stars, &sg.timeDisplay, &sg.quickstats)

	sg.gameMenu.Init(sg.player)
	// sg.shipMenu = NewShipMenu(sg.playerShip)
	// sg.galaxyMenu = NewGalaxyMenu(sg.galaxy, sg.player.SpaceShip)
	sg.crewMenu.Init(sg.playerShip)
	sg.commMenu.Init(sg.playerShip.Comms)
	sg.viewMenu.Init()
	sg.mainMenu.Init()

	sg.Window().AddChildren(&sg.gameMenu, &sg.crewMenu, &sg.commMenu, &sg.viewMenu, &sg.mainMenu)

	sg.gameMenuButton.Init(vec.Dims{10, 1}, vec.Coord{4, 1}, 10, "Game", func() { sg.ActivateMenu(&sg.gameMenu) })
	sg.gameMenuButton.SetupBorder("", "F1")
	sg.shipMenuButton.Init(vec.Dims{10, 1}, vec.Coord{17, 1}, 10, "Ship", nil)
	sg.shipMenuButton.SetupBorder("", "F2")
	sg.galaxyMenuButton.Init(vec.Dims{10, 1}, vec.Coord{30, 1}, 10, "Galaxy", nil)
	sg.galaxyMenuButton.SetupBorder("", "F3")
	sg.crewMenuButton.Init(vec.Dims{10, 1}, vec.Coord{43, 1}, 10, "Crew", func() { sg.ActivateMenu(&sg.crewMenu) })
	sg.crewMenuButton.SetupBorder("", "F4")
	sg.commMenuButton.Init(vec.Dims{10, 1}, vec.Coord{56, 1}, 10, "Communications", func() { sg.ActivateMenu(&sg.commMenu) })
	sg.commMenuButton.SetupBorder("", "F5")
	sg.viewMenuButton.Init(vec.Dims{10, 1}, vec.Coord{69, 1}, 10, "View Mode", func() { sg.ActivateMenu(&sg.viewMenu) })
	sg.viewMenuButton.SetupBorder("", "F6")
	sg.mainMenuButton.Init(vec.Dims{10, 1}, vec.Coord{82, 1}, 10, "Main Menu", func() { sg.ActivateMenu(&sg.mainMenu) })
	sg.mainMenuButton.SetupBorder("", "ESC")
	sg.Window().AddChildren(&sg.gameMenuButton, &sg.shipMenuButton, &sg.galaxyMenuButton, &sg.crewMenuButton, &sg.commMenuButton, &sg.viewMenuButton, &sg.mainMenuButton)

	//setup view mode colour palettes and other runtime jazz
	viewModeData[VIEW_ATMO_PRESSURE].SetTarget(sg.playerShip.LifeSupport.targetPressure)
	viewModeData[VIEW_ATMO_O2].SetTarget(sg.playerShip.LifeSupport.targetO2)
	viewModeData[VIEW_ATMO_TEMP].SetTarget(sg.playerShip.LifeSupport.targetTemp)
	viewModeData[VIEW_ATMO_CO2].SetTarget(sg.playerShip.LifeSupport.targetCO2)

	sg.CenterShip()

	sg.SetKeypressHandler(sg.HandleKeypress)
	sg.SetEventHandler(sg.HandleEvent)
}

func (sg *SpaceshipGame) Update() {
	//simulation!
	for range sg.GetIncrement() {
		sg.galaxy.spaceTime++
		sg.playerShip.Update(sg.GetTime())

		//need starfield shift speed controlled here (currently hardcoded to shift every 100 seconds as long as the ship is moving)
		if sg.playerShip.GetSpeed() != 0 {
			sg.stars.shiftFrequency = 100
		} else {
			sg.stars.shiftFrequency = 0
		}

		for i := range sg.player.MissionLog {
			sg.player.MissionLog[i].Update()
		}
	}

	if sg.GetIncrement() > 0 {
		sg.timeDisplay.UpdateTime()
	}
}

func (sg *SpaceshipGame) HandleEvent(event event.Event) (event_handled bool) {
	switch event.ID() {
	case EV_LOG:
		sg.AddLogMessage(event.(*SpaceLogEvent).message)
		return true
	}
	// switch event.ID {
	// case burl.EV_UPDATE_UI:
	// 	switch event.Message {
	// 	case "inbox":
	// 		sg.commMenu.UpdateInbox()
	// 	case "transmissions":
	// 		sg.commMenu.UpdateTransmissions()
	// 	case "missions":
	// 		sg.gameMenu.UpdateMissions()
	// 	case "stores":
	// 		sg.shipMenu.storesMenu.Update()
	// 	case "ship status":
	// 		sg.quickstats.Update()
	// 	case "ship move":
	// 		if sg.activeMenu == sg.galaxyMenu {
	// 			sg.galaxyMenu.Update()
	// 		}
	// 		sg.quickstats.Update()
	// 	}
	// case burl.EV_LIST_CYCLE:
	// 	switch event.Caller {
	// 	case sg.crewMenu.crewList:
	// 		sg.crewMenu.UpdateCrewDetails()
	// 	}

	return
}

// Activates a menu (crew, rooms, systems, etc). Deactivates menu if menu already active.
func (sg *SpaceshipGame) ActivateMenu(menu Menu) {
	if sg.activeMenu == menu {
		//double-activating means close
		sg.activeMenu.Hide()
		sg.activeMenu = nil
		sg.CenterShip()
		return
	}

	if sg.activeMenu != nil {
		sg.activeMenu.Hide()
	}

	sg.activeMenu = menu
	sg.activeMenu.Show()
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

// returns the number of simulated seconds since launch
func (sg SpaceshipGame) GetTick() int {
	return sg.galaxy.spaceTime - sg.startTime
}

// gets the time from the Galaxy
func (sg SpaceshipGame) GetTime() int {
	return sg.galaxy.spaceTime
}

func (sg SpaceshipGame) Shutdown() {
	sg.SaveShip()
}

func (sg *SpaceshipGame) SaveShip() {
	f, err := os.Create("savefile")
	if err != nil {
		log.Error("Could not open file for saving: " + err.Error())
	}
	defer f.Close()

	enc := gob.NewEncoder(f)
	err = enc.Encode(sg.playerShip)
	if err != nil {
		log.Error("Could not save ship: " + err.Error())
	}
}

func (sg *SpaceshipGame) LoadShip() {
	f, err := os.Open("savefile")
	if err != nil {
		log.Error("Could not open file for loading: " + err.Error())
	}
	defer f.Close()

	s := new(Ship)

	dec := gob.NewDecoder(f)
	err = dec.Decode(s)
	if err != nil {
		log.Error("Could not load ship: " + err.Error())
	}

	//data loaded, now to re-init everything
	s.SetupShip(sg.galaxy)

	//load complete, make the switch!!
	sg.player.SpaceShip = s
	sg.playerShip = s

	sg.SetupUI()
}
