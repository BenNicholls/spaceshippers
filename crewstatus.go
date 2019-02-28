package main

//Effects can have multiple sources, and statuses can have multiple effects.

type EffectID int

const (
	EFFECT_HEAVYBREATHING EffectID = iota
	EFFECT_SUFFOCATING
	EFFECT_DISORIENTED
	EFFECT_SLOW
	EFFECT_POISONED
)

//These effect the way the crew act, effects their stats, etc. Think of them like BUFFS and DEBUFFS
type CrewEffect struct {
	Name        string
	Description string

	Sources map[StatusID]int //map value is a count of how many turns this source has been active

	Duration int
}

func NewCrewEffect(e EffectID) (ce CrewEffect) {
	switch e {
	case EFFECT_HEAVYBREATHING:
		ce = CrewEffect{
			Name:        "Heavy Breathing",
			Description: "Crewman is breathing heavily, trying to catch his/her breath.",
		}
	case EFFECT_SUFFOCATING:
		ce = CrewEffect{
			Name:        "Suffocating",
			Description: "Crewman can't breathe! Get them some air!!!.",
		}
	case EFFECT_DISORIENTED:
		ce = CrewEffect{
			Name:        "Disoriented",
			Description: "Crewman feels a bit woozy. It's so hard to focus sometimes, you know?",
		}
	case EFFECT_SLOW:
		ce = CrewEffect{
			Name:        "Slow",
			Description: "Crewman is moving slowly, and won't be rushed.",
		}
	case EFFECT_POISONED:
		ce = CrewEffect{
			Name:        "Poisoned",
			Description: "Crewman is poisoned, and is losing life.",
		}
	}

	ce.Sources = make(map[StatusID]int)

	return
}

func (ce *CrewEffect) AddSource(s StatusID) {
	if _, ok := ce.Sources[s]; ok {
		return //if source already present, no add
	}

	ce.Sources[s] = 0
}

func (ce *CrewEffect) RemoveSource(s StatusID) {
	delete(ce.Sources, s)
}

//Called once per tick to update durations and such.
func (ce *CrewEffect) Update() {
	for s, _ := range ce.Sources {
		ce.Sources[s] += 1
	}

	ce.Duration++
}

type StatusID int

const (
	STATUS_LOWOXYGEN StatusID = iota
	STATUS_NOOXYGEN
	STATUS_HIGHCO2
	STATUS_CO2_POISONING
)

//These are the sources of CREW EFFECTS.
//example: LOW OXYGEN ENVIRONMENT causes effect HEAVY BREATHING, DISORIENTED, and SLOW
type CrewStatus struct {
	Name        string
	Description string

	Effects  []EffectID //effects this status causes
	Replaces []StatusID //statuses this status replaces if possible. example: NO OXYGEN replaces LOW OXYGEN
}

func NewCrewStatus(s StatusID) (cs CrewStatus) {
	switch s {
	case STATUS_LOWOXYGEN:
		cs = CrewStatus{
			Name:        "Low Oxygen Environment",
			Description: "The Crewman is breathing air with too little oxygen. While not fatal, it makes everything harder.",
			Effects: []EffectID{
				EFFECT_HEAVYBREATHING,
				EFFECT_DISORIENTED,
				EFFECT_SLOW,
			},
			Replaces: []StatusID{
				STATUS_NOOXYGEN,
			},
		}
	case STATUS_NOOXYGEN:
		cs = CrewStatus{
			Name:        "No Oxygen Environment",
			Description: "The Crewman is breathing air with almost no oxygen! Uh oh!",
			Effects: []EffectID{
				EFFECT_HEAVYBREATHING,
				EFFECT_DISORIENTED,
				EFFECT_SLOW,
				EFFECT_SUFFOCATING,
			},
			Replaces: []StatusID{
				STATUS_LOWOXYGEN,
			},
		}
	case STATUS_HIGHCO2:
		cs = CrewStatus{
			Name:        "High CO2 Levels",
			Description: "The Crewman has breathed air with too much carbon dioxide. Eventually leads to CO2 poisoning, which is bad.",
			Effects: []EffectID{
				EFFECT_SLOW,
			},
			Replaces: []StatusID{},
		}
	case STATUS_CO2_POISONING:
		cs = CrewStatus{
			Name:        "Carbon Dioxide Poisoning",
			Description: "The Crewman has respirated a fatal amount of CO2 and is dying!",
			Effects: []EffectID{
				EFFECT_HEAVYBREATHING,
				EFFECT_DISORIENTED,
				EFFECT_POISONED,
			},
			Replaces: []StatusID{},
		}
	}

	return
}
