package statemachine

import (
	"fmt"
	"hps"
)

type Event struct {
	machine    *Machine
	transition *Transition
	from       *State
	to         *State
	time       float64
}

func (e *Event) Machine() hps.IMachine {
	return e.machine
}

func (e *Event) From() *State {
	return e.from
}

func (e *Event) To() *State {
	return e.to
}

func (e *Event) Time() float64 {
	return e.time
}

func (e *Event) Invoke() float64 {
	// clear the saved event
	e.machine.event = nil

	for _, h := range e.from.leaving {
		h(e.from, e.time, e.time-e.from.entered)
	}

	e.machine.state = e.to
	e.to.entered = e.time

	for _, h := range e.to.entering {
		h(e.to, e.time, 0)
	}
	for _, h := range e.machine.changed {
		h(e.machine)
	}

	return e.time
}

func (e *Event) String() string {
	return fmt.Sprintf("%v : %v [ %v -> %v ]", e.time, e.machine.Name(), e.from.Name(), e.to.Name())
}
