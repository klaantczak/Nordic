package model

import (
	nm "hps/networkmachine"
	sm "hps/statemachine"
	"strconv"
)

type Breaker struct {
	m          *nm.Machine
	components []*BreakerComponent
}

func NewBreaker(m *nm.Machine) *Breaker {
	components := []*BreakerComponent{}
	i := 1
	for {
		component, ok := m.GetMachine("Component" + strconv.Itoa(i))
		if !ok {
			break
		}
		components = append(components, NewBreakerComponent(component.(*sm.Machine)))
		i++
	}
	return &Breaker{m, components}
}

func (b *Breaker) Component(index int) *BreakerComponent {
	return b.components[index]
}

func (b *Breaker) Components() []*BreakerComponent {
	return b.components
}

func (b *Breaker) Restore(time float64) {
	for _, component := range b.components {
		component.Restore(time)
	}
}
