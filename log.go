package main

import (
	"github.com/bennicholls/tyumi/event"
	"github.com/bennicholls/tyumi/gfx/ui"
)

var EV_LOG = event.Register("Space Log Message!", event.COMPLEX)

type SpaceLogEvent struct {
	event.EventPrototype

	message string
}

func fireSpaceLogEvent(message string) {
	logEvent := SpaceLogEvent{
		EventPrototype: event.New(EV_LOG),
		message:        message,
	}

	event.Fire(&logEvent)
}

func (sg *SpaceshipGame) AddLogMessage(message string) {
	sg.logOutput.InsertText(ui.JUSTIFY_LEFT, message)
	sg.logOutput.ScrollToBottom()
}
