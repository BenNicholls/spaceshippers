package main

import "github.com/bennicholls/burl-E/burl"

type MissionMenu struct {
	burl.Container

	missionList      *burl.List
	startTimeText    *burl.Textbox
	deadlineTimeText *burl.Textbox
	missionGiverText *burl.Textbox
	stepsList        *burl.List
	missionLogList   *burl.List

	missions []*Mission
}

func NewMissionMenu() *MissionMenu {
	mm := new(MissionMenu)

	mm.Container = *burl.NewContainer(40, 26, 39, 4, 5, true)
	mm.SetTitle("Missions")
	mm.SetVisibility(false)

	return mm
}
