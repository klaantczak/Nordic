package statemachine

type Transition struct {
	to      string
	trigger ITrigger
}

func CreateTransition(to string, trigger ITrigger) *Transition {
	t := &Transition{}
	t.to = to
	t.trigger = trigger
	return t
}

func (t *Transition) To() string {
	return t.to
}

func (t *Transition) Trigger() ITrigger {
	return t.trigger
}

func (t *Transition) SetTrigger(trigger ITrigger) {
	t.trigger = trigger
}

func (t *Transition) ToJSON() map[string]interface{} {
	return map[string]interface{}{
		"to": t.to,
	}
}
