package model

import (
	"hps"
	nm "hps/networkmachine"
	sm "hps/statemachine"
)

type GeneratorBay struct {
	*nm.Machine
	capacity  *hps.Property
	connected *hps.Property
}

func NewGeneratorBay(machine *nm.Machine) *GeneratorBay {
	connected, _ := machine.Property("connected")
	capacity, _ := machine.Property("capacity")
	return &GeneratorBay{machine, capacity, connected}
}

func (gb *GeneratorBay) Generator() (*sm.Machine, bool) {
	for _, m := range gb.Machines() {
		if m.Kind() == "Generator" {
			return m.(*sm.Machine), true
		}
	}
	return nil, false
}

func (gb *GeneratorBay) Breaker() (hps.IMachine, bool) {
	return gb.findMachineByName("Breaker")
}

func (gb *GeneratorBay) Capacity() float64 {
	value, _ := gb.capacity.GetFloat()
	return value
}

func (gb *GeneratorBay) Connected() bool {
	value, _ := gb.connected.GetBool()
	return value
}

func (gb *GeneratorBay) Disconnected() bool {
	value, _ := gb.connected.GetBool()
	return !value
}

func (gb *GeneratorBay) PropertyConnected() *hps.Property {
	return gb.connected
}

func (gb *GeneratorBay) Disconnect() {
	gb.connected.SetValue(false)
}

func (gb *GeneratorBay) Connect() {
	gb.connected.SetValue(true)
}

func (gb *GeneratorBay) findMachineByName(name string) (hps.IMachine, bool) {
	for _, m := range gb.Machines() {
		if m.Name() == name {
			return m, true
		}
	}
	return nil, false
}
