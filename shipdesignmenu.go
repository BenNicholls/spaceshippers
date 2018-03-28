package main

import (
	"io/ioutil"
	"strconv"

	"github.com/bennicholls/burl-E/burl"
	"github.com/veandco/go-sdl2/sdl"
)

var EV_LOADFILE = burl.RegisterCustomEvent() //event.Message will contain local pathname

type ShipDesignMenu struct {
	burl.BaseState

	window *burl.Container

	roomColumn        *burl.Container
	roomLists         *burl.PagedContainer
	installedRoomList *burl.List
	allRoomList       *burl.List
	roomDetails       *burl.Container

	shipView           *burl.TileView
	selectionAnimation *burl.PulseAnimation
	stars              StarField

	shipColumn *burl.Container

	buttons         *burl.Container
	addRemoveButton *burl.Button
	saveButton      *burl.Button
	loadButton      *burl.Button
	returnButton    *burl.Button
	helpText        *burl.Textbox

	dialog Dialog

	ship        *Ship
	roomToAdd   *Room
	roomAddGood bool

	roomTemplateOrder []RoomType //ordering for room templates, so we can later sort/filter

	offX, offY int //shipview camera offsets
}

func NewShipDesignMenu() (sdm *ShipDesignMenu) {
	sdm = new(ShipDesignMenu)

	sdm.window = burl.NewContainer(78, 43, 1, 1, 0, true)
	sdm.window.SetTitle("USE YOUR IMAGINATION")

	sdm.roomColumn = burl.NewContainer(20, 43, 0, 0, 0, true)
	sdm.roomColumn.Add(burl.NewTextbox(20, 1, 0, 0, 0, false, true, "Modules"))
	sdm.roomLists = burl.NewPagedContainer(20, 20, 0, 2, 0, true)

	sdm.allRoomList = burl.NewList(18, 16, 0, 0, 0, true, "No modules exists in the whole universe somehow.")
	sdm.installedRoomList = burl.NewList(18, 16, 0, 0, 0, true, "Ain't got no modules!")

	all := sdm.roomLists.AddPage("All")
	all.Add(sdm.allRoomList)
	installed := sdm.roomLists.AddPage("Installed")
	installed.Add(sdm.installedRoomList)

	sdm.roomDetails = burl.NewContainer(20, 20, 0, 23, 0, true)

	sdm.roomColumn.Add(sdm.roomLists, sdm.roomDetails)

	sdm.shipView = burl.NewTileView(36, 36, 21, 0, 0, false)
	sdm.shipColumn = burl.NewContainer(20, 43, 58, 0, 0, true)
	sdm.buttons = burl.NewContainer(36, 6, 21, 37, 0, true)

	sdm.addRemoveButton = burl.NewButton(8, 1, 0, 0, 0, true, true, "[A]dd Module")
	sdm.loadButton = burl.NewButton(8, 1, 19, 0, 0, true, true, "[L]oad")
	sdm.saveButton = burl.NewButton(8, 1, 28, 0, 0, true, true, "[S]ave")

	sdm.helpText = burl.NewTextbox(36, 4, 0, 2, 0, true, true, "")
	sdm.buttons.Add(sdm.addRemoveButton, sdm.saveButton, sdm.loadButton, sdm.helpText)

	sdm.window.Add(sdm.roomColumn, sdm.shipView, sdm.shipColumn, sdm.buttons)

	sdm.stars = NewStarField(20, sdm.shipView)

	sdm.ship = NewShip("whatever", nil)
	sdm.CenterView()

	sdm.roomTemplateOrder = make([]RoomType, 0)
	for i, _ := range roomTemplates {
		sdm.roomTemplateOrder = append(sdm.roomTemplateOrder, i)
	}

	sdm.UpdateAllRoomList()
	sdm.UpdateRoomDetails()
	sdm.UpdateHelpText()
	sdm.UpdateSelectionAnimation()

	return
}

func (sdm *ShipDesignMenu) CenterView() {
	w, h := sdm.shipView.Dims()
	if sdm.roomToAdd == nil {
		sdm.offX = sdm.ship.x + sdm.ship.width/2 - w/2
		sdm.offY = sdm.ship.y + sdm.ship.height/2 - h/2
	} else {
		sdm.offX = sdm.roomToAdd.X + sdm.roomToAdd.Width/2 - w/2
		sdm.offY = sdm.roomToAdd.Y + sdm.roomToAdd.Height/2 - h/2
	}
}

