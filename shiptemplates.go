package main

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

//this file is going to keep ship template data so we can set up a few kinds of ship
//eventually this will all be read from raw text files, but for now i want to fly a
//spaceship damn it so that will have to wait

var defaultShipTemplates map[ShipType]ShipTemplate

type ShipType int

const (
	SHIPTYPE_CIVILIAN ShipType = iota
	SHIPTYPE_TRANSPORT
	//SHIPTYPE_MINING
	//SHIPTYPE_FIGHTER
	//SHIPTYPE_EXPLORER
	SHIPTYPE_CUSTOM //used for all player-designed ship templates.
)

type ShipTemplate struct {
	Name        string
	Description string
	Shiptype    ShipType
	Rooms       []RoomDef
	CrewNum     int
}

type RoomDef struct {
	RoomType      RoomType
	Rotated       bool
	Width, Height int
	X, Y          int
}

func init() {
	defaultShipTemplates = make(map[ShipType]ShipTemplate)

	defaultShipTemplates[SHIPTYPE_CIVILIAN] = ShipTemplate{
		Name:        "Civilian Craft",
		Description: "The Toyota Camry of spaceships. The Civilian Craft has everything the casual spacegoer might need: door, engine, steering wheel of some variety, snack table, all the cool space things. While it may not look like much, it's cheap to repair, very moddable, and will last forever if you take care of it right!",
		Shiptype:    SHIPTYPE_CIVILIAN,
		Rooms: []RoomDef{
			RoomDef{
				RoomType: ROOM_COCKPIT,
				X:        34,
				Y:        29,
			},
			RoomDef{
				RoomType: ROOM_COMMONAREA,
				X:        26,
				Y:        27,
			},
			RoomDef{
				RoomType: ROOM_ENGINE_SMALL,
				X:        24,
				Y:        28,
			},
		},
		CrewNum: 4,
	}

	defaultShipTemplates[SHIPTYPE_TRANSPORT] = ShipTemplate{
		Name:        "Transport Ship",
		Description: "Used for transporting passengers and cargo. The Transport Ship begins with a larger cargo bay, additional dormitories, and a bulkier engine. Guzzles like an Irishman though.",
		Shiptype:    SHIPTYPE_TRANSPORT,
		Rooms: []RoomDef{
			RoomDef{
				RoomType: ROOM_BRIDGE,
				X:        45,
				Y:        21,
			},
			RoomDef{
				RoomType: ROOM_CORRIDOR,
				Width:    16,
				Height:   4,
				X:        34,
				Y:        29,
			},
			RoomDef{
				RoomType: ROOM_CORRIDOR,
				Rotated:  true,
				X:        36,
				Y:        18,
			},
			RoomDef{
				RoomType: ROOM_QUARTERS,
				X:        39,
				Y:        18,
			},
			RoomDef{
				RoomType: ROOM_QUARTERS,
				X:        39,
				Y:        24,
			},
			RoomDef{
				RoomType: ROOM_CARGOBAY,
				X:        29,
				Y:        20,
			},
			RoomDef{
				RoomType: ROOM_ENGINE_LARGE,
				X:        28,
				Y:        27,
			},
		},
		CrewNum: 6,
	}
}

//we are assuming the ship's defaults are already set
func (s *Ship) SetupFromTemplate(temp ShipTemplate) {
	s.ShipType = temp.Shiptype
	s.Description = temp.Description
	s.ShipType = temp.Shiptype

	for _, r := range temp.Rooms {
		s.AddRoom(CreateRoomFromTemplate(r.RoomType, r.Rotated, r.Width, r.Height), r.X, r.Y)
	}
}

func (st *ShipTemplate) Save() error {
	f, err := os.Create("raws/ship/" + st.Name + ".txt")
	if err != nil {
		return err
	}
	defer f.Close()

	data, err := yaml.Marshal(st)
	if err != nil {
		return err
	}

	f.Write(data)

	return nil
}

func LoadShipTemplate(path string) (ShipTemplate, error) {
	st := ShipTemplate{}

	f, err := os.Open(path)
	if err != nil {
		return st, err
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return st, err
	}

	err = yaml.Unmarshal(data, &st)
	if err != nil {
		return st, err
	}

	return st, nil

}
