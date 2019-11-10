package model

import (
	nm "hps/networkmachine"
)

type Substation struct {
	*nm.Machine
}

func NewSubstation(machine *nm.Machine) *Substation {
	return &Substation{machine}
}

func (s *Substation) HasGeneratorBays() bool {
	for _, b := range s.Machines() {
		if b.Kind() == "Generator Bay" {
			return true
		}
	}
	return false
}

func (s *Substation) HasLoadBays() bool {
	for _, b := range s.Machines() {
		if b.Kind() == "Load Bay" {
			return true
		}
	}
	return false
}

func (s *Substation) GeneratorBays() []*GeneratorBay {
	list := []*GeneratorBay{}
	for _, b := range s.Machines() {
		if b.Kind() == "Generator Bay" {
			list = append(list, b.(*GeneratorBay))
		}
	}
	return list
}

func (s *Substation) LoadBays() []*LoadBay {
	list := []*LoadBay{}
	for _, b := range s.Machines() {
		if b.Kind() == "Load Bay" {
			list = append(list, b.(*LoadBay))
		}
	}
	return list
}

func (s *Substation) LineBay(line string) (*LineBay, bool) {
	for _, b := range s.Machines() {
		if b.Kind() == "Line Bay" || b.Kind() == "Transformer Bay" {
			bay := b.(*LineBay)
			if bay.Line() == line {
				return bay, true
			}
		}
	}
	return nil, false
}

func (s *Substation) LineBays() []*LineBay {
	list := []*LineBay{}
	for _, b := range s.Machines() {
		if b.Kind() == "Line Bay" || b.Kind() == "Transformer Bay" {
			list = append(list, b.(*LineBay))
		}
	}
	return list
}

func (s *Substation) CurrentGeneration() float64 {
	sum := 0.0
	for _, b := range s.WorkingGeneratorBays() {
		sum += b.Capacity()
	}
	return sum
}

func (s *Substation) CurrentLoad() float64 {
	sum := 0.0
	for _, b := range s.WorkingLoadBays() {
		sum += b.Power()
	}
	return sum
}

func (s *Substation) WorkingGeneratorBays() []*GeneratorBay {
	list := []*GeneratorBay{}
	for _, gb := range s.GeneratorBays() {
		if !gb.Connected() {
			continue
		}

		generator, ok := gb.Generator()
		if !ok || generator.State().Name() != "ok" {
			continue
		}

		list = append(list, gb)
	}

	return list
}

func (s *Substation) WorkingLoadBays() []*LoadBay {
	list := []*LoadBay{}
	for _, lb := range s.LoadBays() {
		if !lb.Connected() {
			continue
		}

		load, ok := lb.Load()
		if !ok || load.State().Name() != "ok" {
			continue
		}

		list = append(list, lb)
	}

	return list
}
