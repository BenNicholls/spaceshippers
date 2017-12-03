package main

import (
	"strconv"
)

type GoalStatus int

const (
	goal_INPROGRESS GoalStatus = iota
	goal_COMPLETE
	goal_FAILED
)

type Acheivable interface {
	Update()
	IsComplete() bool
	IsFailed() bool
	GetName() string
	GetDescription() string
}

type Goal struct {
	name        string
	description string
	status      GoalStatus

	success func() bool
	failure func() bool
}

func (g Goal) GetName() string {
	return g.name
}

func (g Goal) GetDescription() string {
	return g.description
}

func (g Goal) IsComplete() bool {
	if g.status == goal_COMPLETE {
		return true
	}

	return false
}

func (g Goal) IsFailed() bool {
	if g.status == goal_FAILED {
		return true
	}

	return false
}

func (g *Goal) Update() {
	if g.status != goal_INPROGRESS {
		return
	}

	if g.failure() {
		g.status = goal_FAILED
	} else if g.success() {
		g.status = goal_COMPLETE
	}
}

//TODO: Implement some sort of "Dirty" flag so we can tell when a
//mission parameter has changed.
type Mission struct {
	Goal

	steps    []Acheivable
	criteria []Acheivable
}

func NewMission(name, desc string) (m *Mission) {
	m = new(Mission)
	m.name = name
	m.description = desc
	m.status = goal_INPROGRESS
	m.steps = make([]Acheivable, 0, 0)
	m.criteria = make([]Acheivable, 0, 0)

	m.success = func() bool {
		for _, s := range m.steps {
			s.Update()
			if !s.IsComplete() {
				return false
			}
		}

		return true
	}

	m.failure = func() bool {
		for _, c := range m.criteria {
			c.Update()
			if c.IsFailed() {
				return true
			}
		}

		return false
	}

	return
}

func (m *Mission) AddStep(s Acheivable) {
	m.steps = append(m.steps, s)
}

func (m *Mission) AddCriteria(c Acheivable) {
	m.criteria = append(m.criteria, c)
}

func GenerateGoToMission(s *Ship, dest, avoid Locatable) (m *Mission) {
	m = NewMission("", "")
	m.name = "Go to " + dest.GetName()
	m.description = "You need to get there buddy."

	m.AddStep(NewGoToStep(s, dest))
	m.AddCriteria(NewStayAwayCriteria(s, avoid, 1e6))

	return
}

///////////////////////////
//STEPS
// - Compose missions out of these! Also see the CRITERIA secion below
//////////////////////////

type GoToStep struct {
	Goal

	ship        *Ship
	destination Locatable
}

func NewGoToStep(s *Ship, d Locatable) (gs *GoToStep) {
	gs = new(GoToStep)
	gs.name = "Go to " + d.GetName()
	gs.description = "Navigate your ship to " + d.GetName() + " safely. Or dangerously! As long as you get there."

	gs.ship = s
	gs.destination = d

	gs.success = func() bool {
		if gs.ship.GetCoords().IsIn(gs.destination) {
			return true
		}

		return false
	}

	gs.failure = func() bool {
		return false
	}

	return
}

////////////////////////////
//CRITERIA
// - Compose missions out of these! Criteria must be met at all times, or else the mission FAILS.
// - Also see the STEPS section above.
////////////////////////////

type StayAwayCriteria struct {
	Goal

	ship        *Ship
	destination Locatable
	distance    float64
}

func NewStayAwayCriteria(s *Ship, d Locatable, dist float64) (sac *StayAwayCriteria) {
	sac = new(StayAwayCriteria)
	sac.name = "Do not approach " + d.GetName()
	sac.description = "Stay at least " + strconv.Itoa(int(dist/1000)) + "km away from " + d.GetName()

	sac.ship = s
	sac.destination = d
	sac.distance = dist

	sac.success = func() bool {
		return false
	}

	sac.failure = func() bool {
		if sac.ship.GetCoords().CalcVector(sac.destination.GetCoords()).Local.Mag() < sac.distance {
			return true
		}

		return false
	}

	return
}
