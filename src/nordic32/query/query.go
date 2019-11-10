package query

import (
	"hps"
	nm "hps/networkmachine"
	sm "hps/statemachine"
	"nordic32/model"
	"strings"
)

func FindAllByKind(machines []hps.IMachine, kind string) []hps.IMachine {
	list := []hps.IMachine{}
	for _, b := range machines {
		if b.Kind() == kind {
			list = append(list, b)
		}
	}
	return list
}

func FindGeneratorBay(substation *nm.Machine) (*nm.Machine, bool) {
	for _, b := range substation.Machines() {
		if b.Kind() == "Generator Bay" {
			return b.(*nm.Machine), true
		}
	}
	return nil, false
}

func FindMachineByName(machines []hps.IMachine, name string) (hps.IMachine, bool) {
	for _, m := range machines {
		if m.Name() == name {
			return m, true
		}
	}
	return nil, false
}

func FindMachineByPath(machines []hps.IMachine, path string) (hps.IMachine, bool) {
	if path == "" {
		return nil, false
	}
	parts := strings.Split(path, "/")

	var ok bool

	machine := hps.IMachine(nil)
	context := machines

	last := len(parts) - 1
	for i, part := range parts {
		machine, ok = FindMachineByName(context, part)
		if !ok || i == last {
			break
		}

		var n nm.INetworkMachine
		n, ok = machine.(nm.INetworkMachine)
		if !ok {
			machine = nil
			break
		}

		context = n.Machines()
	}

	return machine, machine != nil
}

func FindMachinesByKind(machines []hps.IMachine, kind string) []hps.IMachine {
	result := []hps.IMachine{}
	for _, m := range machines {
		if m.Kind() == kind {
			result = append(result, m)
		}
	}
	return result
}

func FindSubstationsNetwork(e hps.IEnvironment) (*model.Network, bool) {
	if machine, ok := FindMachineByPath(e.Machines(), "Substations"); ok {
		return machine.(*model.Network), true
	}
	return nil, false
}

func IsOk(m *sm.Machine) bool {
	return m.State().Name() == "ok"
}

func GetState(m *sm.Machine, name string) *sm.State {
	s, _ := m.GetState(name)
	return s
}

func GetProp(m hps.IMachine, name string) *hps.Property {
	p, ok := m.Property(name)
	if !ok {
		panic("Missing property: " + name)
	}
	return p
}

func GetStrProp(m hps.IMachine, name string) string {
	p, ok := m.Property(name)
	if !ok {
		panic("Missing property: " + name)
	}
	v, ok := p.GetString()
	if !ok {
		panic("Expect the property '" + name + "' to be string")
	}
	return v
}

func GetFloatProp(m hps.IMachine, name string) float64 {
	p, ok := m.Property(name)
	if !ok {
		panic("Missing property: " + name)
	}
	v, ok := p.GetFloat()
	if !ok {
		panic("Expect the property '" + name + "' to be float")
	}
	return v
}

func GetBoolProp(m hps.IMachine, name string) bool {
	p, ok := m.Property(name)
	if !ok {
		panic("Missing property: " + name)
	}
	v, ok := p.GetBool()
	if !ok {
		panic("Expect the property '" + name + "' to be boolean")
	}
	return v
}

func GetTransitionTo(s *sm.State, state string) (*sm.Transition, bool) {
	for _, tr := range s.Transitions() {
		if tr.To() == state {
			return tr, true
		}
	}
	return nil, false
}
