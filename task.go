package main

type Job interface {
    OnEnd()
    OnTick()
    OnInterrupt()
    GetDuration() int //returns spaceTime units
    GetName() string
    GetDescription() string
}

type Task struct {
    name string
    desc string
    duration int
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

func (t Task) OnEnd() {

}

func (t Task) OnTick() {

}

func (t Task) OnInterrupt() {

}

type RepairRoomJob struct {
    Task
    location *Room
}

func NewRepairRoomJob(r *Room) *RepairRoomJob {
    return &RepairRoomJob{Task{"Repair", "Repairing " + r.name, 10}, r}
}

