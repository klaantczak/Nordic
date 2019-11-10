package model

import (
	"hps"
	nm "hps/networkmachine"
	sm "hps/statemachine"
	"strings"
)

func TypeFactory() hps.ITypeFactory {
	return func(m hps.IMachine) hps.IMachine {
		kind := m.Kind()
		if kind == "Generator Bay" {
			return NewGeneratorBay(m.(*nm.Machine))
		}
		if kind == "Line Bay" {
			return NewLineBay(m.(*nm.Machine))
		}
		if kind == "Link" {
			return NewLink(m.(*sm.Machine))
		}
		if kind == "Load Bay" {
			return NewLoadBay(m.(*nm.Machine))
		}
		if strings.HasPrefix(kind, "Substation ") {
			return NewSubstation(m.(*nm.Machine))
		}
		if kind == "Substations" {
			return NewNetwork(m.(*nm.Machine))
		}
		if kind == "Transformer Bay" {
			return NewLineBay(m.(*nm.Machine))
		}
		return m
	}
}
