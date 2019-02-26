package main

type LifeSupportSystem struct {
	SystemStats

	targetPressure float64
	targetO2       float64
	targetTemp     float64

	ship *Ship
}

func NewLifeSupportSystem(s *Ship) *LifeSupportSystem {
	lss := new(LifeSupportSystem)

	lss.ship = s

	lss.targetPressure = 101
	lss.targetO2 = .2095
	lss.targetTemp = 288

	return lss
}

func (lss *LifeSupportSystem) Update(tick int) {
	//pressurize rooms
	//for _, r := range lss.ship.Rooms {
	// if r.atmo.pressure < lss.targetPressure {
	// 	r.atmo.pressure += 0.01
	// } else if r.atmo.pressure > lss.targetPressure {
	// 	r.atmo.pressure -= 0.01
	// }
	//}
}

type Atmosphere struct {
	o2       float64 //
	co2      float64 // these are all percentages.
	n2       float64 //
	pressure float64 // kPa
	volume   float64 // L
	temp     float64 // K
}

//Earth's atmosphere at sea level. v is volume in L
func (a *Atmosphere) InitBreathable(v int) {
	a.o2 = .2095
	a.co2 = .0005
	a.n2 = 1 - a.o2 - a.co2
	a.pressure = 101
	a.temp = 288 //15 degrees celsius approx
	a.volume = float64(v)
}

//Removes a volume of air v (L) from the atmosphere. Returns the volume of gas removed.
func (a *Atmosphere) RemoveVolume(v float64) (removed Atmosphere) {
	removed = *a

	if v >= a.volume {
		removed.volume = a.volume
		a.pressure = 0
		return
	}

	removed.volume = v
	a.pressure = ((a.volume - v) / a.volume) * a.pressure

	return
}

func (a *Atmosphere) Add(a2 Atmosphere) {
	v_ratio := a2.volume / a.volume
	new_pressure := a.pressure + (a2.pressure * v_ratio)

	a.o2 = ((a.o2 * a.pressure) + (a2.o2 * a2.pressure * v_ratio)) / new_pressure
	a.co2 = ((a.co2 * a.pressure) + (a2.co2 * a2.pressure * v_ratio)) / new_pressure
	a.n2 = 1 - a.o2 - a.co2 //nitrogen is just a buffer gas, it is what the other two isn't.

	a.pressure = new_pressure
}
