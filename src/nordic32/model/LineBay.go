package model

import (
	"hps"
	nm "hps/networkmachine"
)

type LineBay struct {
	*nm.Machine
	line *hps.Property
}

func NewLineBay(machine *nm.Machine) *LineBay {
	line, _ := machine.Property("line")
	return &LineBay{machine, line}
}

func (lb *LineBay) Line() string {
	value, _ := lb.line.GetString()
	return value
}

func (lb *LineBay) Breaker() (hps.IMachine, bool) {
	return lb.findMachineByName("Breaker")
}

func (lb *LineBay) PrimaryRelay() (hps.IMachine, bool) {
	return lb.findMachineByName("Primary Relay")
}

func (lb *LineBay) BackupRelay() (hps.IMachine, bool) {
	return lb.findMachineByName("Backup Relay")
}

func (lb *LineBay) Battery() (hps.IMachine, bool) {
	return lb.findMachineByName("Battery")
}

func (lb *LineBay) Wiring() (hps.IMachine, bool) {
	return lb.findMachineByName("Wiring")
}

func (lb *LineBay) findMachineByName(name string) (hps.IMachine, bool) {
	for _, m := range lb.Machines() {
		if m.Name() == name {
			return m, true
		}
	}
	return nil, false
}
