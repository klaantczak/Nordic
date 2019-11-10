package loadflow

import (
	"math"
	"nordic32/model"
)

type LFResult struct {
	Flows  []float64
	States []int
}

func LoadFlow(model *LFModel) *LFResult {
	states := make([]bool, len(model.Nodes)+len(model.Links))
	for i, _ := range states {
		states[i] = true
	}
	return loadflow(model, states)
}

func LoadFlow2(model *LFModel, states []bool) *LFResult {
	return loadflow(model, states)
}

func loadflow(model *LFModel, states []bool) *LFResult {
	for i, n := range model.Nodes {
		if n.Capacity > 0 {
			model.Nodes[i].Power = n.Capacity
		}
	}

	idlookup := make(map[int]int)
	for i, n := range model.Nodes {
		idlookup[n.ID] = i
	}

	statuses := make([]int, len(states))
	for i := 0; i < len(statuses); i++ {
		if !states[i] {
			statuses[i] = StatusDisabled
		}
	}

	for i := 0; i < len(model.Nodes); i++ {
		id := model.Nodes[i].ID
		if !states[i] {
			for j := 0; j < len(model.Links); j++ {
				if model.Links[j].From == id || model.Links[j].To == id {
					statuses[len(model.Nodes)+j] = StatusDisabled
				}
			}
		}
	}

	for i := 0; i < len(model.Nodes); i++ {
		if statuses[i] == StatusDisabled {
			model.Nodes[i].Status = StatusDisabled
		} else {
			model.Nodes[i].Status = StatusEnabled
		}
	}

	for i := 0; i < len(model.Links); i++ {
		if statuses[len(model.Nodes)+i] == StatusDisabled {
			model.Links[i].Status = StatusDisabled
		} else {
			model.Links[i].Status = StatusEnabled
		}
	}

	buildIsland := func(statuses []int, island int) bool {
		next := -1
		for i := 0; i < len(model.Nodes); i++ {
			if statuses[i] == 0 {
				next = i
				break
			}
		}
		if next == -1 {
			return false
		}

		queue := []int{next}
		for i := 0; i < len(queue); i++ {
			statuses[queue[i]] = island
			id := model.Nodes[queue[i]].ID

			for j := 0; j < len(model.Links); j++ {
				status := statuses[len(model.Nodes)+j]
				if status == StatusOverloaded || status == StatusDisabled {
					statuses[len(model.Nodes)+j] = StatusDisabled
					continue
				}

				from := model.Links[j].From
				to := model.Links[j].To

				if from == id {
					statuses[len(model.Nodes)+j] = island
					status := statuses[idlookup[to]]
					if status == 0 {
						queue = append(queue, to)
					}
					continue
				}
				if to == id {
					statuses[len(model.Nodes)+j] = island
					status := statuses[idlookup[from]]
					if status == 0 {
						queue = append(queue, from)
					}
					continue
				}
			}
		}

		return true
	}

	island := 0
	for {
		island += 1

		if !buildIsland(statuses, island) {
			break
		}
	}

	for i := 1; i < island; i++ {
		nodemap := make(map[int]int)

		nodes := []Node{}
		for j := 0; j < len(model.Nodes); j++ {
			if statuses[j] != i {
				continue
			}

			n := model.Nodes[j]
			nodemap[n.ID] = len(nodes)

			power := n.Power
			// TODO: uncomment as compatibility
			// if power > 0.0 {
			// 	power = n.Capacity
			// }
			nodes = append(nodes, Node{
				Power: power,
				Data:  interface{}(n),
			})
		}

		links := []Link{}
		for j := 0; j < len(model.Links); j++ {
			if statuses[len(model.Nodes)+j] == i {
				lnk := model.Links[j]

				links = append(links, Link{
					From: nodemap[lnk.From],
					To:   nodemap[lnk.To],
					X:    lnk.X,
					Data: interface{}(lnk),
				})
			}
		}

		network := Network{nodes, links}

		production, demand, capacity := GetPowerSummary(&network)

		if capacity == 0.0 || demand == 0.0 {
			for _, n := range network.Nodes {
				n.Data.(*LFNode).Status = StatusDisconnected
			}
			for _, n := range network.Links {
				n.Data.(*LFLink).Status = StatusDisconnected
				n.Data.(*LFLink).Flow = 0
			}
		} else if production == -demand {
			// TODO do nothing, backward compatibility, remove
		} else if capacity < -demand {
			for _, n := range network.Nodes {
				n.Data.(*LFNode).Status = StatusCollapsed
			}
			for _, n := range network.Links {
				n.Data.(*LFLink).Status = StatusCollapsed
			}
		} else {
			BalanceNetwork(&network)

			// TODO: remove, compatibility bug
			for _, n := range network.Nodes {
				n.Data.(*LFNode).Power = n.Power
			}

			CalculateLoadFlow(network)

			hasOverloadedLinks := false

			for _, link := range network.Links {
				link.Data.(*LFLink).Flow = link.Flow

				if math.Abs(link.Flow) > link.Data.(*LFLink).Max {
					link.Data.(*LFLink).Status = StatusOverloaded
					hasOverloadedLinks = true
				}
			}

			for i := 0; i < len(model.Links); i++ {
				link := model.Links[i]
				if link.Status == StatusOverloaded {
					statuses[len(model.Nodes)+i] = link.Status
				}
			}

			if hasOverloadedLinks {
				for j := 0; j < len(statuses); j++ {
					if statuses[j] == i {
						statuses[j] = 0
					}
				}

				for {
					if !buildIsland(statuses, island) {
						break
					}

					island += 1
				}
			}
		}
	}

	for _, n := range model.Nodes {
		if n.Capacity > 0 {
			n.Power = n.Capacity
		}
	}

	for i := 0; i < len(model.Nodes); i++ {
		statuses[i] = model.Nodes[i].Status
	}

	for i := 0; i < len(model.Links); i++ {
		link := model.Links[i]
		if link.Status == StatusDisabled {
			link.Flow = 0.0
			statuses[len(model.Nodes)+i] = link.Status
		} else {
			statuses[len(model.Nodes)+i] = link.Status
		}
	}

	flows := make([]float64, len(model.Links))
	for i, lnk := range model.Links {
		flows[i] = lnk.Flow
	}

	return &LFResult{flows, statuses}
}

