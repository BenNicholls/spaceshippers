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
	CO2       burl.Stat //level of CO2 in the blood.
	Dead      bool

	CurrentTask Job

	Statuses map[StatusID]CrewStatus
	Effects  map[EffectID]CrewEffect

	ship *Ship //reference to the ship. THINK: what if a crewman isn't on the ship?
}

func NewCrewman() *Crewman {
	c := new(Crewman)
	c.Vis.Glyph = burl.GLYPH_FACE1
	c.Vis.ForeColour = burl.COL_WHITE
	c.HP = burl.NewStat(100)
	c.Awakeness = burl.NewStat((rand.Intn(4) + 7) * HOUR)
	c.CO2 = burl.NewStat(1000000)
	c.CO2.Set(0)
	c.randomizeName()
	c.Ptype = PERSON_CREWMAN

	c.Statuses = make(map[StatusID]CrewStatus)
	c.Effects = make(map[EffectID]CrewEffect)

	return c
}

func (c *Crewman) randomizeName() {
	c.Name = FIRSTNAMES[rand.Intn(len(FIRSTNAMES))] + " " + LASTNAMES[rand.Intn(len(LASTNAMES))]
}

//general per-tick update function.
func (c *Crewman) Update(spaceTime int) {
	if c.Dead {
		return
	}

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

		//walk around randomly like a doofus.
		if spaceTime%20 == 0 {
			dx, dy := burl.RandomDirection()
			if c.ship.shipMap.GetTile(c.X+dx, c.Y+dy).Empty() {
				c.ship.shipMap.MoveEntity(c.X, c.Y, dx, dy)
				c.Move(dx, dy)
			}
		}
	}

	//do ya damn job
	if c.CurrentTask != nil {
		c.CurrentTask.OnTick()
	} else {
		//job finding code goes here, write the code why don't you
	}

	if c.CO2.GetPct() > 20 {
		c.AddStatus(STATUS_HIGHCO2)
		if c.CO2.IsMax() {
			c.AddStatus(STATUS_CO2_POISONING)
		}
	}

	c.HandleEffects(spaceTime)

	if c.HP.Get() == 0 {
		c.Dead = true
		c.Vis.ForeColour = burl.COL_RED
		if c.CurrentTask != nil {
			c.CurrentTask.OnInterrupt()
		}
		c.CurrentTask = nil
		burl.PushEvent(burl.NewEvent(LOG_EVENT, c.Name+" has died! :("))
	}
}

func (c *Crewman) Breathe() {
	if c.ship == nil { //protection against uninitialized ship? not sure how this would arise but it's probably a good idea.
		burl.LogError("No ship associated with crewman: ", c.Name)
		return
	}

	r := c.ship.GetRoom(c.X, c.Y)
	if r == nil {
		//no room found... person is outside? TODO: spacemen dont like being outside. hurts their eyes. handle that.
		return
	}

	//inhale 350 mL of air from the atmosphere. 500 mL is the normal Tidal Volume for an adult human, 150 mL of which
	//is maintained in the nose/bronchial tubes/breathing... hose.
	v := 0.35
	if c.HasEffect(EFFECT_HEAVYBREATHING) {
		v = 0.7
	}
	if !c.IsAwake() { //crewman can go into Hibernative Naptosis to conserve oxygen
		v = v / 2
	}

	breath := r.atmo.RemoveVolume(v)

	//check CO2 levels. Modify crewman CO2 poisoning counter. TODO: crewman should just straight up suffocate at some point.
	switch {
	case breath.CO2 > 7: //WAY too high. CO2 poisoning in 5-10 mins
		c.CO2.Mod(10000)
	case breath.CO2 > 5: //Very high. CO2 poisoning in around a day
		c.CO2.Mod(55)
	case breath.CO2 > 3: //High. CO2 poisoning in a couple weeks
		c.CO2.Mod(4)
	case breath.CO2 > 1: //A little too much. Get dizzy if breathed for a while. Can't give you CO2 poisoning.
		if c.CO2.GetPct() < 50 {
			c.CO2.Mod(1)
		}
	case breath.CO2 <= 1: //low-ish CO2, slowly brings crewman's CO2 levels back to normal.
		c.CO2.Mod(-10)
	}

	//Check oxygen levels. Gas exchange.
	switch {
	case breath.O2 < 5: //O2 dangerously low
		breath.CO2 += breath.O2
		breath.O2 = 0
		c.AddStatus(STATUS_NOOXYGEN)
	case breath.O2 < 15: //O2 content between 5 - 15 kpa. bad, but not fatal.
		breath.O2 -= 4 //THINK: maybe people can get acclimated to this level of O2? Like sherpas and Joe Sakic do.
		breath.CO2 += 4
		c.AddStatus(STATUS_LOWOXYGEN)
	default:
		breath.O2 -= 6
		breath.CO2 += 6
	}

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

func (c Crewman) HasEffect(e EffectID) bool {
	_, ok := c.Effects[e]
	return ok
}

func (c Crewman) HasStatus(s StatusID) bool {
	_, ok := c.Statuses[s]
	return ok
}

func (c *Crewman) AddStatus(s StatusID) {
	if c.HasStatus(s) {
		return
	}

	status := NewCrewStatus(s)
	c.Statuses[s] = status

	for _, id := range status.Effects {
		c.AddEffect(id, s)
	}

	for _, id := range status.Replaces {
		c.RemoveStatus(id)
	}

	burl.PushEvent(burl.NewEvent(LOG_EVENT, c.Name+" is now affected by "+status.Name))
}

func (c *Crewman) RemoveStatus(s StatusID) {
	if !c.HasStatus(s) {
		return
	}

	status := c.Statuses[s]
	delete(c.Statuses, s)

	for _, e := range status.Effects {
		effect := c.Effects[e]
		effect.RemoveSource(s)

		if len(effect.Sources) == 0 {
			delete(c.Effects, e) //if that was the last source, remove the effect
		} else {
			c.Effects[e] = effect
		}
	}
}

func (c *Crewman) AddEffect(e EffectID, s StatusID) {
	var effect CrewEffect

	if c.HasEffect(e) {
		effect = c.Effects[e]
	} else {
		effect = NewCrewEffect(e)
	}

	effect.AddSource(s)

	c.Effects[e] = effect
}

//Per-turn things for effects
func (c *Crewman) HandleEffects(spaceTime int) {
	var e CrewEffect

	for id := range c.Effects {
		e = c.Effects[id]
		e.Update()

		switch id {
		case EFFECT_SUFFOCATING:
			//person passes out and slowly begins dying after 2 minutes without air
			if e.Duration > 120 {
				if c.IsAwake() {
					c.ConsumeJob(NewSleepJob())
					burl.PushEvent(burl.NewEvent(LOG_EVENT, c.Name+" passed out!"))
				}
				c.Awakeness.Set(0) //make sure they don't wake up from being passed out
				c.HP.Mod(-1)
			}
		case EFFECT_POISONED:
			for s, _ := range e.Sources {
				switch s {
				case STATUS_CO2_POISONING:
					if spaceTime%10 == 0 {
						c.HP.Mod(-1)
					}
				}
			}

		}

		c.Effects[id] = e
	}
}
