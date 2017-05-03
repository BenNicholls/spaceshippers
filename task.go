package main

type Job interface {
	OnEnd()
	OnTick()
	OnInterrupt()
	GetDuration() int
	GetName() string
	GetDescription() string
	SetWorker(w *Crewman)
}

type Task struct {
	name         string
	desc         string
	duration     int
	timeInvested int
	worker       *Crewman //person doing the job, nil if not being done.
}

func (t Task) GetDuration() int {
	return t.duration
}

func (t Task) GetName() string {
	return t.name
}

func (t Task) GetDescription() string {
	return t.desc
}

//when job is finished, release the lowly slave who completed it.
func (t Task) OnEnd() {
	t.worker.CurrentTask = nil
}

func (t *Task) OnTick() {
	t.timeInvested++
}

//TODO: this is a thing that should not be...
func (t Task) OnInterrupt() {
	t.OnEnd()
}

func (t *Task) SetWorker(w *Crewman) {
	t.worker = w
}

type RepairRoomJob struct {
	Task
	location *Room
}

func NewRepairRoomJob(r *Room) *RepairRoomJob {
	return &RepairRoomJob{Task{"Repair", "Repairing " + r.name, 0, 0, nil}, r}
}

func (rj *RepairRoomJob) OnTick() {
	rj.Task.OnTick()

	if rj.timeInvested%rj.location.repairDifficulty == 0 {
		rj.location.state.Mod(1)
		if rj.location.state.IsMax() {
			rj.OnEnd()
		}
	}
}

func (rj *RepairRoomJob) OnEnd() {
	AddMessage("Repair of " + rj.location.name + " by " + rj.worker.Name + " completed.")
	rj.Task.OnEnd()
}

type SleepJob struct {
	Task
	//restfulness int //quality of sleep, determines amount of alertness recovered per tick
}

func NewSleepJob() *SleepJob {
	return &SleepJob{Task{"Sleep", "Sleeping", 0, 0, nil}}
}

func (sj *SleepJob) OnTick() {
	sj.Task.OnTick()

	sj.worker.Awakeness.Mod(4)

	if sj.worker.Awakeness.IsMax() {
		sj.OnEnd()
	}
}