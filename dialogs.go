package main

import (
	"slices"

	"github.com/bennicholls/tyumi"
	"github.com/bennicholls/tyumi/event"
	"github.com/bennicholls/tyumi/gfx"
	"github.com/bennicholls/tyumi/gfx/ui"
	"github.com/bennicholls/tyumi/input"
	"github.com/bennicholls/tyumi/util"
	"github.com/bennicholls/tyumi/vec"
)

type ChooseFileDialog struct {
	tyumi.State
	done bool

	onFileLoaded func(filename string) // function run when dialog exits if a file was successfully chosen

	fileList     ui.List
	okayButton   ui.Button
	cancelButton ui.Button

	filenames  []string
	dirPath    string
	chosenFile string
}

func NewChooseFileDialog(dirPath, ext string, onLoad func(filename string)) (cfd *ChooseFileDialog) {
	cfd = new(ChooseFileDialog)
	cfd.dirPath = dirPath
	cfd.onFileLoaded = onLoad
	cfd.SetKeypressHandler(cfd.HandleKeypress)

	cfd.InitCentered(vec.Dims{20, 29})
	cfd.Window().SetupBorder("Select File!", "")
	cfd.Window().SendEventsToUnfocused = true

	cfd.fileList.Init(vec.Dims{20, 25}, vec.ZERO_COORD, ui.BorderDepth)
	cfd.fileList.EnableBorder()
	cfd.fileList.SetEmptyText("No Files Found!/n/n(Press C or ESCAPE to cancel)")
	cfd.fileList.AcceptInput = true
	cfd.fileList.EnableHighlight()
	cfd.okayButton.Init(vec.Dims{8, 1}, vec.Coord{1, 27}, 1, "[L]oad File", func() {
		if cfd.fileList.Count() > 0 {
			cfd.chosenFile = cfd.filenames[cfd.fileList.GetSelectionIndex()]
			cfd.CreateTimer(20, func() { cfd.done = true })
		}
	})
	cfd.okayButton.EnableBorder()
	cfd.cancelButton.Init(vec.Dims{8, 1}, vec.Coord{11, 27}, 1, "[C]ancel", func() {
		cfd.CreateTimer(20, func() { cfd.done = true })
	})
	cfd.cancelButton.EnableBorder()
	cfd.Window().AddChildren(&cfd.fileList, &cfd.okayButton, &cfd.cancelButton)

	var err error
	cfd.filenames, err = util.GetFileList(dirPath, ext)
	if err != nil {
		cfd.fileList.SetEmptyText("Could not load files!/n/n(See log.txt for details, Press C or ESCAPE to cancel)")
	}

	cfd.fileList.InsertText(ui.JUSTIFY_LEFT, cfd.filenames...)

	return
}

func (cfd *ChooseFileDialog) HandleKeypress(key_event *input.KeyboardEvent) (event_handled bool) {
	switch key_event.Key {
	case input.K_RETURN, input.K_l:
		cfd.okayButton.Press()
	case input.K_ESCAPE, input.K_c:
		cfd.cancelButton.Press()
	}
	return
}

func (cfd *ChooseFileDialog) Done() bool {
	return cfd.done
}

func (cfd *ChooseFileDialog) Shutdown() {
	if cfd.onFileLoaded != nil && cfd.chosenFile != "" {
		cfd.onFileLoaded(cfd.chosenFile)
	}
}

type SaveDialog struct {
	tyumi.State
	done bool

	onFileSaved func(filename string)

	nameInput    ui.InputBox
	fileText     ui.Textbox
	saveButton   ui.Button
	cancelButton ui.Button

	chosenFilename string
	ext            string
	dirPath        string
	filenames      []string //current contents of directory so we can warn against overwrites
}

