package loadflow

import (
	"hps"
	"hps/engine"
	"hps/loggers"
	"nordic32/plugins"
	q "nordic32/query"
)

type ChangedEventHandler func(*LFModel, float64)
type InitialisedEventHandler func(*LFModel)

type Plugin struct {
	logger      hps.ILogger
	initialised []InitialisedEventHandler
	changed     []ChangedEventHandler
	completed   []ChangedEventHandler
	model       *LFModel
}

func NewPlugin(logger hps.ILogger) plugins.IPlugin {
	return &Plugin{
		loggers.NewContextLogger(logger, "loadflow"),
		[]InitialisedEventHandler{},
		[]ChangedEventHandler{},
		[]ChangedEventHandler{},
		nil,
	}
}

func (p *Plugin) Name() string {
	return "loadflow"
}

func (p *Plugin) Init(e hps.IEnvironment) error {
	substations, ok := q.FindSubstationsNetwork(e)
	if !ok {
		p.logger.Print("no substations")
		return nil
	}

	model := NewModelFromNetwork(substations)
	p.model = model

	p.logger.Print("initialising")

	model.Solve()

	for _, h := range p.initialised {
		h(model)
	}

	model.Solved = true

	e.EndingIteration(func(env hps.IEnvironment, evt hps.IEvent) {
		if model.Solved {
			return
		}

		model.Solve()

		for _, h := range p.changed {
			h(model, evt.Time())
		}

		for _, link := range GetTestedListOfOverloadedLines(substations, model) {
			link.Machine.Overload()
		}

		model.Solved = true
	})

	return nil
}

func (p *Plugin) Done(r *engine.SimulationResult) {
	p.logger.Print("completing")

	if p.model != nil {
		for _, h := range p.completed {
			h(p.model, r.Time)
		}
	}
}

func (p *Plugin) Initialised(handler InitialisedEventHandler) {
	p.initialised = append(p.initialised, handler)
}

func (p *Plugin) Changed(handler ChangedEventHandler) {
	p.changed = append(p.changed, handler)
}

func (p *Plugin) Completed(handler ChangedEventHandler) {
	p.completed = append(p.completed, handler)
}
