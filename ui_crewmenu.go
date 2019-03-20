package main

import (
	"github.com/bennicholls/burl-E/burl"
	"github.com/veandco/go-sdl2/sdl"
)

type CrewMenu struct {
	burl.PagedContainer

	detailPage    *burl.Container
	rolePage      *burl.Container
	jobPage       *burl.Container
	passengerPage *burl.Container
	projectPage   *burl.Container

	//detailPage
	crewList    *burl.List
	crewDetails *burl.Container

	detailsNameText     *burl.Textbox
	detailsBirthdayText *burl.Textbox
	detailsRaceText     *burl.Textbox
	detailsJobText      *burl.Textbox
	detailsPicView      *burl.TileView
	detailsStats        *burl.Container
	detailsHPBar        *burl.ProgressBar
	detailsAlertnessBar *burl.ProgressBar
	detailsCO2Bar       *burl.ProgressBar
	detailsStatusText   *burl.Textbox
	detailsEffectsText  *burl.Textbox

	ship *Ship
}

func NewCrewMenu(s *Ship) (cm *CrewMenu) {
	cm = new(CrewMenu)
	cm.ship = s
	cm.PagedContainer = *burl.NewPagedContainer(56, 45, 39, 4, 10, true)
	_, ph := cm.GetPageDims()

	//INIT DETAILS PAGE//
	cm.detailPage = cm.AddPage("Crew Details")
	cm.detailPage.Add(burl.NewTextbox(13, 1, 1, 1, 1, true, true, "Crew List"))
	cm.crewList = burl.NewList(13, ph-4, 1, 3, 1, true, "No crew!")
	cm.crewList.SetHint("PgUp/PgDown to select")
	cm.crewDetails = burl.NewContainer(41, ph, 15, 0, 0, false)
	cm.detailPage.Add(cm.crewList, cm.crewDetails)

	w, _ := cm.crewDetails.Dims()
	cm.detailsNameText = burl.NewTextbox(w-13, 1, 1, 1, 0, false, false, "Name: ")
	cm.detailsBirthdayText = burl.NewTextbox(w-13, 1, 1, 2, 0, false, false, "Age: ")
	cm.detailsRaceText = burl.NewTextbox(w-13, 1, 1, 3, 0, false, false, "Race: ")
	cm.detailsJobText = burl.NewTextbox(w-13, 1, 1, 7, 0, false, false, "Job:")
	cm.detailsPicView = burl.NewTileView(12, 12, 29, 0, 0, false)
	cm.crewDetails.Add(cm.detailsNameText, cm.detailsBirthdayText, cm.detailsRaceText, cm.detailsJobText, cm.detailsPicView)

	cm.detailsStats = burl.NewContainer(w-2, 3, 1, 13, 2, true)
	cm.detailsHPBar = burl.NewProgressBar((w-2)/2, 1, 0, 0, 0, true, true, "HP", burl.COL_RED)
	cm.detailsAlertnessBar = burl.NewProgressBar((w-2)/2, 1, (w-2)/2+1, 0, 0, true, true, "Alertness", burl.COL_GREEN)
	cm.detailsCO2Bar = burl.NewProgressBar((w-2)/2, 1, 0, 2, 0, true, true, "CO2 Buildup", burl.COL_LIGHTGREY)
	cm.detailsStats.Add(cm.detailsAlertnessBar, cm.detailsHPBar, cm.detailsCO2Bar)
	cm.crewDetails.Add(cm.detailsStats)

	cm.detailsStatusText = burl.NewTextbox(16, 23, 1, 18, 6, true, false, "Statuses and crap")
	cm.detailsEffectsText = burl.NewTextbox(w-19, 23, 18, 18, 6, true, false, "Effects and crap")

	cm.crewDetails.Add(cm.detailsStatusText, cm.detailsEffectsText)

	//TODO: OTHER PAGES//
	cm.rolePage = cm.AddPage("Roles")
	cm.jobPage = cm.AddPage("Jobs")
	cm.passengerPage = cm.AddPage("Passengers")
	cm.projectPage = cm.AddPage("Projects")

	cm.SetVisibility(false)
	cm.SetHint("TAB to switch submenus")

	cm.UpdateCrewList()

	return
}

func (cm *CrewMenu) UpdateCrewList() {
	if len(cm.ship.Crew) == 0 {
		cm.crewList.ClearElements()
		return
	}

	i := cm.crewList.GetSelection()
	cm.crewList.ClearElements()
	for _, c := range cm.ship.Crew {
		cm.crewList.Append(c.Name)
	}
	cm.crewList.Select(i)

	cm.InitCrewDetails()
}

//Updates the static data for the crew details page. This is only called when the selected crew member is changed.
func (cm *CrewMenu) InitCrewDetails() {
	if len(cm.ship.Crew) == 0 {
		return
	}

	c := cm.ship.Crew[cm.crewList.GetSelection()]
	cm.detailsNameText.ChangeText("Name: " + c.Name)
	cm.detailsRaceText.ChangeText("Race: " + c.Race)
	cm.detailsBirthdayText.ChangeText("Birthdate: " + GetDateString(c.BirthDate))
	cm.detailsPicView.LoadImageFromXP(c.Pic)

	cm.UpdateCrewDetails()
}

//Updates variable data for crew details page. Called whenever a UI_UPDATE event for the crew is received (most frames)
func (cm *CrewMenu) UpdateCrewDetails() {
	if len(cm.ship.Crew) == 0 {
		return
	}

	c := cm.ship.Crew[cm.crewList.GetSelection()]

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
		for i := 0; i < int(STATUS_MAX); i++ {
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

func (cm *CrewMenu) HandleKeypress(key sdl.Keycode) {
	cm.PagedContainer.HandleKeypress(key)

	switch cm.CurrentPage() {
	case cm.detailPage:
		if key == sdl.K_PAGEUP || key == sdl.K_PAGEDOWN {
			cm.crewList.HandleKeypress(key)
			cm.InitCrewDetails()
		}
	}
}
