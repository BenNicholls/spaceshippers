package main

import (
	"github.com/bennicholls/burl-E/burl"
)

type GasType int

const (
	GAS_O2 GasType = iota
	GAS_CO2
	GAS_N2

	GAS_MAXTYPES
)

type gasData struct {
	name        string
	description string
}

var GasData [GAS_MAXTYPES]gasData

func init() {
	GasData[GAS_O2] = gasData{
		name:        "Oxygen",
		description: "The stuff we breath. It also explodes sometimes",
	}

	GasData[GAS_CO2] = gasData{
		name:        "Carbon Dioxide",
		description: "The waste product of respiration. Poisonous at high concentrations. Rumoured to have destroyed the natural environment of the human race many centuries ago.",
	}

	GasData[GAS_N2] = gasData{
		name:        "Nitrogen",
		description: "Inert gas which comprises most of the standard atmosphere. Doesn't explode or anything.",
	}
}

type Gas struct {
	gasType GasType
	molar   float64
}

func (g Gas) GetName() string {
	return GasData[g.gasType].name
}

func (g Gas) GetDescription() string {
	return GasData[g.gasType].description
}

//returns the molar value of the gas (L * kPa)
func (g Gas) GetAmount() float64 {
	return g.molar
}

//sets molar value of gas to v, keeping volume consistent
func (g *Gas) SetAmount(m float64) {
	g.molar = m
}

func (g *Gas) ChangeAmount(d float64) {
	if g.molar+d < 0 {
		burl.LogError("Attempt to change molar value to negative for", g.GetName(), ", no change made.")
		return
	}
	g.molar += d
}

func (g Gas) GetStorageType() StorageType {
	return STORE_GAS
}

type GasMixture struct {
	gasses map[GasType]Gas
	Volume float64 // L
	Temp   float64 // K

	totalMolar float64
}

//removes all gas from the atmosphere, leaving a vaccuum.
func (gm *GasMixture) InitVaccuum(v float64) {
	gm.Volume = v
	gm.totalMolar = 0
	gm.gasses = nil
}

//Inits atmosphere to the target parameters. v is volume in L. o2 and p are pressures in kPa
func (gm *GasMixture) InitAtmosphere(v, o2, p, t float64) {
	gm.Temp = t
	gm.Volume = v
	gm.AddGas(GAS_O2, o2*v)
	gm.AddGas(GAS_N2, (p-o2)*v)
}

//Initializes atmosphere to standard Earth sea-level values.
func (gm *GasMixture) InitStandardAtmosphere(v float64) {
	gm.InitAtmosphere(v, 21, 101, 288)
}

//AddGas adds an amount of gas g to the mixture. m is the molar value of the gas (volume * kPa).
//volume of the mixture remains constant.
func (gm *GasMixture) AddGas(g GasType, m float64) {
	if gm.gasses == nil {
		gm.gasses = make(map[GasType]Gas)
	}

	if gas, ok := gm.gasses[g]; ok {
		gas.ChangeAmount(m)
		gm.gasses[g] = gas
	} else {
		gm.gasses[g] = Gas{
			gasType: g,
			molar:   m,
		}
	}

	gm.totalMolar += m
}

//removes some gas from the mixture. if not enough of the gas is present in the mixture, removes what is there.
func (gm *GasMixture) RemoveGas(g GasType, m float64) {
	if gm.gasses == nil {
		return
	}

	if gas, ok := gm.gasses[g]; ok {
		if gas.GetAmount() <= m {
			gm.totalMolar -= gas.GetAmount()
			delete(gm.gasses, g)
		} else {
			gm.totalMolar -= m
			gas.ChangeAmount(-m)
			gm.gasses[g] = gas
		}
	}
}

//Removes a volume of gas v (L) from the atmosphere. Returns the volume of gas removed.
func (gm *GasMixture) RemoveVolume(v float64) (removed GasMixture) {
	if gm.totalMolar == 0 { //if gm is a vaccuum, return a vaccuum of size v
		removed.InitVaccuum(v)
		return
	}

	if v >= gm.Volume { //if v is larger than the mixture, return the whole mixture and leave a vaccuum.
		removed = *gm
		gm.InitVaccuum(gm.Volume)
		return
	}

	removed.Volume = v
	for gas := range gm.gasses {
		amount := gm.gasses[gas].GetAmount() * (v / gm.Volume)
		gm.RemoveGas(gas, amount)
		removed.AddGas(gas, amount)
	}
	return
}

func (gm GasMixture) Pressure() float64 {
	return gm.totalMolar / gm.Volume
}

func (gm GasMixture) PartialPressure(g GasType) float64 {
	if gm.gasses == nil {
		return 0
	}

	if _, ok := gm.gasses[g]; !ok {
		return 0
	}

	return gm.gasses[g].GetAmount() / gm.Volume
}

func (gm GasMixture) GetMolarValue(g GasType) float64 {
	if gm.gasses == nil {
		return 0
	}

	return gm.gasses[g].GetAmount()
}

func (gm *GasMixture) AddGasMixture(gm2 GasMixture) {
	for g, gas := range gm2.gasses {
		gm.AddGas(g, gas.molar)
	}
}
