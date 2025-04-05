package main

import (
	"sort"
	"strconv"
	"strings"

	"github.com/bennicholls/tyumi"
	"github.com/bennicholls/tyumi/gfx"
	"github.com/bennicholls/tyumi/gfx/col"
	"github.com/bennicholls/tyumi/gfx/ui"
	"github.com/bennicholls/tyumi/input"
	"github.com/bennicholls/tyumi/log"
	"github.com/bennicholls/tyumi/rl"
	"github.com/bennicholls/tyumi/util"
	"github.com/bennicholls/tyumi/vec"
)

type ShipDesignMenu struct {
	tyumi.Scene

	roomColumn        ui.Element
	roomLists         ui.PageContainer
	installedRoomList ui.List
	allRoomList       ui.List
	roomDetails       ui.Element

	shipView           rl.TileMapView
	stars              StarField
	selectionAnimation gfx.PulseAnimation

	shipColumn     ui.Element
	shipNameText   ui.Textbox
	shipVolumeText ui.Textbox
	shipStatsList  ui.List

	addRemoveButton ui.Button
	saveButton      ui.Button
	loadButton      ui.Button
	helpText        ui.Textbox

	ship             *Ship
	roomToAdd        *Room
	roomToAddElement RoomElement

	roomTemplateOrder []RoomType //ordering for room templates, so we can later sort/filter

	addRoomState    util.StateID
	removeRoomState util.StateID
}

