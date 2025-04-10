package main

import (
	"github.com/bennicholls/tyumi/gfx"
	"github.com/bennicholls/tyumi/gfx/col"
	"github.com/bennicholls/tyumi/gfx/ui"
	"github.com/bennicholls/tyumi/log"
	"github.com/bennicholls/tyumi/vec"
)

// view modes
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
	palette             col.Gradient
	min, max            float64
	target              float64
	cmin, ctarget, cmax col.Colour
	labels              map[float64]string
}

func (vmd *ViewModeData) SetTarget(t float64) {
	vmd.target = t
	vmd.GeneratePalette()
}

func (vmd *ViewModeData) GeneratePalette() {
	if vmd.target == vmd.min || vmd.target == vmd.max {
		vmd.palette = col.GenerateGradient(40, vmd.cmin, vmd.cmax)
	} else {
		targetnum := int(40 * (vmd.target - vmd.min) / (vmd.max - vmd.min))
		vmd.palette = col.GenerateGradient(targetnum, vmd.cmin, vmd.ctarget)
		vmd.palette = append(vmd.palette[:targetnum-1], col.GenerateGradient(41-targetnum, vmd.ctarget, vmd.cmax)...)
		log.Debug(len(vmd.palette))
	}
}

func (vmd ViewModeData) GetColour(val float64) col.Colour {
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
		cmin:        col.BLACK,
		ctarget:     col.GREEN,
		cmax:        col.RED,
		labels: map[float64]string{
			0:   "0 kPa:   Total vaccuum. Very bad for your skin.",
			34:  "34 kPa:  Pressure at the top of Mt. Everest.",
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
		cmin:        col.NAVY,
		ctarget:     col.GREEN,
		cmax:        col.RED,
		labels: map[float64]string{
			0:  "0 kPa:  Zero oxygen. Very very bad.",
			5:  "5 kPa:  Approximate oxygen you get at extreme elevations on Earth. Lowest non-fatal amount for humans.",
			15: "15 kPa: Minimum amount required for normal healthy respiration.",
			21: "21 kPa: Oxygen pressure in a standard Earth atmosphere.",
			50: "50 kPa: Are you trying to explode the ship? This is too much.",
		},
	}
	viewModeData[VIEW_ATMO_TEMP] = ViewModeData{
		name:        "Internal Temperature",
		description: "/nTemperature, in Kelvin (K), of the internal environment./n/nSee Life Support System to dial the thermostat up or down.",
		min:         0,
		target:      290, //approximate default. overwritten on ship setup
		max:         500,
		cmin:        col.BLUE,
		ctarget:     col.GREEN,
		cmax:        col.RED,
		labels: map[float64]string{
			0:   "0K: Absolute zero. If you have this, you may have destroyed the universe.",
			273: "273K: Freezing/melting point of water.",
			288: "288K: Room temperature. Comfortable for humans.",
			373: "373K: Boiling point of water. Uncomfortable for humans.",
			500: "500K: Boiling point of humans (probably). VERY uncomfortable for humans.",
		},
	}
	viewModeData[VIEW_ATMO_CO2] = ViewModeData{
		name:        "Carbon Dioxide Level",
		description: "/nPressure of internal atmosphere's Carbon Dioxide (CO2). CO2 is exhaled by humans every time they breathe! They don't like breathing it back in though. High Levels of CO2 are poisonous./n/nSee Life Support System to manage CO2 elimination.",
		min:         0,
		target:      0, //approximate default. overwritten on ship setup
		max:         10,
		cmin:        col.GREEN,
		ctarget:     col.GREEN,
		cmax:        col.RED,
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
	ui.Element

	modeList            ui.List
	modeDescriptionText ui.Textbox
	paletteView         ui.Element
	paletteLabels       ui.Element
}

func (vm *ViewMenu) Init() {
	vm.Element.Init(menuSize, menuPos, menuDepth)
	vm.EnableBorder()
	vm.Hide()
	vm.AcceptInput = true

	vm.modeList.Init(vec.Dims{18, 30}, vec.Coord{1, 1}, 1)
	vm.modeList.SetupBorder("", "PgUp/PgDown to scroll")
	vm.modeList.SetEmptyText("No Viewmodes Found???")
	vm.modeList.AcceptInput = true
	vm.modeList.ToggleHighlight()
	vm.modeList.OnChangeSelection = vm.UpdateViewModeData

	vm.modeDescriptionText.Init(vec.Dims{18, 12}, vec.Coord{1, 32}, 1, "ViewMode Description", ui.JUSTIFY_CENTER)
	vm.modeDescriptionText.EnableBorder()

	for i := 0; i < VIEWMODE_NUM; i++ {
		vm.modeList.InsertText(ui.JUSTIFY_LEFT, viewModeData[i].name)
	}

	vm.paletteView.Init(vec.Dims{1, 40}, vec.Coord{21, 2}, 1)
	vm.paletteLabels.Init(vec.Dims{32, 42}, vec.Coord{23, 2}, 1)

	vm.AddChildren(&vm.modeList, &vm.modeDescriptionText, &vm.paletteView, &vm.paletteLabels)

	vm.UpdateViewModeData()

	return
}

func (vm *ViewMenu) UpdateViewModeData() {
	mode := vm.GetViewMode()
	vmd := viewModeData[mode]
	vm.modeDescriptionText.ChangeText(vmd.description)

	if mode != VIEW_DEFAULT {
		for i, colour := range vmd.palette {
			vm.paletteView.DrawColours(vec.ZERO_COORD.StepN(vec.DIR_DOWN, i), 0, col.Pair{colour, colour})
		}
	} else {
		vm.paletteView.Clear()
	}

	vm.paletteLabels.RemoveAllChildren()
	for k, v := range vmd.labels {
		pos := int(40 * (k - vmd.min) / (vmd.max - vmd.min))
		vm.paletteLabels.AddChild(ui.NewTextbox(vec.Dims{32, ui.FIT_TEXT}, vec.Coord{0, min(pos, 39)}, 0, "<-- "+v, ui.JUSTIFY_LEFT))
	}

	if mode != VIEW_DEFAULT {
		pos := int(40 * (vmd.target - vmd.min) / (vmd.max - vmd.min))
		vm.paletteView.DrawVisuals(vec.Coord{0, pos}, 0, gfx.NewGlyphVisuals(gfx.GLYPH_DIAMOND, col.Pair{col.WHITE, col.NONE}))
	}
}

func (vm ViewMenu) GetViewMode() int {
	return vm.modeList.GetSelectionIndex()
}
