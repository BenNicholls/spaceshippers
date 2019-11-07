package main

type LifeSupportSystem struct {
	SystemStats

	targetPressure float64
	targetO2       float64
	targetTemp     float64
	targetCO2      float64

	ship *Ship
}

func NewLifeSupportSystem(s *Ship) *LifeSupportSystem {
	lss := new(LifeSupportSystem)

	lss.ship = s

	lss.targetPressure = 101
	lss.targetO2 = 21
	lss.targetTemp = 288
	lss.targetCO2 = 0

	return lss
}

func (lss *LifeSupportSystem) Update(tick int) {

	//co2 scrubbing! for each room, remove CO2 and replace with O2 from stores
	for _, r := range lss.ship.Rooms {
		if r.atmo.PartialPressure(GAS_CO2) > lss.targetCO2 {
			//co2 scub goes here.
		}
	}

	//uses some system to convert stored CO2 back into O2.
}
