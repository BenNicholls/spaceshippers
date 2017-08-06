package main

type MissionStatus int

const (
	mis_INPROGRESS MissionStatus = iota
	mis_COMPLETE
	mis_FAILED
)

type Acheivable interface {
	Update()
	IsComplete() bool
	IsFailed() bool
	GetName() string
	GetDescription() string
}

type MissionLog struct {
	log []*Mission
}

func NewMissionLog() (ml *MissionLog) {
	ml = new(MissionLog)
	ml.log = make([]*Mission, 0)
	return
}

func (ml *MissionLog) Add(m *Mission) {
	ml.log = append(ml.log, m)
}

type Mission struct {
	name        string
	description string

	status MissionStatus
	steps  []*Mission

	success func() bool
	failure func() bool
}

func NewMission(name, desc string) (m *Mission) {
	m = new(Mission)
	m.name = name
	m.description = desc
	m.status = mis_INPROGRESS
	m.steps = make([]*Mission, 0, 0)

	m.success = func() bool { return false }
	m.failure = func() bool { return false }

	return
}
