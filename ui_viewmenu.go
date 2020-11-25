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
	VIEW_ATMO_CO2
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
		description: "/nPressure of internal atmosphere's Oxygen content. Oxygen is important for breathing./n/nSee Life Support System to fanagle with desired oxygen level.",
		min:         0,
		target:      22, //approximate default. overwritten on ship setup
		max:         50,
		cmin:        burl.COL_NAVY,
		ctarget:     burl.COL_GREEN,
		cmax:        burl.COL_RED,
		labels: map[float64]string{
			0:  "0 kPa: Zero oxygen. Very very bad.",
			5:  "5 kPa:  Approximate oxygen you get at extreme elevations on Earth. Lowest non-fatal amount for humans.",
			15: "15 kPa: Minimum amount required for normal healthy respiration.",
			21: "21 kPa: Oxygen pressure in a standard Earth atmosphere.",
			50: "50 kPa: Are you trying to explode the ship? This is too much.",
		},
	}
	viewModeData[VIEW_ATMO_TEMP] = ViewModeData{
		name:        "Internal Temperature",
		description: "/nTemperature, in Kelvin (K), of the internal environment./n/nSee Life Support System to dial the thermostat up or down .",
		min:         0,
		target:      290, //approximate default. overwritten on ship setup
		max:         500,
		cmin:        burl.COL_BLUE,
		ctarget:     burl.COL_GREEN,
		cmax:        burl.COL_RED,
		labels: map[float64]string{
			0:   "0K: Absolute zero. If you have this, you may have destroyed the universe.",
			273: "273K: Freezing/melting point of water.",
			288: "288K: Room temperature. Comfortable for most humans.",
			373: "373K: Boiling point of water. Uncomfortable for humans.",
			500: "500K: Boiling point of humans (probably). VERY uncomfortable for humans.",
		},
	}
	viewModeData[VIEW_ATMO_CO2] = ViewModeData{
		name:        "Carbon Dioxide Level",
		description: "/nPressure of internal atmosphere's Carbon Dioxide (CO2). CO2 is exhaled by humans, every time they breathe!. They don't like breathing it back in though. High Levels of CO2 are poisonous./n/nSee Life Support System to manage CO2 elimination.",
		min:         0,
		target:      0, //approximate default. overwritten on ship setup
		max:         10,
		cmin:        burl.COL_GREEN,
		ctarget:     burl.COL_GREEN,
		cmax:        burl.COL_RED,
		labels: map[float64]string{
			0:  "0 kPa: No CO2 content. Excellent.",
			1:  "1 kPa: Starting to get to be a little much. Humans begin getting dizzy.",
			3:  "3 kPa: Start of ill health effects. Prolonged exposure results in CO2 poisoning.",
			5:  "5 kPa: Highly accelerated CO2 poisoning timeframe.",
			7:  "7 kPa: WAY too high. CO2 poisoning in mere minutes.",
			10: "10 kPa: What in god's name are you doing letting people breathe this?!?",
		},
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

	vm.modeList = burl.NewList(18, 30, 1, 1, 1, true, "No Viewmodes Found???")
	vm.modeList.SetHint("PgUp/PgDown to scroll")

	vm.modeDescriptionText = burl.NewTextbox(18, 12, 1, 32, 1, true, false, "Viewmode Description")

	for i := 0; i < VIEWMODE_NUM; i++ {
		vm.modeList.Append(viewModeData[i].name)
	}

	vm.paletteView = burl.NewTileView(1, 40, 21, 2, 1, false)
	vm.paletteLabels = burl.NewContainer(32, 42, 23, 2, 1, false)

	vm.Add(vm.modeList, vm.modeDescriptionText, vm.paletteView, vm.paletteLabels)

	vm.UpdateViewModeData()

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
		vm.paletteLabels.Add(burl.NewTextbox(32, burl.CalcWrapHeight("<-- "+v, 32), 0, burl.Min(pos, 39), 0, false, false, "<-- "+v))
	}

	if vm.GetViewMode() != VIEW_DEFAULT {
		pos := int(40 * (vmd.target - vmd.min) / (vmd.max - vmd.min))
		vm.paletteView.Draw(0, pos, burl.GLYPH_DIAMOND, burl.COL_WHITE, burl.COL_NONE)
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
