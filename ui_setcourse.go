package main

import "github.com/bennicholls/burl/ui"
import "github.com/veandco/go-sdl2/sdl"
import "github.com/bennicholls/burl/util"

import "strconv"

type Dialog interface {
	ui.UIElem
	HandleInput(key sdl.Keycode)
	Done() bool
}

type SetCourseDialog struct {
	ui.Container

	travelTimeText  *ui.Textbox
	travelSpeedText *ui.Textbox
	arrivalTimeText *ui.Textbox
	fuelGauge       *ui.ProgressBar
	places          *ui.List
	goButton        *ui.Button
	cancelButton    *ui.Button

	ship        *Ship
	destination Locatable
	distance    float64 //distance to destination in meters
	startTime   int
	course      Course
	done        bool
}

func NewSetCourseDialog(s *Ship, d Locatable, time int) *SetCourseDialog {
	sc := new(SetCourseDialog)
	sc.ship = s
	sc.destination = d
	sc.startTime = time
	sc.done = false

	sc.Container = *ui.NewContainer(58, 33, 11, 6, 50, true)
	sc.SetTitle("OFF WE GO!")
	sc.ToggleFocus()

	//left column
	courseLabel := ui.NewTextbox(26, 1, 1, 0, 0, false, true, "Setting Course For:")
	destName := ui.NewTextbox(26, 1, 1, 2, 0, true, true, d.GetName())
	destDescription := ui.NewTextbox(26, 12, 1, 4, 0, true, true, d.GetDescription())
	sc.places = ui.NewList(26, 14, 1, 18, 0, true, "Nothing in orbit! :(")
	sc.Add(courseLabel, destName, destDescription, sc.places)

	sc.distance = s.coords.CalcVector(d.GetCoords()).Distance * METERS_PER_LY
	distanceText := ui.NewTextbox(26, 1, 30, 2, 0, false, false, "Distance: "+strconv.Itoa(int(sc.distance/1000))+" km")
	orbitText := ui.NewTextbox(26, 1, 30, 3, 3, false, false, "Required Speed to Orbit: "+strconv.Itoa(int(d.GetVisitSpeed()/1000))+" km/s")
	shipSpeedText := ui.NewTextbox(26, 1, 30, 4, 3, false, false, "Current Ship Speed: "+strconv.Itoa(s.GetSpeed())+" m/s")
	maxFuelText := ui.NewTextbox(26, 1, 30, 5, 3, false, false, "Fuel Available "+strconv.Itoa(s.Fuel.Get())+" Litres")

	sc.Add(distanceText, orbitText, shipSpeedText, maxFuelText)

	sc.travelTimeText = ui.NewTextbox(26, 1, 30, 7, 3, false, false, "")
	sc.travelSpeedText = ui.NewTextbox(26, 1, 30, 8, 3, false, false, "")
	sc.arrivalTimeText = ui.NewTextbox(26, 1, 30, 9, 3, false, false, "")

	sc.fuelGauge = ui.NewProgressBar(26, 1, 30, 12, 0, true, true, "", 0xFF00FF00)
	sc.fuelGauge.SetProgress(50)

	sc.Add(sc.travelSpeedText, sc.travelTimeText, sc.arrivalTimeText, sc.fuelGauge)

	sc.UpdateCourse()

	sc.goButton = ui.NewButton(20, 1, 33, 18, 0, true, true, "This Looks Good, Let's Go!!")
	sc.goButton.ToggleFocus()
	sc.cancelButton = ui.NewButton(20, 1, 33, 21, 0, true, true, "On Second Thought, Let's Go Somewhere Else")

	sc.Add(sc.goButton, sc.cancelButton)

	return sc
}

func (sc *SetCourseDialog) HandleInput(key sdl.Keycode) {

	switch key {
	case sdl.K_LEFT:
		if sc.fuelGauge.GetProgress() > 10 {
			sc.fuelGauge.ChangeProgress(-5)
			sc.UpdateCourse()
		}
	case sdl.K_RIGHT:
		if sc.fuelGauge.GetProgress() < 100 {
			sc.fuelGauge.ChangeProgress(5)
			sc.UpdateCourse()
		}
	case sdl.K_UP, sdl.K_DOWN:
		sc.goButton.ToggleFocus()
		sc.cancelButton.ToggleFocus()
	case sdl.K_RETURN:
		if sc.goButton.IsFocused() {
			sc.ship.SetCourse(sc.destination, sc.course)
			sc.goButton.Press()
		} else {
			sc.cancelButton.Press()
		}
	}
}

func (sc *SetCourseDialog) UpdateCourse() {

	maxFuel := util.Min(sc.ship.Fuel.Get(), sc.ship.Engine.FuelUse*int(sc.ship.Navigation.CalcMaxBurnTime(sc.destination.GetVisitSpeed(), sc.distance)))
	c := sc.ship.Navigation.ComputeCourse(sc.destination, maxFuel*sc.fuelGauge.GetProgress()/100, sc.startTime)

	sc.fuelGauge.ChangeText("Fuel to burn: " + strconv.Itoa(maxFuel*sc.fuelGauge.GetProgress()/100))
	sc.travelTimeText.ChangeText("Travel Time: " + GetTimeString(c.totalTime))

	speed := sc.ship.GetSpeed() + int(float64(c.accelTime-c.startTime)*sc.ship.Engine.Thrust)
	sc.travelSpeedText.ChangeText("Max Travel Speed: " + strconv.Itoa(speed/1000) + " km/s")
	sc.arrivalTimeText.ChangeText("Arrival Time: " + GetTimeString(c.arrivaltime))

	sc.course = c
}

func (sc SetCourseDialog) Done() bool {
	if sc.goButton.IsFocused() {
		if sc.goButton.PressPulse.IsFinished() {
			return true
		}
	} else {
		if sc.cancelButton.PressPulse.IsFinished() {
			return true
		}
	}

	return false
}
