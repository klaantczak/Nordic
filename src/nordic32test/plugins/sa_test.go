package plugins

import (
	"hps/engine"
	"hps/loggers"
	sm "hps/statemachine"
	"nordic32/plugins/sa"
	"testing"
)

func Test_ParserNoRules(t *testing.T) {
	rules, _ := sa.ReadRulesFromString("[]")
	if len(rules) != 0 {
		t.Errorf("Should parse no rules from empty array")
	}
}

func Test_ParseRules(t *testing.T) {
	rules, _ := sa.ReadRulesFromString(`[
  { "machine" : "test",
    "property" : "param",
    "values" : [
      { "value" : 1, "states" : [
        { "machine" : "test", "state" : "ok" },
        { "machine" : "test", "state" : "fail" }
      ] }
    ] }
	]`)
	if len(rules) != 1 {
		t.Errorf("Should parse one rule from the array")
	}
	if rules[0].Machine != "test" {
		t.Errorf("Expect rules[0].Machine == test")
	}
	if rules[0].Property != "param" {
		t.Errorf("Expect rules[0].Property == param")
	}
	if rules[0].Values[0].Value != 1.0 {
		t.Errorf("Expect rules[0].Values[0].Value == 1.0")
	}
	if rules[0].Values[0].States[0].Machine != "test" {
		t.Errorf("Expect rules[0].Values[0].States[0].Machine == test")
	}
}

func Test_ApplyRules(t *testing.T) {
	e := engine.NewEnvironment(loggers.NewConsoleLogger())

	tm1, _ := sm.NewMachine("tm1", "", []string{"ok"}, "ok")
	v := tm1.AddProperty("v", "float")
	e.AddMachine(tm1)

	tm2, _ := sm.NewMachine("tm2", "", []string{"ok", "fail"}, "ok")
	tm2.AddProperty("v", "float")
	e.AddMachine(tm2)

	rules := []sa.SaRule{
		sa.SaRule{"tm1", "v", []sa.SaValue{sa.SaValue{1.0, []sa.SaState{sa.SaState{"tm2", "ok"}}}}},
		sa.SaRule{"tm1", "v", []sa.SaValue{sa.SaValue{2.0, []sa.SaState{sa.SaState{"tm2", "fail"}}}}},
	}

	err := sa.ApplyRules(rules, e)
	if err != nil {
		t.Errorf("ERROR: %s", err)
	}

	tm2.SetState("fail", 0)
	if value, _ := v.GetFloat(); value != 2.0 {
		t.Errorf("Expect tm1.v == 2")
	}

	tm2.SetState("ok", 0)
	if value, _ := v.GetFloat(); value != 1.0 {
		t.Errorf("Expect tm1.v == 1")
	}
}
