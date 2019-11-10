package plugins

import (
	"fmt"
	"hps"
	"hps/engine"
	sm "hps/statemachine"
	"nordic32/model"
	q "nordic32/query"
	"strings"
)

type PowerSubsystemChangesPlugin struct {
	logger hps.ILogger
}

func NewPowerSubsystemChangesPlugin(logger hps.ILogger) IPlugin {
	return &PowerSubsystemChangesPlugin{logger}
}

func (p *PowerSubsystemChangesPlugin) Name() string {
	return "PowerSubsystemChangesPlugin"
}

func (p *PowerSubsystemChangesPlugin) Init(e hps.IEnvironment) error {
	substations, _ := q.FindSubstationsNetwork(e)
	modifications := []string{}
	states := make([]string, 1000)
	length := trackPowerNetworkModifications(substations, func(i int, connected bool, ok bool, name string, verb string) {
		modifications = append(modifications, fmt.Sprintf("%s of %s", verb, name))
		state := ""
		if connected {
			if ok {
				state = "0 (ok)"
			} else {
				state = "1 (failure)"
			}
		} else {
			state = "X (disconnected)"
		}
		states[i] = fmt.Sprintf("%s=%s", name, state)
	})
	states = states[0:length]

	snapshot := ""
	e.StartingIteration(func(env hps.IEnvironment, evt hps.IEvent) {
		modifications = []string{}
		snapshot = strings.Join(states, ", ")
	})

	e.EndingIteration(func(env hps.IEnvironment, evt hps.IEvent) {
		if len(modifications) == 0 {
			return
		}

		p.logger.Printf("REPORT: State BEFORE event: %s", snapshot)
		p.logger.Printf("REPORT: EVENT: %s", strings.Join(modifications, ", "))
		p.logger.Printf("REPORT: State AFTER event: %s", strings.Join(states, ", "))
	})

	return nil
}

func (p *PowerSubsystemChangesPlugin) Done(r *engine.SimulationResult) {
}

func trackPowerNetworkModifications(substations *model.Network, modified func(int, bool, bool, string, string)) int {
	i := 0
	for _, s := range substations.Substations() {
		for _, bay := range s.GeneratorBays() {
			idx := i
			i++

			connected := bay.PropertyConnected()
			generator, _ := bay.Generator()
			connected.Changed(func(p *hps.Property, oldValue interface{}) {
				v := bay.Connected()
				modified(idx, v, q.IsOk(generator), bay.Name(), connectionVerb(v))
			})
			q.GetState(generator, "fail").Entering(func(p *sm.State, time float64, duration float64) {
				modified(idx, bay.Connected(), q.IsOk(generator), generator.Name(), "failure")
			})
			q.GetState(generator, "ok").Entering(func(p *sm.State, time float64, duration float64) {
				modified(idx, bay.Connected(), q.IsOk(generator), generator.Name(), "recovery")
			})
			modified(idx, true, true, generator.Name(), "init")
		}

		for _, bay := range s.LoadBays() {
			idx := i
			i++
			connected := bay.PropertyConnected()
			load, _ := bay.Load()
			connected.Changed(func(p *hps.Property, oldValue interface{}) {
				v := bay.Connected()
				modified(idx, v, q.IsOk(load), load.Name(), connectionVerb(v))
			})
			q.GetState(load, "fail").Entering(func(p *sm.State, time float64, duration float64) {
				modified(idx, bay.Connected(), q.IsOk(load), load.Name(), "failure")
			})
			q.GetState(load, "ok").Entering(func(p *sm.State, time float64, duration float64) {
				modified(idx, bay.Connected(), q.IsOk(load), load.Name(), "recovery")
			})
			modified(idx, true, true, load.Name(), "init")
		}
	}

	for _, m := range substations.Links() {
		idx := i
		i++
		m.StateFail().Entering(func(p *sm.State, time float64, duration float64) {
			modified(idx, m.Connected(), m.Ok(), m.Name(), "failure")
		})
		m.StateOk().Entering(func(p *sm.State, time float64, duration float64) {
			modified(idx, m.Connected(), m.Ok(), m.Name(), "recovery")
		})
		m.PropertyConnected().Changed(func(p *hps.Property, oldValue interface{}) {
			v := m.Connected()
			modified(idx, v, m.Ok(), m.Name(), connectionVerb(v))
		})
		modified(idx, true, true, m.Name(), "init")
	}

	return i
}

func connectionVerb(v bool) string {
	if v {
		return "connection"
	}
	return "disconnection"
}
