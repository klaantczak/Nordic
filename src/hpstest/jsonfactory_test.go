package hpstest

import (
	jf "hps/jsonfactory"
	sm "hps/statemachine"
	"hps/statemachine/triggers"
	"hps/tools"
	"math"
	"testing"
)

func Test_LoadDeterministicTriggers(t *testing.T) {
	factory := jf.NewFactory(nil, nil)

	err := factory.LoadJSON([]byte(`{
    "machines" : [
      {
        "name" : "test",
        "type" : "state-machine",
        "structure" : {
          "states": [ "s1", "s2", "s3" ],
          "initial": "s1",
          "transitions": {
            "s1": { "s2": [ { "type": "deterministic", "parameter": 1 } ] },
            "s2": { "s3": [ { "type": "deterministic", "parameter": 1.5 } ] },
            "s3": { "s1": [ { "type": "deterministic", "parameter": 0.019230769230769232 } ] }
          }
        }
      }
    ]
  }`))

	if err != nil {
		t.Errorf("cannot load state machines into factory, %v", err)
		return
	}

	m, err := factory.Create("test", "test")
	if err != nil {
		t.Errorf("error: %v", err)
		return
	}

	stm, _ := m.(*sm.Machine)

	if stm.State().Name() != "s1" {
		t.Errorf("Expecting state 's1' but it is '%v'", stm.State().Name())
		return
	}

	expected := 0.0
	actual := 0.0

	e, _ := stm.Event(actual)
	actual = e.Time()
	expected += 1
	if expected != actual {
		t.Errorf("The event time should be %v but it is %v", expected, actual)
		return
	}

	e.Invoke()

	if stm.State().Name() != "s2" {
		t.Errorf("Expecting state 's2' but it is '%v'", stm.State().Name())
		return
	}

	e, _ = stm.Event(actual)
	actual = e.Time()
	expected += 1.5
	if expected != actual {
		t.Errorf("The event time should be %v but it is %v", expected, actual)
	}

	e.Invoke()

	oneWeek, err := tools.ParseDuration("1 week")
	if err != nil {
		t.Errorf("Error: %v", err)
		return
	}

	if stm.State().Name() != "s3" {
		t.Errorf("Expecting state 's3' but it is '%v'", stm.State().Name())
		return
	}

	e, _ = stm.Event(actual)
	actual = e.Time()
	expected += oneWeek
	if expected != actual {
		t.Errorf("The event time should be %v but it is %v", expected, actual)
	}
}

func Test_LoadStateMachines(t *testing.T) {
	factory := jf.NewFactory(nil, nil)
	err := factory.LoadJSON([]byte(`{
    "machines" : [
      {
        "name" : "deterministic",
        "type" : "state-machine",
        "structure" : {
          "states": [ "ok", "fail" ],
          "initial": "ok",
          "transitions": {
            "ok": { "fail": [ { "type": "deterministic", "parameter": 1 } ] },
            "fail": { "ok": [ { "type": "deterministic", "parameter": 1 } ] }
          }
        }
      },
      {
        "name" : "probabilistic",
        "type" : "state-machine",
        "structure" : {
          "states": [ "ok", "fail" ],
          "initial": "ok",
          "transitions": {
            "ok": { "fail": [ { "type": "probabilistic",
              "distribution": "exponential", "parameter": 0.1 } ] },
            "fail": { "ok": [ { "type": "probabilistic",
              "distribution": "exponential", "parameter": 0.1 } ] }
          }
        }
      },
      {
        "name" : "parameterised",
        "type" : "state-machine",
        "properties" : { "trigger" : { "type" : "ITrigger" } },
        "structure" : {
          "states": [ "ok", "fail" ],
          "initial": "ok",
          "transitions": {
            "ok": { "fail": [ { "type": "property", "property": "trigger" } ] },
            "fail": { "ok": [ { "type": "property", "property": "trigger" } ] }
          }
        }
      }
    ]
  }`))
	if err != nil {
		t.Errorf("Error: %v", err)
		return
	}

	// test deterministic state machine
	{
		m, _ := factory.Create("deterministic", "test")
		dm, _ := m.(*sm.Machine)

		if dm.State().Name() != "ok" {
			t.Errorf("Initial state should be 'ok'.")
		}

		e, _ := dm.Event(0.0)
		dmt := e.Invoke()

		e, _ = dm.Event(dmt)
		dmt = e.Invoke()

		if dmt != 2 {
			t.Errorf("Expeting current time of deterministic machine to be 2, not %v", dmt)
		}
	}

	// test probabilistic state machine
	{
		m, _ := factory.Create("probabilistic", "test")
		pbm, _ := m.(*sm.Machine)

		if pbm.State().Name() != "ok" {
			t.Errorf("Initial state should be 'ok'.")
		}

		pbmt := 0.0
		for i := 0; i < 1000; i++ {
			e, _ := pbm.Event(pbmt)
			pbmt += e.Invoke()
		}
		// average time is 0.1 +/- 10%
		avg := pbmt / 1000.0
		if math.Abs(avg-0.1) < 0.01 {
			t.Errorf("Expeting average transition time to be 0.1(+/-0.01), not %v", avg)
		}
	}

	// test parameterised state machine
	{
		m, _ := factory.Create("parameterised", "test")
		prm, _ := m.(*sm.Machine)

		trigger, ok := prm.Property("trigger")
		if !ok {
			t.Errorf("Machine should have the property 'Trigger'")
			return
		}
		trigger.SetValue(triggers.NewDeterministicTrigger(1))

		e, _ := m.Event(0.0)
		et := e.Invoke()

		e, _ = m.Event(et)
		prmt := e.Invoke()

		if prmt != 2 {
			t.Errorf("Expeting current time of parameterised machine to be 2, not %v", prmt)
		}
	}
}

func Test_LoadNetwork(t *testing.T) {
	factory := jf.NewFactory(nil, nil)
	err := factory.LoadJSON([]byte(`{
    "networks" : [
      {
        "name" : "sample",
        "machines" : [
          { "name" : "t1", "machine" : "test", "properties" : { "name" : "t1" } },
          { "name" : "t2", "machine" : "test", "properties" : { "name" : "t2" } },
          { "name" : "t3", "machine" : "test", "properties" : { "name" : "t3" } }
        ]
      }
    ],
    "machines" : [
      {
        "name" : "test",
        "type" : "state-machine",
        "properties" : { "name" : { "type" : "string" } },
        "structure" : {
          "states": [ "ok", "fail" ],
          "initial": "ok",
          "transitions": {
            "ok": { "fail": [ { "type": "deterministic", "parameter": 1 } ] },
            "fail": { "ok": [ { "type": "deterministic", "parameter": 1 } ] }
          }
        }
      }
    ]
  }`))
	if err != nil {
		t.Errorf("Error: %v", err)
		return
	}

	ms, err := factory.CreateNetwork("sample")
	if err != nil {
		t.Errorf("Error: %v", err)
		return
	}

	if len(ms) != 3 {
		t.Errorf("Expecting 3 machines in the network, but it has %v", len(ms))
		return
	}

	property, _ := ms[0].Property("name")
	value, _ := property.Value().(string)
	if value != "t1" {
		t.Errorf("Expecting property value to be 't1', but it is %v", value)
		return
	}
}
