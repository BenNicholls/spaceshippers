package main

import (
	"github.com/bennicholls/tyumi"
	"github.com/bennicholls/tyumi/event"
	"github.com/bennicholls/tyumi/gfx/ui"
	"github.com/bennicholls/tyumi/input"
	"github.com/bennicholls/tyumi/vec"
)

var spaceEvents map[int]SpaceEvent

// SpaceEvent is a happening -- an event presented to the player that
// they must deal with. There might be decisions to make, or it might just be
// a notification that something has begun or changed that they will have to
// handle.
type SpaceEvent struct {
	Title  string
	ID     int
	Unique bool //whether the event should only ever be presented to the player once.

	Pic         string //picture (currently in rexpaint .xp output format) TODO: create special "picture" resource type.
	Description string

	Choices []EventChoice
}

type EventChoice struct {
	Text   string
	Result func()
}

type SpaceEventDialog struct {
	tyumi.Scene
	done bool

	eventPicView    ui.Image
	titleText       ui.Textbox
	descriptionText ui.Textbox
	choiceList      ui.List

	choiceButtons []ui.Button

	event SpaceEvent
}

func NewSpaceEventDialog(e SpaceEvent) (sed *SpaceEventDialog) {
	sed = new(SpaceEventDialog)
	sed.InitCentered(vec.Dims{48, 37})
	sed.Window().EnableBorder()
	sed.SetKeypressHandler(sed.HandleKeypress)
	sed.Window().SendEventsToUnfocused = true

	sed.eventPicView.Init(vec.ZERO_COORD, 0, e.Pic)
	sed.titleText.Init(vec.Dims{48, 1}, vec.Coord{0, 16}, ui.BorderDepth, e.Title, ui.JUSTIFY_CENTER)
	sed.titleText.EnableBorder()
	sed.descriptionText.Init(vec.Dims{48, 8}, vec.Coord{0, 18}, ui.BorderDepth, e.Description, ui.JUSTIFY_CENTER)
	sed.descriptionText.EnableBorder()
	sed.choiceList.Init(vec.Dims{20, 10}, vec.Coord{14, 27}, ui.BorderDepth)
	sed.choiceList.EnableBorder()
	sed.choiceList.SetEmptyText("No choices, how'd this happen")
	sed.choiceList.ToggleHighlight()
	sed.choiceList.AcceptInput = true

	sed.choiceButtons = make([]ui.Button, len(e.Choices))
	for i := range e.Choices {
		sed.choiceButtons[i].Init(vec.Dims{20, 3}, vec.ZERO_COORD, 1, "/n"+e.Choices[i].Text+"/n", nil)
		sed.choiceList.Insert(&sed.choiceButtons[i])
	}

	sed.Window().AddChildren(&sed.eventPicView, &sed.titleText, &sed.descriptionText, &sed.choiceList)

	sed.event = e

	return
}

func (sed *SpaceEventDialog) HandleKeypress(key_event *input.KeyboardEvent) (event_handled bool) {
	if key_event.Key == input.K_RETURN {
		sed.choiceButtons[sed.choiceList.GetSelectionIndex()].Press()
		sed.CreateTimer(20, func() {
			sed.done = true
			sed.event.Choices[sed.choiceList.GetSelectionIndex()].Result()
		})
		event_handled = true
	}

	return
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
		Pic:         "res/art/anomaly.xp",
		Description: "While on a routine non-descript operation near Earth, your craft, " + sg.playerShip.Name + ", became entangled in a SPACETIME ANOMALY OF SOME DESCRIPTION and was hurled into the depths of space! Ship damaged, crew rattled, beverages spilled! It is now your job to somehow traverse the galaxy and find your way home!/n/nGood luck Captain! You'll need it!",
		Choices: []EventChoice{
			EventChoice{
				Text: "Holy Moly! Time for an adventure!",
				Result: func() {
					sg.player.AddMission(GenerateGoToMission(sg.playerShip, sg.galaxy.GetEarth(), nil))
					// welcomeMessage := "Hi Captain! Welcome to " + sg.playerShip.GetName() + "! I am the Ship Computer Interactive Parameter-Parsing Intelligence Entity, but you can call me SCIPPIE! "
					// sg.OpenDialog(NewCommDialog("SCIPPIE", sg.player.Name+", Captain of "+sg.playerShip.GetName(), "res/art/scippie.xp", welcomeMessage))
				},
			},
			EventChoice{
				Text: "Oh no! I'm too scared!",
				Result: func() {
					event.FireSimple(tyumi.EV_QUIT)
				},
			},
		},
	}
}
