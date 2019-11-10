package tools

import (
	"encoding/json"
	"fmt"
	"hps"
	nm "hps/networkmachine"
	sm "hps/statemachine"
)

type EventLog struct {
	e      hps.IEnvironment
	time   float64
	logger hps.ILogger
}

// NewEventLog creates event log which outputs to stdout.
func NewEventLog(logger hps.ILogger) *EventLog {
	el := &EventLog{}
	el.logger = logger
	return el
}

func (el *EventLog) AttachTraces(e hps.IEnvironment) error {
	el.time = 0.0

	e.StartingIteration(func(env hps.IEnvironment, evt hps.IEvent) {
		el.time = evt.Time()
	})

	machineTracers([]string{"machines"}, e.Machines(), func(delta map[string]interface{}) {
		b, _ := json.Marshal(struct {
			Time  float64     `json:"time"`
			Delta interface{} `json:"delta"`
		}{el.time, delta})
		el.logger.Print(string(b))
	})

	return nil
}

func (el *EventLog) Printf(format string, v ...interface{}) {
	b, _ := json.Marshal(struct {
		Time    float64 `json:"time"`
		Message string  `json:"message"`
	}{el.time, fmt.Sprintf(format, v...)})
	el.logger.Print(string(b))

}

func (el *EventLog) Print(v ...interface{}) {
	b, _ := json.Marshal(struct {
		Time    float64 `json:"time"`
		Message string  `json:"message"`
	}{el.time, fmt.Sprint(v...)})
	el.logger.Print(string(b))
}

func createTracer(path []string, tag string) func(interface{}) map[string]interface{} {
	root := map[string]interface{}{}
	leaf := root
	for _, p := range path {
		newLeaf := map[string]interface{}{}
		leaf[p] = (interface{})(newLeaf)
		leaf = newLeaf
	}
	return func(value interface{}) map[string]interface{} {
		leaf[tag] = value
		return root
	}
}

func propertyTracers(path []string, machine hps.IMachine, callback func(map[string]interface{})) {
	for _, p := range machine.Properties() {
		tracer := createTracer(append(path, "properties", p.Name()), "value")
		p.Changed(func(p *hps.Property, oldValue interface{}) {
			if v, ok := p.Value().(hps.IJSON); ok {
				callback(tracer(v.ToJSON()))
			} else {
				callback(tracer(p.Value()))
			}
		})
	}
}

func machineTracers(path []string, machines []hps.IMachine, callback func(map[string]interface{})) {
	for _, m := range machines {
		machineTracer(append(path, m.Name()), m, callback)
	}
}

func machineTracer(path []string, machine hps.IMachine, callback func(map[string]interface{})) {
	propertyTracers(path, machine, callback)

	switch m := machine.(type) {
	case sm.IStateMachine:
		stateTracer := createTracer(path, "state")
		m.Changed(func(p *sm.Machine) {
			callback(stateTracer(m.State().Name()))
		})
	case nm.INetworkMachine:
		machineTracers(append(path, "content", "machines"), m.Machines(), callback)
	}
}
