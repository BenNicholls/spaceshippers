package main

import "github.com/bennicholls/burl-E/burl"

type CommsMenu struct {
	burl.PagedContainer

	inboxPage         *burl.Container
	contactsPage      *burl.Container
	transmissionsPage *burl.Container
	logsPage          *burl.Container
}

func NewCommsMenu() (cm *CommsMenu) {
	cm = new(CommsMenu)

	cm.PagedContainer = *burl.NewPagedContainer(40, 27, 39, 4, 5, true)
	cm.SetTitle("Comm Panel")
	cm.SetVisibility(false)

	cm.inboxPage = burl.NewContainer(40, 23, 0, 3, 1, true)
	cm.inboxPage.Add(burl.NewTextbox(10, 1, 2, 2, 1, true, true, "inbox page"))

	cm.contactsPage = burl.NewContainer(40, 23, 0, 3, 1, true)
	cm.contactsPage.Add(burl.NewTextbox(10, 1, 2, 2, 1, true, true, "contacts page"))

	cm.transmissionsPage = burl.NewContainer(40, 23, 0, 3, 1, true)
	cm.transmissionsPage.Add(burl.NewTextbox(20, 1, 2, 2, 1, true, true, "transmissions page"))

	cm.logsPage = burl.NewContainer(40, 23, 0, 3, 1, true)
	cm.logsPage.Add(burl.NewTextbox(10, 1, 2, 2, 1, true, true, "logs page"))

	cm.AddPage("Inbox", cm.inboxPage)
	cm.AddPage("Contacts", cm.contactsPage)
	cm.AddPage("Transmissions", cm.transmissionsPage)
	cm.AddPage("Log", cm.logsPage)

	return
}
