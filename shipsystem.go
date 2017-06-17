package main

import "github.com/bennicholls/burl/core"
import "github.com/bennicholls/burl/util"
import "math"

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
			g := ps.ship.coords.CalcVector(ps.ship.Destination.GetCoords())
			dx, dy := g.local.ToVector().Get()
			vx, vy := ps.ship.Heading.ToRect().Get()
			ax, ay := 0.0, 0.0

			//check if we are currently heading in the wrong direction
			if dx * vx < 0 {
				ax = -vx
			} 

			if dy * vy < 0 {
				ay = -vy
			} 

			//if heading in the (vaguely) right direction, steam towards object and brake if close
			if dy * vy >= 0 && dx * vx >= 0 {
				ax, ay = dx, dy

				//braking code
				t := (ps.ship.Heading.R - float64(ps.ship.Destination.GetVisitSpeed())) / float64(ps.ship.Engine.Thrust)
				decelDistance := ps.ship.Heading.R*t - float64(ps.ship.Engine.Thrust)*t*t/2
				if g.local.Mag()-ps.ship.Destination.GetVisitDistance() < int(decelDistance) {
					ax = -ax
					ay = -ay
				}
			}		
			impulse := util.Vec2Polar{R: float64(ps.Thrust), Phi: math.Atan2(ay, ax)}

			ps.ship.Heading = ps.ship.Heading.Add(impulse)
			//ps.ship.Fuel.Mod(-ps.FuelUse)
		}
	}
}
