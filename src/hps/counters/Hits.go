package counters

import (
	sm "hps/statemachine"
)

// Hits contains counter and provides handlers that increases it for various HPS
// events.
type Hits struct {
	Value int
}

// StateHandler returns StateMachine State change handler that increases this
// hits counter.
func (h *Hits) StateHandler() sm.StateEventHandler {
	return func(s *sm.State, time float64, duration float64) {
		h.Value++
	}
}
