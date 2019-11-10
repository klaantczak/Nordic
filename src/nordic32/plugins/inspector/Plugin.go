package inspector

import (
	"hps"
	"hps/engine"
	"hps/loggers"
	sm "hps/statemachine"
	"nordic32/model"
	"nordic32/plugins"
	"nordic32/query"
)

type Item struct {
	Link *model.Link
	Max  float64
}

type Plugin struct {
	logger hps.ILogger
	rnd    hps.IRnd
	table  []Item
}

func NewPlugin(logger hps.ILogger, rnd hps.IRnd) plugins.IPlugin {
	return &Plugin{
		loggers.NewContextLogger(logger, "inspector"),
		rnd,
		[]Item{}}
}

func (p *Plugin) Name() string {
	return "inspector"
}

func (p *Plugin) Init(e hps.IEnvironment) error {
	substations, ok := query.FindSubstationsNetwork(e)
	if !ok {
		p.logger.Print("no substations")
		return nil
	}

	links := substations.Links()

	for _, link := range links {
		p.table = append(p.table, Item{link, link.Max()})
	}

	for _, m := range e.Machines() {
		if m.Kind() == "Inspector" {
			inspector := m.(*sm.Machine)
			working, _ := inspector.GetState("working")
			working.Entering(func(s *sm.State, time float64, duration float64) {
				for _, em := range p.table {
					if em.Link.Max() != em.Max {
						em.Link.SetMax(em.Max)
					}
				}
			})
		}
	}

	return nil
}

func (p *Plugin) Done(r *engine.SimulationResult) {
}
