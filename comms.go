package main

import (
	"math/rand"

	"github.com/bennicholls/burl-E/burl"
)

type CommSystem struct {
	Freq        int //how often it scans for transmissions. this should be user modifiable and increase power usage
	Range       int //max range for incoming transmissions
	Sensitivity int //likeliness to pickup transmission

	Inbox         []CommMessage //list of messages
	Transmissions []CommMessage //list of intercepted transmissions
}

func NewCommSystem() (cs *CommSystem) {
	cs = new(CommSystem)

	cs.Freq = HOUR
	cs.Range = 1e6
	cs.Sensitivity = 5

	cs.Inbox = make([]CommMessage, 0, 100)
	cs.Transmissions = make([]CommMessage, 0, 100)

	return
}

type CommMessage struct {
	title   string
	sender  *Person
	date    int //date of message in SpaceTime format
	message string
}

func (cs *CommSystem) Update(tick int) {
	if tick%cs.Freq == 0 {
		if rand.Intn(100) < cs.Sensitivity {
			cs.AddRandomTransmission(tick)
		}
	}
}

func (cs *CommSystem) AddRandomTransmission(tick int) {
	t := rand.Intn(100)
	trans := CommMessage{}
	trans.date = tick

	switch {
	case t < 1:
		//1% chance to get a birthday message from Mumsy
		trans.sender = &Person{"Mom", PERSON_CONTACT, DEFAULT_PIC}
		trans.title = "Happy Birthday"
		trans.message = "Glad to find you out there, don't know your frequency exactly. Anyways, Happy Birthday son. Love you!"
		if len(cs.Inbox) != cap(cs.Inbox) {
			cs.Inbox = append(cs.Inbox, trans)
		}
		burl.PushEvent(burl.NewEvent(burl.UPDATE_UI_EVENT, "inbox"))
		burl.PushEvent(burl.NewEvent(LOG_EVENT, "A new message has been received! Check your inbox."))
	case t < 10:
		//9% chance to win a radio contest
		trans.sender = &Person{"1781.2 NOVA-FM", PERSON_CONTACT, DEFAULT_PIC}
		trans.title = "You have won!"
		trans.message = "By transdimensional FM radio scanning, we've determined that you are our 10 billionth listener! That is great!"
		if len(cs.Transmissions) != cap(cs.Transmissions) {
			cs.Transmissions = append(cs.Transmissions, trans)
		}
		burl.PushEvent(burl.NewEvent(burl.UPDATE_UI_EVENT, "transmissions"))
		burl.PushEvent(burl.NewEvent(LOG_EVENT, "A transmission has been decoded."))
	default:
		trans.sender = &Person{"Unknown", PERSON_CONTACT, DEFAULT_PIC}
		trans.title = "--indecipherable--"
		trans.message = "--there is a message here, but it is too faint or corrupted to decode--"
		if len(cs.Transmissions) != cap(cs.Transmissions) {
			cs.Transmissions = append(cs.Transmissions, trans)
		}
		burl.PushEvent(burl.NewEvent(burl.UPDATE_UI_EVENT, "transmissions"))
		burl.PushEvent(burl.NewEvent(LOG_EVENT, "A garbled transmission was intercepted."))
	}
}
