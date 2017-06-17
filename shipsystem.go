package main

import "github.com/bennicholls/burl/core"
import "github.com/bennicholls/burl/util"

type PropulsionSystem struct {
	ship        *Ship
	RepairState core.Stat //0 = broken. NOTE: Do systems break, or do rooms break? Think on this.
	Thrust      int       //acceleration provided by the ship in m/s^2
	FuelUse     int       //fuel used in 1 second while on
	Firing      bool
}

func NewPropulsionSystem(s *Ship) *PropulsionSystem {
	ps := new(PropulsionSystem)
	ps.ship = s
	ps.RepairState = core.NewStat(100)
	ps.Thrust = 10
	ps.FuelUse = 2
	ps.Firing = false

	return ps
}

func (ps *PropulsionSystem) Update() {
	if ps.Firing && ps.ship.Destination != nil {
		if ps.ship.Fuel.Get()-ps.FuelUse < 0 {
			ps.Firing = false
		} else {
			//check to see if we should be coasting
			if ps.ship.Navigation.Course.R == 0 {
				return
			}
			impulse := util.Vec2Polar{R: float64(ps.Thrust), Phi: ps.ship.Navigation.Course.Phi}
			ps.ship.Velocity = ps.ship.Velocity.Add(impulse)
			ps.ship.Fuel.Mod(-ps.FuelUse)
		}
	}
}

type NavigationSystem struct {
	ship *Ship
	navRate int //how often the nevigation system adjusts the ship's course
	Course util.Vec2Polar //ship's current thrust vector. Course.R is zero, means we don't thrust at all.
}

func NewNavigationSystem(s *Ship) *NavigationSystem {
	n := new(NavigationSystem)
	n.ship = s
	n.navRate = 100

	return n
}

func (ns *NavigationSystem) Update() {
	if ns.ship.Engine.Firing {
		g := ns.ship.coords.CalcVector(ns.ship.Destination.GetCoords())
		dx, dy := g.local.ToVector().Get()
		vx, vy := ns.ship.Velocity.ToRect().Get()
		ax, ay := 0.0, 0.0

		//check if we are currently heading in the wrong direction
		if dx*vx < 0 {
			ax = -vx
		}

		if dy*vy < 0 {
			ay = -vy
		}

		//if heading in the (vaguely) right direction, steam towards object and brake if close
		//TODO: possibly split braking code into an x and y portion? eh, probably overkill.
		if dy*vy >= 0 && dx*vx >= 0 {
			ax, ay = dx, dy

			//braking code
			t := (ns.ship.Velocity.R - float64(ns.ship.Destination.GetVisitSpeed())) / float64(ns.ship.Engine.Thrust)
			decelDistance := ns.ship.Velocity.R*t - float64(ns.ship.Engine.Thrust)*t*t/2
			if g.local.Mag()-ns.ship.Destination.GetVisitDistance() < int(decelDistance) {
				ax = -ax
				ay = -ay
			} else if ns.ship.Fuel.Get()*ns.ship.Engine.Thrust < ns.ship.GetSpeed() {
				ax, ay = 0, 0
			}
		}

		ns.Course = util.Vec2{X: ax, Y: ay}.ToPolar()
	}
}