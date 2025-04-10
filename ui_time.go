package main

import (
	"github.com/bennicholls/tyumi/gfx"
	"github.com/bennicholls/tyumi/gfx/ui"
	"github.com/bennicholls/tyumi/vec"
)

type TimeDisplay struct {
	ui.Element

	timeText     ui.Textbox
	speedDisplay ui.Element

	galaxy *Galaxy
	speed  int
}

func (td *TimeDisplay) Init(pos vec.Coord, g *Galaxy) {
	td.Element.Init(vec.Dims{16, 3}, pos, menuDepth)
	td.TreeNode.Init(td)
	td.SetupBorder("", "+/-: change speed")
	td.speedDisplay.Init(vec.Dims{4, 1}, vec.Coord{12, 0}, ui.BorderDepth)
	td.speedDisplay.EnableBorder()
	td.timeText.Init(vec.Dims{16, 1}, vec.Coord{0, 2}, ui.BorderDepth, "", ui.JUSTIFY_CENTER)
	td.timeText.EnableBorder()

	td.AddChildren(&td.timeText, &td.speedDisplay)
	td.AddChild(ui.NewTextbox(vec.Dims{11, 1}, vec.ZERO_COORD, 0, "Simulation Speed: ", ui.JUSTIFY_CENTER))

	td.galaxy = g
	td.UpdateTime()

	return
}

func (td *TimeDisplay) UpdateTime() {
	td.timeText.ChangeText(GetTimeString(td.galaxy.spaceTime) + " " + GetDateString(td.galaxy.spaceTime))
}

func (td *TimeDisplay) UpdateSpeed(simSpeed int) {
	if td.speed == simSpeed {
		return
	}

	td.speed = simSpeed
	td.Updated = true
}

func (td *TimeDisplay) Render() {
	for i := range 4 {
		if i < td.speed {
			td.speedDisplay.DrawGlyph(vec.Coord{i, 0}, 0, gfx.GLYPH_TRIANGLE_RIGHT)
		} else {
			td.speedDisplay.DrawGlyph(vec.Coord{i, 0}, 0, gfx.GLYPH_UNDERSCORE)
		}
	}
}