func NewSaveDialog(dirPath, ext, default_filename string, onSave func(filename string)) (sd *SaveDialog) {
	sd = new(SaveDialog)
	sd.dirPath = dirPath
	sd.ext = ext

	sd.InitCentered(vec.Dims{26, 10})
	sd.SetKeypressHandler(sd.HandleKeypress)
	sd.Window().SetupBorder("Choose Save Name", "[ENTER]/[ESC]")

	sd.Window().AddChild(ui.NewTextbox(vec.Dims{5, 1}, vec.Coord{1, 2}, 0, "Filename:", ui.JUSTIFY_LEFT))
	sd.nameInput.Init(vec.Dims{17, 1}, vec.Coord{7, 2}, 0, 0)
	sd.nameInput.EnableBorder()
	sd.nameInput.Focus()
	sd.nameInput.ChangeText(default_filename)
	sd.nameInput.OnTextChanged = sd.UpdateFileText

	sd.fileText.Init(vec.Dims{24, 3}, vec.Coord{1, 4}, 0, "Input filename.", ui.JUSTIFY_CENTER)

	sd.saveButton.Init(vec.Dims{5, 1}, vec.Coord{7, 8}, 2, "Save", func() {
		sd.chosenFilename = sd.dirPath + sd.nameInput.InputtedText()
		sd.CreateTimer(20, func() { sd.done = true })
	})
	sd.saveButton.EnableBorder()
	sd.cancelButton.Init(vec.Dims{5, 1}, vec.Coord{14, 8}, 2, "Cancel", func() {
		sd.CreateTimer(20, func() { sd.done = true })
	})
	sd.cancelButton.EnableBorder()

	sd.Window().AddChildren(&sd.nameInput, &sd.fileText, &sd.saveButton, &sd.cancelButton)

	sd.filenames, _ = util.GetFileList(dirPath, "")
	sd.UpdateFileText()

	return
}

func (sd *SaveDialog) HandleKeypress(key_event *input.KeyboardEvent) (event_handled bool) {
	switch key_event.Key {
	case input.K_RETURN:
		if sd.nameInput.InputtedText() != "" {
			sd.saveButton.Press()
			return true
		}
	case input.K_ESCAPE:
		sd.cancelButton.Press()
		return true
	}

	return
}

func (sd *SaveDialog) UpdateFileText() {
	if filename := sd.nameInput.InputtedText(); filename == "" {
		sd.fileText.ChangeText("Please input filename.")
	} else {
		name := filename + sd.ext
		sd.fileText.ChangeText("Will save the file as " + sd.dirPath + name)
		if slices.Contains(sd.filenames, name) {
			sd.fileText.AppendText("/nThis will OVERWRITE the current file!")
		}
	}
}

func (sd *SaveDialog) Shutdown() {
	if sd.onFileSaved != nil && sd.chosenFilename != "" {
		sd.onFileSaved(sd.chosenFilename)
	}
}

func (sd *SaveDialog) Done() bool {
	return sd.done
}

type CommDialog struct {
	tyumi.State

	okayButton ui.Button
	done       bool
}

func NewCommDialog(from, to, picFile, message string) (cd *CommDialog) {
	cd = new(CommDialog)
	cd.Init()

	cd.Window().AddChild(ui.NewImage(vec.ZERO_COORD, 0, picFile))
	cd.Window().AddChild(ui.NewTextbox(vec.Dims{35, 5}, vec.Coord{13, 3}, 0, message, ui.JUSTIFY_LEFT))
	cd.Window().AddChild(ui.NewTextbox(vec.Dims{35, 1}, vec.Coord{13, 0}, 0, "FROM: "+from, ui.JUSTIFY_LEFT))
	cd.Window().AddChild(ui.NewTextbox(vec.Dims{35, 1}, vec.Coord{13, 1}, 0, "TO: "+to, ui.JUSTIFY_LEFT))
	cd.okayButton.MoveTo(vec.Coord{12 + (36-cd.okayButton.Size().W)/2, 10})

	return
}

func NewSimpleCommDialog(message string) (cd *CommDialog) {
	cd = new(CommDialog)
	cd.Init()

	cd.Window().AddChild(ui.NewTextbox(vec.Dims{48, 5}, vec.Coord{0, 1}, 0, message, ui.JUSTIFY_CENTER))
	cd.okayButton.CenterHorizontal()

	return
}

func (cd *CommDialog) Init() {
	cd.State.InitCentered(vec.Dims{48, 12})
	cd.Window().EnableBorder()

	cd.okayButton.Init(vec.Dims{6, 1}, vec.Coord{0, 10}, 1, "Sounds Good!", nil)
	cd.okayButton.EnableBorder()
	cd.okayButton.Focus()
	//cd.okayButton.OnPressAnimation.(*gfx.PulseAnimation).Blocking = true

	cd.Window().AddChild(&cd.okayButton)

	cd.Events().Listen(gfx.EV_ANIMATION_COMPLETE)
	cd.SetEventHandler(cd.HandleEvent)
}

func (cd *CommDialog) HandleEvent(game_event event.Event) (event_handled bool) {
	if game_event.ID() == gfx.EV_ANIMATION_COMPLETE {
		cd.done = true
		return true
	}

	return
}

func (cd CommDialog) Done() bool {
	return cd.done
}
