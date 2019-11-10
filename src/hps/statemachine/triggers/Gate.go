package triggers

import (
	sm "hps/statemachine"
	"math"
)

type GateTrigger struct {
	machine *sm.Machine
	state   *sm.State
	opened  bool
	comment string
}

func Gate(machine *sm.Machine, state *sm.State) *GateTrigger {
	t := &GateTrigger{}
	t.machine = machine
	t.state = state
	_ = sm.ITrigger(t)
	return t
}

func (t *GateTrigger) Open() {
	if t.machine.State() == t.state {
		t.opened = true
	}
}

func (t *GateTrigger) Time(time float64) float64 {
	if t.opened {
		return time
	}
	return math.Inf(1)
}

func (t *GateTrigger) SetMachine(machine *sm.Machine) {
}

func (t *GateTrigger) StateHandler(reactivate bool) sm.StateEventHandler {
	return func(s *sm.State, time float64, duration float64) {
		t.Open()
		if reactivate {
			t.machine.Reactivate(time)
		}
	}
}

func (t *GateTrigger) SetComment(comment string) {
	t.comment = comment
}