func (sdm *ShipDesignMenu) UpdateHelpText() {
	if sdm.roomToAdd != nil {
		sdm.helpText.ChangeText("ADDING MODULE: " + sdm.roomToAdd.Name + "/n/n Press ARROW KEYS to move, R to rotate, ENTER to add module to ship, and ESCAPE to cancel.")
	} else if sdm.roomLists.CurrentIndex() == 0 {
		sdm.helpText.ChangeText("Welcome to the Ship Designer!/n/n Use PGUP/PGDOWN to select a module to add. Press TAB to see all modules currently installed.")
	} else {
		sdm.helpText.ChangeText("Welcome to the Ship Designer!/n/n Use PGUP/PGDOWN to select a module to remove. Press TAB to see all available modules.")
	}
}

func (sdm *ShipDesignMenu) HandleKeypress(key sdl.Keycode) {
	if sdm.dialog != nil {
		sdm.dialog.HandleKeypress(key)
		return
	}

	if sdm.roomToAdd == nil {
		switch key {
		case sdl.K_a:
			if sdm.roomLists.CurrentIndex() == 0 {
				sdm.roomToAdd = CreateRoomFromTemplate(sdm.roomTemplateOrder[sdm.allRoomList.GetSelection()], false)
				if sdm.ship.volume == 0 {
					sdm.roomToAdd.X = sdm.ship.shipMap.Width/2 - sdm.roomToAdd.Width/2
					sdm.roomToAdd.Y = sdm.ship.shipMap.Height/2 - sdm.roomToAdd.Height/2
				} else {
					sdm.roomToAdd.X = sdm.ship.x + sdm.ship.width/2 - sdm.roomToAdd.Width/2
					sdm.roomToAdd.Y = sdm.ship.y + sdm.ship.height/2 - sdm.roomToAdd.Height/2
				}
				sdm.CenterView()
				sdm.UpdateRoomState()
				sdm.UpdateHelpText()
				sdm.addRemoveButton.Press()
			}
		case sdl.K_r:
			if sdm.roomLists.CurrentIndex() == 1 && len(sdm.installedRoomList.Elements) > 0 {
				room := sdm.ship.Rooms[sdm.installedRoomList.GetSelection()]
				sdm.ship.RemoveRoom(room)
				sdm.UpdateInstalledRoomList()
				sdm.addRemoveButton.Press()
			}
		case sdl.K_l:
			sdm.loadButton.Press()
			sdm.dialog = NewChooseFileDialog("raws/ship/")
		case sdl.K_TAB:
			sdm.roomLists.HandleKeypress(key)
			sdm.UpdateRoomDetails()
			sdm.UpdateHelpText()
			sdm.UpdateSelectionAnimation()
			if sdm.roomLists.CurrentIndex() == 0 {
				sdm.addRemoveButton.ChangeText("[A]dd Module")
			} else {
				sdm.addRemoveButton.ChangeText("[R]emove Module")
			}

		case sdl.K_PAGEUP:
			if sdm.roomLists.CurrentIndex() == 0 {
				sdm.allRoomList.Prev()
			} else {
				sdm.installedRoomList.Prev()
				sdm.UpdateSelectionAnimation()
			}
			sdm.UpdateRoomDetails()
		case sdl.K_PAGEDOWN:
			if sdm.roomLists.CurrentIndex() == 0 {
				sdm.allRoomList.Next()
			} else {
				sdm.installedRoomList.Next()
				sdm.UpdateSelectionAnimation()
			}
			sdm.UpdateRoomDetails()
		}
	} else {
		switch key {
		case sdl.K_r:
			sdm.roomToAdd.Rotate()
			sdm.UpdateRoomState()
		case sdl.K_RETURN:
			sdm.AddRoomToShip()
			sdm.UpdateHelpText()
		case sdl.K_UP:
			sdm.roomToAdd.Y -= 1
			sdm.UpdateRoomState()
		case sdl.K_DOWN:
			sdm.roomToAdd.Y += 1
			sdm.UpdateRoomState()
		case sdl.K_LEFT:
			sdm.roomToAdd.X -= 1
			sdm.UpdateRoomState()
		case sdl.K_RIGHT:
			sdm.roomToAdd.X += 1
			sdm.UpdateRoomState()
		case sdl.K_ESCAPE:
			sdm.roomToAdd = nil
			sdm.UpdateHelpText()
		}
	}
}

