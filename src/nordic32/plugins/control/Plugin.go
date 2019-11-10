package control

import (
	"hps"
	"hps/engine"
	"hps/loggers"
	sm "hps/statemachine"
	"hps/statemachine/triggers"
	"hps/tools"
	"nordic32/plugins"
	"nordic32/query"
)

type Plugin struct {
	logger    hps.ILogger
	frequency string
	quering   string
}

func NewPlugin(logger hps.ILogger) plugins.IPlugin {
	p := &Plugin{loggers.NewContextLogger(logger, "control"), "", ""}
	_ = plugins.IPlugin(p)
	return p
}

func (p *Plugin) Name() string {
	return "control"
}

func (p *Plugin) Init(e hps.IEnvironment) error {
	machine := addControlMachine(e)

	trigger := func() {
		if machine.State().Name() == "working" {
			return
		}

		start, _ := machine.Property("start")
		value := start.Value()
		if _, ok := value.(*triggers.IdleTrigger); !ok {
			return
		}

		duration2h, _ := tools.ParseDuration("2h")
		start.SetValue(triggers.NewDeterministicTrigger(duration2h))
		machine.ClearEvent()
	}

	network, ok := query.FindSubstationsNetwork(e)
	if !ok {
		p.logger.Print("no substations")
		return nil
	}

	for _, link := range network.Links() {
		link.PropertyConnected().Changed(func(p *hps.Property, oldValue interface{}) {
			trigger()
		})
		link.Changed(func(p *sm.Machine) {
			trigger()
		})
	}
	for _, substation := range network.Substations() {
		for _, bay := range substation.GeneratorBays() {
			bay.PropertyConnected().Changed(func(p *hps.Property, oldValue interface{}) {
				trigger()
			})
			generator, _ := bay.Generator()
			generator.Changed(func(p *sm.Machine) {
				trigger()
			})
		}
		for _, bay := range substation.LoadBays() {
			bay.PropertyConnected().Changed(func(p *hps.Property, oldValue interface{}) {
				trigger()
			})
			load, _ := bay.Load()
			load.Changed(func(p *sm.Machine) {
				trigger()
			})
		}
	}

	working, _ := machine.GetState("working")
	working.Entering(func(p *sm.State, time float64, duration float64) {
		start, _ := machine.Property("start")
		start.SetValue(triggers.NewIdleTrigger())
		machine.ClearEvent()

		// do control
		actions := BuildActions(network)

		if len(actions) > 0 {
			ApplyActions(actions, BuildGraphFromNetwork(network))
		}
	})

	return nil
}

func (p *Plugin) SetFrequency(frequency string) {
	p.frequency = frequency
}

func (p *Plugin) SetQuerying(quering string) {
	p.quering = quering
}

func (p *Plugin) Done(r *engine.SimulationResult) {
}
