package main

import (
	"github.com/bennicholls/tyumi/gfx/col"
	"github.com/bennicholls/tyumi/gfx/ui"
	"github.com/bennicholls/tyumi/vec"
)

type CrewMenu struct {
	ui.PageContainer

	detailPage    *ui.Page
	rolePage      *ui.Page
	jobPage       *ui.Page
	passengerPage *ui.Page
	projectPage   *ui.Page

	//detailPage
	crewList    ui.List
	crewDetails ui.Element

	detailsNameText     ui.Textbox
	detailsBirthdayText ui.Textbox
	detailsRaceText     ui.Textbox
	detailsJobText      ui.Textbox
	detailsPicView      ui.Image
	detailsStats        ui.Element
	detailsHPBar        ui.ProgressBar
	detailsAlertnessBar ui.ProgressBar
	detailsCO2Bar       ui.ProgressBar
	detailsStatusText   ui.Textbox
	detailsEffectsText  ui.Textbox

	ship *Ship
}

func (cm *CrewMenu) Init(s *Ship) {
	cm.PageContainer.Init(menuSize, menuPos, menuDepth)
	cm.EnableBorder()
	cm.AcceptInput = true
	cm.Hide()

	//INIT DETAILS PAGE//
	cm.detailPage = cm.CreatePage("Crew Details")
	ph := cm.detailPage.Size().H
	cm.detailPage.AddChild(ui.NewTitleTextbox(vec.Dims{13, 1}, vec.Coord{1, 1}, 1, "Crew List"))
	cm.crewList.Init(vec.Dims{13, ph - 4}, vec.Coord{1, 3}, 1)
	cm.crewList.SetupBorder("", "PgUp/PgDown to select")
	cm.crewList.SetEmptyText("No Crew!")
	cm.crewList.OnChangeSelection = cm.InitCrewDetails

	cm.crewDetails.Init(vec.Dims{41, ph}, vec.Coord{15, 0}, 0)
	cm.detailPage.AddChildren(&cm.crewList, &cm.crewDetails)

	w := cm.crewDetails.Size().W
	cm.detailsNameText.Init(vec.Dims{w - 13, 1}, vec.Coord{1, 1}, 0, "Name: ", ui.JUSTIFY_LEFT)
	cm.detailsBirthdayText.Init(vec.Dims{w - 13, 1}, vec.Coord{1, 2}, 0, "Age: ", ui.JUSTIFY_LEFT)
	cm.detailsRaceText.Init(vec.Dims{w - 13, 1}, vec.Coord{1, 3}, 0, "Race: ", ui.JUSTIFY_LEFT)
	cm.detailsJobText.Init(vec.Dims{w - 13, 1}, vec.Coord{1, 7}, 0, "Job: ", ui.JUSTIFY_LEFT)
	cm.detailsPicView.Init(vec.Coord{29, 0}, 0, "")
	cm.crewDetails.AddChildren(&cm.detailsNameText, &cm.detailsBirthdayText, &cm.detailsRaceText, &cm.detailsJobText, &cm.detailsPicView)

	cm.detailsStats.Init(vec.Dims{w - 2, 3}, vec.Coord{1, 13}, 2)
	cm.detailsStats.EnableBorder()
	cm.detailsHPBar.Init(vec.Dims{(w - 2) / 2, 1}, vec.ZERO_COORD, ui.BorderDepth, col.RED, "HP")
	cm.detailsHPBar.EnableBorder()
	cm.detailsAlertnessBar.Init(vec.Dims{(w - 2) / 2, 1}, vec.Coord{(w-2)/2 + 1, 0}, ui.BorderDepth, col.GREEN, "Alertness")
	cm.detailsHPBar.EnableBorder()
	cm.detailsCO2Bar.Init(vec.Dims{(w - 2) / 2, 1}, vec.Coord{0, 2}, ui.BorderDepth, col.LIGHTGREY, "CO2 Buildup")
	cm.detailsHPBar.EnableBorder()
	cm.detailsStats.AddChildren(&cm.detailsAlertnessBar, &cm.detailsHPBar, &cm.detailsCO2Bar)
	cm.crewDetails.AddChild(&cm.detailsStats)

	cm.detailsStatusText.Init(vec.Dims{16, 23}, vec.Coord{1, 18}, 6, "Statuses and crap", ui.JUSTIFY_CENTER)
	cm.detailsEffectsText.Init(vec.Dims{w - 19, 23}, vec.Coord{18, 18}, 6, "Effects and crap", ui.JUSTIFY_CENTER)

	cm.crewDetails.AddChildren(&cm.detailsStatusText, &cm.detailsEffectsText)

	//TODO: OTHER PAGES//
	cm.rolePage = cm.CreatePage("Roles")
	cm.jobPage = cm.CreatePage("Jobs")
	cm.passengerPage = cm.CreatePage("Passengers")
	cm.projectPage = cm.CreatePage("Projects")

	cm.ship = s
	cm.UpdateCrewList()

	return
}

func (cm *CrewMenu) UpdateCrewList() {
	i := cm.crewList.GetSelectionIndex()
	cm.crewList.RemoveAll()
	if len(cm.ship.Crew) == 0 {
		return
	}

	for _, c := range cm.ship.Crew {
		cm.crewList.InsertText(ui.JUSTIFY_LEFT, c.Name)
	}
	cm.crewList.Select(i)

	cm.InitCrewDetails()
}

// Updates the static data for the crew details page. This is only called when the selected crew member is changed.
func (cm *CrewMenu) InitCrewDetails() {
	if len(cm.ship.Crew) == 0 {
		return
	}

	c := cm.ship.Crew[cm.crewList.GetSelectionIndex()]
	cm.detailsNameText.ChangeText("Name: " + c.Name)
	cm.detailsRaceText.ChangeText("Race: " + c.Race)
	cm.detailsBirthdayText.ChangeText("Birthdate: " + GetDateString(c.BirthDate))
	cm.detailsPicView.LoadImage(c.Pic)

	cm.UpdateCrewDetails()
}

// Updates variable data for crew details page. Called whenever a UI_UPDATE event for the crew is received (most frames)
func (cm *CrewMenu) UpdateCrewDetails() {
	if len(cm.ship.Crew) == 0 {
		return
	}

	c := cm.ship.Crew[cm.crewList.GetSelectionIndex()]

	cm.detailsHPBar.SetProgress(c.HP.GetPct())
	cm.detailsAlertnessBar.SetProgress(c.Awakeness.GetPct())
	cm.detailsCO2Bar.SetProgress(c.CO2.GetPct())
	jobstring := c.Name + " is "
	if c.CurrentTask != nil {
		jobstring += c.CurrentTask.GetDescription()
	} else if c.Dead {
		jobstring += "dead. :("
	} else {
		jobstring += "idling."
	}
	cm.detailsJobText.ChangeText(jobstring)

	//TODO: this could be MUCH better probably. Currently kind of hacked together.
	statusString := "STATUSES:/n/n"
	effectString := "EFFECTS:/n/n"
	if len(c.Statuses) == 0 {
		statusString += " - " + "No statuses. Boooring!"
	} else {
		for i := range int(STATUS_MAX) {
			if s, ok := c.Statuses[StatusID(i)]; ok {
				statusString += " - " + s.Name + "/n"
				effectString += " -> "
				for _, e := range s.Effects {
					effectString += c.Effects[e].Abbreviation + " "
				}
				effectString += "/n"
			}
		}
	}
	cm.detailsStatusText.ChangeText(statusString)
	cm.detailsEffectsText.ChangeText(effectString)
}
