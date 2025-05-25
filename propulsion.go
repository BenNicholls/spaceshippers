package main

import (
	"github.com/bennicholls/tyumi/event"
	"github.com/bennicholls/tyumi/rl"
)

var EV_SHIPMOVE = event.Register("Ship Moved", event.SIMPLE)

type PropulsionSystem struct {
	SystemStats

	ship        *Ship
	RepairState rl.Stat[int] //0 = broken. NOTE: Do systems break, or do rooms break? Think on this.
	Thrust      float64      //acceleration provided by the engines in m/s^2.
	FuelUse     float64      //fuel used in 1 second while on
	Firing      bool
}

func NewPropulsionSystem(s *Ship) *PropulsionSystem {
	ps := new(PropulsionSystem)

	ps.InitStats()
	ps.UpdateEngineStats()
	ps.ship = s
	ps.RepairState = rl.NewBasicStat(100)
	ps.Firing = false

	return ps
}

func (ps *PropulsionSystem) UpdateEngineStats() {
	ps.Thrust = float64(ps.GetStat(STAT_SUBLIGHT_THRUST))   //TODO: add calculation for ship size here!
	ps.FuelUse = float64(ps.GetStat(STAT_SUBLIGHT_FUELUSE)) //TODO: afterburner calc here??
}

func (ps *PropulsionSystem) Update(tick int) {
	ps.UpdateEngineStats()

	if ps.Firing && ps.ship.destination != nil {
		if ps.ship.Storage.GetItemVolume("Fuel")-ps.FuelUse < 0 {
			ps.Firing = false
			fireSpaceLogEvent("Out of fuel! What a catastrophe!")
		} else {
			switch ps.ship.Navigation.CurrentCourse.Phase {
			case phase_ACCEL:
				ps.ship.Velocity.R += ps.Thrust
				ps.ship.Storage.Remove(&Item{
					Name:        "Fuel",
					StorageType: STORE_LIQUID,
					Volume:      ps.FuelUse,
				})
			case phase_BRAKE:
				ps.ship.Velocity.R -= ps.Thrust
				ps.ship.Storage.Remove(&Item{
					Name:        "Fuel",
					StorageType: STORE_LIQUID,
					Volume:      ps.FuelUse,
				})
			case phase_COAST:
			}
		}
	}

	rectVelo := ps.ship.Velocity.ToRect()
	if rectVelo.NonZero() {
		ps.ship.Coords.moveLocal(rectVelo.X, rectVelo.Y)
		event.FireSimple(EV_SHIPMOVE)
	}
}
