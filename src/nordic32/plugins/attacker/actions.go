package attacker

import (
	"hps"
	sm "hps/statemachine"
	"nordic32/model"
)

func rndOne(rnd hps.IRnd, n int) int {
	return int(rnd.Next() * float64(n))
}

func actions(m *Attacker, network *model.Network, rnd hps.IRnd) {
	m.StateDisconnectLoad().Entering(func(p *sm.State, time float64, duration float64) {
		if substations := network.SubstationsWithLoads(); len(substations) > 0 {
			substation := substations[rndOne(rnd, len(substations))]
			if bays := substation.LoadBays(); len(bays) > 0 {
				bay := bays[rndOne(rnd, len(bays))]
				bay.Disconnect()
			}
		}
	})

	m.StateDisconnectGenerator().Entering(func(p *sm.State, time float64, duration float64) {
		if substations := network.SubstationsWithGenerators(); len(substations) > 0 {
			substation := substations[rndOne(rnd, len(substations))]
			if bays := substation.GeneratorBays(); len(bays) > 0 {
				bay := bays[rndOne(rnd, len(bays))]
				bay.Disconnect()
			}
		}
	})

	m.StateDisconnectLine().Entering(func(p *sm.State, time float64, duration float64) {
		if substations := network.Substations(); len(substations) > 0 {
			substation := substations[rndOne(rnd, len(substations))]
			if bays := substation.LineBays(); len(bays) > 0 {
				bay := bays[rndOne(rnd, len(bays))]
				if link, ok := network.Link(bay.Line()); ok {
					link.Disconnect()
				}
			}
		}
	})
}
