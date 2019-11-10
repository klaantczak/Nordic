package statemachine

import (
	"encoding/json"
	"fmt"
	"hps"
)

type ITrigger interface {
	Time(float64) float64
	SetMachine(*Machine)
	SetComment(string)
}

type TriggerJSONUnmarshaller func([]byte) (ITrigger, error)

var TriggerJSONUnmarshallers = map[string]TriggerJSONUnmarshaller{}

func RegisterTriggerJSONUnmarshaller(triggerType string, unmarshaller TriggerJSONUnmarshaller) {
	TriggerJSONUnmarshallers[triggerType] = unmarshaller
}

func init() {
	hps.RegisterPropertyValueJSONUnmarshaller("Trigger", func(data []byte) (interface{}, error) {
		trigger, err := UnmarshalTrigger(data)
		return interface{}(trigger), err
	})
}

func UnmarshalTrigger(data []byte) (ITrigger, error) {
	properties := map[string]interface{}{}
	err := json.Unmarshal(data, &properties)
	if err != nil {
		return nil, err
	}

	triggerTypeProperty, ok := properties["type"]
	if !ok {
		return nil, fmt.Errorf("Expecting the \"type\" property in the trigger's JSON.")
	}

	triggerType, ok := triggerTypeProperty.(string)
	if !ok {
		return nil, fmt.Errorf("Expecting the \"type\" property to be string.")
	}

	unmarshaller, ok := TriggerJSONUnmarshallers[triggerType]
	if !ok {
		return nil, fmt.Errorf("No marshaller is registered for this trigger type.")
	}

	return unmarshaller(data)
}
