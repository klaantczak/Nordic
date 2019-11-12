package vpnattacker

import (
	"hps"
	"hps/engine"
	"hps/loggers"
	sm "hps/statemachine"
	"nordic32/plugins"
	q "nordic32/query"
)

type Plugin struct {
	logger hps.ILogger
	rnd    hps.IRnd
}

func NewPlugin(logger hps.ILogger, rnd hps.IRnd) plugins.IPlugin {
	return &Plugin{loggers.NewContextLogger(logger, "vpnattacker"), rnd}
}

func (p *Plugin) Name() string {
	return "vpnattacker"
}

func (p *Plugin) Init(e hps.IEnvironment) error {
	network, ok := q.FindSubstationsNetwork(e)
	if !ok {
		p.logger.Print("no substations")
		return nil
	}

	countVpnAttackers := 0
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
