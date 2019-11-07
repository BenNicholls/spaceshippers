package main

import "github.com/bennicholls/burl-E/burl"

type PropulsionSystem struct {
	SystemStats

	ship        *Ship
	RepairState burl.Stat //0 = broken. NOTE: Do systems break, or do rooms break? Think on this.
	Thrust      float64   //acceleration provided by the engines in m/s^2.
	FuelUse     float64   //fuel used in 1 second while on
	Firing      bool
}

func NewPropulsionSystem(s *Ship) *PropulsionSystem {
	ps := new(PropulsionSystem)

	ps.InitStats()
	ps.UpdateEngineStats()
	ps.ship = s
	ps.RepairState = burl.NewStat(100)
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
			burl.PushEvent(burl.NewEvent(LOG_EVENT, "Out of fuel! What a catastrophe!"))
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

	x, y := ps.ship.Velocity.ToRect().Get()
	if x != 0 || y != 0 {
		ps.ship.Coords.moveLocal(x, y)
		burl.PushEvent(burl.NewEvent(burl.EV_UPDATE_UI, "ship move"))
		burl.PushEvent(burl.NewEvent(burl.EV_UPDATE_UI, "ship status"))
	}
}
