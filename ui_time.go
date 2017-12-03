package main

import (
	"github.com/bennicholls/burl-E/burl"
)

type TimeDisplay struct {
	burl.Container

	timeText *burl.Textbox
	dateText *burl.Textbox

	galaxy *Galaxy
}

func NewTimeDisplay(g *Galaxy) (td *TimeDisplay) {
	td = new(TimeDisplay)

	td.Container = *burl.NewContainer(10, 2, 1, 1, 10, true)
	td.timeText = burl.NewTextbox(10, 1, 0, 0, 0, false, true, "")
	td.dateText = burl.NewTextbox(10, 1, 0, 1, 0, false, true, "")

	td.Add(td.timeText, td.dateText)

	td.galaxy = g
	td.Update()

	return
}

func (td *TimeDisplay) Update() {
	td.timeText.ChangeText(GetTimeString(td.galaxy.spaceTime))
	td.dateText.ChangeText(GetDateString(td.galaxy.spaceTime))
}
