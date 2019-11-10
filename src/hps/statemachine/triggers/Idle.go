package triggers

import (
	"encoding/json"
	"fmt"
	sm "hps/statemachine"
	"math"
)

type IdleTrigger struct {
	comment string
}

func NewIdleTrigger() *IdleTrigger {
	t := &IdleTrigger{}
	_ = sm.ITrigger(t)
	return t
}

func (t *IdleTrigger) Time(time float64) float64 {
	return math.Inf(1)
}

func (t *IdleTrigger) GetComment() string {
	return t.comment
}

func (t *IdleTrigger) SetComment(comment string) {
	t.comment = comment
}

func (t *IdleTrigger) SetMachine(machine *sm.Machine) {
}

func (t *IdleTrigger) UnmarshalJSON(data []byte) error {
	item := struct {
		Type    string `json:"type"`
		Comment string `json:"comment,omitempty"`
	}{}
	err := json.Unmarshal(data, item)
	if err != nil {
		return nil
	}
	if item.Type != "idle" {
		return fmt.Errorf("Expect transition of the idle type.")
	}
	*t = IdleTrigger{}
	t.comment = item.Comment
	return nil
}

func (t *IdleTrigger) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type    string `json:"type"`
		Comment string `json:"comment,omitempty"`
	}{"idle", t.comment})
}
