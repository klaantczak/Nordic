package statemachine

import (
	"bytes"
	"encoding/json"
	"fmt"
	"hps"
)

type MachineEventHandler func(p *Machine)

type Machine struct {
	name       string
	kind       string
	state      *State
	states     []*State
	event      *Event
	changed    []MachineEventHandler
	properties map[string]*hps.Property
}

func NewMachine(name string, kind string, states []string, state string) (*Machine, error) {
	m := &Machine{}
	m.name = name
	m.kind = kind

	m.states = []*State{}
	for _, name := range states {
		s := &State{m, name, []*Transition{}, 0, []StateEventHandler{}, []StateEventHandler{}}
		if name == state {
			m.state = s
		}
		m.states = append(m.states, s)
	}
	if m.state == nil {
		return nil, fmt.Errorf("initial state is incorrect")
	}

	m.properties = map[string]*hps.Property{}
	_ = hps.IMachine(m)
	_ = IStateMachine(m)
	return m, nil
}

func (m *Machine) Name() string {
	return m.name
}

func (m *Machine) Kind() string {
	return m.kind
}

func (m *Machine) States() []*State {
	list := []*State{}
	for _, s := range m.states {
		list = append(list, s)
	}
	return list
}

func (m *Machine) State() *State {
	return m.state
}

// SetState modifies the current state of the machine. Returns whether
// state was successfuly changed or not.
func (m *Machine) SetState(name string, time float64) bool {
	state, ok := m.GetState(name)
	if !ok {
		return false
	}

	if state.Name() == m.state.Name() {
		return false
	}

	for _, h := range m.state.leaving {
		h(m.state, time, time-m.state.entered)
	}

	state.entered = time
	m.event = nil
	m.state = state

	for _, h := range state.entering {
		h(state, time, 0)
	}

	for _, handler := range m.changed {
		handler(m)
	}
	return true
}

func (m *Machine) GetState(name string) (*State, bool) {
	for _, s := range m.states {
		if s.Name() == name {
			return s, true
		}
	}
	return nil, false
}

func (m *Machine) Event(time float64) (hps.IEvent, bool) {
	if m.event == nil {
		var e *Event = nil
		for _, t := range m.state.transitions {
			te := &Event{}
			te.machine = m
			te.transition = t
			te.from = m.state
			te.to, _ = m.GetState(t.to)
			te.time = t.Trigger().Time(time)
			if e == nil || e.time > te.time {
				e = te
			}
		}
		m.event = e
	}
	return m.event, m.event != nil
}

func (m *Machine) Reactivate(time float64) {
	m.event = nil
	m.Event(time)
}

func (m *Machine) ClearEvent() {
	// TODO THIS IS FOR BACKWARD COMPATIBILITY ONLY
	m.event = nil
}

func (m *Machine) Properties() []*hps.Property {
	list := []*hps.Property{}
	for _, p := range m.properties {
		list = append(list, p)
	}
	return list
}

func (m *Machine) AddProperty(name, dataType string) *hps.Property {
	p := hps.NewProperty(name, dataType)
	m.properties[name] = p
	return p
}

func (m *Machine) Property(name string) (*hps.Property, bool) {
	p, ok := m.properties[name]
	return p, ok
}

func (m *Machine) Changed(handler MachineEventHandler) {
	m.changed = append(m.changed, handler)
}

