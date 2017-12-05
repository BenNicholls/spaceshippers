package main

import "math/rand"
import "github.com/bennicholls/burl-E/burl"

var FIRSTNAMES []string
var LASTNAMES []string

func init() {
	FIRSTNAMES = []string{"Armund", "Bort", "Chet", "Danzig", "Elton", "Francine", "Geralt", "Hooper", "Ingrid", "Jassy", "Klepta", "Liam", "Mumpy", "Ninklas", "Oliver", "Pernissa", "Quentin", "Rosalinda", "Shlupp", "Timmy", "Ursula", "Vivica", "Wendel", "Xavier", "Yuppie", "Zelda"}
	LASTNAMES = []string{"Andleman", "Bunchlo", "Cogsworth", "Doofer", "Encelada", "Fink", "Gusto", "Humber", "Illiamson", "Jasprex", "Klefbom", "Lorax", "Munkleberg", "Ning", "Olberson", "Pinzip", "Quaker", "Ruffsborg", "Shlemko", "Thrace", "Undergarb", "Von Satan", "White", "Xom", "Yillian", "Zaphod"}
}

type Crewman struct {
	burl.BurlEntity
	Person

	//defining characteristics of various types
	HP        burl.Stat
	Awakeness burl.Stat

	CurrentTask Job
}

func NewCrewman() *Crewman {
	c := new(Crewman)
	c.Vis.Glyph = burl.GLYPH_FACE1
	c.Vis.ForeColour = burl.COL_WHITE
	c.HP = burl.NewStat(100)
	c.Awakeness = burl.NewStat((rand.Intn(4) + 7) * HOUR)
	c.randomizeName()
	c.Ptype = PERSON_CREWMAN
	return c
}

func (c *Crewman) randomizeName() {
	c.Name = FIRSTNAMES[rand.Intn(len(FIRSTNAMES))] + " " + LASTNAMES[rand.Intn(len(LASTNAMES))]
}

//general per-tick update function.
func (c *Crewman) Update() {
	//increase sleepy. if too sleepy, drop what your doing and go to sleep.
	if c.IsAwake() {
		c.Awakeness.Mod(-1)
	}
	if c.Awakeness.Get() == 0 {
		if c.CurrentTask != nil {
			c.CurrentTask.OnInterrupt()
		}
		c.ConsumeJob(NewSleepJob())
	}

	if c.CurrentTask != nil {
		c.CurrentTask.OnTick()
	} else {
		//job finding code goes here, write the code why don't you
	}
}

func (c Crewman) GetStatus() string {

	if c.Awakeness.GetPct() < 15 {
		return "Tired"
	} else if c.HP.GetPct() > 80 {
		return "Great"
	} else if c.HP.GetPct() > 50 {
		return "Fine"
	} else if c.HP.GetPct() > 20 {
		return "Struggling"
	} else if c.HP.GetPct() > 0 {
		return "Near Death"
	} else {
		return "Dead"
	}
}

func (c Crewman) IsAwake() bool {
	if c.CurrentTask != nil && c.CurrentTask.GetName() == "Sleep" {
		return false
	}
	return true
}

func (c *Crewman) ConsumeJob(j Job) {
	if c.CurrentTask != nil {
		c.CurrentTask.OnInterrupt()
	}

	c.CurrentTask = j
	c.CurrentTask.SetWorker(c)
}
