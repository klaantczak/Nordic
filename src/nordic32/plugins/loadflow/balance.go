package loadflow

// Returns network production, consumption and its limit.
func GetPowerSummary(n *Network) (float64, float64, float64) {
	consumption := 0.0
	production := 0.0
	capacity := 0.0

	for _, v := range n.Nodes {
		if v.Power > 0 {
			production += v.Power
			capacity += v.Data.(*LFNode).Capacity
		} else {
			consumption += v.Power
		}
	}

	return production, consumption, capacity
}

// Simple network balancer for fully connected networks with all generators at
// maximum level.
func BalanceNetwork(nwt *Network) bool {
	production, consumption, capacity := GetPowerSummary(nwt)

	if capacity == 0 || consumption == 0 {
		return false
	}

	// TODO: Should be removed
	if consumption == production {
		return false
	}

	r := -consumption / capacity

	if r > 1 {
		// generators cannot satisfy the load
		return false
	}

	for i, n := range nwt.Nodes {
		if n.Power > 0 {
			nwt.Nodes[i].Power = r * n.Data.(*LFNode).Capacity
		}
	}

	return true
}
