package main

import (
	"strconv"

	"github.com/bennicholls/burl-E/burl"
	"github.com/veandco/go-sdl2/sdl"
)

type ShipDesignMenu struct {
	burl.BaseState

	window *burl.Container

	roomColumn  *burl.Container
	roomList    *burl.List
	roomDetails *burl.Container

	shipView   *burl.TileView
	shipColumn *burl.Container
	buttons    *burl.Container
	helpText   *burl.Textbox

	dialog Dialog

	stars StarField

	ship *Ship

	roomToAdd   *Room
	roomAddGood bool

	offX, offY int //shipview camera offsets
}

func NewShipDesignMenu() (sdm *ShipDesignMenu) {
	sdm = new(ShipDesignMenu)

	sdm.window = burl.NewContainer(78, 43, 1, 1, 0, true)
	sdm.window.SetTitle("USE YOUR IMAGINATION")

	sdm.roomColumn = burl.NewContainer(20, 43, 0, 0, 0, true)
	sdm.roomColumn.Add(burl.NewTextbox(20, 1, 0, 0, 0, true, true, "Installed Modules"))
	sdm.roomList = burl.NewList(20, 20, 0, 2, 0, true, "Ain't got no modules! Press [A] to add some Mr. No-Ship!!")
	sdm.roomList.Highlight = false

	sdm.roomDetails = burl.NewContainer(20, 20, 0, 23, 0, true)

	sdm.roomColumn.Add(sdm.roomList, sdm.roomDetails)

	sdm.shipView = burl.NewTileView(36, 36, 21, 0, 0, false)
	sdm.shipColumn = burl.NewContainer(20, 43, 58, 0, 0, true)
	sdm.buttons = burl.NewContainer(36, 6, 21, 37, 0, true)
	sdm.helpText = burl.NewTextbox(36, 6, 21, 37, 1, true, true, "")
	sdm.helpText.SetVisibility(false)

	sdm.window.Add(sdm.roomColumn, sdm.shipView, sdm.shipColumn, sdm.buttons, sdm.helpText)

	sdm.stars = NewStarField(20, sdm.shipView)

	sdm.ship = NewShip("whatever", nil)
	sdm.CenterView()

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

func (sdm *ShipDesignMenu) ActivateHelpText(help string) {
	sdm.helpText.ChangeText(help)
	sdm.helpText.SetVisibility(true)
}

func (sdm *ShipDesignMenu) HandleKeypress(key sdl.Keycode) {
	if sdm.dialog != nil {
		sdm.dialog.HandleKeypress(key)
		return
	}

	switch key {
	case sdl.K_a:
		if sdm.roomToAdd == nil {
			sdm.roomToAdd = CreateRoomFromTemplate(ROOM_ENGINE_MEDIUM)
			sdm.roomToAdd.X = sdm.ship.shipMap.Width/2 - sdm.roomToAdd.Width/2
			sdm.roomToAdd.Y = sdm.ship.shipMap.Height/2 - sdm.roomToAdd.Height/2
			sdm.CenterView()
			sdm.UpdateRoomState()
			sdm.ActivateHelpText("ADDING MODULE: " + sdm.roomToAdd.Name + "/n/n Press ARROW KEYS to move, ENTER to add module to ship, and ESCAPE to cancel.")
		}
	case sdl.K_RETURN:
		if sdm.roomToAdd != nil {
			sdm.AddRoomToShip()
			sdm.helpText.SetVisibility(false)
		}
	case sdl.K_UP:
		if sdm.roomToAdd != nil {
			sdm.roomToAdd.Y -= 1
			sdm.UpdateRoomState()
		}
	case sdl.K_DOWN:
		if sdm.roomToAdd != nil {
			sdm.roomToAdd.Y += 1
			sdm.UpdateRoomState()
		}
	case sdl.K_LEFT:
		if sdm.roomToAdd != nil {
			sdm.roomToAdd.X -= 1
			sdm.UpdateRoomState()
		}
	case sdl.K_RIGHT:
		if sdm.roomToAdd != nil {
			sdm.roomToAdd.X += 1
			sdm.UpdateRoomState()
		}
	case sdl.K_ESCAPE:
		if sdm.roomToAdd != nil {
			sdm.roomToAdd = nil
			sdm.helpText.SetVisibility(false)
		}
	}
}

func (sdm *ShipDesignMenu) Update() {
	sdm.Tick++

	if sdm.Tick%10 == 0 {
		sdm.stars.Shift()
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
	sdm.UpdateRoomList()
	sdm.UpdateRoomDetails()
}

func (sdm *ShipDesignMenu) UpdateRoomList() {
	sdm.roomList.ClearElements()
	for _, r := range sdm.ship.Rooms {
		sdm.roomList.Append(r.Name)
	}
}

func (sdm *ShipDesignMenu) UpdateRoomDetails() {
	sdm.roomDetails.ClearElements()
	if len(sdm.roomList.Elements) != 0 {
		room := sdm.ship.Rooms[sdm.roomList.GetSelection()]
		sdm.roomDetails.Add(burl.NewTextbox(20, 1, 0, 0, 0, true, true, room.Name))
		sdm.roomDetails.Add(burl.NewTextbox(20, 2, 0, 2, 0, false, true, room.Description))
		sdm.roomDetails.Add(burl.NewTextbox(20, 1, 0, 5, 0, false, false, "Dims: ("+strconv.Itoa(room.Width)+"x"+strconv.Itoa(room.Height)+")"))
		sdm.roomDetails.Add(burl.NewTextbox(20, 1, 0, 7, 0, false, false, "STATS:"))
		for i, s := range room.Stats {
			sdm.roomDetails.Add(burl.NewTextbox(20, 1, 2, 8+i, 0, false, false, s.GetName()+": "+strconv.Itoa(s.Modifier)))
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
	}

	sdm.window.Render()

	if sdm.dialog != nil {
		sdm.dialog.Render()
	}
}
