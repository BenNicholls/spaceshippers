package main

import "github.com/bennicholls/burl-E/burl"

type MissionMenu struct {
	burl.Container

	missionList     *burl.List
	descriptionText *burl.Textbox
	statusText      *burl.Textbox
	goalList        *burl.List
	criteriaList    *burl.List

	missions *[]Mission
}

func NewMissionMenu(m *[]Mission) *MissionMenu {
	mm := new(MissionMenu)
	mm.missions = m

	mm.Container = *burl.NewContainer(40, 26, 39, 4, 15, true)
	mm.SetTitle("Missions")
	mm.SetVisibility(false)
	mm.missionList = burl.NewList(10, 26, 0, 0, 0, false, "No Missions To Do!")
	mm.statusText = burl.NewTextbox(28, 1, 11, 1, 1, true, true, "")
	mm.descriptionText = burl.NewTextbox(28, 4, 11, 3, 1, true, true, "")

	mm.goalList = burl.NewList(28, 8, 11, 8, 1, true, "Nothing to do???")
	mm.goalList.SetTitle("TO DO")
	mm.goalList.ToggleHighlight()

	mm.criteriaList = burl.NewList(28, 8, 11, 17, 1, true, "Nothing to do???")
	mm.criteriaList.SetTitle("CRITERIA")
	mm.criteriaList.ToggleHighlight()

	mm.Add(mm.missionList, mm.descriptionText, mm.statusText, mm.goalList, mm.criteriaList)

	mm.Update()

	return mm
}

func (mm *MissionMenu) Update() {
	//update mission list
	mm.missionList.ClearElements()
	for _, mis := range *mm.missions {
		mm.missionList.Add(burl.NewTextbox(10, 2, 0, 0, 0, false, false, mis.name))
	}
	mm.missionList.Calibrate()

	if len(*mm.missions) != 0 {
		m := (*mm.missions)[mm.missionList.GetSelection()]

		switch m.status {
		case goal_COMPLETE:
			mm.statusText.ChangeText("Mission Successful!")
		case goal_FAILED:
			mm.statusText.ChangeText("Mission Failed!")
		case goal_INPROGRESS:
			mm.statusText.ChangeText("Mission in Progress.")
		}

		mm.descriptionText.ChangeText(m.description)

		//TODO: should not be rebuilding these lists all the damn time
		mm.goalList.ClearElements()
		for _, g := range m.steps {
			step := burl.NewContainer(29, 3, 0, 0, 0, false)
			step.Add(burl.NewTextbox(29, 1, 0, 0, 0, false, false, g.GetName()))
			step.Add(burl.NewTextbox(28, 2, 1, 1, 0, false, false, "- "+g.GetDescription()))
			mm.goalList.Add(step)
		}

		mm.criteriaList.ClearElements()
		for _, c := range m.criteria {
			criteria := burl.NewContainer(29, 3, 0, 0, 0, false)
			criteria.Add(burl.NewTextbox(29, 1, 0, 0, 0, false, false, c.GetName()))
			criteria.Add(burl.NewTextbox(28, 2, 1, 1, 0, false, false, "- "+c.GetDescription()))
			mm.criteriaList.Add(criteria)
		}

	} else {
		mm.descriptionText.ChangeText("")
		mm.goalList.ClearElements()
		mm.criteriaList.ClearElements()
	}
}
