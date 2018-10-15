package main

import (
	"github.com/bennicholls/burl-E/burl"
	"github.com/veandco/go-sdl2/sdl"
)

type CommMenu struct {
	burl.PagedContainer

	inboxPage *burl.Container
	inboxList *burl.List

	contactsPage *burl.Container

	transmissionsPage *burl.Container
	transmissionsList *burl.List

	logsPage *burl.Container

	comms *CommSystem
}

func NewCommsMenu(comm *CommSystem) (cm *CommMenu) {
	cm = new(CommMenu)

	cm.comms = comm

	cm.PagedContainer = *burl.NewPagedContainer(56, 45, 39, 4, 10, true)
	cm.SetVisibility(false)
	cm.SetHint("TAB to switch submenus")

	_, ph := cm.GetPageDims()

	cm.inboxPage = cm.AddPage("Messages")
	cm.inboxList = burl.NewList(56, ph-2, 0, 0, 2, false, "NO INBOX MESSAGES")
	cm.inboxPage.Add(cm.inboxList)
	cm.UpdateInbox()

	cm.contactsPage = cm.AddPage("Contacts")
	cm.contactsPage.Add(burl.NewTextbox(10, 1, 2, 2, 1, true, true, "contacts page"))

	cm.transmissionsPage = cm.AddPage("Radios")
	cm.transmissionsList = burl.NewList(56, 32, 0, 10, 0, true, "NO TRANSMISSIONS")
	cm.transmissionsPage.Add(cm.transmissionsList)

	cm.logsPage = cm.AddPage("Logs")
	cm.logsPage.Add(burl.NewTextbox(10, 1, 2, 2, 1, true, true, "logs page"))

	return
}

func (sg *SpaceshipGame) HandleKeypressCommMenu(key sdl.Keycode) {
	sg.commMenu.HandleKeypress(key)

	switch sg.commMenu.CurrentIndex() {
	case 0: //Inbox
		sg.commMenu.inboxList.HandleKeypress(key)
		if key == sdl.K_RETURN && len(sg.commMenu.comms.Inbox) > 0 {
			s := sg.commMenu.inboxList.GetSelection()
			msg := sg.commMenu.comms.Inbox[s]
			burl.OpenDialog(NewCommDialog(msg.sender.Name, "You", msg.sender.Pic, msg.message))
		}
	case 2: //Transmissions
		sg.commMenu.transmissionsList.HandleKeypress(key)
		if key == sdl.K_RETURN && len(sg.commMenu.comms.Transmissions) > 0 {
			s := sg.commMenu.transmissionsList.GetSelection()
			msg := sg.commMenu.comms.Transmissions[s]
			burl.OpenDialog(NewCommDialog(msg.sender.Name, "You", msg.sender.Pic, msg.message))
		}
	}
}

func (cm *CommMenu) Update() {
	switch cm.CurrentIndex() {
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
	//build message list
	cm.inboxList.ClearElements()
	w, _ := cm.inboxList.Dims()
	for _, m := range cm.comms.Inbox {
		message := burl.NewContainer(w, 3, 0, 0, 0, false)
		message.Add(burl.NewTextbox(w, 1, 0, 0, 0, false, false, m.title))
		message.Add(burl.NewTextbox(w/2, 1, 0, 1, 0, false, false, "From: "+m.sender.Name))
		message.Add(burl.NewTextbox(w/2, 1, w/2, 1, 0, false, false, "Date: "+GetDateString(m.date)))
		message.Add(burl.NewTextbox(w, 1, 0, 2, 0, false, false, m.message[:40]+"..."))
		cm.inboxList.Add(message)
	}
}

func (cm *CommMenu) UpdateContacts() {

}

func (cm *CommMenu) UpdateTransmissions() {
	//build message list
	cm.transmissionsList.ClearElements()
	w, _ := cm.transmissionsList.Dims()
	for _, m := range cm.comms.Transmissions {
		message := burl.NewContainer(w, 3, 0, 0, 0, false)
		message.Add(burl.NewTextbox(w, 1, 0, 0, 0, false, false, m.title))
		message.Add(burl.NewTextbox(w/2, 1, 0, 1, 0, false, false, "From: "+m.sender.Name))
		message.Add(burl.NewTextbox(w/2, 1, w/2, 1, 0, false, false, "Date: "+GetDateString(m.date)))
		message.Add(burl.NewTextbox(w, 1, 0, 2, 0, false, false, m.message[:40]+"..."))
		cm.transmissionsList.Add(message)
	}
}

func (cm *CommMenu) UpdateLogs() {

}
