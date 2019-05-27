package main

import "github.com/bennicholls/burl-E/burl"

type QuickStatsWindow struct {
	burl.Container

	hullBar   *burl.ProgressBar
	fuelBar   *burl.ProgressBar
	energyBar *burl.ProgressBar
	courseBar *burl.ProgressBar

	alertText *burl.Textbox

	ship *Ship
}

func NewQuickStatsWindow(x, y int, s *Ship) (qsw *QuickStatsWindow) {
	qsw = new(QuickStatsWindow)
	qsw.Container = *burl.NewContainer(39, 3, x, y, 9, true)

	qsw.alertText = burl.NewTextbox(39, 1, 0, 0, 1, true, true, "All Optimal")
	qsw.hullBar = burl.NewProgressBar(9, 1, 0, 2, 1, true, true, "HULL", burl.COL_RED)
	qsw.fuelBar = burl.NewProgressBar(9, 1, 10, 2, 1, true, true, "FUEL", burl.COL_PURPLE)
	qsw.energyBar = burl.NewProgressBar(9, 1, 20, 2, 1, true, true, "ENERGY", burl.COL_BLUE)
	qsw.courseBar = burl.NewProgressBar(9, 1, 30, 2, 1, true, true, "COURSE", burl.COL_GREEN)

	qsw.Add(qsw.alertText)
	qsw.Add(qsw.hullBar, qsw.fuelBar, qsw.energyBar, qsw.courseBar)

	qsw.ship = s

	qsw.Update()

	return
}

func (qsw *QuickStatsWindow) Update() {
	qsw.hullBar.SetProgress(qsw.ship.Hull.GetPct())
	qsw.fuelBar.SetProgress(int(100 * qsw.ship.Storage.GetItemVolume("Fuel") / float64(qsw.ship.Storage.GetStat(STAT_FUEL_STORAGE))))
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
