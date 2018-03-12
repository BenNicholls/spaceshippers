package main

import (
	"github.com/bennicholls/burl-E/burl"
	"github.com/veandco/go-sdl2/sdl"
)

var spaceEvents map[int]SpaceEvent

//SpaceEvent is a happening -- an event presented to the player that
//they must deal with. There might be decisions to make, or it might just be
//a notification that something has begun or changed that they will have to
//handle.
type SpaceEvent struct {
	Title  string
	ID     int
	Unique bool //whether the event should only ever be presented to the player once.

	Pic         string //picture (currently in rexpaint .csv output format) TODO: create special "picture" resource type.
	Description string

	Choices []EventChoice
}

type EventChoice struct {
	Text   string
	Result func()
}

type SpaceEventDialog struct {
	burl.BaseState

	container       *burl.Container
	eventPicView    *burl.TileView
	titleText       *burl.Textbox
	descriptionText *burl.Textbox
	choiceList      *burl.List

	choiceButtons []*burl.Button

	event SpaceEvent
	done  bool
}

func NewSpaceEventDialog(e SpaceEvent) (sed *SpaceEventDialog) {
	sed = new(SpaceEventDialog)
	sed.container = burl.NewContainer(48, 37, 0, 0, 50, true)
	sed.container.CenterInConsole()

	sed.eventPicView = burl.NewTileView(48, 15, 0, 0, 0, true)
	sed.eventPicView.LoadImageFromCSV(e.Pic)
	sed.titleText = burl.NewTextbox(48, 1, 0, 16, 0, true, true, e.Title)
	sed.descriptionText = burl.NewTextbox(48, 8, 0, 18, 0, true, true, e.Description)
	sed.choiceList = burl.NewList(20, 10, 14, 27, 0, true, "no choices, how'd this happen")

	sed.choiceButtons = make([]*burl.Button, len(e.Choices))
	for i := range e.Choices {
		sed.choiceButtons[i] = burl.NewButton(20, 3, 0, 0, 1, false, true, "/n"+e.Choices[i].Text+"/n")
		sed.choiceList.Add(sed.choiceButtons[i])
	}

	sed.container.Add(sed.eventPicView, sed.titleText, sed.descriptionText, sed.choiceList)

	sed.event = e

	return
}

func (sed *SpaceEventDialog) HandleKeypress(key sdl.Keycode) {
	if key == sdl.K_RETURN {
		sed.choiceButtons[sed.choiceList.GetSelection()].HandleKeypress(key)
	} else {
		sed.choiceList.HandleKeypress(key)
	}
}

func (sed *SpaceEventDialog) HandleEvent(event *burl.Event) {
	switch event.ID {
	case burl.EV_ANIMATION_DONE:
		if event.Caller == sed.choiceButtons[sed.choiceList.GetSelection()] {
			sed.event.Choices[sed.choiceList.GetSelection()].Result()
			sed.container.ToggleVisible()
			sed.done = true
		}
	}
}

func (sed *SpaceEventDialog) Render() {
	sed.container.Render()
}

func (sed *SpaceEventDialog) Done() bool {
	return sed.done
}

func (sg *SpaceshipGame) LoadSpaceEvents() {

	spaceEvents = make(map[int]SpaceEvent)

	//Story Event 1
	spaceEvents[1] = SpaceEvent{
		Title:       "Trapped in space!",
		ID:          1,
		Unique:      true,
		Pic:         "res/art/anomaly.csv",
		Description: "While on a routine non-descript operation near Earth, your craft, " + sg.playerShip.Name + ", became entagled in a SPACETIME ANOMALY OF SOME DESCRIPTION and was hurled into the depths of the space! Ship damaged, crew rattled, beverages spilled! It is now your job to somehow traverse the galaxy and find your way home!/n/nGood luck Captain! You'll need it!",
		Choices: []EventChoice{
			EventChoice{
				Text: "Holy Moly! Time for an adventure!",
				Result: func() {
					sg.AddMission(GenerateGoToMission(sg.playerShip, sg.galaxy.GetEarth(), nil))
					welcomeMessage := "Hi Captain! Welcome to " + sg.playerShip.GetName() + "! I am the Ship Computer Interactive Parameter-Parsing Intelligence Entity, but you can call me SCIPPIE! "
					sg.dialog = NewCommDialog("SCIPPIE", sg.player.Name+", Captain of "+sg.playerShip.GetName(), "res/art/scippie.csv", welcomeMessage)
				},
			},
		},
	}
}
