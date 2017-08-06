package main

import "github.com/bennicholls/burl-E/burl"

func (sg *SpaceshipGame) SetupCrewMenu() {
	sg.crewMenu = burl.NewContainer(20, 27, 59, 4, 3, true)
	sg.crewMenu.SetTitle("Crew Roster")
	sg.crewMenu.SetVisibility(false)
	sg.crewMenu.ToggleFocus()
	w, h := sg.crewMenu.Dims()
	sg.crewList = burl.NewList(w, h, 0, 0, 0, false, "")
	for _, c := range sg.playerShip.Crew {
		sg.crewList.Append(c.Name)
	}
	sg.crewDetails = burl.NewContainer(w, 3*h/4, 0, h/4+1, 0, true)
	sg.crewDetails.SetTitle("Crew Detail")
	sg.crewDetails.SetVisibility(false)
	sg.crewMenu.Add(sg.crewList, sg.crewDetails)
}

func (sg *SpaceshipGame) UpdateCrewList() {
	i := sg.crewList.GetSelection()
	sg.crewList.ClearElements()
	for _, c := range sg.playerShip.Crew {
		sg.crewList.Append(c.Name)
	}
	sg.crewList.Select(i)
}

func (sg *SpaceshipGame) UpdateCrewDetails() {
	c := sg.playerShip.Crew[sg.crewList.GetSelection()]
	w, _ := sg.crewDetails.Dims()

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

	sg.crewDetails.Add(name, hp, awake, status, job)
}

//Toggles the crew detail view
//TODO: this needs to reshape the crewlist to be constrained above the detail
//view, but we can't do that until we add the ability to reshape ui elements in burl.
func (sg *SpaceshipGame) ToggleCrewDetails() {
	if sg.crewMenu.IsVisible() {
		sg.crewDetails.ToggleVisible()
		if sg.crewDetails.IsVisible() {
			sg.UpdateCrewDetails()
		}
	}
}
