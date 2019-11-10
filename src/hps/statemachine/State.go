package statemachine

type StateEventHandler func(p *State, time float64, duration float64)

type State struct {
	machine     *Machine
	name        string
	transitions []*Transition
	entered     float64
	entering    []StateEventHandler
	leaving     []StateEventHandler
}

func (s *State) Name() string {
	return s.name
}

func (s *State) Transitions() []*Transition {
	return s.transitions
}

func (s *State) AddTransition(to string, trigger ITrigger) {
	t := &Transition{}
	t.to = to
	t.trigger = trigger
	s.transitions = append(s.transitions, t)
	trigger.SetMachine(s.machine)
}

func (s *State) Entering(handler StateEventHandler) *State {
	s.entering = append(s.entering, handler)
	return s
}

func (s *State) Leaving(handler StateEventHandler) *State {
	s.leaving = append(s.leaving, handler)
	return s
}

func (s *State) ToJSON() map[string]interface{} {
	transitions := map[string]interface{}{}
	for _, transition := range s.transitions {
		transitions[transition.to] = transition.ToJSON()
	}

	return map[string]interface{}{
		"name":        s.Name(),
		"transitions": transitions,
	}
}
