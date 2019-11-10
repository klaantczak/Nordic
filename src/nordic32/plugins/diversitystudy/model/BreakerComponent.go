package model

import (
	sm "hps/statemachine"
)

type BreakerComponent struct {
	m *sm.Machine
}

func NewBreakerComponent(m *sm.Machine) *BreakerComponent {
	return &BreakerComponent{m}
}

func (bc *BreakerComponent) StateCompromisedFail() *sm.State {
	s, _ := bc.m.GetState("compromised-fail")
	return s
}

func (bc *BreakerComponent) StateCompromisedOk() *sm.State {
	s, _ := bc.m.GetState("compromised-ok")
	return s
}

func (bc *BreakerComponent) StateFail() *sm.State {
	s, _ := bc.m.GetState("fail")
	return s
}

func (bc *BreakerComponent) StateOk() *sm.State {
	s, _ := bc.m.GetState("ok")
	return s
}

func (bc *BreakerComponent) Failed() bool {
	s := bc.m.State().Name()
	return s == "fail" || s == "compromised-fail"
}

func (bc *BreakerComponent) Restore(time float64) {
	bc.m.SetState("ok", time)
}

func (bc *BreakerComponent) Compromise(time float64) {
	if bc.Failed() {
		bc.m.SetState("compromised-fail", time)
	} else {
		bc.m.SetState("compromised-ok", time)
	}
}