func (sdm *ShipDesignMenu) UpdateSelectionAnimation() {
	sdm.shipView.RemoveAnimation(sdm.selectionAnimation)
	if len(sdm.installedRoomList.Elements) != 0 {
		room := sdm.ship.Rooms[sdm.installedRoomList.GetSelection()]
		sdm.selectionAnimation = burl.NewPulseAnimation(0, 0, room.Width, room.Height, 50, 0, true)
		sdm.shipView.AddAnimation(sdm.selectionAnimation)
		if sdm.roomLists.CurrentIndex() == 1 {
			sdm.selectionAnimation.Activate()
		}
	}
}

func (sdm *ShipDesignMenu) Update() {
	if sdm.dialog != nil && sdm.dialog.Done() {
		sdm.dialog = nil
	}

	sdm.Tick++

	if sdm.Tick%10 == 0 {
		sdm.stars.Shift()
	}
}

func (sdm *ShipDesignMenu) HandleEvent(e *burl.Event) {
	switch e.ID {
	case EV_LOADFILE:
		temp, err := LoadShipTemplate("raws/ship/" + e.Message)
		if err != nil {
			burl.LogError(err.Error())
		} else {
			sdm.ship = NewShip("whatever", nil)
			sdm.ship.SetupFromTemplate(temp)
			sdm.UpdateInstalledRoomList()
			sdm.UpdateRoomDetails()
			sdm.CenterView()
		}
	}
}

func (sdm *ShipDesignMenu) UpdateRoomState() {
	if sdm.roomToAdd != nil {
		sdm.roomAddGood = sdm.ship.CheckRoomValidAdd(sdm.roomToAdd, sdm.roomToAdd.X, sdm.roomToAdd.Y)
	}
}

func (sdm *ShipDesignMenu) AddRoomToShip() {
	sdm.UpdateRoomState()
	if sdm.roomAddGood {
		sdm.ship.AddRoom(sdm.roomToAdd, sdm.roomToAdd.X, sdm.roomToAdd.Y)
		sdm.roomToAdd = nil
	}
	sdm.UpdateInstalledRoomList()
	sdm.UpdateRoomDetails()
}

func (sdm *ShipDesignMenu) UpdateInstalledRoomList() {
	sdm.installedRoomList.ClearElements()
	for _, r := range sdm.ship.Rooms {
		sdm.installedRoomList.Append(r.Name)
	}
	sdm.installedRoomList.CheckSelection()
	sdm.UpdateSelectionAnimation()
	sdm.UpdateRoomDetails()
}

func (sdm *ShipDesignMenu) UpdateAllRoomList() {
	sdm.allRoomList.ClearElements()
	for _, temp := range sdm.roomTemplateOrder {
		sdm.allRoomList.Append(roomTemplates[temp].name)
	}
}

func (sdm *ShipDesignMenu) UpdateRoomDetails() {
	sdm.roomDetails.ClearElements()

	var room *Room

	switch sdm.roomLists.CurrentIndex() {
	case 0: //All
		room = CreateRoomFromTemplate(sdm.roomTemplateOrder[sdm.allRoomList.GetSelection()], false)
	case 1: //Installed modules
		if len(sdm.installedRoomList.Elements) != 0 {
			room = sdm.ship.Rooms[sdm.installedRoomList.GetSelection()]
		}
	}

	if room != nil {
		sdm.roomDetails.Add(burl.NewTextbox(20, 1, 0, 0, 0, true, true, room.Name))
		sdm.roomDetails.Add(burl.NewTextbox(20, 3, 0, 2, 0, false, true, room.Description))
		sdm.roomDetails.Add(burl.NewTextbox(20, 1, 0, 6, 0, false, false, "Dims: ("+strconv.Itoa(room.Width)+"x"+strconv.Itoa(room.Height)+")"))
		sdm.roomDetails.Add(burl.NewTextbox(20, 1, 0, 8, 0, false, false, "STATS:"))
		for i, s := range room.Stats {
			sdm.roomDetails.Add(burl.NewTextbox(20, 1, 2, 9+i, 0, false, false, s.GetName()+": "+strconv.Itoa(s.Modifier)))
		}
	}
}

