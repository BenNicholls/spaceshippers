package main

import (
	"math/rand"

	"github.com/bennicholls/tyumi/gfx"
	"github.com/bennicholls/tyumi/gfx/col"
	"github.com/bennicholls/tyumi/log"
	"github.com/bennicholls/tyumi/rl"
	"github.com/bennicholls/tyumi/vec"
)

var FIRSTNAMES []string
var LASTNAMES []string

func init() {
	FIRSTNAMES = []string{"Armund", "Bort", "Chet", "Danzig", "Elton", "Francine", "Geralt", "Hooper", "Ingrid", "Jassy", "Klepta", "Liam", "Mumpy", "Ninklas", "Oliver", "Pernissa", "Quentin", "Rosalinda", "Shlupp", "Timmy", "Ursula", "Vivica", "Wendel", "Xavier", "Yuppie", "Zelda"}
	LASTNAMES = []string{"Andleman", "Bunchlo", "Cogsworth", "Doofer", "Encelada", "Fink", "Gusto", "Humber", "Illiamson", "Jasprex", "Klefbom", "Lorax", "Munkleberg", "Ning", "Olberson", "Pinzip", "Quaker", "Ruffsborg", "Shlemko", "Thrace", "Undergarb", "Von Satan", "White", "Xom", "Yillian", "Zaphod"}
}

var ENTITY_CREWMAN = rl.RegisterEntityType(rl.EntityData{
	Name:    "Crewman",
	Desc:    "A loyal member of the crew, devoted to exploring space with their best buds!",
	Glyph:   gfx.GLYPH_FACE1,
	Colours: col.Pair{col.WHITE, col.NONE},
})

type Crewman struct {
	rl.Entity
	Person

	Updated bool // true if crew menu UI needs to re-render

	//defining characteristics of various types
	HP        rl.Stat[int]
	Awakeness rl.Stat[int]
	CO2       rl.Stat[int] //level of CO2 in the blood.
	Dead      bool

	CurrentTask Job

	Statuses map[StatusID]CrewStatus
	Effects  map[EffectID]CrewEffect

	ship *Ship //reference to the ship. THINK: what if a crewman isn't on the ship?
}

func NewCrewman() *Crewman {
	c := new(Crewman)
	c.Init(ENTITY_CREWMAN)
	c.HP = rl.NewBasicStat(100)
	c.Awakeness = rl.NewBasicStat((rand.Intn(4) + 7) * HOUR)
	c.CO2 = rl.NewBasicStat(1000000)
	c.CO2.Set(0)
	c.randomizeName()
	c.Ptype = PERSON_CREWMAN
	c.Pic = DEFAULT_PIC
	c.BirthDate = rand.Intn(30 * CYCLE) //born sometime during the first 30 cycles of the digital age. so they'll be at least twenty?
	c.Race = "Human"

	c.Statuses = make(map[StatusID]CrewStatus)
	c.Effects = make(map[EffectID]CrewEffect)

	return c
}

func (c *Crewman) randomizeName() {
	c.Name = FIRSTNAMES[rand.Intn(len(FIRSTNAMES))] + " " + LASTNAMES[rand.Intn(len(LASTNAMES))]
}

func (c Crewman) GetVisuals() (v gfx.Visuals) {
	v = c.Entity.GetVisuals()
	if c.Dead {
		v.Colours.Fore = col.RED
	} else if !c.IsAwake() {
		v.Colours.Fore = col.YELLOW
	}

	return
}

// general per-tick update function.
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
		c.Updated = true
		if c.Awakeness.GetPct() < 10 {
			c.AddStatus(STATUS_SLEEPY)
		}
		if c.Awakeness.Get() == 0 {
			c.ConsumeJob(NewSleepJob())
		}

		//walk around randomly like a doofus.
		if spaceTime%20 == 0 {
			dir := vec.RandomDirection()
			if c.ship.shipMap.GetTile(c.Position().Step(dir)).IsPassable() {
				c.ship.shipMap.MoveEntity(c.Position(), c.Position().Step(dir))
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
		if c.CurrentTask != nil {
			c.CurrentTask.OnInterrupt()
			c.CurrentTask = nil
		}
		fireSpaceLogEvent(c.Name + " has died! :(")
	}
}

func (c *Crewman) Breathe() {
	if c.ship == nil { //protection against uninitialized ship? not sure how this would arise but it's probably a good idea.
		log.Error("No ship associated with crewman: ", c.Name)
		return
	}

	r := c.ship.GetRoom(c.Position())
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
	case breath.PartialPressure(GAS_CO2) > 7: //WAY too high. CO2 poisoning in 5-10 mins
		c.CO2.Mod(10000)
	case breath.PartialPressure(GAS_CO2) > 5: //Very high. CO2 poisoning in around a day
		c.CO2.Mod(55)
	case breath.PartialPressure(GAS_CO2) > 3: //High. CO2 poisoning in a couple weeks
		c.CO2.Mod(4)
	case breath.PartialPressure(GAS_CO2) > 1: //A little too much. Get dizzy if breathed for a while. Can't give you CO2 poisoning.
		if c.CO2.GetPct() < 50 {
			c.CO2.Mod(1)
		}
	case breath.PartialPressure(GAS_CO2) <= 1: //low-ish CO2, slowly brings crewman's CO2 levels back to normal.
		c.CO2.Mod(-10)
	}

	//Check oxygen levels.
	var oxygenIntake float64
	switch {
	case breath.PartialPressure(GAS_O2) < 5: //O2 dangerously low
		oxygenIntake = breath.GetMolarValue(GAS_O2)
		c.AddStatus(STATUS_NOOXYGEN)
	case breath.PartialPressure(GAS_O2) < 15: //O2 content between 5 - 15 kpa. bad, but not fatal.
		oxygenIntake = 4 * breath.Volume //THINK: maybe people can get acclimated to this level of O2? Like sherpas and Joe Sakic do.
		c.AddStatus(STATUS_LOWOXYGEN)
	default:
		oxygenIntake = 6 * breath.Volume
	}

	//gas exchange
	breath.RemoveGas(GAS_O2, oxygenIntake)
	breath.AddGas(GAS_CO2, oxygenIntake)

	//exhale.
	r.atmo.AddGasMixture(breath)
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
	c.Updated = true

	for _, id := range status.Effects {
		c.AddEffect(id, s)
	}

	for _, id := range status.Replaces {
		c.RemoveStatus(id)
	}
}

func (c *Crewman) RemoveStatus(s StatusID) {
	if !c.HasStatus(s) {
		return
	}

	status := c.Statuses[s]
	delete(c.Statuses, s)
	c.Updated = true

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
	c.Updated = true
}

// Per-turn things for effects
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
					fireSpaceLogEvent(c.Name + " passed out!")
				}
				c.Awakeness.Set(0) //make sure they don't wake up from being passed out
				c.HP.Mod(-1)
			}
		case EFFECT_POISONED:
			for s := range e.Sources {
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
