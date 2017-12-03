package main

import "github.com/bennicholls/burl-E/burl"
import "strconv"

type ShipStatsWindow struct {
	burl.Container

	name        *burl.Textbox
	speed       *burl.Textbox
	fuel        *burl.ProgressBar
	location    *burl.Textbox
	destination *burl.Textbox

	playerShip *Ship
}

func NewShipStatsWindow(ship *Ship) *ShipStatsWindow {
	ss := new(ShipStatsWindow)
	ss.playerShip = ship

	ss.Container = *burl.NewContainer(26, 12, 1, 32, 10, true)
	ss.name = burl.NewTextbox(26, 1, 0, 0, 1, false, true, ss.playerShip.name)
	ss.speed = burl.NewTextbox(26, 1, 0, 2, 1, false, false, "Speed: "+strconv.Itoa(ss.playerShip.GetSpeed()))
	ss.fuel = burl.NewProgressBar(26, 1, 0, 3, 1, false, false, "", burl.COL_GREEN)
	ss.location = burl.NewTextbox(26, 1, 0, 10, 1, false, false, "")
	ss.destination = burl.NewTextbox(26, 1, 0, 11, 1, false, false, "")

	ss.Add(ss.name, ss.speed, ss.fuel, ss.location, ss.destination)

	return ss
}

func (ss *ShipStatsWindow) Update() {
	ss.name.ChangeText(ss.playerShip.name)

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
	if ss.playerShip.CurrentLocation != nil {
		locString += ss.playerShip.CurrentLocation.GetName()
	} else {
		locString += "NO LOCATION. HOW'D YOU DO THIS."
	}
	if ss.playerShip.Destination != nil {
		dstString += ss.playerShip.Destination.GetName()
	} else {
		dstString += "NO DESTINATION. Let's go somewhere!!"
	}
	ss.location.ChangeText(locString)
	ss.destination.ChangeText(dstString)
}
