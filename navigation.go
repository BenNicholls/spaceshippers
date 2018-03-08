package main

import "math"
import "github.com/bennicholls/burl-E/burl"

type NavigationSystem struct {
	ship          *Ship
	galaxy        *Galaxy
	CurrentCourse Course //computed course for the ship to take
}

func NewNavigationSystem(s *Ship, g *Galaxy) *NavigationSystem {
	n := new(NavigationSystem)
	n.ship = s
	n.galaxy = g

	return n
}

func (ns *NavigationSystem) Update(tick int) {
	if ns.ship.Engine.Firing {
		//keep ship on course magically. Kepler, Newton, all the physicists, I am SO SORRY.
		targetVec := ns.ship.Coords.CalcVector(ns.ship.destination.GetCoords()).Local.ToPolar()
		ns.ship.Velocity.Phi = targetVec.Phi

		//phase transitions
		switch ns.CurrentCourse.Phase {
		case phase_ACCEL:
			if tick > ns.CurrentCourse.AccelTime {
				if tick > ns.CurrentCourse.BrakeTime {
					ns.CurrentCourse.Phase = phase_BRAKE
				} else {
					ns.CurrentCourse.Phase = phase_COAST
				}
			}
		case phase_COAST:
			if tick > ns.CurrentCourse.BrakeTime {
				ns.CurrentCourse.Phase = phase_BRAKE
			}
		}

		burl.PushEvent(burl.NewEvent(burl.EV_UPDATE_UI, "ship status"))
	}

	//If we're moving, we need to check our locations/destinations.
	if ns.ship.GetSpeed() > 0 {
		//change location if we move away,
		if !ns.ship.Coords.IsIn(ns.ship.currentLocation) {
			ns.ship.currentLocation = ns.galaxy.GetLocation(ns.ship.Coords)
			burl.PushEvent(burl.NewEvent(burl.EV_UPDATE_UI, "ship status"))
		}

		//change destination/location when we arrive!
		if ns.ship.Coords.IsIn(ns.ship.destination) {
			ns.ship.currentLocation = ns.ship.destination
			burl.PushEvent(burl.NewEvent(burl.EV_UPDATE_UI, "ship status"))
			if ns.ship.Velocity.R < ns.ship.destination.GetVisitSpeed() {
				//if going slow enough while in range, stop the boat. TODO: put parking orbit code here.
				ns.ship.destination = nil
				ns.ship.Velocity.R = 0
				ns.ship.Engine.Firing = false
				ns.CurrentCourse.Done = true

				burl.PushEvent(burl.NewEvent(LOG_EVENT, "We have arrived at "+ns.ship.currentLocation.GetName()))
			}
		}
	}
}

type CoursePhase int

const (
	phase_ACCEL CoursePhase = iota
	phase_COAST
	phase_BRAKE
)

//Plan of action for the nav system
type Course struct {
	//precomputed factors, we'll see if they turn out to be right
	FuelUse   int //amount of fuel the plan uses
	TotalTime int //time the course takes

	StartTime   int
	StartPos    burl.Vec2
	AccelTime   int //time to stop accelerating
	BrakeTime   int //time to start braking
	Arrivaltime int //time of arrival

	Phase CoursePhase
	Done  bool
}

//Computes the course parameters subject to a fuel limit
//These formulae are the product of way too much algebra don't screw with them.
func (ns NavigationSystem) ComputeStraightCourse(V_f, B, D float64) (t_a, t_c, t_d int) {
	V_i := ns.ship.Velocity.R
	T := ns.ship.Engine.Thrust

	c := 0.0
	//limit fuel use to the maximum possible
	max_burnTime := ns.CalcMaxBurnTime(V_f, D)
	if B < max_burnTime {
		//add some coast time in the middle of the course if we're not going all-out
		c = (2*D + (V_f*V_f+V_i*V_i)/(2*T) - V_f*V_i/T - B*(V_f+V_i) - T*B*B/2) / (V_f + V_i + T*B)
	}

	t_d = burl.RoundFloatToInt(math.Sqrt((V_f*V_f+V_i*V_i)/2+T*T*c*c/4+T*D)/T - V_f/T - c/2)
	t_a = burl.RoundFloatToInt(math.Sqrt((V_f*V_f+V_i*V_i)/2+T*T*c*c/4+T*D)/T - V_i/T - c/2)
	t_c = burl.RoundFloatToInt(c)

	return
}

func (ns NavigationSystem) CalcMaxBurnTime(V_f, D float64) float64 {
	V_i := ns.ship.Velocity.R
	T := ns.ship.Engine.Thrust
	return 2*math.Sqrt((V_f*V_f+V_i*V_i)/2+T*D)/T - (V_f+V_i)/T
}

//Calculates the required burntime required to brake the ship down to V_f. If V_f > ship speed, returns 0.
func (ns NavigationSystem) CalcMinBurnTime(V_f float64) float64 {
	if V_f > ns.ship.Velocity.R {
		return 0
	} else {
		return (ns.ship.Velocity.R - V_f) / ns.ship.Engine.Thrust * float64(ns.ship.Engine.FuelUse)
	}
}

//Returns a Course based on the parameters.
func (ns NavigationSystem) ComputeCourse(d Locatable, fuelToUse, tick int) (course Course) {
	course.StartTime = tick
	V_f := d.GetVisitSpeed()
	B := float64(fuelToUse) / float64(ns.ship.Engine.FuelUse) //burn time available
	D := ns.ship.Coords.CalcVector(d.GetCoords()).Local.Mag() + d.GetVisitDistance()

	t_a, t_c, t_d := ns.ComputeStraightCourse(V_f, B, D)

	course.AccelTime = tick + t_a
	course.BrakeTime = course.AccelTime + t_c
	course.TotalTime = t_a + t_c + t_d
	course.FuelUse = (t_a + t_d) * ns.ship.Engine.FuelUse
	course.Arrivaltime = tick + course.TotalTime

	return
}