func NewShipDesignMenu() (sdm *ShipDesignMenu) {
	sdm = new(ShipDesignMenu)
	sdm.ship = NewShip("Unsaved Ship", nil)

	sdm.InitBordered()
	sdm.Window().SetupBorder("USE YOUR IMAGINATION", "")
	windowStyle := ui.DefaultBorderStyle
	windowStyle.TitleJustification = ui.JUSTIFY_CENTER
	sdm.Window().Border.SetStyle(ui.BORDER_STYLE_CUSTOM, windowStyle)
	sdm.Window().SendEventsToUnfocused = true

	sdm.roomColumn.Init(vec.Dims{20, 52}, vec.ZERO_COORD, ui.BorderDepth)
	sdm.roomColumn.EnableBorder()
	sdm.roomColumn.AddChild(ui.NewTitleTextbox(vec.Dims{20, 1}, vec.ZERO_COORD, ui.BorderDepth, "Modules"))

	sdm.roomLists.Init(vec.Dims{20, 20}, vec.Coord{0, 2}, ui.BorderDepth)
	sdm.roomLists.EnableBorder()
	sdm.roomLists.OnPageChanged = func() {
		sdm.UpdateRoomDetails()
		if sdm.roomLists.GetPageIndex() == 1 {
			sdm.ChangeState(sdm.removeRoomState)
		} else {
			sdm.ChangeState(util.STATE_NONE)
		}
	}
	sdm.roomLists.AcceptInput = true
	allPage := sdm.roomLists.CreatePage("All")
	sdm.allRoomList.Init(allPage.Size(), vec.ZERO_COORD, 0)
	sdm.allRoomList.SetEmptyText("No modules exist in the whole universe somehow.")
	sdm.allRoomList.ToggleHighlight()
	sdm.allRoomList.OnChangeSelection = sdm.OnRoomSelectionChange
	sdm.allRoomList.AcceptInput = true
	allPage.AddChild(&sdm.allRoomList)

	installedPage := sdm.roomLists.CreatePage("Installed")
	sdm.installedRoomList.Init(installedPage.Size(), vec.ZERO_COORD, 0)
	sdm.installedRoomList.SetEmptyText("Ain't got no modules!")
	sdm.installedRoomList.ToggleHighlight()
	sdm.installedRoomList.OnChangeSelection = sdm.OnRoomSelectionChange
	sdm.installedRoomList.AcceptInput = true
	installedPage.AddChild(&sdm.installedRoomList)

	sdm.roomDetails.Init(vec.Dims{20, 29}, vec.Coord{0, 23}, ui.BorderDepth)
	roomName := ui.NewTitleTextbox(vec.Dims{20, 1}, vec.ZERO_COORD, ui.BorderDepth, "Room Name")
	roomName.SetLabel("Room Name")
	roomDesc := ui.NewTextbox(vec.Dims{18, 5}, vec.Coord{1, 2}, 0, "Room Description", ui.JUSTIFY_CENTER)
	roomDesc.SetLabel("Room Description")
	roomDims := ui.NewTextbox(vec.Dims{20, 1}, vec.Coord{0, 8}, 0, "DIMS: ???", ui.JUSTIFY_LEFT)
	roomDims.SetLabel("Room Dimensions")
	roomStats := ui.NewList(vec.Dims{18, 18}, vec.Coord{2, 11}, 0)
	roomStats.SetLabel("Room Stats")
	sdm.roomDetails.AddChild(ui.NewTextbox(vec.Dims{20, 1}, vec.Coord{0, 10}, 0, "STATS:", ui.JUSTIFY_LEFT))
	sdm.roomDetails.AddChildren(roomName, roomDesc, roomDims, roomStats)
	sdm.roomColumn.AddChildren(&sdm.roomLists, &sdm.roomDetails)

	sdm.shipView.Init(vec.Dims{52, 45}, vec.Coord{21, 0}, 1, &sdm.ship.shipMap)
	sdm.shipView.SetDefaultVisuals(gfx.Visuals{Mode: gfx.DRAW_NONE, Colours: col.Pair{col.WHITE, col.BLACK}})
	sdm.stars.Init(sdm.shipView.Size(), vec.Coord{21, 0}, 0, 20, 10)
	sdm.shipColumn.Init(vec.Dims{20, 52}, vec.Coord{74, 0}, ui.BorderDepth)
	sdm.shipColumn.EnableBorder()
	sdm.shipColumn.AddChild(ui.NewTitleTextbox(vec.Dims{20, 1}, vec.ZERO_COORD, ui.BorderDepth, "Ship Details"))
	sdm.shipColumn.AddChild(ui.NewTextbox(vec.Dims{5, 1}, vec.Coord{0, 3}, 0, "Ship Name:", ui.JUSTIFY_LEFT))
	sdm.shipColumn.AddChild(ui.NewTextbox(vec.Dims{5, 1}, vec.Coord{0, 4}, 0, "Volume:", ui.JUSTIFY_LEFT))
	sdm.shipNameText.Init(vec.Dims{14, 2}, vec.Coord{6, 3}, 0, "", ui.JUSTIFY_LEFT)
	sdm.shipVolumeText.Init(vec.Dims{14, 1}, vec.Coord{6, 4}, 0, strconv.Itoa(sdm.ship.volume), ui.JUSTIFY_LEFT)
	sdm.shipStatsList.Init(vec.Dims{20, 30}, vec.Coord{0, 22}, ui.BorderDepth)
	sdm.shipStatsList.EnableBorder()
	sdm.shipStatsList.SetEmptyText("Ship has no stats to display. Try adding some modules!")
	sdm.shipColumn.AddChildren(&sdm.shipNameText, &sdm.shipVolumeText, &sdm.shipStatsList)
	sdm.Window().AddChildren(&sdm.shipView, &sdm.stars, &sdm.shipColumn)

	var buttons ui.Element
	buttons.Init(vec.Dims{52, 6}, vec.Coord{21, 46}, ui.BorderDepth)
	buttons.EnableBorder()
	sdm.addRemoveButton.Init(vec.Dims{8, 1}, vec.ZERO_COORD, ui.BorderDepth, "[A]dd Module", nil)
	sdm.addRemoveButton.EnableBorder()
	sdm.loadButton.Init(vec.Dims{8, 1}, vec.Coord{35, 0}, ui.BorderDepth, "[L]oad", func() {
		sdm.OpenDialog(NewChooseFileDialog("raws/ship/", ".shp", sdm.LoadShip))
	})
	sdm.loadButton.EnableBorder()
	sdm.saveButton.Init(vec.Dims{8, 1}, vec.Coord{44, 0}, ui.BorderDepth, "[S]ave", func() {
		if sdm.ship.Name != "Unsaved Ship" {
			sdm.OpenDialog(NewSaveDialog("raws/ship/", ".shp", sdm.ship.Name, sdm.SaveShip))
		} else {
			sdm.OpenDialog(NewSaveDialog("raws/ship/", ".shp", "", sdm.SaveShip))
		}
	})
	sdm.saveButton.EnableBorder()
	sdm.helpText.Init(vec.Dims{52, 4}, vec.Coord{0, 2}, ui.BorderDepth, "", ui.JUSTIFY_CENTER)
	sdm.helpText.EnableBorder()
	buttons.AddChildren(&sdm.addRemoveButton, &sdm.loadButton, &sdm.saveButton, &sdm.helpText)

	sdm.Window().AddChildren(&sdm.roomColumn, &buttons)

	sdm.selectionAnimation = gfx.NewPulseAnimation(vec.Rect{}, 1, 60, col.Pair{col.LIGHTGREY, col.LIGHTGREY})
	sdm.selectionAnimation.Repeat = true
	sdm.shipView.AddAnimation(&sdm.selectionAnimation)

	sdm.SetKeypressHandler(sdm.HandleKeypress)

	sdm.CenterView()

	sdm.roomTemplateOrder = make([]RoomType, 0)
	for i := range roomTemplates {
		sdm.roomTemplateOrder = append(sdm.roomTemplateOrder, i)
	}
	sdm.roomToAddElement.Init(vec.ZERO_COORD, 10, nil)
	sdm.roomToAddElement.Hide()
	sdm.shipView.AddChild(&sdm.roomToAddElement)

	sdm.addRoomState = sdm.RegisterState(util.State{
		OnEnter: func(prev util.StateID) {
			sdm.roomToAdd = CreateRoomFromTemplate(sdm.roomTemplateOrder[sdm.allRoomList.GetSelectionIndex()], false)
			sdm.roomToAddElement.SetRoom(sdm.roomToAdd)
			sdm.CenterView()
			sdm.roomToAddElement.Center()
			sdm.roomToAddElement.Show()
			sdm.UpdateRoomState()
			sdm.allRoomList.AcceptInput = false
		},
		OnLeave: func(next util.StateID) {
			sdm.roomToAdd = nil
			sdm.roomToAddElement.Hide()
			sdm.allRoomList.AcceptInput = true
		},
	})

	sdm.removeRoomState = sdm.RegisterState(util.State{
		OnEnter: func(prev util.StateID) {
			sdm.addRemoveButton.ChangeText("[R]emove Module")
			sdm.UpdateInstalledRoomList()
			sdm.UpdateSelectionAnimation()
		},
		OnLeave: func(next util.StateID) {
			sdm.addRemoveButton.ChangeText("[A]dd Module")
			sdm.selectionAnimation.Stop()
		},
	})

	sdm.OnStateChange = func(prev, next util.StateID) { sdm.UpdateHelpText(next) }

	sdm.UpdateAllRoomList()
	sdm.UpdateRoomDetails()
	sdm.UpdateHelpText(sdm.CurrentState())
	sdm.UpdateShipDetails()

	return
}

