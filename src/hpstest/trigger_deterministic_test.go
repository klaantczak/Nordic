package hpstest

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"hps/statemachine/triggers"
	"testing"
)

func Test_DeterministicTrigger_Time(t *testing.T) {
	trigger := triggers.NewDeterministicTrigger(1.0)
	time := trigger.Time(2.0)
	assert.InEpsilon(t, time, 3.0, 0.001)
}

func Test_DeterministicTrigger_Marshal(t *testing.T) {
	var trigger *triggers.DeterministicTrigger
	var data []byte
	var err error

	// marshal parameter without comment
	trigger = triggers.NewDeterministicTrigger(0.123)
	data, err = json.Marshal(trigger)
	assert.Nil(t, err)
	assert.Equal(t, string(data), `{"type":"deterministic","parameter":0.123}`)

	// marshal parameter with comment
	trigger = triggers.NewDeterministicTrigger(0.123)
	trigger.SetComment("ok")
	data, err = json.Marshal(trigger)
	assert.Nil(t, err)
	assert.Equal(t, string(data), `{"type":"deterministic","parameter":0.123,"comment":"ok"}`)
}

func Test_DeterministicTrigger_Unmarshal(t *testing.T) {
	var trigger triggers.DeterministicTrigger
	var err error

	// unmarshal with type and parameter
	err = json.Unmarshal([]byte(`{"type":"deterministic","parameter":0.123}`), &trigger)
	assert.Nil(t, err)
	assert.InEpsilon(t, trigger.GetParameter(), 0.123, 0.001)

	// unmarshal with type, parameter and comment
	err = json.Unmarshal([]byte(`{"type":"deterministic","parameter":0.456,"comment":"ok"}`), &trigger)
	assert.Nil(t, err)
	assert.InEpsilon(t, trigger.GetParameter(), 0.456, 0.001)
	assert.Equal(t, trigger.GetComment(), "ok")
}
