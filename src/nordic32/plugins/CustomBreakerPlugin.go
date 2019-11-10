package plugins

import (
	"hps"
	"hps/engine"
	"hps/loggers"
	sm "hps/statemachine"
	q "nordic32/query"
)

type CustomBreakerPlugin struct {
	logger hps.ILogger
}

func NewCustomBreakerPlugin(logger hps.ILogger) IPlugin {
	return &CustomBreakerPlugin{loggers.NewContextLogger(logger, "CustomBreakerPlugin")}
}

func (p *CustomBreakerPlugin) Name() string {
	return "CustomBreakerPlugin"
}

func (p *CustomBreakerPlugin) Init(e hps.IEnvironment) error {
	network, _ := q.FindSubstationsNetwork(e)

	countBreakers := 0

	for _, substation := range network.Substations() {
		for _, lineBay := range substation.LineBays() {
			generalBreaker, _ := lineBay.Breaker()
			if breaker, ok := generalBreaker.(*sm.Machine); ok {
				countBreakers++

				link, _ := network.Link(lineBay.Line())

				failSt, _ := breaker.GetState("fail")

				failSt.Entering(func(p *sm.State, time float64, duration float64) {
					link.Disconnect()
				})
			}
		}
	}

	if countBreakers == 0 {
		p.logger.Print("no breakers was found")
	} else {
		p.logger.Printf("attached handlers to %d breakers", countBreakers)
	}

	return nil
}

func (p *CustomBreakerPlugin) Done(r *engine.SimulationResult) {
}