func (sdm *ShipDesignMenu) LoadShip(filename string) {
	if filename == "" {
		return
	}

	template, err := LoadShipTemplate("raws/ship/" + filename)
	if err != nil {
		log.Error(err.Error())
		return
	}

	sdm.ship = NewShip("whatever", nil)
	sdm.ship.SetupFromTemplate(template)
	sdm.ship.Name = template.Name
	sdm.shipView.SetTileMap(sdm.ship.shipMap)
	sdm.UpdateInstalledRoomList()
	sdm.UpdateRoomDetails()
	sdm.UpdateShipDetails()
	sdm.CenterView()

	switch sdm.CurrentState() {
	case sdm.addRoomState:
		sdm.UpdateRoomState()
	case sdm.removeRoomState:
		sdm.UpdateSelectionAnimation()
	}
}

func (sdm *ShipDesignMenu) SaveShip(filename string) {
	sdm.ship.Name = strings.TrimSuffix(filename, ".shp")
	template := sdm.ship.CreateShipTemplate()
	err := template.Save()
	if err != nil {
		log.Error(err)
	}
	sdm.UpdateShipDetails()
}

func (sdm *ShipDesignMenu) CenterView() {
	if sdm.ship.Bounds().Area() > 0 {
		sdm.shipView.CenterOnTileMapCoord(sdm.ship.Bounds().Center())
	} else {
		sdm.shipView.CenterTileMap()
	}
}

