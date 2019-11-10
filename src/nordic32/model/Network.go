package model

import (
	nm "hps/networkmachine"
	"strings"
)

type Network struct {
	*nm.Machine
}

func NewNetwork(machine *nm.Machine) *Network {
	return &Network{machine}
}

func (n *Network) Links() []*Link {
	list := []*Link{}
	for _, m := range n.Machines() {
		if link, ok := m.(*Link); ok {
			list = append(list, link)
		}
	}
	return list
}

func (n *Network) Link(name string) (*Link, bool) {
	for _, m := range n.Machines() {
		if link, ok := m.(*Link); ok && m.Name() == name {
			return link, true
		}
	}
	return nil, false
}

func (n *Network) Substation(name string) (*Substation, bool) {
	for _, m := range n.Machines() {
		if strings.HasPrefix(m.Kind(), "Substation ") && m.Name() == name {
			return m.(*Substation), true
		}
	}
	return nil, false
}

func (n *Network) Substations() []*Substation {
	list := []*Substation{}
	for _, m := range n.Machines() {
		if strings.HasPrefix(m.Kind(), "Substation ") {
			list = append(list, m.(*Substation))
		}
	}
	return list
}

func (n *Network) SubstationsWithGenerators() []*Substation {
	list := []*Substation{}
	for _, m := range n.Substations() {
		if m.HasGeneratorBays() {
			list = append(list, m)
		}
	}
	return list
}

func (n *Network) SubstationsWithLoads() []*Substation {
	list := []*Substation{}
	for _, m := range n.Substations() {
		if m.HasLoadBays() {
			list = append(list, m)
		}
	}
	return list
}
