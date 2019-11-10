package diversitystudy

import (
	"hps"
	sm "hps/statemachine"
	"hps/statemachine/triggers"
	n32model "nordic32/model"
	"nordic32/plugins/diversitystudy/model"
)

func setBreakerMtbf(breaker *model.Breaker, mtbf float64) {
	for _, component := range breaker.Components() {
		okState := component.StateOk()
		okState.Transitions()[0].SetTrigger(triggers.NewProbabilisticTrigger(1 / mtbf))
	}
}

func attachBreakerHandlers(e hps.IEnvironment, breaker *model.Breaker, fail func()) {
	componentFailed := func(p *sm.State, time float64, duration float64) {
		for _, component := range breaker.Components() {
			if !component.Failed() {
				return
			}
		}
		fail()
	}

	for _, component := range breaker.Components() {
		component.StateFail().Entering(componentFailed)
		component.StateCompromisedFail().Entering(componentFailed)
	}

}

func attachLineBreakerHandlers(e hps.IEnvironment, breaker *model.Breaker, link *n32model.Link) {
	attachBreakerHandlers(e, breaker, func() { link.Disconnect() })
}

func attachGeneratorBreakerHandlers(e hps.IEnvironment, breaker *model.Breaker, bay *n32model.GeneratorBay) {
	attachBreakerHandlers(e, breaker, func() { bay.Disconnect() })
}

func attachLoadBreakerHandlers(e hps.IEnvironment, breaker *model.Breaker, bay *n32model.LoadBay) {
	attachBreakerHandlers(e, breaker, func() { bay.Disconnect() })
}
