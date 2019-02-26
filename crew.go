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
	burl.EntityPrototype
	Person

	//defining characteristics of various types
	HP        burl.Stat
	Awakeness burl.Stat

	CurrentTask Job

	ship *Ship //reference to the ship. THINK: what if a crewman isn't on the ship?
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
func (c *Crewman) Update(spaceTime int) {

	//breath every 5 seconds (average respiratory rate for adult human)
	//TODO: respiratory rate should change based on exertion/age/physiology/whatever?
	if spaceTime%5 == 0 {
		c.Breathe()
	}

	//increase sleepy. if too sleepy, drop what you're doing and go to sleep.
	if c.IsAwake() {
		c.Awakeness.Mod(-1)
		burl.PushEvent(burl.NewEvent(burl.EV_UPDATE_UI, "crew"))
		if c.Awakeness.Get() == 0 {
			c.ConsumeJob(NewSleepJob())
		}
	}

	//do ya damn job
	if c.CurrentTask != nil {
		c.CurrentTask.OnTick()
	} else {
		//job finding code goes here, write the code why don't you
	}
}

func (c *Crewman) Breathe() {
	if c.ship == nil { //protection against initialized ship? not sure how this would arise.
		burl.LogError("No ship associated with crewman: ", c.Name)
		return
	}

	r := c.ship.GetRoom(c.X, c.Y)

	//inhale 350 mL of air from the atmosphere. 500 mL is the normal Tidal Volume for an adult human, 150 mL of which
	//is maintained in the nose/bronchial tubes/breathing... hose.
	//TODO: this volume should change depending on exertion level? should also increase if the
	//person is oxygen-starved.
	breath := r.atmo.RemoveVolume(0.35)

	//breathing code goes here. exchange o2 for co2, etc.

	//exhale.
	r.atmo.Add(breath)
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
