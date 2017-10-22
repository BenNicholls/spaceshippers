package main

import "github.com/bennicholls/burl-E/burl"

type MissionMenu struct {
	burl.Container

	missionList      *burl.List
	descriptionText  *burl.Textbox
	statusText       *burl.Textbox
	// startTimeText    *burl.Textbox
	// deadlineTimeText *burl.Textbox
	// missionGiverText *burl.Textbox
	// stepsList        *burl.List
	// missionLogList   *burl.List

	missions *[]Mission
}

func NewMissionMenu(m *[]Mission) *MissionMenu {
	mm := new(MissionMenu)
	mm.missions = m

	mm.Container = *burl.NewContainer(40, 26, 39, 4, 5, true)
	mm.SetTitle("Missions")
	mm.SetVisibility(false)
	mm.missionList = burl.NewList(15, 26, 0, 0, 0, true, "No Missions To Do!")
	mm.descriptionText = burl.NewTextbox(24, 4, 16, 1, 0, true, true, "")
	mm.statusText = burl.NewTextbox(10, 1, 27, 12, 0, true, true, "")
	mm.Add(mm.missionList, mm.descriptionText, mm.statusText)

	mm.UpdateMissionList()

	return mm
}

func (mm *MissionMenu) UpdateMissionList() {
	mm.missionList.ClearElements()
	for _, mis := range *mm.missions {
		mm.missionList.Add(burl.NewTextbox(15, 2, 0, 0, 0, false, false, mis.name))
	}
	mm.missionList.Calibrate()
	mm.Update()
}

func (mm *MissionMenu) Update() {
	if len(*mm.missions) != 0 {
		m := (*mm.missions)[mm.missionList.GetSelection()]
		mm.descriptionText.ChangeText(m.description)
		switch m.status {
		case goal_COMPLETE:
			mm.statusText.ChangeText("Mission Successful!")
		case goal_FAILED:
			mm.statusText.ChangeText("Mission Failed!")
		case goal_INPROGRESS:
			mm.statusText.ChangeText("Mission in Progress.")
		}
	} else {
		mm.descriptionText.ChangeText("")
	}
}
