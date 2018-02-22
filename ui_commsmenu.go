package main

import "github.com/bennicholls/burl-E/burl"

type CommsMenu struct {
	burl.PagedContainer

	inboxPage *burl.Container
	inboxList *burl.List

	contactsPage *burl.Container

	transmissionsPage *burl.Container
	transmissionsList *burl.List

	logsPage *burl.Container

	comms *CommSystem
}

func NewCommsMenu(comm *CommSystem) (cm *CommsMenu) {
	cm = new(CommsMenu)

	cm.comms = comm

	cm.PagedContainer = *burl.NewPagedContainer(40, 27, 39, 4, 5, true)
	cm.SetTitle("Comm Panel")
	cm.SetVisibility(false)

	cm.inboxPage = cm.AddPage("Inbox")
	cm.inboxList = burl.NewList(38, 21, 0, 0, 2, false, "NO INBOX MESSAGES")
	cm.inboxPage.Add(cm.inboxList)
	cm.UpdateInbox()

	cm.contactsPage = cm.AddPage("Contacts")
	cm.contactsPage.Add(burl.NewTextbox(10, 1, 2, 2, 1, true, true, "contacts page"))

	cm.transmissionsPage = cm.AddPage("Transmissions")
	cm.transmissionsList = burl.NewList(37, 12, 0, 11, 0, false, "NO TRANSMISSIONS")
	cm.transmissionsPage.Add(cm.transmissionsList)

	cm.logsPage = cm.AddPage("Logs")
	cm.logsPage.Add(burl.NewTextbox(10, 1, 2, 2, 1, true, true, "logs page"))

	return
}

func (cm *CommsMenu) Update() {
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

func (cm *CommsMenu) UpdateInbox() {
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

func (cm *CommsMenu) UpdateContacts() {

}

func (cm *CommsMenu) UpdateTransmissions() {
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

func (cm *CommsMenu) UpdateLogs() {

}
