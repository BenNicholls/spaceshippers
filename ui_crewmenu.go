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

	ship *Ship
}

func NewCrewMenu(s *Ship) (cm *CrewMenu) {
	cm = new(CrewMenu)
	cm.ship = s
	cm.PagedContainer = *burl.NewPagedContainer(56, 45, 39, 4, 10, true)
	_, ph := cm.GetPageDims()

	cm.detailPage = cm.AddPage("Crew Details")
	cm.rolePage = cm.AddPage("Roles")
	cm.jobPage = cm.AddPage("Jobs")
	cm.passengerPage = cm.AddPage("Passengers")
	cm.projectPage = cm.AddPage("Projects")

	cm.crewList = burl.NewList(18, ph-2, 1, 1, 1, true, "No crew!")
	cm.crewList.SetHint("PgUp/PgDown to select")
	cm.crewDetails = burl.NewContainer(36, ph, 20, 0, 0, true)

	cm.detailPage.Add(cm.crewList, cm.crewDetails)

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

	cm.UpdateCrewDetails()
}

func (cm *CrewMenu) UpdateCrewDetails() {
	if len(cm.ship.Crew) == 0 {
		cm.crewDetails.ClearElements()
		return
	}

	c := cm.ship.Crew[cm.crewList.GetSelection()]
	w, _ := cm.crewDetails.Dims()

	name := burl.NewTextbox(w, 1, 0, 0, 0, false, true, c.Name)
	hp := burl.NewProgressBar(w, 1, 0, 3, 0, false, true, "HP: Lots", burl.COL_RED)
	hp.SetProgress(c.HP.GetPct())
	awake := burl.NewProgressBar(w, 1, 0, 4, 0, false, true, "Awakeness: Lots", burl.COL_GREEN)
	awake.SetProgress(c.Awakeness.GetPct())
	status := burl.NewTextbox(w, 1, 0, 6, 0, false, false, c.Name+" is "+c.GetStatus())
	jobstring := c.Name + " is "
	if c.CurrentTask != nil {
		jobstring += c.CurrentTask.GetDescription()
	} else {
		jobstring += "idiling."
	}
	job := burl.NewTextbox(w, 1, 0, 7, 0, false, false, jobstring)

	cm.crewDetails.Add(name, hp, awake, status, job)
}

func (cm *CrewMenu) HandleKeypress(key sdl.Keycode) {
	cm.PagedContainer.HandleKeypress(key)

	switch cm.CurrentPage() {
	case cm.detailPage:
		cm.crewList.HandleKeypress(key)
	}
}
