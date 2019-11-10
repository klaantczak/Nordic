package triggers

import (
	"encoding/json"
	"fmt"
	"hps"
	sm "hps/statemachine"
	"math"
)

type ProbabilisticTrigger struct {
	parameter float64
	comment   string
}

func NewProbabilisticTrigger(parameter float64) *ProbabilisticTrigger {
	t := &ProbabilisticTrigger{}
	t.parameter = parameter
	_ = sm.ITrigger(t)
	return t
}

func (t *ProbabilisticTrigger) Time(time float64) float64 {
	return time + math.Log(1-hps.RndNext())/(-t.parameter)
}

func (t *ProbabilisticTrigger) GetComment() string {
	return t.comment
}

func (t *ProbabilisticTrigger) SetComment(comment string) {
	t.comment = comment
}

func (t *ProbabilisticTrigger) GetParameter() float64 {
	return t.parameter
}

func (t *ProbabilisticTrigger) SetMachine(machine *sm.Machine) {
}

func (t *ProbabilisticTrigger) UnmarshalJSON(data []byte) error {
	trigger := struct {
		Type         string  `json:"type"`
		Distribution string  `json:"distribution"`
		Parameter    float64 `json:"parameter"`
		Comment      string  `json:"comment,omitempty"`
	}{"", "", 0.0, ""}
	err := json.Unmarshal(data, &trigger)
	if err != nil {
		return err
	}
	if trigger.Type != "probabilistic" || trigger.Distribution != "exponential" {
		return fmt.Errorf("Expect the trigger type to be \"probabilistic\" and the distribution is \"exponential\".")
	}

	*t = ProbabilisticTrigger{}
	t.parameter = trigger.Parameter
	t.comment = trigger.Comment
	return nil
}

func (t *ProbabilisticTrigger) MarshalJSON() ([]byte, error) {
	trigger := struct {
		Type         string  `json:"type"`
		Distribution string  `json:"distribution"`
		Parameter    float64 `json:"parameter"`
		Comment      string  `json:"comment,omitempty"`
	}{"probabilistic", "exponential", t.parameter, t.comment}
	return json.Marshal(trigger)
}
