package modificator

import (
	"hps"
	"hps/engine"
	"hps/loggers"
	sm "hps/statemachine"
	"nordic32/plugins"
	"nordic32/query"
)

type Plugin struct {
	logger hps.ILogger
	rnd    hps.IRnd
}

func NewPlugin(logger hps.ILogger, rnd hps.IRnd) plugins.IPlugin {
	return &Plugin{loggers.NewContextLogger(logger, "modificator"), rnd}
}

func (p *Plugin) Name() string {
	return "modificator"
}

func (p *Plugin) Init(e hps.IEnvironment) error {
	substations, ok := query.FindSubstationsNetwork(e)
	if !ok {
		p.logger.Print("no substations")
		return nil
	}

	countModificators := 0
	for _, m := range e.Machines() {
		if m.Kind() == "Modificator" {
			actions(m.(*sm.Machine), substations, p.rnd)
			countModificators++
		}
	}

	if countModificators == 0 {
		p.logger.Print("no modificators")
	} else {
		p.logger.Printf("added handlers to %d modificators", countModificators)
	}

	return nil
}

func (p *Plugin) Done(r *engine.SimulationResult) {
}
