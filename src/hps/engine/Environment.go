package engine

import (
	"bytes"
	"encoding/json"
	"fmt"
	"hps"
	"strings"
	"time"
)

// SimulationResult contains status and exit time of the simulation.
type SimulationResult struct {
	Status   int
	Time     float64
	Events   int
	Duration time.Duration
}

// Environment is a top-level element of the simulation machine tree. It
// contains machines, events, and methods for running the simulation.
type Environment struct {
	logger            hps.ILogger
	machines          []hps.IMachine
	startingIteration []hps.EnvironmentEventHandler
	endingIteration   []hps.EnvironmentEventHandler
}

// NewEnvironment creates the environmenwith the default settings.
func NewEnvironment(logger hps.ILogger) *Environment {
	e := &Environment{
		logger,
		[]hps.IMachine{},
		[]hps.EnvironmentEventHandler{},
		[]hps.EnvironmentEventHandler{},
	}
	_ = hps.IEnvironment(e)
	return e
}

// AddMachine adds machine into the environment.
func (e *Environment) AddMachine(m hps.IMachine) {
	e.machines = append(e.machines, m)
}

// RemoveMachine removes machine from the environment.
func (e *Environment) RemoveMachine(m hps.IMachine) {
	for i, em := range e.machines {
		if em == m {
			e.machines = append(e.machines[0:i], e.machines[i+1:]...)
			break
		}
	}
}

// Machines returns list of machines in the environment.
func (e *Environment) Machines() []hps.IMachine {
	return e.machines
}

// Machine returns first machine that matches the given name.
func (e *Environment) Machine(name string) (hps.IMachine, bool) {
	for _, m := range e.machines {
		if m.Name() == name {
			return m, true
		}
	}
	return nil, false
}

// Event gets the next event from the environment machines.
func (e *Environment) Event(time float64) (hps.IEvent, bool) {
	var event hps.IEvent
	hasEvent := false

	for _, m := range e.machines {
		me, mhe := m.Event(time)

		if mhe && (event == nil || me.Time() < event.Time()) {
			event = me
			hasEvent = true
		}
	}

	return event, hasEvent
}

// StartingIteration adds new handler to be invoked when the iteration starts.
func (e *Environment) StartingIteration(handler hps.EnvironmentEventHandler) {
	e.startingIteration = append(e.startingIteration, handler)
}

// EndingIteration adds new handler to be invoked when the iteration ends.
func (e *Environment) EndingIteration(handler hps.EnvironmentEventHandler) {
	e.endingIteration = append(e.endingIteration, handler)
}

// Run runs the simulation. Simulation stops when one the specified limits is
// reached.
func (e *Environment) Run(limits *Limits) *SimulationResult {
	msgLimits := []string{}
	if limits.Duration != 0 {
		msgLimits = append(msgLimits, fmt.Sprintf("duration = %v", limits.Duration))
	}
	if limits.Time != 0 {
		msgLimits = append(msgLimits, fmt.Sprintf("time = %v", limits.Time))
	}
	if limits.Events != 0 {
		msgLimits = append(msgLimits, fmt.Sprintf("events = %v", limits.Events))
	}
	if limits.Predicate != nil {
		msgLimits = append(msgLimits, "predicate")
	}
	e.logger.Printf("run with limits: %s", strings.Join(msgLimits, ", "))

	count := 0
	currentTime := 0.0
	started := time.Now().Unix()

	createAndLogResult := func(statis int) *SimulationResult {
		duration := time.Duration(time.Now().Unix()-started) * time.Second

		result := &SimulationResult{statis, currentTime, count, duration}

		msg := fmt.Sprintf("duration = %v, events = %d, time = %v",
			result.Duration,
			result.Events,
			result.Time)

		switch result.Status {
		case COMPLETED_BY_DURATION:
			e.logger.Printf("duration limit is reached, %s", msg)
		case COMPLETED_BY_TIME:
			e.logger.Printf("time limit is reached, %s", msg)
		case COMPLETED_BY_ITERATIONS:
			e.logger.Printf("iterations limit is reached, %s", msg)
		case COMPLETED_BY_IDLE:
			e.logger.Printf("no more events, %s", msg)
		}

		return result
	}

	for {
		event, ok := e.Event(currentTime)
		if !ok {
			return createAndLogResult(COMPLETED_BY_IDLE)
		}

		currentTime = event.Time()

		if limits.Events > 0 && count > limits.Events {
			return createAndLogResult(COMPLETED_BY_ITERATIONS)
		}

		if limits.Time > 0 && currentTime > limits.Time {
			return createAndLogResult(COMPLETED_BY_TIME)
		}

		if limits.Duration > 0 && float64(time.Now().Unix()-started) >= float64(limits.Duration) {
			return createAndLogResult(COMPLETED_BY_DURATION)
		}

		if limits.Predicate != nil && !limits.Predicate() {
			return createAndLogResult(COMPLETED_BY_PREDICATE)
		}

		for _, h := range e.startingIteration {
			h(e, event)
		}

		count++

		event.Invoke()

		for _, h := range e.endingIteration {
			h(e, event)
		}
	}
}

// MarshalJSON returns the environment and its machines serialised to JSON.
func (e *Environment) MarshalJSON() ([]byte, error) {
	machines, _ := e.marshalMachinesJSON()
	return json.Marshal(struct {
		Machines *json.RawMessage `json:"machines"`
	}{(*json.RawMessage)(&machines)})
}

func (e *Environment) marshalMachinesJSON() ([]byte, error) {
	buf := bytes.Buffer{}
	buf.WriteString("{")
	for i, machine := range e.machines {
		if i != 0 {
			buf.WriteString(",")
		}
		key, _ := json.Marshal(machine.Name())
		buf.Write(key)
		buf.WriteString(":")
		value, _ := json.Marshal(machine)
		buf.Write(value)
	}
	buf.WriteString("}")
	return buf.Bytes(), nil
}
