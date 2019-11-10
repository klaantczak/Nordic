package control

import (
	"hps"
	sm "hps/statemachine"
	"nordic32/model"
	"nordic32/plugins/loadflow"
	"nordic32/query"
	"strings"
)

type Node struct {
	Id         int
	Name       string
	Node       *model.Substation
	Generators float64
	Load       float64
	Power      float64
	Capacity   float64
}

type Link struct {
	Id   int
	From int
	To   int
	X    float64
	Max  float64
}

type Load struct {
	Power   float64
	Bay     *model.LoadBay
	Node    *Node
	Enabled bool
}

func (load *Load) Enable() {
	if load.Enabled {
		return
	}

	load.Node.Load += load.Power
	load.Enabled = true
}

func (load *Load) Disable() {
	if !load.Enabled {
		return
	}

	load.Node.Load -= load.Power
	load.Enabled = false
}

type Network struct {
	Nodes []*Node
	Links []*Link
	Loads []*Load
}

type Action struct {
	Command    string
	Bay        hps.IMachine
	Substation *model.Substation
	Link       *model.Link
}

/**
 * Creates model base to find shed loads.
 *
 * Model inclues generators, loads and links that are not in failed state.
 */
func buildModelBase(substations *model.Network) *Network {
	nodes := []*Node{}
	loads := []*Load{}
	links := []*Link{}

	index := 0
	for _, n := range substations.Substations() {
		// assume that any disconnected generator, can be connected:
		gsum := 0.0
		for _, b := range n.GeneratorBays() {
			if g, ok := b.Generator(); ok {
				if query.IsOk(g) {
					gsum += b.Capacity()
				}
			} else {
				panic("Expect generator")
			}
		}

		node := Node{
			Id:         index,
			Name:       n.Name(),
			Generators: gsum,
			Load:       0.0,
			Power:      0.0,
			Capacity:   0.0,
			Node:       n,
		}

		nodes = append(nodes, &node)

		// assume that any disconnected (normally, by control) load can be
		// connected:
		for _, b := range n.LoadBays() {
			l, _ := b.Load()
			if query.IsOk(l) {
				loads = append(loads, &Load{
					Power: b.Power(),
					Bay:   b,
					Node:  &node,
				})
			}
		}

		index++
	}

	// assume that any disconnected line (normally, by protection) can be
	// connected:

	for _, n := range substations.Links() {
		if n.Ok() {
			xValue := n.X()
			maxValue := n.Max()
			fromName := n.From()
			toName := n.To()

			fromId := 0
			for _, n := range nodes {
				if n.Name == fromName {
					fromId = n.Id
				}
			}

			toId := 0
			for _, n := range nodes {
				if n.Name == toName {
					toId = n.Id
				}
			}

			item := &Link{
				Id:   index,
				From: fromId,
				To:   toId,
				X:    xValue,
				Max:  maxValue,
			}

			n.PropertyMax().Changed(func(p *hps.Property, oldValue interface{}) {
				value, _ := p.GetFloat()
				item.Max = value
			})

			links = append(links, item)

			index++
		}
	}

	return &Network{nodes, links, loads}
}

func shedLoadFlowOk(substations *model.Network, model *Network) bool {
	lfnodes := []*loadflow.LFNode{}
	lflinks := []*loadflow.LFLink{}

	for _, n := range model.Nodes {
		power := n.Generators - n.Load

		capacity := power
		if capacity < 0 {
			capacity = 0
		}

		nodeType := ""
		if 0 < capacity {
			nodeType = "producer"
		} else if 0 > power {
			nodeType = "consumer"
		} else {
			nodeType = "distributor"
		}

		lfnodes = append(lfnodes, &loadflow.LFNode{
			ID:       n.Id,
			Name:     n.Name,
			Power:    n.Generators - n.Load,
			Capacity: capacity,
			Type:     nodeType,
		})
	}

	for _, n := range model.Links {
		lflinks = append(lflinks, &loadflow.LFLink{
			ID:   n.Id,
			From: n.From,
			To:   n.To,
			X:    n.X,
			Max:  n.Max,
		})
	}

	result := loadflow.LoadFlow(&loadflow.LFModel{lfnodes, lflinks, true})

	for _, v := range result.States {
		if v == -4 {
			return false
		}
	}

	return true
}

