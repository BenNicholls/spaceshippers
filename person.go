package main

type PersonType int

var DEFAULT_PIC = "res/art/scippie.xp"

const (
	PERSON_PLAYER PersonType = iota
	PERSON_CREWMAN
	PERSON_CONTACT
)

type Person struct {
	Name      string
	Ptype     PersonType
	Pic       string //string to a picture file
	BirthDate int    //birthdate in SpaceTime format
	Race      string //eventually this should be a RaceType or something
}

//Creates a person with given name.
func NewPersonContact(name string) (p *Person) {
	p = new(Person)

	p.Name = name
	p.Ptype = PERSON_CONTACT
	p.Pic = DEFAULT_PIC

	return
}