func buildStates(model *LFModel) []bool {
	states := make([]bool, len(model.Nodes)+len(model.Links))

	for i := 0; i < len(states); i++ {
		states[i] = true
	}

	for i, n := range model.Nodes {
		if n.Power > 0 {
			states[i] = 0 != n.Capacity
		} else if n.Power < 0 {
			states[i] = 0 != n.Power
		} else {
			states[i] = true
		}
	}

	for i, lnk := range model.Links {
		states[len(model.Nodes)+i] = lnk.Enabled
	}

	return states
}

func getMostOverloadedLine(network *model.Network, model *LFModel, states []bool, disconnected []int) int {
	mergedStates := make([]bool, len(states))
	for k, v := range states {
		mergedStates[k] = v
	}

	for _, v := range disconnected {
		mergedStates[v] = false
	}

	result := loadflow(model, mergedStates)

	maxOverloadedIndex := -1
	maxOverloadedRatio := 0.0

	for i, v := range result.States {
		if v != -4 {
			continue
		}

		machine := model.Links[i-len(model.Nodes)].Machine

		ratio := 0.0
		if i < len(result.Flows) {
			ratio = math.Abs(result.Flows[i]) / machine.Max()
		}

		if maxOverloadedIndex == -1 || maxOverloadedRatio < ratio {
			maxOverloadedIndex = i
			maxOverloadedRatio = ratio
		}
	}

	return maxOverloadedIndex
}

func GetTestedListOfOverloadedLines(network *model.Network, model *LFModel) []*LFLink {
	states := buildStates(model)

	disconnected := []int{}

	for {
		line := getMostOverloadedLine(network, model, states, disconnected)
		if line != -1 {
			disconnected = append(disconnected, line)
		} else {
			break
		}
	}

	result := make([]*LFLink, len(disconnected))
	for i, d := range disconnected {
		result[i] = model.Links[d-len(model.Nodes)]
	}

	return result
}
