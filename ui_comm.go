package main

import "github.com/bennicholls/burl-E/burl"
import "github.com/veandco/go-sdl2/sdl"

type CommDialog struct {
	burl.Container

	senderText    *burl.Textbox
	recipientText *burl.Textbox
	senderPic     *burl.TileView
	messageText   *burl.Textbox
	okayButton    *burl.Button
}

func NewCommDialog(from, to, picFile, message string) (cd *CommDialog) {
	cd = new(CommDialog)
	cd.Container = *burl.NewContainer(48, 12, 1, 1, 50, true)
	cd.CenterInConsole()

	cd.senderPic = burl.NewTileView(12, 12, 0, 0, 0, false)
	cd.senderPic.LoadImageFromCSV(picFile)

	cd.messageText = burl.NewTextbox(35, 5, 13, 3, 0, false, false, message)

	cd.senderText = burl.NewTextbox(35, 1, 13, 0, 0, false, false, "FROM: "+from)
	cd.recipientText = burl.NewTextbox(35, 1, 13, 1, 0, false, false, "TO: "+to)
	cd.okayButton = burl.NewButton(6, 1, 0, 10, 0, true, true, "Sounds Good!")
	w, _ := cd.Dims()
	cd.okayButton.CenterX(w, 12)
	cd.okayButton.ToggleFocus()

	cd.Add(cd.senderPic, cd.senderText, cd.recipientText, cd.messageText, cd.okayButton)

	return
}

func (cd *CommDialog) HandleInput(key sdl.Keycode) {
	switch key {
	case sdl.K_RETURN:
		if cd.okayButton.IsFocused() {
			cd.okayButton.Press()
		}
	}
}

func (cd CommDialog) Done() bool {
	if cd.okayButton.PressPulse.IsFinished() {
		return true
	}

	return false
}