// Loads the machine from the JSON definition. The definition is the JavaScript
// object with the following properties:
//
// - `name`, string, required, uniquely identifies the machine in some context
// - `machine`, string, optional, type of the machine
// - `states`, []string, required, non-empty, list of machine states, first
//   counts as initial unless the `state` property is specified.
// - `state`, string, optional, initial/current state.
// - `properties`, map of properties, optional
// - `transitions`, map of transitions, optional
//
// Properties are defined as map. Map keys correspond to property names and
// values contain property definitions. Example:
//
//     "properties" : {
//       "failureRate" : { "type" : "ITrigger" }
//     }
//
// Transitions are defined as two nested maps, containing the list of triggers
// that may perform the transition. Example:
//
//    "transitions" : {
//      "ok" : { "fail" : [ { "type" : "idle" } ] },
//      "fail" : { "ok" : [ { "type" : ... }, { ... }, ... ] }
//    }
func (m *Machine) UnmarshalJSON(data []byte) error {
	machine := struct {
		Name        string           `json:"name"`
		Machine     string           `json:"machine"`
		States      []string         `json:"states"`
		State       string           `json:"state"`
		Properties  *json.RawMessage `json:"properties"`
		Transitions *json.RawMessage `json:"transitions"`
	}{}
	err := json.Unmarshal(data, &machine)
	if err != nil {
		return err
	}
	if machine.Name == "" {
		return fmt.Errorf("The property \"name\" is required.")
	}
	if len(machine.States) == 0 {
		return fmt.Errorf("The property \"states\" is required and should be non-empty.")
	}

	*m = Machine{}
	m.name = machine.Name
	m.kind = machine.Machine
	m.states = []*State{}
	m.properties = map[string]*hps.Property{}

	if machine.Properties != nil {
		properties := map[string]struct {
			Type  string           `json:"type"`
			Value *json.RawMessage `json:"value"`
		}{}
		err = json.Unmarshal([]byte(*machine.Properties), &properties)
		if err != nil {
			return err
		}
		for k, v := range properties {
			property := m.AddProperty(k, v.Type)
			if v.Value != nil {
				value, err := hps.PropertyValueUnmarshalJSON(property.Type(), []byte(*v.Value))
				if err != nil {
					return err
				}
				property.SetValue(value)
			}
		}
	}

	if len(machine.States) == 0 {
		return fmt.Errorf("machine has no states")
	}

	for _, name := range machine.States {
		s := &State{m, name, []*Transition{}, 0, []StateEventHandler{}, []StateEventHandler{}}
		m.states = append(m.states, s)
	}

	if machine.State == "" {
		m.state = m.states[0]
	} else {
		for _, s := range m.states {
			if s.name == machine.State {
				m.state = s
				break
			}
		}
		if m.state == nil {
			return fmt.Errorf("initial state is incorrect")
		}
	}

	transitions := map[string]map[string][]*json.RawMessage{}
	if machine.Transitions != nil {
		err := json.Unmarshal([]byte(*machine.Transitions), &transitions)
		if err != nil {
			return err
		}

		for from, tmp1 := range transitions {
			for to, tmp2 := range tmp1 {
				for _, item := range tmp2 {
					fromState, ok := m.GetState(from)
					if !ok {
						return fmt.Errorf("State is not defined.")
					}
					toState, ok := m.GetState(to)
					if !ok {
						return fmt.Errorf("State is not defined.")
					}

					transition, err := UnmarshalTrigger([]byte(*item))
					if err != nil {
						return err
					}
					fromState.AddTransition(toState.Name(), transition)
				}
			}
		}
	}

	return nil
}

func (m *Machine) MarshalJSON() ([]byte, error) {
	states, _ := m.MarshalStatesJSON()
	properties, _ := m.MarshalPropertiesJSON()
	return json.Marshal(struct {
		Name       string           `json:"name"`
		Machine    string           `json:"machine"`
		States     *json.RawMessage `json:"states"`
		State      string           `json:"state"`
		Properties *json.RawMessage `json:"properties"`
	}{m.Name(), m.Kind(), (*json.RawMessage)(&states), m.state.Name(), (*json.RawMessage)(&properties)})
}

func (m *Machine) MarshalPropertiesJSON() ([]byte, error) {
	buf := bytes.Buffer{}
	buf.WriteString("{")
	for i, property := range m.Properties() {
		if i != 0 {
			buf.WriteString(",")
		}
		key, _ := json.Marshal(property.Name())
		buf.Write(key)
		buf.WriteString(":")
		value, _ := json.Marshal(property)
		buf.Write(value)
	}
	buf.WriteString("}")
	return buf.Bytes(), nil
}

func (m *Machine) MarshalStatesJSON() ([]byte, error) {
	buf := bytes.Buffer{}
	buf.WriteString("{")
	for i, state := range m.States() {
		if i != 0 {
			buf.WriteString(",")
		}
		key, _ := json.Marshal(state.Name())
		buf.Write(key)
		buf.WriteString(":")
		value, _ := json.Marshal(state.ToJSON())
		buf.Write(value)
	}
	buf.WriteString("}")
	return buf.Bytes(), nil
}
