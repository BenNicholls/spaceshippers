package main

import (
	"github.com/bennicholls/burl-E/burl"
	"github.com/bennicholls/tyumi/gfx/col"
	"github.com/bennicholls/tyumi/gfx/ui"
	"github.com/bennicholls/tyumi/vec"
)

type QuickStatsWindow struct {
	ui.Element

	hullBar   ui.ProgressBar
	fuelBar   ui.ProgressBar
	energyBar ui.ProgressBar
	courseBar ui.ProgressBar

	alertText ui.Textbox

	ship *Ship
}

func (qsw *QuickStatsWindow) Init(s *Ship) {
	qsw.Element.Init(vec.Dims{39, 3}, vec.Coord{39, 50}, menuDepth)
	qsw.EnableBorder()

	qsw.alertText.Init(vec.Dims{39, 1}, vec.ZERO_COORD, ui.BorderDepth, "All Optimal", ui.JUSTIFY_CENTER)
	qsw.alertText.EnableBorder()
	qsw.hullBar.Init(vec.Dims{9, 1}, vec.Coord{0, 2}, ui.BorderDepth, col.RED, "HULL")
	qsw.hullBar.EnableBorder()
	qsw.fuelBar.Init(vec.Dims{9, 1}, vec.Coord{10, 2}, ui.BorderDepth, col.PURPLE, "FUEL")
	qsw.fuelBar.EnableBorder()
	qsw.energyBar.Init(vec.Dims{9, 1}, vec.Coord{20, 2}, ui.BorderDepth, col.BLUE, "ENERGY")
	qsw.energyBar.EnableBorder()
	qsw.courseBar.Init(vec.Dims{9, 1}, vec.Coord{30, 2}, ui.BorderDepth, col.GREEN, "COURSE")
	qsw.courseBar.EnableBorder()

	qsw.AddChild(&qsw.alertText)
	qsw.AddChildren(&qsw.hullBar, &qsw.fuelBar, &qsw.energyBar, &qsw.courseBar)

	qsw.ship = s
	qsw.Update()

	return
}

func (qsw *QuickStatsWindow) Update() {
	qsw.hullBar.SetProgress(qsw.ship.Hull.GetPct())
	//qsw.fuelBar.SetProgress(int(100 * qsw.ship.Storage.GetItemVolume("Fuel") / float64(qsw.ship.Storage.GetStat(STAT_FUEL_STORAGE))))
	//qsw.powerBar.ChangeProgress(qsw.ship.PowerSystem.GetPowerUsagePct())
	qsw.courseBar.SetProgress(qsw.ship.Navigation.GetCurrentProgress())
}

type ShipStatsWindow struct {
	burl.Container

	name        *burl.Textbox
	speed       *burl.Textbox
	fuel        *burl.ProgressBar
	location    *burl.Textbox
	destination *burl.Textbox

	playerShip *Ship
}

// func NewShipStatsWindow(x, y int, ship *Ship) *ShipStatsWindow {
// 	ss := new(ShipStatsWindow)
// 	ss.playerShip = ship

// 	ss.Container = *burl.NewContainer(26, 10, x, y, 10, true)
// 	ss.name = burl.NewTextbox(26, 1, 0, 0, 1, false, true, ss.playerShip.Name)
// 	ss.speed = burl.NewTextbox(26, 1, 0, 2, 1, false, false, "Speed: "+strconv.Itoa(ss.playerShip.GetSpeed()))
// 	ss.fuel = burl.NewProgressBar(26, 1, 0, 3, 1, false, false, "", burl.COL_GREEN)
// 	ss.location = burl.NewTextbox(26, 1, 0, 8, 1, false, false, "")
// 	ss.destination = burl.NewTextbox(26, 1, 0, 9, 1, false, false, "")

// 	ss.Add(ss.name, ss.speed, ss.fuel, ss.location, ss.destination)

// 	return ss
// }

// func (ss *ShipStatsWindow) Update() {
// 	ss.name.ChangeText(ss.playerShip.Name)

// 	speed := float64(ss.playerShip.GetSpeed())
// 	switch {
// 	case speed < 1000:
// 		ss.speed.ChangeText("Speed: " + strconv.FormatFloat(speed, 'f', 0, 64) + " m/s")
// 	case speed < 100000000:
// 		ss.speed.ChangeText("Speed: " + strconv.FormatFloat(speed/1000, 'f', 2, 64) + " km/s")
// 	default:
// 		ss.speed.ChangeText("Speed: " + strconv.FormatFloat(speed/float64(METERS_PER_LY), 'f', 4, 64) + "c")
// 	}

// 	ss.fuel.ChangeText("Fuel: " + ss.playerShip.Fuel.String() + " Litres")
// 	ss.fuel.SetProgress(ss.playerShip.Fuel.GetPct())

// 	locString := "Location: "
// 	dstString := "Destination: "
// 	if ss.playerShip.currentLocation != nil {
// 		locString += ss.playerShip.currentLocation.GetName()
// 	} else {
// 		locString += "NO LOCATION. HOW'D YOU DO THIS."
// 	}
// 	if ss.playerShip.destination != nil {
// 		dstString += ss.playerShip.destination.GetName()
// 	} else {
// 		dstString += "NO DESTINATION. Let's go somewhere!!"
// 	}
// 	ss.location.ChangeText(locString)
// 	ss.destination.ChangeText(dstString)
// }