func shedLoadEnableOneByOne(substations *model.Network, model *Network) {
	for _, load := range model.Loads {
		load.Enable()

		if shedLoadFlowOk(substations, model) {
			continue
		}

		load.Disable()
	}
}

func shedLoad(substations *model.Network) ([]*Load, []*Load) {
	model := buildModelBase(substations)

	// optimistic check
	for _, load := range model.Loads {
		load.Enable()
	}

	if shedLoadFlowOk(substations, model) {
		return model.Loads, []*Load{}
	}

	for _, load := range model.Loads {
		load.Disable()
	}

	// enable one by one
	shedLoadEnableOneByOne(substations, model)

	enabled := []*Load{}
	disabled := []*Load{}
	for _, load := range model.Loads {
		if load.Enabled {
			enabled = append(enabled, load)
		} else {
			disabled = append(disabled, load)
		}
	}

	return enabled, disabled
}

func BuildActions(substations *model.Network) []Action {
	actions := []Action{}

	enabledLoads, disabledLoads := shedLoad(substations)

	for _, load := range disabledLoads {
		if load.Bay.Connected() {
			actions = append(actions, Action{
				Command:    "disconnect-load",
				Bay:        load.Bay,
				Substation: load.Node.Node,
			})
		}
	}

	for _, load := range enabledLoads {
		if load.Bay.Disconnected() {
			actions = append(actions, Action{
				Command:    "connect-load",
				Bay:        load.Bay,
				Substation: load.Node.Node,
			})
		}
	}

	// Control should find the configuration without overloaded links
	// so as a result of the control, the network should be balanced
	// properly so all the links may be connected.
	//
	// The link can be disconnected by protection or as attack target.
	//
	// Connecting failed link does not make anny effect.
	for _, link := range substations.Links() {
		if link.Disconnected() {
			actions = append(actions, Action{
				Command: "connect-link",
				Link:    link,
			})
		}
	}

	// Control also turns on the disconnected generators. Connecting
	// failed generator does not make any effect.
	for _, substation := range substations.SubstationsWithGenerators() {
		for _, bay := range substation.GeneratorBays() {
			if !bay.Connected() {
				actions = append(actions, Action{
					Command:    "connect-generator",
					Bay:        bay,
					Substation: substation,
				})
			}
		}
	}

	return actions
}

func ApplyActions(actions []Action, graph Graph) {
	reachable := func(graph Graph, name string) bool {
		if _, ok := graph["National Control Centre"]; !ok {
			// in case of control center not in the graph, the control is direct
			return true
		}
		return connected(graph, "National Control Centre", name)
	}

	connect := func(m hps.IMachine) {
		connected, _ := m.Property("connected")
		connected.SetValue(true)
	}

	disconnect := func(m hps.IMachine) {
		connected, _ := m.Property("connected")
		connected.SetValue(false)
	}

	for _, action := range actions {
		switch action.Command {
		case "disconnect-load":
			if reachable(graph, action.Substation.Name()) {
				disconnect(action.Bay)
			}
		case "connect-load":
			if reachable(graph, action.Substation.Name()) {
				connect(action.Bay)
			}
		case "connect-link":
			// connecting link requires both substations being connected
			if reachable(graph, action.Link.From()) && reachable(graph, action.Link.To()) {
				action.Link.Connect()
			}

		case "connect-generator":
			if reachable(graph, action.Substation.Name()) {
				connect(action.Bay)
			}
		}
	}
}

func BuildGraphFromNetwork(substations *model.Network) Graph {
	graph := Graph{}

	machines := []hps.IMachine{}
	for _, m := range substations.Machines() {
		kind := m.Kind()
		if strings.HasPrefix(kind, "Substation ") || kind == "Control Centre" || kind == "Data Centre" {
			machines = append(machines, m)
		}
	}

	for _, m := range machines {
		node := GraphNode{m.Name(), false, map[string]bool{}}
		if sm, ok := m.(sm.IStateMachine); ok {
			node.Enabled = sm.State().Name() == "ok"
		}
		graph[m.Name()] = node
	}

	for _, m := range substations.Machines() {
		if m.Kind() == "Data Link" {
			from, _ := m.Property("from")
			fromValue, _ := from.GetString()
			to, _ := m.Property("to")
			toValue, _ := to.GetString()

			graph[fromValue].Links[toValue] = m.(*sm.Machine).State().Name() == "ok"
			graph[toValue].Links[fromValue] = m.(*sm.Machine).State().Name() == "ok"
		}
	}

	return graph
}
