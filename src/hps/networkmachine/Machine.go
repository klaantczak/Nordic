package networkmachine

import (
	"bytes"
	"encoding/json"
	"fmt"
	"hps"
)

type Machine struct {
	name       string
	kind       string
	properties map[string]*hps.Property
	machines   []hps.IMachine
}

func NewMachine(name string, kind string) *Machine {
	m := &Machine{}
	m.name = name
	m.kind = kind
	m.properties = map[string]*hps.Property{}
	m.machines = []hps.IMachine{}
	_ = hps.IMachine(m)
	_ = INetworkMachine(m)
	return m
}

func (m *Machine) Name() string {
	return m.name
}

func (m *Machine) Kind() string {
	return m.kind
}

func (m *Machine) Event(time float64) (hps.IEvent, bool) {
	var event hps.IEvent = nil

	for _, mm := range m.machines {
		if e, ok := mm.Event(time); ok && (event == nil || event.Time() > e.Time()) {
			event = e
		}
	}

	return event, event != nil
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

func (m *Machine) AddMachine(machine hps.IMachine) error {
	if _, ok := m.GetMachine(machine.Name()); ok {
		return fmt.Errorf("A machine named \"%s\" already exists in this network.")
	}
	m.machines = append(m.machines, machine)
	return nil
}

func (m *Machine) GetMachine(name string) (hps.IMachine, bool) {
	for _, sm := range m.machines {
		if sm.Name() == name {
			return sm, true
		}
	}
	return nil, false
}

func (m *Machine) Machines() []hps.IMachine {
	list := []hps.IMachine{}
	for _, v := range m.machines {
		list = append(list, v)
	}
	return list
}

func (m *Machine) MarshalJSON() ([]byte, error) {
	machines, _ := m.MarshalMachinesJSON()
	content, _ := json.Marshal(struct {
		Machines *json.RawMessage `json:"machines"`
	}{(*json.RawMessage)(&machines)})

	properties, _ := m.MarshalPropertiesJSON()

	return json.Marshal(struct {
		Name       string           `json:"name"`
		Machine    string           `json:"machine"`
		Content    *json.RawMessage `json:"content"`
		Properties *json.RawMessage `json:"properties"`
	}{m.Name(), m.Kind(), (*json.RawMessage)(&content), (*json.RawMessage)(&properties)})
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

func (m *Machine) MarshalMachinesJSON() ([]byte, error) {
	buf := bytes.Buffer{}
	buf.WriteString("{")
	for i, machine := range m.Machines() {
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
