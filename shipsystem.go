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
	Braking     bool
}

func NewPropulsionSystem(s *Ship) *PropulsionSystem {
	ps := new(PropulsionSystem)
	ps.ship = s
	ps.RepairState = core.NewStat(100)
	ps.Thrust = 10
	ps.FuelUse = 2
	ps.Firing = false
	ps.Braking = false

	return ps
}

func (ps *PropulsionSystem) Update() {
	if ps.Firing && ps.ship.Destination != nil {
		if ps.ship.Fuel.Get()-ps.FuelUse < 0 {
			ps.Firing = false
		} else {
			impulse := util.Vec2Polar{R: float64(ps.Thrust), Phi: ps.ship.Course.Phi}
			if ps.Braking {
				impulse.Phi += math.Pi
			}
			ps.ship.Heading = ps.ship.Heading.Add(impulse)
			ps.ship.Fuel.Mod(-ps.FuelUse)
		}
	}
}
