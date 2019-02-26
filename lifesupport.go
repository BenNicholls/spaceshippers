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
	lss.targetO2 = 21
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
	O2     float64 //
	CO2    float64 // kpa
	N2     float64 //
	Volume float64 // L
	Temp   float64 // K

	pressure float64 // derived value
}

//Inits atmosphere to the target parameters. v is volume in L
func (a *Atmosphere) Init(v, o2, p, t float64) {
	a.O2 = o2
	a.CO2 = 0
	a.N2 = p - a.O2
	a.Temp = t
	a.Volume = v

	a.CalcPressure()
}

func (a *Atmosphere) CalcPressure() {
	a.pressure = a.O2 + a.CO2 + a.N2
}

//removes all gas from the atmosphere, leaving a vaccuum.
func (a *Atmosphere) InitVaccuum(v float64) {
	a.O2 = 0
	a.CO2 = 0
	a.N2 = 0
	a.Volume = v
	a.pressure = 0
}

//Initializes atmosphere to standard Earth sea-level values.
func (a *Atmosphere) InitStandard(v float64) {
	a.O2 = 21
	a.CO2 = 0
	a.N2 = 80
	a.Volume = v
	a.Temp = 288
	a.pressure = 101
}

//Removes a volume of air v (L) from the atmosphere. Returns the volume of gas removed.
func (a *Atmosphere) RemoveVolume(v float64) (removed Atmosphere) {
	removed = *a

	if v >= a.Volume {
		removed.Volume = a.Volume
		a.InitVaccuum(a.Volume)
		return
	}

	removed.Volume = v

	v_ratio := 1 - (v / a.Volume)
	a.O2 *= v_ratio
	a.CO2 *= v_ratio
	a.N2 *= v_ratio
	a.CalcPressure()

	return
}

func (a *Atmosphere) Add(a2 Atmosphere) {
	a.O2 = ((a.O2 * a.Volume) + (a2.O2 * a2.Volume)) / a.Volume
	a.CO2 = ((a.CO2 * a.Volume) + (a2.CO2 * a2.Volume)) / a.Volume
	a.N2 = ((a.N2 * a.Volume) + (a2.N2 * a2.Volume)) / a.Volume

	a.CalcPressure()
}
