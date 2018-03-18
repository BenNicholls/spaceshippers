package main

type PersonType int

var DEFAULT_PIC = "res/art/scippie.xp"

const (
	PERSON_PLAYER PersonType = iota
	PERSON_CREWMAN
	PERSON_CONTACT
)

type Person struct {
	Name  string
	Ptype PersonType
	Pic   string //string to a picture file
}