func (sdm *ShipDesignMenu) Render() {
	sdm.stars.Draw()
	sdm.ship.DrawToTileView(sdm.shipView, sdm.offX, sdm.offY)

	if sdm.roomToAdd != nil {
		for i := 0; i < sdm.roomToAdd.Width*sdm.roomToAdd.Height; i++ {

			x := i%sdm.roomToAdd.Width + sdm.roomToAdd.X - sdm.offX
			y := i/sdm.roomToAdd.Width + sdm.roomToAdd.Y - sdm.offY
			w, h := sdm.shipView.Dims()

			if burl.CheckBounds(x, y, w, h) {
				tv := sdm.roomToAdd.RoomMap.GetTile(i%sdm.roomToAdd.Width, i/sdm.roomToAdd.Width).GetVisuals()
				if sdm.roomAddGood {
					sdm.shipView.Draw(x, y, tv.Glyph, tv.ForeColour, burl.COL_GREEN)
				} else {
					sdm.shipView.Draw(x, y, tv.Glyph, tv.ForeColour, burl.COL_RED)
				}
			}
		}
	} else if sdm.roomLists.CurrentIndex() == 1 && len(sdm.installedRoomList.Elements) != 0 {
		room := sdm.ship.Rooms[sdm.installedRoomList.GetSelection()]
		sdm.selectionAnimation.MoveTo(room.X-sdm.offX, room.Y-sdm.offY)
	}

	sdm.window.Render()

	if sdm.dialog != nil {
		sdm.dialog.Render()
	}
}

type ChooseFileDialog struct {
	burl.BaseState

	container    *burl.Container
	fileList     *burl.List
	okayButton   *burl.Button
	cancelButton *burl.Button

	filenames []string

	dirPath string
}

func NewChooseFileDialog(dirPath string) (cfd *ChooseFileDialog) {
	cfd = new(ChooseFileDialog)
	cfd.dirPath = dirPath

	cfd.container = burl.NewContainer(20, 29, 0, 0, 50, true)
	cfd.container.CenterInConsole()
	cfd.container.ToggleFocus()
	cfd.container.SetTitle("Select file!")

	cfd.fileList = burl.NewList(20, 25, 0, 0, 0, true, "No Files Found!/n/n(Press C or ESCAPE to cancel)")
	cfd.fileList.ToggleFocus()
	cfd.okayButton = burl.NewButton(8, 1, 1, 27, 1, true, true, "[L]oad File")
	cfd.cancelButton = burl.NewButton(8, 1, 11, 27, 2, true, true, "[C]ancel")

	cfd.container.Add(cfd.fileList, cfd.okayButton, cfd.cancelButton)

	cfd.filenames = make([]string, 0)
	dirContents, err := ioutil.ReadDir(cfd.dirPath)
	if err != nil {
		cfd.fileList.ChangeEmptyText("Could not load files!/n/n(See log.txt for details, Press C or ESCAPE to cancel)")
		burl.LogError(err.Error())
	} else {
		for i, file := range dirContents {
			if !file.IsDir() {
				cfd.filenames = append(cfd.filenames, dirContents[i].Name())
				cfd.fileList.Append(file.Name())
			}
		}
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

func (cfd *ChooseFileDialog) Render() {
	cfd.container.Render()
}

func (cfd *ChooseFileDialog) Done() bool {
	if cfd.okayButton.PressPulse.IsFinished() {
		burl.PushEvent(burl.NewEvent(EV_LOADFILE, cfd.filenames[cfd.fileList.GetSelection()]))
		cfd.container.ToggleVisible()
		return true
	} else if cfd.cancelButton.PressPulse.IsFinished() {
		cfd.container.ToggleVisible()
		return true
	}

	return false
}
