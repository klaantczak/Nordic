package triggers

import (
	"encoding/json"
	"fmt"
	sm "hps/statemachine"
)

type DeterministicTrigger struct {
	parameter float64
	comment   string
}

func NewDeterministicTrigger(parameter float64) *DeterministicTrigger {
	t := &DeterministicTrigger{}
	t.parameter = parameter
	_ = sm.ITrigger(t)
	return t
}

func (t *DeterministicTrigger) Time(time float64) float64 {
	return time + t.parameter
}

func (t *DeterministicTrigger) GetComment() string {
	return t.comment
}

func (t *DeterministicTrigger) SetComment(comment string) {
	t.comment = comment
}

func (t *DeterministicTrigger) GetParameter() float64 {
	return t.parameter
}

func (t *DeterministicTrigger) SetMachine(machine *sm.Machine) {
}

func (t *DeterministicTrigger) UnmarshalJSON(data []byte) error {
	trigger := struct {
		Type      string  `json:"type"`
		Parameter float64 `json:"parameter"`
		Comment   string  `json:"comment,omitempty"`
	}{"", 0.0, ""}
	err := json.Unmarshal(data, &trigger)
	if err != nil {
		return err
	}
	if trigger.Type != "deterministic" {
		return fmt.Errorf("Expect the trigger type to be \"deterministic\".")
	}

	*t = DeterministicTrigger{}
	t.parameter = trigger.Parameter
	t.comment = trigger.Comment
	return nil
}

func (t *DeterministicTrigger) MarshalJSON() ([]byte, error) {
	trigger := struct {
		Type      string  `json:"type"`
		Parameter float64 `json:"parameter"`
		Comment   string  `json:"comment,omitempty"`
	}{"deterministic", t.parameter, t.comment}
	return json.Marshal(trigger)
}
