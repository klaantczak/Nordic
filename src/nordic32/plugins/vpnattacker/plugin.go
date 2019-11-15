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
	return &Plugin{loggers.NewContextLogger(logger, "vpn-attacker"), rnd}
}

func (p *Plugin) Name() string {
	return "vpn-attacker"
}

func (p *Plugin) Init(e hps.IEnvironment) error {
	network, ok := q.FindSubstationsNetwork(e)
	if !ok {
		p.logger.Print("no substations")
		return nil
	}

	attackers := []*Attacker{}
	for _, machine := range q.FindAllByKind(e.Machines(), "Vpn Attacker") {
		m, ok := machine.(*sm.Machine)
		if !ok {
			p.logger.Printf("The attacker %s is not a statemachine", machine.Name())
			return nil
		}
		attacker, ok := ToAttacker(m)
		if !ok {
			p.logger.Printf("The attacker %s is not a valid attacker", machine.Name())
			return nil
		}
		attackers = append(attackers, attacker)
	}

	if len(attackers) == 0 {
		p.logger.Print("network has no attackers")
		return nil
	}

	for _, attacker := range attackers {
		actions(attacker, network, p.rnd)
	}

	p.logger.Printf("Initialised %d attacker(s)", len(attackers))
	return nil
}

func (p *Plugin) Done(r *engine.SimulationResult) {
}
