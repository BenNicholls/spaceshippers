package main

import (
	"io"
	"os"

	"github.com/bennicholls/tyumi/vec"
	"gopkg.in/yaml.v2"
)

type ShipTemplate struct {
	Name        string
	Description string
	Rooms       []RoomDef
	CrewNum     int
}

// TO DO: get rid of yaml, move to json or something
type RoomDef struct {
	RoomType RoomType
	Rotated  bool
	Size     vec.Dims  `yaml:",flow"`
	Position vec.Coord `yaml:",flow"`
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

	data, err := io.ReadAll(f)
	if err != nil {
		return st, err
	}

	err = yaml.Unmarshal(data, &st)
	if err != nil {
		return st, err
	}

	return st, nil
}

// we are assuming the ship's defaults are already set
func (s *Ship) SetupFromTemplate(template ShipTemplate) {
	s.Description = template.Description

	for _, room := range template.Rooms {
		s.AddRoom(room.Position, CreateRoomFromTemplate(room.RoomType, room.Rotated, room.Size))
	}

	for range template.CrewNum {
		s.AddCrewman(NewCrewman())
	}
}

func (s *Ship) CreateShipTemplate() (st ShipTemplate) {
	st = ShipTemplate{
		Name:        s.Name,
		Description: s.Description,
		CrewNum:     4,
	}

	st.Rooms = make([]RoomDef, len(s.Rooms))
	for i, room := range s.Rooms {
		st.Rooms[i] = RoomDef{
			RoomType: room.Roomtype,
			Rotated:  room.Rotated,
			Size:     vec.Dims{room.Width, room.Height},
			Position: room.pos,
		}
		if room.Rotated {
			st.Rooms[i].Size.W, st.Rooms[i].Size.H = st.Rooms[i].Size.H, st.Rooms[i].Size.W
		}
	}

	return
}
