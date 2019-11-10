package loadflow

import (
	"hps"
	sm "hps/statemachine"
	"nordic32/model"
)

const StatusEnabled = 1
const StatusDisabled = -1
const StatusDisconnected = -2
const StatusCollapsed = -3
const StatusOverloaded = -4

type LFModel struct {
	Nodes  []*LFNode
	Links  []*LFLink
	Solved bool
}

func NewModelFromNetwork(network *model.Network) *LFModel {
	model := &LFModel{[]*LFNode{}, []*LFLink{}, false}

	for index, node := range network.Substations() {
		item := &LFNode{
			index,
			node.Name(),
			0,
			"distributor",
			[]*LFGenerator{},
			[]*LFLoad{},
			0,
			StatusEnabled}
		createItem(model, item, node)
	}

	for index, m := range network.Links() {
		link := &LFLink{}
		link.ID = index
		link.From = model.FindNodeByName(m.From()).ID
		link.To = model.FindNodeByName(m.To()).ID
		link.X = m.X()
		link.Max = m.Max()
		link.Enabled = true
		link.Machine = m

		createLinkHandler(model, link, m)

		model.Links = append(model.Links, link)
	}

	return model
}

func (m *LFModel) FindNodeByName(name string) *LFNode {
	for _, n := range m.Nodes {
		if n.Name == name {
			return n
		}
	}
	return nil
}

func (m *LFModel) Solve() {
	loadflow(m, buildStates(m))
	m.Solved = true

}

func createItem(m *LFModel, n *LFNode, s *model.Substation) {
	updater := func() {
		generation := 0.0
		for _, g := range n.Generators {
			connected, _ := g.Connected.GetBool()
			g.Enabled = connected && g.Machine.State().Name() == "ok"
			if g.Enabled {
				generation += g.Value
			}
		}

		load := 0.0
		for _, l := range n.Loads {
			connected, _ := l.Connected.GetBool()
			l.Enabled = connected && l.Machine.State().Name() == "ok"
			if l.Enabled {
				load += l.Value
			}
		}

		delta := generation - load

		if delta > 0.0 {
			n.Capacity = delta
			n.Power = delta
			n.Type = "producer"
		} else if delta < 0.0 {
			n.Capacity = 0.0
			n.Power = delta
			n.Type = "consumer"
		} else {
			n.Capacity = 0.0
			n.Power = 0.0
			n.Type = "distributor"
		}

		m.Solved = false
	}

	for _, bay := range s.GeneratorBays() {
		generator, _ := bay.Generator()
		n.Generators = append(n.Generators, &LFGenerator{generator, bay.PropertyConnected(), bay.Capacity(), true})
	}

	for _, g := range n.Generators {
		g.Connected.Changed(func(p *hps.Property, oldValue interface{}) {
			updater()
		})
		g.Machine.Changed(func(m *sm.Machine) {
			updater()
		})
	}

	for _, bay := range s.LoadBays() {
		load, _ := bay.Load()
		n.Loads = append(n.Loads, &LFLoad{load, bay.PropertyConnected(), bay.Power(), true})
	}

	for _, l := range n.Loads {
		l.Connected.Changed(func(p *hps.Property, oldValue interface{}) {
			updater()
		})
		l.Machine.Changed(func(m *sm.Machine) {
			updater()
		})
	}

	updater()

	m.Nodes = append(m.Nodes, n)
}

func createLinkHandler(model *LFModel, link *LFLink, machine *model.Link) {
	machine.StateFail().Entering(func(p *sm.State, time float64, duration float64) {
		// TODO: should it be machine.Connected() && machine.Enaled()?
		link.Enabled = false
		model.Solved = false
	})

	machine.StateOk().Entering(func(p *sm.State, time float64, duration float64) {
		// TODO: should it be machine.Connected() && machine.Enaled()?
		link.Enabled = true
		model.Solved = false
	})

	machine.PropertyConnected().Changed(func(p *hps.Property, oldValue interface{}) {
		// TODO: should it be machine.Connected() && machine.Enaled()?
		link.Enabled = machine.Connected()
		model.Solved = false
	})

	machine.PropertyMax().Changed(func(p *hps.Property, oldValue interface{}) {
		link.Max = machine.Max()
		model.Solved = false
	})
}