func (sdm *ShipDesignMenu) HandleKeypress(key_event *input.KeyboardEvent) (event_handled bool) {
	//general inputs valid for all states
	switch key_event.Key {
	case input.K_l:
		sdm.loadButton.Press()
		return true
	case input.K_s:
		sdm.UpdateInstalledRoomList()
		if sdm.installedRoomList.Count() != 0 {
			sdm.saveButton.Press()
			return true
		}
	case input.K_ESCAPE:
		if sdm.CurrentState() != util.STATE_NONE {
			sdm.ChangeState(util.STATE_NONE)
		} else {
			var startMenu StartMenu
			startMenu.Init()
			tyumi.ChangeScene(&startMenu)
		}

		return true
	}

	// state specific inputs
	switch sdm.CurrentState() {
	case sdm.addRoomState:
		if dir := key_event.Direction(); dir != vec.DIR_NONE {
			sdm.roomToAddElement.Move(dir.X, dir.Y)
			sdm.UpdateRoomState()
		}
		switch key_event.Key {
		case input.K_r:
			sdm.roomToAdd.Rotate()
			sdm.roomToAddElement.Resize(sdm.roomToAdd.Size())
			sdm.UpdateRoomState()
		case input.K_RETURN, input.K_a:
			sdm.AddRoomToShip()
			sdm.CenterView()
		}
	case sdm.removeRoomState:
		if key_event.Key == input.K_r {
			if roomIndex := sdm.installedRoomList.GetSelectionIndex(); roomIndex != -1 {
				room := sdm.ship.Rooms[sdm.installedRoomList.GetSelectionIndex()]
				sdm.ship.RemoveRoom(room)
				sdm.UpdateInstalledRoomList()
				sdm.UpdateShipDetails()
				sdm.CenterView()
				sdm.UpdateSelectionAnimation()
				sdm.addRemoveButton.Press()
			}
		}
	default: // STATE_NONE, default state
		switch key_event.Key {
		case input.K_a:
			if sdm.roomLists.GetPageIndex() != 0 {
				return
			}

			sdm.ChangeState(sdm.addRoomState)
			sdm.addRemoveButton.Press()
			return true
		}
	}

	return
}

func (sdm *ShipDesignMenu) AddRoomToShip() {
	sdm.UpdateRoomState()
	if sdm.roomToAddElement.valid {
		sdm.ship.AddRoom(sdm.roomToAdd.pos, sdm.roomToAdd)
		sdm.UpdateShipDetails()
		sdm.ChangeState(util.STATE_NONE)
	}
}

func (sdm *ShipDesignMenu) UpdateHelpText(state util.StateID) {
	switch state {
	case sdm.addRoomState:
		sdm.helpText.ChangeText("ADDING MODULE/n/n Press ARROW KEYS to move, R to rotate, A or ENTER to add module to ship, and ESCAPE to cancel.")
	case sdm.removeRoomState:
		sdm.helpText.ChangeText("REMOVING MODULE/n/n Use UP/DOWN to select a module to remove, and R to remove the module. Press TAB to see all available modules.")
	default:
		sdm.helpText.ChangeText("Welcome to the Ship Designer!/n/n Use UP/DOWN to select a module to add. Press TAB to see all modules currently installed.")
	}
}

func (sdm *ShipDesignMenu) OnRoomSelectionChange() {
	if sdm.CurrentState() == sdm.removeRoomState {
		sdm.UpdateSelectionAnimation()
	}
	sdm.UpdateRoomDetails()
}

func (sdm *ShipDesignMenu) UpdateSelectionAnimation() {
	if sdm.installedRoomList.Count() != 0 {
		room := sdm.ship.Rooms[sdm.installedRoomList.GetSelectionIndex()]
		pos := sdm.shipView.MapCoordToViewCoord(room.pos)
		sdm.selectionAnimation.SetArea(vec.Rect{pos, room.Size()})
		sdm.selectionAnimation.Start()
	} else {
		sdm.selectionAnimation.Stop()
	}
}

