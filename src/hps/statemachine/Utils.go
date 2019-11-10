package statemachine

import (
	"hps"
)

func Cast(machines []hps.IMachine) ([]*Machine, bool) {
	list := []*Machine{}
	for _, m := range machines {
		if nm, ok := m.(*Machine); ok {
			list = append(list, nm)
		} else {
			return nil, false
		}
	}
	return list, true
}
