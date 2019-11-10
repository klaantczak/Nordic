package diversitystudy

import (
	"hps"
	"hps/engine"
	"hps/loggers"
	nm "hps/networkmachine"
	sm "hps/statemachine"
	"nordic32/plugins"
	"nordic32/plugins/diversitystudy/model"
	q "nordic32/query"
	"os"
	"strconv"
)

type Plugin struct {
	logger hps.ILogger
	rnd    hps.IRnd
}

func NewPlugin(logger hps.ILogger, rnd hps.IRnd) plugins.IPlugin {
	return &Plugin{loggers.NewContextLogger(logger, "diversitystudy"), rnd}
}

func (p *Plugin) Name() string {
	return "diversitystudy"
}

func (p *Plugin) Init(e hps.IEnvironment) error {
	network, ok := q.FindSubstationsNetwork(e)
	if !ok {
		p.logger.Print("no substations")
		return nil
	}

	mtbfstr := os.ExpandEnv("$BREAKER_COMPONENT_MTBF")
	mtbf := 0.0
	if mtbfstr != "" {
		mtbf, _ = strconv.ParseFloat(mtbfstr, 64)
		p.logger.Printf("set MTBF to %v", mtbf)
	}

	breakerCounter := 0
	for _, substation := range network.Substations() {
		for _, lineBay := range substation.LineBays() {
			generalBreaker, _ := lineBay.Breaker()
			if compositeBreaker, ok := generalBreaker.(*nm.Machine); ok {
				breaker := model.NewBreaker(compositeBreaker)
				link, _ := network.Link(lineBay.Line())

				if mtbfstr != "" {
					setBreakerMtbf(breaker, mtbf)
				}

				attachLineBreakerHandlers(e, breaker, link)

				breakerCounter++
			}
		}

		for _, generatorBay := range substation.GeneratorBays() {
			generalBreaker, _ := generatorBay.Breaker()

			if compositeBreaker, ok := generalBreaker.(*nm.Machine); ok {
				breaker := model.NewBreaker(compositeBreaker)

				if mtbfstr != "" {
					setBreakerMtbf(breaker, mtbf)
				}

				attachGeneratorBreakerHandlers(e, breaker, generatorBay)

				breakerCounter++
			}
		}

		for _, loadBay := range substation.LoadBays() {
			generalBreaker, _ := loadBay.Breaker()

			if compositeBreaker, ok := generalBreaker.(*nm.Machine); ok {
				breaker := model.NewBreaker(compositeBreaker)

				if mtbfstr != "" {
					setBreakerMtbf(breaker, mtbf)
				}

				attachLoadBreakerHandlers(e, breaker, loadBay)

				breakerCounter++
			}
		}
	}

	if breakerCounter == 0 {
		p.logger.Print("no multicomponent breakers")
	} else {
		p.logger.Printf("attached to %d breakers", breakerCounter)
	}

	for _, m := range e.Machines() {
		if m.Kind() == "Breaker Attacker" {
			attacker := model.NewAttacker(m.(*sm.Machine))
			attachAttackerHandlers(attacker, network, p.rnd)
			p.logger.Print("attached to attacker")
		} else {
			p.logger.Print("no breaker attacker")
		}
		if m.Kind() == "Breaker Inspector" {
			attachInspectorHandlers(m.(*sm.Machine), network, p.rnd)
			p.logger.Print("attached to inspector")
		} else {
			p.logger.Print("no breaker inspector")
		}
	}
	return nil
}

func (p *Plugin) Done(r *engine.SimulationResult) {
}
