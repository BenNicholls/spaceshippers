package main

import "github.com/bennicholls/burl/core"

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
			ps.ship.Speed += ps.Thrust
			ps.ship.Fuel.Mod(-ps.FuelUse)
		}
	}
}
