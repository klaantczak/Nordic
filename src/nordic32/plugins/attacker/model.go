package attacker

import (
	sm "hps/statemachine"
)

type Attacker struct {
	m                   *sm.Machine
	disconnectLoad      *sm.State
	disconnectGenerator *sm.State
	disconnectLine      *sm.State
}

func ToAttacker(machine *sm.Machine) (*Attacker, bool) {
	disconnectLoad, ok := machine.GetState("disconnectLoad")
	if !ok {
		return nil, false
	}
	disconnectGenerator, ok := machine.GetState("disconnectGenerator")
	if !ok {
		return nil, false
	}
	disconnectLine, ok := machine.GetState("disconnectLine")
	if !ok {
		return nil, false
	}
	return &Attacker{
		machine,
		disconnectLoad,
		disconnectGenerator,
		disconnectLine,
	}, true
}

func (a *Attacker) StateDisconnectLoad() *sm.State {
	return a.disconnectLoad
}

func (a *Attacker) StateDisconnectGenerator() *sm.State {
	return a.disconnectGenerator
}

func (a *Attacker) StateDisconnectLine() *sm.State {
	return a.disconnectLine
}
