package sa

import (
	"encoding/json"
	"fmt"
	"hps"
	"hps/engine"
	sm "hps/statemachine"
	"io/ioutil"
	"nordic32/plugins"
	q "nordic32/query"
)

// [
//    { machine : 'machine 1/machine 2/...',
//      property : 'param',
//      values : [ { value : 1, states : [ { machine : 'machine 1/...', state : 'ok' }, ... ] ] }
// ]

type SaState struct {
	Machine string `json:"machine"`
	State   string `json:"state"`
}

type SaValue struct {
	Value  float64   `json:"value"`
	States []SaState `json:"states"`
}

type SaRule struct {
	Machine  string    `json:"machine"`
	Property string    `json:"property"`
	Values   []SaValue `json:"values"`
}

type Plugin struct {
	file     string
	fileFlag *string
	e        hps.IEnvironment
}

func NewPlugin() plugins.IPlugin {
	return &Plugin{}
}

func (p *Plugin) Name() string {
	return "stochastic associations"
}

func (p *Plugin) Init(e hps.IEnvironment) error {
	p.e = e

	if p.fileFlag == nil {
		return nil
	}

	p.file = *p.fileFlag

	if p.file == "" {
		return nil
	}

	rules, err := ReadRulesFromFile(p.file)
	if err != nil {
		return err
	}

	ApplyRules(rules, e)

	return nil
}

func (p *Plugin) Done(r *engine.SimulationResult) {
}

// Reads the stochastic associations from string
func ReadRulesFromString(text string) ([]SaRule, error) {
	var rules []SaRule
	err := json.Unmarshal([]byte(text), &rules)
	return rules, err
}

func ReadRulesFromFile(path string) ([]SaRule, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return []SaRule{}, err
	}

	return ReadRulesFromString(string(content))
}

type SaMachineState struct {
	Machine *sm.Machine
	State   *sm.State
}

func ApplyRules(rules []SaRule, e hps.IEnvironment) error {
	for _, rule := range rules {
		machine, ok := q.FindMachineByPath(e.Machines(), rule.Machine)
		if !ok {
			return fmt.Errorf("Cannot find machine by path '%s'.", rule.Machine)
		}

		property, _ := machine.Property(rule.Property)

		trigger := func(p *hps.Property, value float64, mss []SaMachineState) {
			ok := true
			for _, ms := range mss {
				if ms.Machine.State() != ms.State {
					ok = false
					break
				}
			}
			if ok {
				property.SetValue(value)
			}
		}

		for _, value := range rule.Values {
			mss := []SaMachineState{}
			for _, st := range value.States {
				m, _ := q.FindMachineByPath(e.Machines(), st.Machine)
				state, _ := m.(*sm.Machine).GetState(st.State)
				mss = append(mss, SaMachineState{m.(*sm.Machine), state})
			}

			for _, ms := range mss {
				ms.State.Entering(func(p *sm.State, time float64, duration float64) {
					trigger(property, value.Value, mss)
				})
			}
		}
	}

	return nil
}
