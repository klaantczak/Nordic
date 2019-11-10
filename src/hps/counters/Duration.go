package counters

import (
	sm "hps/statemachine"
)

type Duration struct {
	Value float64
}

func (d *Duration) LeavingStateHandler() sm.StateEventHandler {
	return func(s *sm.State, time float64, duration float64) {
		d.Value += duration
	}
}
