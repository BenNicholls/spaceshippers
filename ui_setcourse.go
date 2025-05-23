package main

import (
	"fmt"

	"github.com/bennicholls/tyumi"
	"github.com/bennicholls/tyumi/gfx/col"
	"github.com/bennicholls/tyumi/gfx/ui"
	"github.com/bennicholls/tyumi/input"
	"github.com/bennicholls/tyumi/vec"
)

type SetCourseDialog struct {
	tyumi.Scene

	done bool

	travelTimeText  ui.Textbox
	travelSpeedText ui.Textbox
	arrivalTimeText ui.Textbox
	fuelGauge       ui.ProgressBar
	places          ui.List
	goButton        ui.Button
	cancelButton    ui.Button

	ship        *Ship
	destination Locatable
	distance    float64 //distance to destination in meters
	startTime   int
	course      Course
}

func NewSetCourseDialog(s *Ship, d Locatable, time int) (sc *SetCourseDialog) {
	sc = new(SetCourseDialog)
	sc.ship = s
	sc.destination = d
	sc.startTime = time

	sc.InitCentered(vec.Dims{58, 33})
	sc.Window().SetupBorder("OFF WE GO!", "")
	windowStyle := ui.DefaultBorderStyle
	windowStyle.TitleJustification = ui.JUSTIFY_CENTER
	windowStyle.Colours = col.Pair{col.ORANGE, col.DARKGREY}
	sc.Window().Border.SetStyle(ui.BORDER_STYLE_CUSTOM, windowStyle)
	sc.SetKeypressHandler(sc.HandleKeypress)

	//left column
	courseLabel := ui.NewTextbox(vec.Dims{26, 1}, vec.Coord{1, 1}, 0, "Setting Course For:", ui.JUSTIFY_CENTER)
	destName := ui.NewTitleTextbox(vec.Dims{26, 1}, vec.Coord{1, 2}, 0, d.GetName())
	destDescription := ui.NewTitleTextbox(vec.Dims{26, 13}, vec.Coord{1, 4}, 0, d.GetDescription())
	sc.places.Init(vec.Dims{26, 14}, vec.Coord{1, 18}, 0)
	sc.places.EnableBorder()
	sc.places.SetEmptyText("Nothing in orbit! :()")
	sc.Window().AddChildren(courseLabel, destName, destDescription, &sc.places)

	sc.distance = s.Coords.CalcVector(d.GetCoords()).Distance * METERS_PER_LY
	distanceText := ui.NewTextbox(vec.Dims{26, 1}, vec.Coord{30, 2}, 0, fmt.Sprintf("Distance: %.0f km", sc.distance/1000), ui.JUSTIFY_LEFT)
	orbitText := ui.NewTextbox(vec.Dims{26, 1}, vec.Coord{30, 3}, 0, fmt.Sprintf("Required Speed to Orbit: %.0f km/s", d.GetVisitSpeed()/1000), ui.JUSTIFY_LEFT)
	shipSpeedText := ui.NewTextbox(vec.Dims{26, 1}, vec.Coord{30, 4}, 0, fmt.Sprintf("Current Ship Speed: %d m/s", s.GetSpeed()), ui.JUSTIFY_LEFT)
	maxFuelText := ui.NewTextbox(vec.Dims{26, 1}, vec.Coord{30, 5}, 0, fmt.Sprintf("Fuel Available: %.0f Litres", s.Storage.GetItemVolume("Fuel")), ui.JUSTIFY_LEFT)
	fuelUseText := ui.NewTextbox(vec.Dims{26, 1}, vec.Coord{30, 6}, 0, fmt.Sprintf("Fuel Use Rate: %.2f Litres per second", s.Engine.FuelUse), ui.JUSTIFY_LEFT)
	engineThrustText := ui.NewTextbox(vec.Dims{26, 1}, vec.Coord{30, 7}, 0, fmt.Sprintf("Total Engine Thrust: %.1f m/s/s", s.Engine.Thrust), ui.JUSTIFY_LEFT)

	sc.Window().AddChildren(distanceText, orbitText, shipSpeedText, maxFuelText, fuelUseText, engineThrustText)

	sc.travelTimeText.Init(vec.Dims{26, 1}, vec.Coord{30, 9}, 0, "", ui.JUSTIFY_LEFT)
	sc.travelSpeedText.Init(vec.Dims{26, 1}, vec.Coord{30, 10}, 0, "", ui.JUSTIFY_LEFT)
	sc.arrivalTimeText.Init(vec.Dims{26, 1}, vec.Coord{30, 11}, 0, "", ui.JUSTIFY_LEFT)
	sc.fuelGauge.Init(vec.Dims{26, 1}, vec.Coord{30, 14}, 0, col.GREEN, "")
	sc.fuelGauge.SetupBorder("", "<-/->")
	sc.fuelGauge.SetProgress(50)
	sc.Window().AddChildren(&sc.travelSpeedText, &sc.travelTimeText, &sc.arrivalTimeText, &sc.fuelGauge)

	sc.goButton.Init(vec.Dims{20, 1}, vec.Coord{33, 20}, 1, "This Looks Good, Let's Go!!", func() {
		sc.ship.SetCourse(sc.destination, sc.course)
		fireSpaceLogEvent("Setting course for " + sc.destination.GetName())
		sc.CreateTimer(20, func() { sc.done = true })
	})
	sc.goButton.EnableBorder()
	sc.cancelButton.Init(vec.Dims{20, 1}, vec.Coord{33, 23}, 1, "On Second Thought, nevermind.", func() {
		fireSpaceLogEvent("Course selection cancelled.")
		sc.CreateTimer(20, func() { sc.done = true })
	})
	sc.cancelButton.EnableBorder()
	sc.Window().AddChildren(&sc.goButton, &sc.cancelButton)

	sc.Window().SetTabbingOrder(&sc.goButton, &sc.cancelButton)
	sc.goButton.Focus()

	sc.UpdateCourse()

	return sc
}

