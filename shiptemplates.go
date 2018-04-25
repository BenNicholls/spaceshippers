package main

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type ShipTemplate struct {
	Name        string
	Description string
	Rooms       []RoomDef
	CrewNum     int
}

type RoomDef struct {
	RoomType      RoomType
	Rotated       bool
	Width, Height int
	X, Y          int
}

func (st *ShipTemplate) Save() error {
	f, err := os.Create("raws/ship/" + st.Name + ".shp")
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

//we are assuming the ship's defaults are already set
func (s *Ship) SetupFromTemplate(temp ShipTemplate) {
	s.Description = temp.Description

	for _, r := range temp.Rooms {
		s.AddRoom(CreateRoomFromTemplate(r.RoomType, r.Rotated, r.Width, r.Height), r.X, r.Y)
	}
}

func (s *Ship) CreateShipTemplate() (st ShipTemplate) {
	st = ShipTemplate{
		Name:        s.Name,
		Description: s.Description,
	}

	st.Rooms = make([]RoomDef, len(s.Rooms))
	for i, room := range s.Rooms {
		st.Rooms[i] = RoomDef{
			RoomType: room.Roomtype,
			Rotated:  room.Rotated,
			Width:    room.Width,
			Height:   room.Height,
			X:        room.X,
			Y:        room.Y,
		}
		if room.Rotated {
			st.Rooms[i].Width, st.Rooms[i].Height = st.Rooms[i].Height, st.Rooms[i].Width
		}
	}

	st.CrewNum = 4

	return
}
