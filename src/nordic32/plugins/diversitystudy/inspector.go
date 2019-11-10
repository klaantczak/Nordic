package diversitystudy

import (
	"hps"
	nm "hps/networkmachine"
	sm "hps/statemachine"
	n32model "nordic32/model"
	"nordic32/plugins/diversitystudy/model"
)

func attachInspectorHandlers(m *sm.Machine, network *n32model.Network, rnd hps.IRnd) {
	inspection, _ := m.GetState("inspection")
	inspection.Entering(func(p *sm.State, time float64, duration float64) {
		for _, substation := range network.Substations() {
			for _, lineBay := range selectBays(substation) {
				machine, _ := lineBay.Breaker()
				model.NewBreaker(machine.(*nm.Machine)).Restore(time)
			}
		}
	})
}
