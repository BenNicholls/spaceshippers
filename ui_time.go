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

func NewTimeDisplay(g *Galaxy) (td *TimeDisplay) {
	td = new(TimeDisplay)

	td.Container = *burl.NewContainer(26, 1, 1, 43, 10, true)
	td.timeText = burl.NewTextbox(20, 1, 5, 0, 0, false, false, "")
	td.speedDisplay = burl.NewTileView(4, 1, 0, 0, 0, true)

	td.Add(td.timeText, td.speedDisplay)

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
