package hpstest

import (
	"encoding/json"
	"fmt"
	sm "hps/statemachine"
	"hps/statemachine/triggers"
	"hps/tools"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_StateMachine_UnmarshalError(t *testing.T) {
	var m sm.Machine
	assert.NotNil(t, json.Unmarshal([]byte(``), &m))
	assert.NotNil(t, json.Unmarshal([]byte(`{}`), &m))
	assert.NotNil(t, json.Unmarshal([]byte(`[]`), &m))
	assert.NotNil(t, json.Unmarshal([]byte(`""`), &m))
	assert.NotNil(t, json.Unmarshal([]byte(`1`), &m))
	assert.NotNil(t, json.Unmarshal([]byte(`{"name":"ok"}`), &m))
	assert.NotNil(t, json.Unmarshal([]byte(`{"states":["ok"]}`), &m))
	assert.NotNil(t, json.Unmarshal([]byte(`{"name":"ok","states"[]}`), &m))
}

func Test_StateMachine_UnmarshalOk(t *testing.T) {
	var m sm.Machine
	assert.Nil(t, json.Unmarshal([]byte(`{"name":"ok","states":["ok"]}`), &m))
}

func Test_DeterministicStateMachine(t *testing.T) {
	var m sm.Machine
	err := json.Unmarshal([]byte(`{
		"name" : "test",
		"machine" : "default",
		"states" : [ "ok", "fail" ],
		"transitions" : {
			"ok" : { "fail" : [ { "type" : "deterministic", "parameter" : 1 } ] },
			"fail" : { "ok" : [ { "type" : "deterministic", "parameter" : 1 } ] }
		}
	}`), &m)
	assert.Nil(t, err)

	ok, _ := m.GetState("ok")
	fail, _ := m.GetState("fail")

	assert.Equal(t, 2, len(m.States()), "Machine should have two states.")
	assert.Equal(t, "ok", m.State().Name(), "Initial state should be 'ok'.")
	assert.Equal(t, 1, len(ok.Transitions()), "Transition from 'ok' should be defined.")
	assert.Equal(t, "fail", ok.Transitions()[0].To(), "Transition from 'ok' to 'fail' should be defined.")
	assert.Equal(t, 1, len(fail.Transitions()), "Transition from 'fail' should be defined.")
	assert.Equal(t, "ok", fail.Transitions()[0].To(), "Transition from 'fail' to 'ok' should be defined.")

	time := 0.0

	event, _ := m.Event(time)
	assert.NotNil(t, event)

	smEvent, _ := event.(*sm.Event)
	assert.Equal(t, "ok", smEvent.From().Name())
	assert.Equal(t, "fail", smEvent.To().Name())

	time = event.Invoke()
	assert.Equal(t, "fail", m.State().Name())

	event, _ = m.Event(time)
	assert.NotNil(t, event)

	smEvent, _ = event.(*sm.Event)
	assert.Equal(t, "fail", smEvent.From().Name())
	assert.Equal(t, "ok", smEvent.To().Name())
}

func Test_StateMachineWithParameter(t *testing.T) {
	var m sm.Machine
	err := json.Unmarshal([]byte(`{
		"name" : "test",
		"machine" : "default",
		"states" : [ "ok", "fail" ],
		"properties" : {
			"failure_rate" : {
				"type" : "Trigger",
				"value" : { "type" : "deterministic", "parameter" : 2 }
			}
		},
		"transitions" : {
			"ok" : { "fail" : [ { "type" : "property", "property" : "failure_rate" } ] },
			"fail" : { "ok" : [ { "type" : "deterministic", "parameter" : 1 } ] }
		}
	}`), &m)
	assert.Nil(t, err)

	event, ok := m.Event(0.0)
	assert.Equal(t, true, ok)
	assert.NotNil(t, event)
}

func Test_StateMachineWithParameterToJson(t *testing.T) {
	var m sm.Machine
	err := json.Unmarshal([]byte(`{
		"name" : "test",
		"machine" : "default",
		"states" : [ "ok", "fail" ],
		"properties" : { "failure_rate" : { "type" : "ITrigger" } },
		"transitions" : {
			"ok" : { "fail" : [ { "type" : "property", "property" : "failure_rate" } ] },
			"fail" : { "ok" : [ { "type" : "deterministic", "parameter" : 2 } ] }
		}
	}`), &m)
	assert.Nil(t, err)

	duration2h, _ := tools.ParseDuration("2h")
	property, _ := m.Property("failure_rate")
	property.SetValue(triggers.NewDeterministicTrigger(duration2h))

	jsonBin, _ := json.MarshalIndent(&m, "", " ")
	jsonStr := string(jsonBin)
	expected := strings.Join([]string{
		"{",
		" \"name\": \"test\",",
		" \"machine\": \"default\",",
		" \"states\": {",
		"  \"ok\": {",
		"   \"name\": \"ok\",",
		"   \"transitions\": {",
		"    \"fail\": {",
		"     \"to\": \"fail\"",
		"    }",
		"   }",
		"  },",
		"  \"fail\": {",
		"   \"name\": \"fail\",",
		"   \"transitions\": {",
		"    \"ok\": {",
		"     \"to\": \"ok\"",
		"    }",
		"   }",
		"  }",
		" },",
		" \"state\": \"ok\",",
		" \"properties\": {",
		"  \"failure_rate\": {",
		"   \"name\": \"failure_rate\",",
		"   \"value\": {",
		"    \"type\": \"deterministic\",",
		"    \"parameter\": 0.00022831050228310502",
		"   },",
		"   \"required\": true",
		"  }",
		" }",
		"}",
	}, "\n")
	if !assert.JSONEq(t, expected, jsonStr) {
		fmt.Println(jsonStr)
	}
}
