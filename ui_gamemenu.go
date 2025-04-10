package main

import (
	"github.com/bennicholls/tyumi/gfx/ui"
	"github.com/bennicholls/tyumi/vec"
)

type GameMenu struct {
	ui.PageContainer

	playerPage *ui.Page

	//MISSIONS PAGE----------
	missionsPage *ui.Page

	missionList            ui.List
	missionDescriptionText ui.Textbox
	missionStatusText      ui.Textbox
	missionGoalList        ui.List
	missionCriteriaList    ui.List
	//-----------------------

	player *Player
}

func (gm *GameMenu) Init(p *Player) {
	gm.PageContainer.Init(menuSize, menuPos, menuDepth)
	gm.EnableBorder()
	gm.Hide()
	gm.player = p
	gm.AcceptInput = true

	gm.playerPage = gm.CreatePage("Player")

	gm.missionsPage = gm.CreatePage("Missions")
	ph := gm.missionsPage.Size().H
	gm.missionList.Init(vec.Dims{16, ph - 2}, vec.Coord{1, 1}, 1)
	gm.missionList.SetupBorder("", "PgUp/PgDown")
	gm.missionList.SetEmptyText("No Missions To Do!")
	gm.missionList.ToggleHighlight()
	gm.missionList.AcceptInput = true
	gm.missionList.OnChangeSelection = gm.UpdateMissionView
	gm.missionList.SetPadding(1)

	gm.missionStatusText.Init(vec.Dims{37, 1}, vec.Coord{19, 0}, ui.BorderDepth, "Mission Status", ui.JUSTIFY_CENTER)
	gm.missionStatusText.EnableBorder()
	gm.missionDescriptionText.Init(vec.Dims{37, 4}, vec.Coord{19, 2}, ui.BorderDepth, "Description goes here!", ui.JUSTIFY_CENTER)
	gm.missionDescriptionText.EnableBorder()

	gm.missionGoalList.Init(vec.Dims{37, 21}, vec.Coord{19, 7}, ui.BorderDepth)
	gm.missionGoalList.SetupBorder("TO DO", "")
	gm.missionGoalList.SetEmptyText("Nothing to do???")

	gm.missionCriteriaList.Init(vec.Dims{37, 13}, vec.Coord{19, 29}, ui.BorderDepth)
	gm.missionCriteriaList.SetupBorder("CRITERIA", "")
	gm.missionCriteriaList.SetEmptyText("No criteria, do it however you want buddy.")

	gm.missionsPage.AddChildren(&gm.missionList, &gm.missionDescriptionText, &gm.missionStatusText, &gm.missionGoalList, &gm.missionCriteriaList)

	return
}

func (gm *GameMenu) UpdateMissions() {
	gm.missionList.RemoveAll()
	for _, mis := range gm.player.MissionLog {
		gm.missionList.InsertText(ui.JUSTIFY_CENTER, mis.name)
	}

	gm.UpdateMissionView()
}

func (gm *GameMenu) UpdateMissionView() {
	missionIndex := gm.missionList.GetSelectionIndex()
	if missionIndex == -1 { // no missions
		gm.missionDescriptionText.ChangeText("")
		gm.missionGoalList.RemoveAll()
		gm.missionCriteriaList.RemoveAll()
		return
	}

	m := gm.player.MissionLog[missionIndex]

	switch m.status {
	case goal_COMPLETE:
		gm.missionStatusText.ChangeText("Mission Successful!")
	case goal_FAILED:
		gm.missionStatusText.ChangeText("Mission Failed!")
	case goal_INPROGRESS:
		gm.missionStatusText.ChangeText("Mission in Progress.")
	}

	gm.missionDescriptionText.ChangeText(m.description)

	gm.missionGoalList.RemoveAll()
	for _, g := range m.steps {
		var step ui.Element
		step.Init(vec.Dims{27, 3}, vec.ZERO_COORD, 0)
		step.AddChild(ui.NewTextbox(vec.Dims{27, 1}, vec.Coord{0, 0}, 0, g.GetName(), ui.JUSTIFY_LEFT))
		step.AddChild(ui.NewTextbox(vec.Dims{1, 1}, vec.Coord{1, 1}, 0, "- ", ui.JUSTIFY_LEFT))
		step.AddChild(ui.NewTextbox(vec.Dims{25, 2}, vec.Coord{2, 1}, 0, g.GetDescription(), ui.JUSTIFY_LEFT))
		gm.missionGoalList.Insert(&step)
	}

	gm.missionCriteriaList.RemoveAll()
	if len(m.criteria) != 0 {
		for _, c := range m.criteria {
			var criteria ui.Element
			criteria.Init(vec.Dims{27, 3}, vec.Coord{3, 0}, 0)
			criteria.AddChild(ui.NewTextbox(vec.Dims{27, 1}, vec.Coord{0, 0}, 0, c.GetName(), ui.JUSTIFY_LEFT))
			criteria.AddChild(ui.NewTextbox(vec.Dims{26, 2}, vec.Coord{1, 1}, 0, "- "+c.GetDescription(), ui.JUSTIFY_LEFT))
			gm.missionCriteriaList.Insert(&criteria)
		}
	}
}