func (sdm *ShipDesignMenu) UpdateRoomState() {
	if sdm.roomToAdd != nil {
		sdm.roomToAdd.pos = sdm.roomToAddElement.Bounds().Coord.Add(sdm.shipView.GetCameraOffset())
		sdm.roomToAddElement.SetValid(sdm.ship.CheckRoomValidAdd(sdm.roomToAdd))
	}
}

func (sdm *ShipDesignMenu) UpdateInstalledRoomList() {
	selection := sdm.installedRoomList.GetSelectionIndex()
	sdm.installedRoomList.RemoveAll()
	for _, r := range sdm.ship.Rooms {
		sdm.installedRoomList.InsertText(ui.JUSTIFY_LEFT, r.Name)
	}
	sdm.installedRoomList.Select(selection)
	sdm.UpdateRoomDetails()
}

func (sdm *ShipDesignMenu) UpdateAllRoomList() {
	sdm.allRoomList.RemoveAll()
	for _, temp := range sdm.roomTemplateOrder {
		sdm.allRoomList.InsertText(ui.JUSTIFY_LEFT, roomTemplates[temp].name)
	}
}

func (sdm *ShipDesignMenu) UpdateRoomDetails() {
	var room *Room
	switch sdm.roomLists.GetPageIndex() {
	case 0: //All
		if sdm.allRoomList.Count() > 0 {
			room = CreateRoomFromTemplate(sdm.roomTemplateOrder[sdm.allRoomList.GetSelectionIndex()], false)
		}
	case 1: //Installed modules
		if sdm.installedRoomList.Count() != 0 {
			room = sdm.ship.Rooms[sdm.installedRoomList.GetSelectionIndex()]
		}
	}

	if room != nil {
		sdm.roomDetails.Show()
		ui.GetLabelled[*ui.Textbox](sdm.Window(), "Room Name").ChangeText(room.Name)
		ui.GetLabelled[*ui.Textbox](sdm.Window(), "Room Description").ChangeText(room.Description)
		ui.GetLabelled[*ui.Textbox](sdm.Window(), "Room Dimensions").ChangeText("Dims: (" + strconv.Itoa(room.Width) + "x" + strconv.Itoa(room.Height) + ")")
		statList := ui.GetLabelled[*ui.List](sdm.Window(), "Room Stats")
		statList.RemoveAll()
		for _, s := range room.Stats {
			statList.InsertText(ui.JUSTIFY_LEFT, s.GetName()+": "+strconv.Itoa(s.Modifier))
		}
	} else {
		sdm.roomDetails.Hide()
	}
}

func (sdm *ShipDesignMenu) UpdateShipDetails() {
	sdm.shipNameText.ChangeText(sdm.ship.Name)
	sdm.shipVolumeText.ChangeText(strconv.Itoa(sdm.ship.volume))

	sdm.shipStatsList.RemoveAll()

	for _, sys := range sdm.ship.Systems {
		var statStrings []string
		for _, stat := range sys.GetAllStats() {
			statStrings = append(statStrings, stat.GetName()+": "+strconv.Itoa(stat.Modifier))
		}
		sort.Strings(statStrings)
		sdm.shipStatsList.InsertText(ui.JUSTIFY_LEFT, statStrings...)
	}
}

type RoomElement struct {
	rl.TileMapView

	valid bool
	room  *Room
}

func (re *RoomElement) Init(pos vec.Coord, depth int, room *Room) {
	if room != nil {
		re.TileMapView.Init(room.Size(), pos, depth, &room.RoomMap)
	} else {
		re.TileMapView.Init(vec.Dims{1, 1}, pos, depth, nil)
	}

	re.TreeNode.Init(re)
}

func (re *RoomElement) SetRoom(room *Room) {
	if re.room == room {
		return
	}

	re.room = room
	re.Resize(room.Size())
	re.SetTileMap(&re.room.RoomMap)
	re.valid = false
}

func (re *RoomElement) SetValid(valid bool) {
	if re.valid == valid {
		return
	}

	re.valid = valid
	re.Updated = true
}

func (re *RoomElement) Render() {
	if re.room == nil {
		return
	}

	if re.valid {
		re.SetDefaultColours(col.Pair{col.NONE, col.GREEN})
	} else {
		re.SetDefaultColours(col.Pair{col.NONE, col.RED})
	}

	re.TileMapView.Render()
}
