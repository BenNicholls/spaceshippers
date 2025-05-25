package main

import (
	"github.com/bennicholls/tyumi"
	"github.com/bennicholls/tyumi/event"
	"github.com/bennicholls/tyumi/gfx/ui"
	"github.com/bennicholls/tyumi/input"
	"github.com/bennicholls/tyumi/vec"
)

type CommMenu struct {
	ui.PageContainer

	inboxPage *ui.Page
	inboxList ui.List

	contactsPage *ui.Page

	transmissionsPage *ui.Page
	transmissionsList ui.List

	logsPage *ui.Page

	comms *CommSystem
}

func (cm *CommMenu) Init(comm *CommSystem) {
	cm.comms = comm

	cm.PageContainer.Init(menuSize, menuPos, menuDepth)
	cm.EnableBorder()
	cm.Hide()
	cm.AcceptInput = true
	cm.OnPageChanged = cm.UpdateCurrentPage
	cm.Listen(EV_INBOXMESSAGERECEIVED, EV_TRANSMISSIONRECEIVED)
	cm.SuppressDuplicateEvents(event.KeepFirst)
	cm.SetEventHandler(cm.handleEvent)

	cm.inboxPage = cm.CreatePage("Messages")
	ph := cm.inboxPage.Size().H
	cm.inboxList.Init(vec.Dims{56, ph - 2}, vec.ZERO_COORD, 2)
	cm.inboxList.SetEmptyText("NO INBOX MESSAGES")
	cm.inboxList.AcceptInput = true
	cm.inboxPage.AddChild(&cm.inboxList)
	cm.UpdateInbox()

	cm.contactsPage = cm.CreatePage("Contacts")
	cm.contactsPage.AddChild(ui.NewTitleTextbox(vec.Dims{10, 1}, vec.Coord{2, 2}, 1, "contacts page"))

	cm.transmissionsPage = cm.CreatePage("Radios")
	cm.transmissionsList.Init(vec.Dims{56, 32}, vec.Coord{0, 10}, ui.BorderDepth)
	cm.transmissionsList.EnableBorder()
	cm.transmissionsList.SetEmptyText("NO TRANSMISSIONS")
	cm.transmissionsList.ToggleHighlight()
	cm.transmissionsList.AcceptInput = true
	cm.transmissionsPage.AddChild(&cm.transmissionsList)
	cm.UpdateTransmissions()

	cm.logsPage = cm.CreatePage("Logs")
	cm.logsPage.AddChild(ui.NewTitleTextbox(vec.Dims{10, 1}, vec.Coord{2, 2}, 1, "logs page"))

	return
}

func (cm *CommMenu) HandleKeypress(key_event *input.KeyboardEvent) (event_handled bool) {
	if key_event.Handled() || key_event.PressType == input.KEY_RELEASED {
		return
	}

	switch cm.GetPageIndex() {
	case 0: // Inbox
		if key_event.Key == input.K_RETURN {
			if s := cm.inboxList.GetSelectionIndex(); s != -1 {
				msg := cm.comms.Inbox[s]
				tyumi.OpenDialog(NewCommDialog(msg.sender.Name, "You", msg.sender.Pic, msg.message))
			}
		}
	case 2: // Transmissions
		if key_event.Key == input.K_RETURN {
			if s := cm.transmissionsList.GetSelectionIndex(); s != -1 {
				msg := cm.comms.Transmissions[s]
				tyumi.OpenDialog(NewCommDialog(msg.sender.Name, "You", msg.sender.Pic, msg.message))
				tyumi.OpenDialog(NewCommDialog(msg.sender.Name, "You", msg.sender.Pic, msg.message))
			}
		}
	}

	return
}

func (cm *CommMenu) handleEvent(e event.Event) (event_handled bool) {
	switch e.ID() {
	case EV_INBOXMESSAGERECEIVED:
		if cm.GetPageIndex() == 0 {
			cm.UpdateInbox()
		}
		event_handled = true
	case EV_TRANSMISSIONRECEIVED:
		if cm.GetPageIndex() == 2 {
			cm.UpdateTransmissions()
		}
		event_handled = true
	}

	return
}

func (cm *CommMenu) UpdateCurrentPage() {
	switch cm.GetPageIndex() {
	case 0: //Inbox
		cm.UpdateInbox()
	case 1: //Contacts
		cm.UpdateContacts()
	case 2: //Transmissions
		cm.UpdateTransmissions()
	case 3: //Logs
		cm.UpdateLogs()
	}
}

func (cm *CommMenu) UpdateInbox() {
	if cm.GetPageIndex() != 0 {
		return
	}

	//build message list
	cm.inboxList.RemoveAll()
	w := cm.inboxList.Size().W
	for _, m := range cm.comms.Inbox {
		var message ui.Element
		message.Init(vec.Dims{w, 3}, vec.ZERO_COORD, 0)
		message.AddChild(ui.NewTextbox(vec.Dims{w, 1}, vec.Coord{0, 0}, 0, m.title, ui.JUSTIFY_LEFT))
		message.AddChild(ui.NewTextbox(vec.Dims{w / 2, 1}, vec.Coord{0, 1}, 0, "From: "+m.sender.Name, ui.JUSTIFY_LEFT))
		message.AddChild(ui.NewTextbox(vec.Dims{w / 2, 1}, vec.Coord{w / 2, 1}, 0, "Date: "+GetDateString(m.date), ui.JUSTIFY_LEFT))
		message.AddChild(ui.NewTextbox(vec.Dims{w, 1}, vec.Coord{0, 2}, 0, m.message[:40]+"...", ui.JUSTIFY_LEFT))
		cm.inboxList.Insert(&message)
	}
}

func (cm *CommMenu) UpdateContacts() {
	if cm.GetPageIndex() != 1 {
		return
	}
}

func (cm *CommMenu) UpdateTransmissions() {
	if cm.GetPageIndex() != 2 {
		return
	}

	//build message list
	cm.transmissionsList.RemoveAll()
	w := cm.transmissionsList.Size().W
	for _, m := range cm.comms.Transmissions {
		var message ui.Element
		message.Init(vec.Dims{w, 3}, vec.ZERO_COORD, 0)
		message.AddChild(ui.NewTextbox(vec.Dims{w, 1}, vec.Coord{0, 0}, 0, m.title, ui.JUSTIFY_LEFT))
		message.AddChild(ui.NewTextbox(vec.Dims{w / 2, 1}, vec.Coord{0, 1}, 0, "From: "+m.sender.Name, ui.JUSTIFY_LEFT))
		message.AddChild(ui.NewTextbox(vec.Dims{w / 2, 1}, vec.Coord{w / 2, 1}, 0, "Date: "+GetDateString(m.date), ui.JUSTIFY_LEFT))
		message.AddChild(ui.NewTextbox(vec.Dims{w, 1}, vec.Coord{0, 2}, 0, m.message[:40]+"...", ui.JUSTIFY_LEFT))
		cm.transmissionsList.Insert(&message)
	}
}

func (cm *CommMenu) UpdateLogs() {
	if cm.GetPageIndex() != 3 {
		return
	}

}
