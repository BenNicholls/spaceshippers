package main

import "math"
import "github.com/bennicholls/burl/util"

type NavigationSystem struct {
	ship          *Ship
	galaxy        *Galaxy
	currentCourse Course //computed course for the ship to take
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
		targetVec := ns.ship.coords.CalcVector(ns.ship.Destination.GetCoords()).local.ToPolar()
		ns.ship.Velocity.Phi = targetVec.Phi

		//phase transitions
		switch ns.currentCourse.phase {
		case phase_ACCEL:
			if tick > ns.currentCourse.accelTime {
				if tick > ns.currentCourse.brakeTime {
					ns.currentCourse.phase = phase_BRAKE
				} else {
					ns.currentCourse.phase = phase_COAST
				}
			}
		case phase_COAST:
			if tick > ns.currentCourse.brakeTime {
				ns.currentCourse.phase = phase_BRAKE
			}
		}
	}

	//change location if we move away,
	if !ns.ship.coords.IsIn(ns.ship.CurrentLocation) {
		ns.ship.CurrentLocation = ns.galaxy.GetLocation(ns.ship.coords)
	}

	//change destination/location when we arrive!
	if ns.ship.coords.IsIn(ns.ship.Destination) {
		ns.ship.CurrentLocation = ns.ship.Destination
		if ns.ship.Velocity.R < ns.ship.Destination.GetVisitSpeed() {
			//if going slow enough while in range, stop the boat. TODO: put parking orbit code here.
			ns.ship.Destination = nil
			ns.ship.Velocity.R = 0
			ns.ship.Engine.Firing = false
			ns.currentCourse.done = true
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
	fuelUse   int //amount of fuel the plan uses
	totalTime int //time the course takes

	startTime   int
	startPos    util.Vec2
	accelTime   int //time to stop accelerating
	brakeTime   int //time to start braking
	arrivaltime int //time of arrival

	phase CoursePhase
	done  bool
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

	t_d = util.RoundFloatToInt(math.Sqrt((V_f*V_f+V_i*V_i)/2+T*T*c*c/4+T*D)/T - V_f/T - c/2)
	t_a = util.RoundFloatToInt(math.Sqrt((V_f*V_f+V_i*V_i)/2+T*T*c*c/4+T*D)/T - V_i/T - c/2)
	t_c = util.RoundFloatToInt(c)

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
	course.startTime = tick
	V_f := d.GetVisitSpeed()
	B := float64(fuelToUse) / float64(ns.ship.Engine.FuelUse) //burn time available
	D := ns.ship.coords.CalcVector(d.GetCoords()).local.Mag() + d.GetVisitDistance()

	t_a, t_c, t_d := ns.ComputeStraightCourse(V_f, B, D)

	course.accelTime = tick + t_a
	course.brakeTime = course.accelTime + t_c
	course.totalTime = t_a + t_c + t_d
	course.fuelUse = (t_a + t_d) * ns.ship.Engine.FuelUse

	return
}
