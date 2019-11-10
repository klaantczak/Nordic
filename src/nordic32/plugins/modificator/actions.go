package modificator

import (
	"hps"
	sm "hps/statemachine"
	"math"
	"nordic32/model"
	"nordic32/plugins/loadflow"
)

func rndOne(rnd hps.IRnd, n int) int {
	return int(rnd.Next() * float64(n))
}

func actions(m *sm.Machine, network *model.Network, rnd hps.IRnd) {
	modifyThreshold, _ := m.GetState("modifyThreshold")

	substations := network.Substations()
	links := network.Links()
	model := loadflow.NewModelFromNetwork(network)

	modifyThreshold.Entering(func(s *sm.State, time float64, duration float64) {
		substation := substations[rndOne(rnd, len(substations))]
		bays := substation.LineBays()
		bay := bays[rndOne(rnd, len(bays))]
		link, _ := network.Link(bay.Line())

		idx := -1
		for i, v := range links {
			if v == link {
				idx = i
				break
			}
		}

		result := loadflow.LoadFlow(model)
		link.SetMax(math.Abs(result.Flows[idx]) * 1.1)
	})
}
