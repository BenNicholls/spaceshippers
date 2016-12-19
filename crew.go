package main

import "math/rand"

var FIRSTNAMES []string
var LASTNAMES []string

func init() {
    FIRSTNAMES = []string{"Armund", "Bort", "Chet", "Danzig", "Elton", "Francine", "Geralt", "Hooper", "Ingrid", "Jassy", "Klepta", "Liam", "Mumpy", "Ninklas", "Oliver", "Pernissa", "Quentin", "Rosalinda", "Shlupp", "Timmy", "Ursula", "Vivica", "Wendel", "Xavier", "Yuppie", "Zelda"}
    LASTNAMES = []string{"Andleman", "Bunchlo", "Cogsworth", "Doofer", "Encelada", "Fink", "Gusto", "Humber", "Illiamson", "Jasprex", "Klefbom", "Lorax", "Munkleberg", "Ning", "Olberson", "Pinzip", "Quaker", "Ruffsborg", "Shlemko", "Thrace", "Undergarb", "Von Satan", "White", "Xom", "Yillian", "Zaphod"}
}

type Crewman struct {
    Name string
    HP int
    CurrentTask Job
}

func NewCrewman() *Crewman {
    c := new(Crewman)
    c.HP = 100
    c.randomizeName()
    return c
}

func (c *Crewman) randomizeName() {
    c.Name = FIRSTNAMES[rand.Intn(len(FIRSTNAMES))] + " " + LASTNAMES[rand.Intn(len(LASTNAMES))]
}

//general per-tick update function. 
func (c *Crewman) Update() {
}

func (c Crewman) GetStatus() string {
    if c.HP > 80 {
        return "Great"
    } else if c.HP > 50 {
        return "Fine"
    } else if c.HP > 20 {
        return "Struggling"
    } else  if c.HP > 0 {
        return "Near Death"
    } else {
        return "Dead"
    }
}
