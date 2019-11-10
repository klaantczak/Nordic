package protection

import (
	"hps"
	nm "hps/networkmachine"
	sm "hps/statemachine"
	"nordic32/model"
	q "nordic32/query"
	"strconv"
)

func actions(network *model.Network, m *model.Link) {
	query := createDisconnectionQuery(network, m)
	m.PropertyOverloaded().Changed(func(p *hps.Property, oldValue interface{}) {
		value, _ := p.GetBool()
		if value && query() {
			m.Disconnect()
		}
	})
}

func createDisconnectionQuery(network *model.Network, m *model.Link) func() bool {
	from, ok := network.Substation(m.From())
	if !ok {
		panic("Substation " + m.From() + " could not be found.")
	}

	fromBay, ok := from.LineBay(m.Name())
	if !ok {
		panic("Error while processing machine " + m.Name() + ": " +
			"Could not find corresponding line bay for machine " + from.Name())
	}

	to, ok := network.Substation(m.To())
	if !ok {
		panic("Substation " + m.To() + " could not be found.")
	}

	toBay, ok := to.LineBay(m.Name())
	if !ok {
		panic("Error while processing machine " + m.Name() + ": " +
			"Could not find corresponding line bay for machine " + to.Name())
	}

	fromQuery := createWorkingBayQuery(fromBay)
	toQuery := createWorkingBayQuery(toBay)

	return func() bool {
		return fromQuery() || toQuery()
	}
}

func createWorkingBayQuery(m *model.LineBay) func() bool {
	breaker, _ := m.Breaker()
	primaryRelay, _ := m.PrimaryRelay()
	backupRelay, _ := m.BackupRelay()
	battery, _ := m.Battery()
	wiring, _ := m.Wiring()

	return func() bool {
		return breakerIsOk(breaker) && (ok(primaryRelay) || ok(backupRelay)) && ok(battery) && ok(wiring)
	}
}

func ok(m hps.IMachine) bool {
	return m.(*sm.Machine).State().Name() == "ok"
}

func breakerIsOk(m hps.IMachine) bool {
	if simpleBreaker, ok := m.(sm.IStateMachine); ok {
		return simpleBreaker.State().Name() == "ok"
	}

	// multi-channel compromisable breaker
	// TODO move this code to the diversity study plugin
	if twoChannelBreaker, ok := m.(*nm.Machine); ok {
		machines := twoChannelBreaker.Machines()

		i := 1
		isOk := true
		for {
			component, ok := q.FindMachineByName(machines, "Component"+strconv.Itoa(i))
			if !ok {
				break
			}
			state := component.(*sm.Machine).State().Name()
			isOk = isOk && (state == "ok" || state == "compromised-ok")
			i++
		}

		return isOk
	}

	return true
}
