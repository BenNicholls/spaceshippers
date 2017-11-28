package main

import "github.com/bennicholls/burl-E/burl"

type CrewMenu struct {
	burl.Container

	crewList    *burl.List
	crewDetails *burl.Container

	ship *Ship
}

func NewCrewMenu(s *Ship) (cm *CrewMenu) {
	cm = new(CrewMenu)
	cm.ship = s

	cm.Container = *burl.NewContainer(20, 27, 59, 4, 3, true)
	cm.SetTitle("Crew Roster")
	cm.SetVisibility(false)
	cm.ToggleFocus()
	w, h := cm.Dims()
	cm.crewList = burl.NewList(w, h, 0, 0, 0, false, "")
	for _, c := range cm.ship.Crew {
		cm.crewList.Append(c.Name)
	}
	cm.crewDetails = burl.NewContainer(w, 3*h/4, 0, h/4+1, 0, true)
	cm.crewDetails.SetTitle("Crew Detail")
	cm.crewDetails.SetVisibility(false)
	cm.Add(cm.crewList, cm.crewDetails)

	return
}

func (cm *CrewMenu) UpdateCrewList() {
	i := cm.crewList.GetSelection()
	cm.crewList.ClearElements()
	for _, c := range cm.ship.Crew {
		cm.crewList.Append(c.Name)
	}
	cm.crewList.Select(i)
}

func (cm *CrewMenu) UpdateCrewDetails() {
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

//Toggles the crew detail view
//TODO: this needs to reshape the crewlist to be constrained above the detail
//view, but we can't do that until we add the ability to reshape ui elements in burl.
func (cm *CrewMenu) ToggleCrewDetails() {
	if cm.IsVisible() {
		cm.crewDetails.ToggleVisible()
		if cm.crewDetails.IsVisible() {
			cm.UpdateCrewDetails()
		}
	}
}
