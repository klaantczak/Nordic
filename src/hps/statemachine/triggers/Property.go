package triggers

import (
	"encoding/json"
	"fmt"
	sm "hps/statemachine"
)

type PropertyTrigger struct {
	property string
	machine  *sm.Machine
	comment  string
}

func NewPropertyTrigger(property string) *PropertyTrigger {
	t := &PropertyTrigger{}
	t.property = property
	_ = sm.ITrigger(t)
	return t
}

func (t *PropertyTrigger) GetProperty() string {
	return t.property
}

func (t *PropertyTrigger) Time(time float64) float64 {
	property, ok := t.machine.Property(t.property)
	if !ok {
		panic("Machine does not have the property \"" + t.property + "\"")
	}
	value := property.Value()
	trigger, ok := value.(sm.ITrigger)
	if !ok {
		panic("Expect property value to be a trigger.")
	}
	return trigger.Time(time)
}

func (t *PropertyTrigger) SetComment(comment string) {
	t.comment = comment
}

func (t *PropertyTrigger) SetMachine(machine *sm.Machine) {
	t.machine = machine
}

func (t *PropertyTrigger) UnmarshalJSON(data []byte) error {
	trigger := struct {
		Type     string `json:"type"`
		Property string `json:"property"`
		Comment  string `json:"comment,omitempty"`
	}{"", "", ""}
	err := json.Unmarshal(data, &trigger)
	if err != nil {
		return err
	}
	if trigger.Type != "property" {
		return fmt.Errorf("Expect the trigger type to be \"property\".")
	}

	*t = PropertyTrigger{}
	t.property = trigger.Property
	t.comment = trigger.Comment
	return nil
}

func (t *PropertyTrigger) MarshalJSON() ([]byte, error) {
	trigger := struct {
		Type     string `json:"type"`
		Property string `json:"property"`
		Comment  string `json:"comment,omitempty"`
	}{"property", t.property, t.comment}
	return json.Marshal(trigger)
}
