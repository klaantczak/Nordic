package hpstest

import (
	"encoding/json"
	"hps"
	"hps/counters"
	"hps/engine"
	sm "hps/statemachine"
	"testing"

	"hps/loggers"

	"github.com/stretchr/testify/assert"
)

func Test_RunSimulation(t *testing.T) {
	var m sm.Machine
	err := json.Unmarshal([]byte(`{
		"name" : "Example",
		"states" : [ "idle", "working" ],
		"transitions" : {
			"idle" : { "working" : [ { "type" : "deterministic", "parameter" : 1 } ] },
			"working" : { "idle" : [ { "type" : "deterministic", "parameter" : 1 } ] }
		}
	}`), &m)
	assert.Nil(t, err)

	e := engine.NewEnvironment(loggers.NewConsoleLogger())
	e.AddMachine(hps.IMachine(&m))

	hits := counters.Hits{}
	working, _ := m.GetState("working")
	working.Entering(hits.StateHandler())

	limits := engine.NewLimits()
	limits.Time = 1.0
	e.Run(limits)

	assert.NotEqual(t, 0, hits.Value)
}

func Test_SimpleGenerator(t *testing.T) {
	var m sm.Machine
	err := json.Unmarshal([]byte(`{
		"name" : "test",
		"machine" : "default",
		"states" : [ "ok" ],
		"transitions" : {
			"ok" : { "ok" : [ { "type" : "deterministic", "parameter" : 2 } ] }
		}
	}`), &m)
	assert.Nil(t, err)

	e := engine.NewEnvironment(loggers.NewConsoleLogger())
	e.AddMachine(hps.IMachine(&m))

	duration := counters.Duration{}
	ok, _ := m.GetState("ok")
	ok.Leaving(duration.LeavingStateHandler())

	limits := engine.NewLimits()
	limits.Events = 10
	e.Run(limits)

	assert.True(t, duration.Value > 10)
}
