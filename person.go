package main

type PersonType int

const (
	PERSON_PLAYER PersonType = iota
	PERSON_CREWMAN
	PERSON_CONTACT
)

type Person struct {
	Name  string
	Ptype PersonType
}
