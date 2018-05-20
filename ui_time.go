package main

import (
	"github.com/bennicholls/burl-E/burl"
)

type TimeDisplay struct {
	burl.Container

	timeText     *burl.Textbox
	speedDisplay *burl.TileView

	galaxy *Galaxy
}

func NewTimeDisplay(x, y int, g *Galaxy) (td *TimeDisplay) {
	td = new(TimeDisplay)

	td.Container = *burl.NewContainer(16, 3, x, y, 10, true)
	td.Container.SetHint("+/-: change speed")
	td.speedDisplay = burl.NewTileView(4, 1, 12, 0, 0, true)
	td.timeText = burl.NewTextbox(16, 1, 0, 2, 0, true, true, "")

	td.Add(td.timeText, td.speedDisplay)
	td.Add(burl.NewTextbox(11, 1, 0, 0, 0, false, true, "Simulation Speed: "))

	td.galaxy = g
	td.UpdateTime()

	return
}

func (td *TimeDisplay) UpdateTime() {
	td.timeText.ChangeText(GetTimeString(td.galaxy.spaceTime) + " " + GetDateString(td.galaxy.spaceTime))
}

func (td *TimeDisplay) UpdateSpeed(simSpeed int) {
	td.speedDisplay.Reset()
	for i := 0; i < 4; i++ {
		if i < simSpeed {
			td.speedDisplay.Draw(i, 0, burl.GLYPH_TRIANGLE_RIGHT, burl.COL_WHITE, burl.COL_BLACK)
		} else {
			td.speedDisplay.Draw(i, 0, burl.GLYPH_UNDERSCORE, burl.COL_WHITE, burl.COL_BLACK)
		}
	}
}
