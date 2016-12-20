package main

//TODO: custom minimal value support. right now just uses zero
type Stat struct {
    val int
    max int
}

//make a new stat with value at max
func NewStat(v int) Stat {
    return Stat{v, v}
}

func (s Stat) Get() int {
    return s.val
}

func (s Stat) GetMax() int {
    return s.max
}

func (s Stat) GetPct() int {
    return int(100*(float32(s.val)/float32(s.max)))
}

func (s Stat) IsMax() bool {
    if s.val == s.max {
        return true
    }
    return false 
}

//takes a delta
func (s *Stat) Mod(d int) {
    s.val += d

    if s.val > s.max {
        s.val = s.max
    } else if s.val < 0 {
        s.val = 0
    }
}