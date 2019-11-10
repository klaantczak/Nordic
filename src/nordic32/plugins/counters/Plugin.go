package counters

import (
	"hps"
	"hps/engine"
	"hps/loggers"
	"nordic32/plugins"
	"nordic32/plugins/loadflow"
)

type Plugin struct {
	logger                hps.ILogger
	lastCalculatedLoad    float64
	totalLoad             float64
	e                     hps.IEnvironment
	loadChangingEventTime float64
}

func NewPlugin(logger hps.ILogger, loadflowPlugin *loadflow.Plugin) plugins.IPlugin {
	p := &Plugin{loggers.NewContextLogger(logger, "counters"), 0.0, 0.0, nil, 0.0}
	loadflowPlugin.Initialised(func(model *loadflow.LFModel) {
		load := p.calculateLoad(model)
		p.lastCalculatedLoad = load

		p.logger.Printf("initial load is %f", load)
	})

	loadflowPlugin.Changed(func(model *loadflow.LFModel, time float64) {
		load := p.calculateLoad(model)

		if p.lastCalculatedLoad != load {
			p.totalLoad += (time - p.loadChangingEventTime) * p.lastCalculatedLoad
			p.loadChangingEventTime = time
			p.lastCalculatedLoad = load
		}
	})

	loadflowPlugin.Completed(func(model *loadflow.LFModel, time float64) {
		p.totalLoad += (time - p.loadChangingEventTime) * p.lastCalculatedLoad
		p.logger.Printf("Total load: %f", p.totalLoad)
	})

	return p
}

func (p *Plugin) Name() string {
	return "counters"
}

func (p *Plugin) Init(e hps.IEnvironment) error {
	p.e = e
	return nil
}

func (p *Plugin) Done(r *engine.SimulationResult) {
}

func (p *Plugin) TotalLoad() float64 {
	return p.totalLoad
}

func (p *Plugin) calculateLoad(model *loadflow.LFModel) float64 {
	load := 0.0
	for _, node := range model.Nodes {
		if node.Status == loadflow.StatusEnabled {
			for _, l := range node.Loads {
				if l.Enabled {
					load += l.Value
				}
			}
		}
	}
	return load
}
