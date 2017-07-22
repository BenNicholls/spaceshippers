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

type MissionLog []*Mission 

type Mission struct {
	name string
	description string

	status MissionStatus
	steps []*Mission

	success func() bool
	failure func() bool
}

func NewMission(name, desc string) *Mission {
	m := new(Mission)
	m.name = name
	m.description = desc
	m.status = mis_INPROGRESS
	m.steps = make([]*Mission, 0, 0)

	m.success = func() bool {return false}
	m.failure = func() bool {return false}

	return m
}
