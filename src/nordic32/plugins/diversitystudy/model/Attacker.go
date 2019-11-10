package model

import (
	sm "hps/statemachine"
)

type Attacker struct {
	m *sm.Machine
}

func NewAttacker(m *sm.Machine) *Attacker {
	return &Attacker{m}
}

func (a *Attacker) StateAttack() *sm.State {
	s, _ := a.m.GetState("attack")
	return s
}

func (a *Attacker) Target() string {
	p, _ := a.m.Property("target")
	v, _ := p.GetString()
	return v
}
