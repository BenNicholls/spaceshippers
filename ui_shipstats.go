package main

import "github.com/bennicholls/burl-E/burl"
import "strconv"

type QuickStatsWindow struct {
	burl.Container

	hullBar   *burl.ProgressBar
	fuelBar   *burl.ProgressBar
	energyBar *burl.ProgressBar
	voyageBar *burl.ProgressBar

	alertText *burl.Textbox

	ship *Ship
}

func NewQuickStatsWindow(x, y int, s *Ship) (qsw *QuickStatsWindow) {
	qsw = new(QuickStatsWindow)
	qsw.Container = *burl.NewContainer(23, 3, x, y, 9, true)

	qsw.alertText = burl.NewTextbox(23, 1, 0, 0, 1, true, true, "All Optimal")
	qsw.hullBar = burl.NewProgressBar(5, 1, 0, 2, 1, true, true, "HULL", burl.COL_RED)
	qsw.fuelBar = burl.NewProgressBar(5, 1, 6, 2, 1, true, true, "FUEL", burl.COL_PURPLE)
	qsw.energyBar = burl.NewProgressBar(5, 1, 12, 2, 1, true, true, "ENERGY", burl.COL_BLUE)
	qsw.voyageBar = burl.NewProgressBar(5, 1, 18, 2, 1, true, true, "VOYAGE", burl.COL_GREEN)

	qsw.Add(qsw.alertText)
	qsw.Add(qsw.hullBar, qsw.fuelBar, qsw.energyBar, qsw.voyageBar)

	qsw.ship = s

	qsw.Update()

	return
}

func (qsw *QuickStatsWindow) Update() {
	qsw.hullBar.ChangeProgress(qsw.ship.Hull.GetPct())
	qsw.fuelBar.ChangeProgress(qsw.ship.Fuel.GetPct())
	//qsw.powerBar.ChangeProgress(qsw.ship.PowerSystem.GetPowerUsagePct())
	//qsw.powerBar.ChangeProgress(qsw.ship.Navigation.GetCurrentProgress())
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

func NewShipStatsWindow(x, y int, ship *Ship) *ShipStatsWindow {
	ss := new(ShipStatsWindow)
	ss.playerShip = ship

	ss.Container = *burl.NewContainer(26, 10, x, y, 10, true)
	ss.name = burl.NewTextbox(26, 1, 0, 0, 1, false, true, ss.playerShip.Name)
	ss.speed = burl.NewTextbox(26, 1, 0, 2, 1, false, false, "Speed: "+strconv.Itoa(ss.playerShip.GetSpeed()))
	ss.fuel = burl.NewProgressBar(26, 1, 0, 3, 1, false, false, "", burl.COL_GREEN)
	ss.location = burl.NewTextbox(26, 1, 0, 8, 1, false, false, "")
	ss.destination = burl.NewTextbox(26, 1, 0, 9, 1, false, false, "")

	ss.Add(ss.name, ss.speed, ss.fuel, ss.location, ss.destination)

	return ss
}

func (ss *ShipStatsWindow) Update() {
	ss.name.ChangeText(ss.playerShip.Name)

	speed := float64(ss.playerShip.GetSpeed())
	switch {
	case speed < 1000:
		ss.speed.ChangeText("Speed: " + strconv.FormatFloat(speed, 'f', 0, 64) + " m/s")
	case speed < 100000000:
		ss.speed.ChangeText("Speed: " + strconv.FormatFloat(speed/1000, 'f', 2, 64) + " km/s")
	default:
		ss.speed.ChangeText("Speed: " + strconv.FormatFloat(speed/float64(METERS_PER_LY), 'f', 4, 64) + "c")
	}

	ss.fuel.ChangeText("Fuel: " + ss.playerShip.Fuel.String() + " Litres")
	ss.fuel.SetProgress(ss.playerShip.Fuel.GetPct())

	locString := "Location: "
	dstString := "Destination: "
	if ss.playerShip.currentLocation != nil {
		locString += ss.playerShip.currentLocation.GetName()
	} else {
		locString += "NO LOCATION. HOW'D YOU DO THIS."
	}
	if ss.playerShip.destination != nil {
		dstString += ss.playerShip.destination.GetName()
	} else {
		dstString += "NO DESTINATION. Let's go somewhere!!"
	}
	ss.location.ChangeText(locString)
	ss.destination.ChangeText(dstString)
}
