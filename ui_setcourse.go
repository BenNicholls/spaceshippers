package main

import "github.com/bennicholls/burl-E/burl"
import "github.com/veandco/go-sdl2/sdl"
import "strconv"

type SetCourseDialog struct {
	burl.StatePrototype

	travelTimeText  *burl.Textbox
	travelSpeedText *burl.Textbox
	arrivalTimeText *burl.Textbox
	fuelGauge       *burl.ProgressBar
	places          *burl.List
	goButton        *burl.Button
	cancelButton    *burl.Button

	ship        *Ship
	destination Locatable
	distance    float64 //distance to destination in meters
	startTime   int
	course      Course
}

func NewSetCourseDialog(s *Ship, d Locatable, time int) *SetCourseDialog {
	sc := new(SetCourseDialog)
	sc.ship = s
	sc.destination = d
	sc.startTime = time

	sc.Window = burl.NewContainer(58, 33, 1, 1, 50, true)
	sc.Window.CenterInConsole()
	sc.Window.SetTitle("OFF WE GO!")
	sc.Window.ToggleFocus()

	//left column
	courseLabel := burl.NewTextbox(26, 1, 1, 0, 0, false, true, "Setting Course For:")
	destName := burl.NewTextbox(26, 1, 1, 2, 1, true, true, d.GetName())
	destDescription := burl.NewTextbox(26, 13, 1, 4, 1, true, true, d.GetDescription())
	sc.places = burl.NewList(26, 14, 1, 18, 1, true, "Nothing in orbit! :(")
	sc.Window.Add(courseLabel, destName, destDescription, sc.places)

	sc.distance = s.Coords.CalcVector(d.GetCoords()).Distance * METERS_PER_LY
	distanceText := burl.NewTextbox(26, 1, 30, 2, 0, false, false, "Distance: "+strconv.Itoa(int(sc.distance/1000))+" km")
	orbitText := burl.NewTextbox(26, 1, 30, 3, 3, false, false, "Required Speed to Orbit: "+strconv.Itoa(int(d.GetVisitSpeed()/1000))+" km/s")
	shipSpeedText := burl.NewTextbox(26, 1, 30, 4, 3, false, false, "Current Ship Speed: "+strconv.Itoa(s.GetSpeed())+" m/s")
	maxFuelText := burl.NewTextbox(26, 1, 30, 5, 3, false, false, "Fuel Available: "+strconv.Itoa(s.Fuel.Get())+" Litres")
	fuelUseText := burl.NewTextbox(26, 1, 30, 6, 3, false, false, "Fuel Use Rate: "+strconv.Itoa(s.Engine.FuelUse)+" Litres per second")
	engineThrustText := burl.NewTextbox(26, 1, 30, 7, 3, false, false, "Total Engine Thrust: "+strconv.Itoa(int(s.Engine.Thrust))+" m/s/s")

	sc.Window.Add(distanceText, orbitText, shipSpeedText, maxFuelText, fuelUseText, engineThrustText)

	sc.travelTimeText = burl.NewTextbox(26, 1, 30, 9, 3, false, false, "")
	sc.travelSpeedText = burl.NewTextbox(26, 1, 30, 10, 3, false, false, "")
	sc.arrivalTimeText = burl.NewTextbox(26, 1, 30, 11, 3, false, false, "")

	sc.fuelGauge = burl.NewProgressBar(26, 1, 30, 14, 1, true, true, "", burl.COL_GREEN)
	sc.fuelGauge.SetProgress(50)

	sc.Window.Add(sc.travelSpeedText, sc.travelTimeText, sc.arrivalTimeText, sc.fuelGauge)

	sc.UpdateCourse()

	sc.goButton = burl.NewButton(20, 1, 33, 20, 1, true, true, "This Looks Good, Let's Go!!")
	sc.goButton.ToggleFocus()
	sc.cancelButton = burl.NewButton(20, 1, 33, 23, 2, true, true, "On Second Thought, nevermind.")

	sc.Window.Add(sc.goButton, sc.cancelButton)

	return sc
}

func (sc *SetCourseDialog) HandleKeypress(key sdl.Keycode) {
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
	maxFuel := burl.Min(sc.ship.Fuel.Get(), sc.ship.Engine.FuelUse*int(sc.ship.Navigation.CalcMaxBurnTime(sc.destination.GetVisitSpeed(), sc.distance)))
	c := sc.ship.Navigation.ComputeCourse(sc.destination, maxFuel*sc.fuelGauge.GetProgress()/100, sc.startTime)

	sc.fuelGauge.ChangeText("Fuel to burn: " + strconv.Itoa(maxFuel*sc.fuelGauge.GetProgress()/100))
	sc.travelTimeText.ChangeText("Travel Time: " + GetDurationString(c.TotalTime))

	speed := sc.ship.GetSpeed() + int(float64(c.AccelTime-c.StartTime)*sc.ship.Engine.Thrust)
	sc.travelSpeedText.ChangeText("Max Travel Speed: " + strconv.Itoa(speed/1000) + " km/s")
	sc.arrivalTimeText.ChangeText("Arrival Time: " + GetTimeString(c.Arrivaltime) + ", " + GetDateString(c.Arrivaltime))

	sc.course = c
}

func (sc *SetCourseDialog) Done() bool {
	if sc.goButton.IsFocused() {
		if sc.goButton.PressPulse.IsFinished() {
			burl.PushEvent(burl.NewEvent(LOG_EVENT, "Setting course for "+sc.destination.GetName()))
			return true
		}
	} else {
		if sc.cancelButton.PressPulse.IsFinished() {
			burl.PushEvent(burl.NewEvent(LOG_EVENT, "Course selection cancelled."))
			return true
		}
	}

	return false
}
