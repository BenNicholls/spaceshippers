package main

import "github.com/bennicholls/burl/ui"

type MissionMenu struct {
	ui.Container

	missionList *ui.List
	startTimeText *ui.Textbox
	deadlineTimeText *ui.Textbox
	missionGiverText *ui.Textbox
	stepsList *ui.List
	missionLogList *ui.List

	missions []*Mission
}

func NewMissionMenu() *MissionMenu {
	mm := new(MissionMenu)

	mm.Container = *ui.NewContainer(40, 26, 39, 4, 5, true)
	mm.SetTitle("Missions")
	mm.SetVisibility(false)
	
	return mm
}