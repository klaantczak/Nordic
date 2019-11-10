package hpstest

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"hps/statemachine/triggers"
	"testing"
)

func Test_ProbabilisticTrigger_Create(t *testing.T) {
	_ = triggers.NewProbabilisticTrigger(1.0)
}

func Test_ProbabilisticTrigger_Marshal(t *testing.T) {
	var trigger *triggers.ProbabilisticTrigger
	var data []byte
	var err error

	// marshal parameter without comment
	trigger = triggers.NewProbabilisticTrigger(0.123)
	data, err = json.Marshal(trigger)
	assert.Nil(t, err)
	assert.Equal(t, string(data), `{"type":"probabilistic","distribution":"exponential","parameter":0.123}`)

	// marshal parameter with comment
	trigger = triggers.NewProbabilisticTrigger(0.123)
	trigger.SetComment("ok")
	data, err = json.Marshal(trigger)
	assert.Nil(t, err)
	assert.Equal(t, string(data), `{"type":"probabilistic","distribution":"exponential","parameter":0.123,"comment":"ok"}`)
}

func Test_ProbabilisticTrigger_Unmarshal(t *testing.T) {
	var trigger triggers.ProbabilisticTrigger
	var err error

	// unmarshal with type and parameter
	err = json.Unmarshal([]byte(`{"type":"probabilistic","distribution":"exponential","parameter":0.123}`), &trigger)
	assert.Nil(t, err)
	assert.InEpsilon(t, trigger.GetParameter(), 0.123, 0.001)

	// unmarshal with type, parameter and comment
	err = json.Unmarshal([]byte(`{"type":"probabilistic","distribution":"exponential","parameter":0.456,"comment":"ok"}`), &trigger)
	assert.Nil(t, err)
	assert.InEpsilon(t, trigger.GetParameter(), 0.456, 0.001)
	assert.Equal(t, trigger.GetComment(), "ok")
}
