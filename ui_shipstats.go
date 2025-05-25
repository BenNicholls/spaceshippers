package main

import (
	"github.com/bennicholls/tyumi/event"
	"github.com/bennicholls/tyumi/gfx/col"
	"github.com/bennicholls/tyumi/gfx/ui"
	"github.com/bennicholls/tyumi/vec"
)

type QuickStatsWindow struct {
	ui.Element

	hullBar   ui.ProgressBar
	fuelBar   ui.ProgressBar
	energyBar ui.ProgressBar
	courseBar ui.ProgressBar

	alertText ui.Textbox

	ship *Ship
}

func (qsw *QuickStatsWindow) Init(s *Ship) {
	qsw.Element.Init(vec.Dims{39, 3}, vec.Coord{39, 50}, menuDepth)
	qsw.EnableBorder()
	qsw.Listen(EV_SHIPMOVE, EV_STORAGEITEMCHANGED)
	qsw.SuppressDuplicateEvents(event.KeepFirst)
	qsw.SetEventHandler(qsw.handleEvent)

	qsw.alertText.Init(vec.Dims{39, 1}, vec.ZERO_COORD, ui.BorderDepth, "All Optimal", ui.JUSTIFY_CENTER)
	qsw.alertText.EnableBorder()
	qsw.hullBar.Init(vec.Dims{9, 1}, vec.Coord{0, 2}, ui.BorderDepth, col.RED, "HULL")
	qsw.hullBar.EnableBorder()
	qsw.fuelBar.Init(vec.Dims{9, 1}, vec.Coord{10, 2}, ui.BorderDepth, col.PURPLE, "FUEL")
	qsw.fuelBar.EnableBorder()
	qsw.energyBar.Init(vec.Dims{9, 1}, vec.Coord{20, 2}, ui.BorderDepth, col.BLUE, "ENERGY")
	qsw.energyBar.EnableBorder()
	qsw.courseBar.Init(vec.Dims{9, 1}, vec.Coord{30, 2}, ui.BorderDepth, col.GREEN, "COURSE")
	qsw.courseBar.EnableBorder()

	qsw.AddChild(&qsw.alertText)
	qsw.AddChildren(&qsw.hullBar, &qsw.fuelBar, &qsw.energyBar, &qsw.courseBar)

	qsw.ship = s
	qsw.hullBar.SetProgress(qsw.ship.Hull.GetPct())
	qsw.fuelBar.SetProgress(int(100 * qsw.ship.Storage.GetItemVolume("Fuel") / float64(qsw.ship.Storage.GetStat(STAT_LIQUID_STORAGE))))
	qsw.courseBar.SetProgress(0)

	return
}

func (qsw *QuickStatsWindow) handleEvent(e event.Event) (event_handled bool) {
	switch e.ID() {
	case EV_SHIPMOVE:
		qsw.courseBar.SetProgress(qsw.ship.Navigation.GetCurrentProgress())
	case EV_STORAGEITEMCHANGED:
		qsw.fuelBar.SetProgress(int(100 * qsw.ship.Storage.GetItemVolume("Fuel") / float64(qsw.ship.Storage.GetStat(STAT_LIQUID_STORAGE))))
	}

	return true
}
