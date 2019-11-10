package protection

import (
	"hps"
	"hps/engine"
	"hps/loggers"
	"nordic32/plugins"
	q "nordic32/query"
)

type Plugin struct {
	logger hps.ILogger
}

func NewPlugin(logger hps.ILogger) plugins.IPlugin {
	return &Plugin{loggers.NewContextLogger(logger, "protection")}
}

func (p *Plugin) Name() string {
	return "protection"
}

func (p *Plugin) Init(e hps.IEnvironment) error {
	network, ok := q.FindSubstationsNetwork(e)
	if !ok {
		p.logger.Print("no substations")
		return nil
	}

	links := network.Links()
	for _, m := range links {
		actions(network, m)
	}

	p.logger.Printf("Added protection actions to %d link(s)", len(links))
	return nil
}

func (p *Plugin) Done(r *engine.SimulationResult) {
}
