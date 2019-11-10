package triggers

import (
	"encoding/json"
	sm "hps/statemachine"
)

func init() {
	sm.RegisterTriggerJSONUnmarshaller("idle", func(data []byte) (sm.ITrigger, error) {
		var trigger IdleTrigger
		err := json.Unmarshal(data, &trigger)
		return &trigger, err
	})
	sm.RegisterTriggerJSONUnmarshaller("deterministic", func(data []byte) (sm.ITrigger, error) {
		var trigger DeterministicTrigger
		err := json.Unmarshal(data, &trigger)
		return &trigger, err
	})
	sm.RegisterTriggerJSONUnmarshaller("probabilistic", func(data []byte) (sm.ITrigger, error) {
		var trigger ProbabilisticTrigger
		err := json.Unmarshal(data, &trigger)
		return &trigger, err
	})
	sm.RegisterTriggerJSONUnmarshaller("property", func(data []byte) (sm.ITrigger, error) {
		var trigger PropertyTrigger
		err := json.Unmarshal(data, &trigger)
		return &trigger, err
	})
}
