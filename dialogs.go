package main

import (
	"github.com/bennicholls/burl-E/burl"
	"github.com/veandco/go-sdl2/sdl"
)

var EV_LOADFILE = burl.RegisterCustomEvent() //event.Message will contain local pathname
var EV_SAVEFILE = burl.RegisterCustomEvent() //event.Message will contain local pathname

type ChooseFileDialog struct {
	burl.StatePrototype

	fileList     *burl.List
	okayButton   *burl.Button
	cancelButton *burl.Button

	filenames []string

	dirPath string
}

func NewChooseFileDialog(dirPath, ext string) (cfd *ChooseFileDialog) {
	cfd = new(ChooseFileDialog)
	cfd.dirPath = dirPath

	cfd.Window = burl.NewContainer(20, 29, 0, 0, 50, true)
	cfd.Window.CenterInConsole()
	cfd.Window.ToggleFocus()
	cfd.Window.SetTitle("Select file!")

	cfd.fileList = burl.NewList(20, 25, 0, 0, 0, true, "No Files Found!/n/n(Press C or ESCAPE to cancel)")
	cfd.fileList.ToggleFocus()
	cfd.okayButton = burl.NewButton(8, 1, 1, 27, 1, true, true, "[L]oad File")
	cfd.cancelButton = burl.NewButton(8, 1, 11, 27, 2, true, true, "[C]ancel")

	cfd.Window.Add(cfd.fileList, cfd.okayButton, cfd.cancelButton)

	var err error
	cfd.filenames, err = burl.GetFileList(dirPath, ext)
	if err != nil {
		cfd.fileList.ChangeEmptyText("Could not load files!/n/n(See log.txt for details, Press C or ESCAPE to cancel)")
	}

	for _, name := range cfd.filenames {
		cfd.fileList.Append(name)
	}

	return
}

func (cfd *ChooseFileDialog) HandleKeypress(key sdl.Keycode) {
	cfd.fileList.HandleKeypress(key)
	switch key {
	case sdl.K_RETURN, sdl.K_l:
		if len(cfd.fileList.Elements) != 0 {
			cfd.okayButton.Press()
		}
	case sdl.K_ESCAPE, sdl.K_c:
		cfd.cancelButton.Press()
	}
}

func (cfd *ChooseFileDialog) Done() bool {
	if cfd.okayButton.PressPulse.IsFinished() {
		burl.PushEvent(burl.NewEvent(EV_LOADFILE, cfd.filenames[cfd.fileList.GetSelection()]))
		return true
	} else if cfd.cancelButton.PressPulse.IsFinished() {
		return true
	}

	return false
}

type SaveDialog struct {
	burl.StatePrototype

	nameInput    *burl.Inputbox
	fileText     *burl.Textbox
	saveButton   *burl.Button
	cancelButton *burl.Button

	ext       string
	dirPath   string
	filenames []string //current contents of directory so we can warn against overwrites
}

func NewSaveDialog(dirPath, ext, def string) (sd *SaveDialog) {
	sd = new(SaveDialog)
	sd.dirPath = dirPath
	sd.ext = ext

	sd.Window = burl.NewContainer(26, 10, 0, 0, 50, true)
	sd.Window.CenterInConsole()
	sd.Window.SetTitle("Choose Save Name")
	sd.Window.ToggleFocus()

	sd.Window.Add(burl.NewTextbox(5, 1, 1, 2, 0, false, false, "Filename:"))
	sd.nameInput = burl.NewInputbox(17, 1, 7, 2, 0, true)
	sd.nameInput.ChangeText(def)
	sd.nameInput.ToggleFocus()

	sd.fileText = burl.NewTextbox(24, 3, 1, 4, 0, false, true, "Input filename.")

	sd.saveButton = burl.NewButton(5, 1, 7, 8, 2, true, true, "Save")
	sd.saveButton.ToggleFocus()
	sd.cancelButton = burl.NewButton(5, 1, 14, 8, 1, true, true, "Cancel")
	sd.cancelButton.ToggleFocus()

	sd.Window.Add(sd.nameInput, sd.fileText, sd.saveButton, sd.cancelButton)

	sd.filenames, _ = burl.GetFileList(dirPath, "")
	sd.UpdateFileText()

	return
}

func (sd *SaveDialog) HandleKeypress(key sdl.Keycode) {
	switch key {
	case sdl.K_RETURN:
		if sd.nameInput.GetText() != "" {
			sd.saveButton.Press()
		}
	case sdl.K_ESCAPE:
		sd.cancelButton.Press()
	default:
		sd.nameInput.HandleKeypress(key)
		sd.UpdateFileText()
	}
}

func (sd *SaveDialog) UpdateFileText() {
	if sd.nameInput.GetText() == "" {
		sd.fileText.ChangeText("Please input filename.")
	} else {
		name := sd.nameInput.GetText() + sd.ext
		sd.fileText.ChangeText("Will save the file as " + sd.dirPath + name)
		for _, filename := range sd.filenames {
			if name == filename {
				sd.fileText.AppendText("/nThis will OVERWRITE the current file!")
				break
			}
		}
	}
}

func (sd *SaveDialog) Done() bool {
	if sd.saveButton.PressPulse.IsFinished() {
		burl.PushEvent(burl.NewEvent(EV_SAVEFILE, sd.nameInput.GetText()+sd.ext))
		return true
	} else if sd.cancelButton.PressPulse.IsFinished() {
		return true
	}

	return false
}

type CommDialog struct {
	burl.StatePrototype

	senderText    *burl.Textbox
	recipientText *burl.Textbox
	senderPic     *burl.TileView
	messageText   *burl.Textbox
	okayButton    *burl.Button
}

func NewCommDialog(from, to, picFile, message string) (cd *CommDialog) {
	cd = new(CommDialog)
	cd.Window = burl.NewContainer(48, 12, 1, 1, 50, true)
	cd.Window.CenterInConsole()

	cd.okayButton = burl.NewButton(6, 1, 0, 10, 1, true, true, "Sounds Good!")
	cd.okayButton.ToggleFocus()
	w, _ := cd.Window.Dims()

	if from == "" && to == "" && picFile == "" {
		//special dialog version with just a message.
		cd.messageText = burl.NewTextbox(48, 5, 0, 1, 0, false, true, message)
		cd.okayButton.CenterX(w, 0)
	} else {
		cd.senderPic = burl.NewTileView(12, 12, 0, 0, 0, false)
		cd.senderPic.LoadImageFromXP(picFile)
		cd.messageText = burl.NewTextbox(35, 5, 13, 3, 0, false, false, message)
		cd.senderText = burl.NewTextbox(35, 1, 13, 0, 0, false, false, "FROM: "+from)
		cd.recipientText = burl.NewTextbox(35, 1, 13, 1, 0, false, false, "TO: "+to)
		cd.Window.Add(cd.senderPic, cd.senderText, cd.recipientText)
		cd.okayButton.CenterX(w, 12)
	}

	cd.Window.Add(cd.messageText, cd.okayButton)

	return
}

func (cd *CommDialog) HandleKeypress(key sdl.Keycode) {
	cd.okayButton.HandleKeypress(key)
}

func (cd CommDialog) Done() bool {
	return cd.okayButton.PressPulse.IsFinished()
}
