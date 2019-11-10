package control

type GraphNode struct {
	Name    string
	Enabled bool
	Links   map[string]bool
}

type Graph map[string]GraphNode

func connected(nodes Graph, node1 string, node2 string) bool {
	n1, ok := nodes[node1]
	if !ok {
		panic("node1 is not in the graph")
	}

	n2, ok := nodes[node2]
	if !ok {
		panic("node2 is not in the graph")
	}

	if !n1.Enabled || n2.Enabled {
		return false
	}

	visited := []string{node1}

	index := 0

	for index != len(visited) {
		links := nodes[visited[index]].Links

		names := []string{}
		for k, _ := range links {
			names = append(names, k)
		}

		for idx := 0; idx < len(names); idx++ {
			name := names[idx]

			if _, ok := links[name]; !ok {
				continue
			}

			if node2 == name {
				return true
			}

			node := nodes[name]

			if !node.Enabled {
				continue
			}

			visit := true
			for _, n := range visited {
				if n == node.Name {
					visit = false
					break
				}
			}

			if visit {
				visited = append(visited, node.Name)
			}
		}

		index++
	}

	return false
}
