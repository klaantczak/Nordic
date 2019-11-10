package model

import (
	"hps"
	nm "hps/networkmachine"
	sm "hps/statemachine"
)

type LoadBay struct {
	*nm.Machine
	power     *hps.Property
	connected *hps.Property
}

func NewLoadBay(machine *nm.Machine) *LoadBay {
	connected, _ := machine.Property("connected")
	power, _ := machine.Property("power")
	return &LoadBay{machine, power, connected}
}

func (lb *LoadBay) Load() (*sm.Machine, bool) {
	for _, mm := range lb.Machines() {
		if mm.Kind() == "Load" {
			return mm.(*sm.Machine), true
		}
	}
	return nil, false
}

func (lb *LoadBay) Breaker() (hps.IMachine, bool) {
	return lb.findMachineByName("Breaker")
}

func (lb *LoadBay) Power() float64 {
	value, _ := lb.power.GetFloat()
	return value
}

func (lb *LoadBay) Connected() bool {
	value, _ := lb.connected.GetBool()
	return value
}

func (lb *LoadBay) Disconnected() bool {
	value, _ := lb.connected.GetBool()
	return !value
}

func (lb *LoadBay) PropertyConnected() *hps.Property {
	return lb.connected
}

func (lb *LoadBay) Disconnect() {
	lb.connected.SetValue(false)
}

func (lb *LoadBay) Connect() {
	lb.connected.SetValue(true)
}

func (lb *LoadBay) findMachineByName(name string) (hps.IMachine, bool) {
	for _, m := range lb.Machines() {
		if m.Name() == name {
			return m, true
		}
	}
	return nil, false
}
