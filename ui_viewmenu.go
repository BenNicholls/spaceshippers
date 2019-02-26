package main

import (
	"github.com/bennicholls/burl-E/burl"
	"github.com/veandco/go-sdl2/sdl"
)

//view modes
const (
	VIEW_DEFAULT int = iota
	VIEW_ATMO_PRESSURE
	VIEW_ATMO_O2
	VIEW_ATMO_TEMP
	VIEWMODE_NUM
)

type ViewModeData struct {
	name                string
	description         string
	palette             burl.Palette
	min, max            float64
	target              float64
	cmin, ctarget, cmax uint32
	labels              map[float64]string
}

func (vmd *ViewModeData) SetTarget(t float64) {
	vmd.target = t
	vmd.GeneratePalette()
}

func (vmd *ViewModeData) GeneratePalette() {
	if vmd.target == vmd.min || vmd.target == vmd.max {
		vmd.palette = burl.GeneratePalette(40, vmd.cmin, vmd.cmax)
	} else {
		targetnum := int(40 * (vmd.target - vmd.min) / (vmd.max - vmd.min))
		vmd.palette = burl.GeneratePalette(targetnum, vmd.cmin, vmd.ctarget)
		vmd.palette.Add(burl.GeneratePalette(41-targetnum, vmd.ctarget, vmd.cmax))
	}
}

func (vmd ViewModeData) GetColour(val float64) uint32 {
	return vmd.palette[int(float64(len(vmd.palette))*(val-vmd.min)/(vmd.max-vmd.min))]
}

var viewModeData []ViewModeData

func init() {
	viewModeData = make([]ViewModeData, VIEWMODE_NUM)
	viewModeData[VIEW_DEFAULT] = ViewModeData{
		name:        "Default View Mode",
		description: "/nNormal view mode./n/nKeeps the flashing red panic squares off your screen so you can relax.",
	}
	viewModeData[VIEW_ATMO_PRESSURE] = ViewModeData{
		name:        "Atmospheric Pressure",
		description: "/nPressure, in kPa, of the internal atmosphere. Too low, everyone dies. Too high, ship blows up. Also everyone dies./n/nSee Life Support System to fanagle with the pressure settings. If you dare.",
		min:         0,
		target:      100, //approximate default. overwritten on ship setup
		max:         500,
		cmin:        burl.COL_BLACK,
		ctarget:     burl.COL_GREEN,
		cmax:        burl.COL_RED,
		labels: map[float64]string{
			0:   "0 kPa: Total vaccuum. Very bad for your skin.",
			34:  "34 kPa: Pressure at the top of Mt. Everest.",
			101: "100 kPa: Approximate air pressue at sea level on our beloved Earth.",
			500: "500 kPa: The pressure of a very strong fist punch. Good for Mike Tyson, bad for air that you want to breathe.",
		},
	}
	viewModeData[VIEW_ATMO_O2] = ViewModeData{
		name:        "Oxygen Level",
		description: "/nPercentage of internal atmosphere occupied by Oxygen. Oxygen is important for breathing./n/nSee Life Support System to fanagle with desired oxygen level.",
		min:         0,
		target:      22, //approximate default. overwritten on ship setup
		max:         50,
		cmin:        burl.COL_NAVY,
		ctarget:     burl.COL_GREEN,
		cmax:        burl.COL_RED,
	}
	viewModeData[VIEW_ATMO_TEMP] = ViewModeData{
		name:        "Internal Temperature",
		description: "/nTemperature, in Kelvin (K), of the internal environment. Humans like a temperature around 290 K or so./n/nSee Life Support System to dial up or down the thermostat.",
		min:         0,
		target:      290, //approximate default. overwritten on ship setup
		max:         1000,
		cmin:        burl.COL_BLUE,
		ctarget:     burl.COL_GREEN,
		cmax:        burl.COL_RED,
	}
}

type ViewMenu struct {
	burl.Container

	modeList            *burl.List
	modeDescriptionText *burl.Textbox
	paletteView         *burl.TileView
	paletteLabels       *burl.Container
}

func NewViewMenu() (vm *ViewMenu) {
	vm = new(ViewMenu)
	vm.Container = *burl.NewContainer(56, 45, 39, 4, 10, true)

	vm.SetVisibility(false)

	vm.modeList = burl.NewList(25, 32, 1, 1, 1, true, "No Viewmodes Found???")
	vm.modeList.SetHint("PgUp/PgDown to scroll")

	vm.modeDescriptionText = burl.NewTextbox(25, 10, 1, 34, 1, true, false, "Viewmode Description")

	for i := 0; i < VIEWMODE_NUM; i++ {
		vm.modeList.Append(viewModeData[i].name)
	}

	vm.paletteView = burl.NewTileView(1, 40, 28, 2, 1, false)
	vm.paletteLabels = burl.NewContainer(25, 42, 30, 2, 1, false)

	vm.UpdateViewModeData()

	vm.Add(vm.modeList, vm.modeDescriptionText, vm.paletteView, vm.paletteLabels)

	return
}

func (vm *ViewMenu) UpdateViewModeData() {
	mode := vm.modeList.GetSelection()
	vm.modeDescriptionText.ChangeText(viewModeData[mode].description)

	vm.paletteView.Reset()
	vm.paletteView.DrawPalette(0, 0, viewModeData[vm.GetViewMode()].palette, burl.VERTICAL)

	vm.paletteLabels.ClearElements()
	vmd := viewModeData[vm.GetViewMode()]
	for k, v := range vmd.labels {
		pos := int(40 * (k - vmd.min) / (vmd.max - vmd.min))
		vm.paletteLabels.Add(burl.NewTextbox(25, burl.CalcWrapHeight("<-- "+v, 25), 0, burl.Min(pos, 39), 0, false, false, "<-- "+v))
	}

	if vmd.target != vmd.min && vmd.target != vmd.max {
		pos := int(40 * (vmd.target - vmd.min) / (vmd.max - vmd.min))
		vm.paletteLabels.Add(burl.NewTextbox(25, 1, 0, burl.Min(pos, 39), 0, false, false, "<-- TARGET"))
	}
}

func (vm *ViewMenu) HandleKeypress(key sdl.Keycode) {
	if key == sdl.K_PAGEUP || key == sdl.K_PAGEDOWN {
		vm.modeList.HandleKeypress(key)
		vm.UpdateViewModeData()
	}
}

func (vm ViewMenu) GetViewMode() int {
	return vm.modeList.GetSelection()
}