func (sc *SetCourseDialog) HandleKeypress(key_event *input.KeyboardEvent) (event_handled bool) {
	switch key_event.Direction() {
	case vec.DIR_LEFT:
		if fuel := sc.fuelGauge.GetProgress(); fuel > 10 {
			sc.fuelGauge.SetProgress(fuel - 5)
			sc.UpdateCourse()
			return true
		}
	case vec.DIR_RIGHT:
		if fuel := sc.fuelGauge.GetProgress(); fuel < 100 {
			sc.fuelGauge.SetProgress(fuel + 5)
			sc.UpdateCourse()
			return true
		}
	}

	return
}

func (sc *SetCourseDialog) UpdateCourse() {
	maxFuel := min(sc.ship.Storage.GetItemVolume("Fuel"), sc.ship.Engine.FuelUse*(sc.ship.Navigation.CalcMaxBurnTime(sc.destination.GetVisitSpeed(), sc.distance)))
	c := sc.ship.Navigation.ComputeCourse(sc.destination, maxFuel*sc.fuelGauge.GetProgressNormalized(), sc.startTime)

	sc.fuelGauge.ChangeText("Fuel to burn: " + fmt.Sprint(maxFuel*sc.fuelGauge.GetProgressNormalized()))
	sc.travelTimeText.ChangeText(fmt.Sprintf("Travel Time: %s", GetDurationString(c.TotalTime)))

	speed := sc.ship.GetSpeed() + int(float64(c.AccelTime-c.StartTime)*sc.ship.Engine.Thrust)
	sc.travelSpeedText.ChangeText(fmt.Sprintf("Max Travel Speed: %d km/s", speed/1000))
	sc.arrivalTimeText.ChangeText(fmt.Sprintf("Arrival Time: %s, %s", GetTimeString(c.Arrivaltime), GetDateString(c.Arrivaltime)))

	sc.course = c
}

func (sc *SetCourseDialog) Done() bool {
	return sc.done
}
