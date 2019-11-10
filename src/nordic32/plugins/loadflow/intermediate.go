package loadflow

import (
	"encoding/json"
)

type PNode struct {
	Id       int     `json:"id"`
	Capacity float64 `json:"capacity"`
	Power    float64 `json:"power"`
	Type     string  `json:"type"`
}

type PLink struct {
	Id   int     `json:"id"`
	From int     `json:"from"`
	To   int     `json:"to"`
	X    float64 `json:"x"`
	Max  float64 `json:"max"`
}

type PNetwork struct {
	Nodes []PNode `json:"nodes"`
	Links []PLink `json:"links"`
}

func ParseNetwork(text string) *LFModel {
	n := PNetwork{}
	json.Unmarshal([]byte(text), &n)

	nodes := make([]*LFNode, len(n.Nodes))
	for i, nd := range n.Nodes {
		nodes[i] = &LFNode{
			ID:       nd.Id,
			Capacity: nd.Capacity,
			Power:    nd.Power,
			Type:     nd.Type,
		}
	}

	links := make([]*LFLink, len(n.Links))
	for i, lnk := range n.Links {
		links[i] = &LFLink{
			ID:   lnk.Id,
			From: lnk.From,
			To:   lnk.To,
			X:    lnk.X,
			Max:  lnk.Max,
		}
	}

	return &LFModel{nodes, links, false}
}

func ParseFlags(text string) []bool {
	flags := []bool{}
	json.Unmarshal([]byte(text), &flags)
	return flags
}

func ParseInts(text string) []int {
	numbers := []int{}
	json.Unmarshal([]byte(text), &numbers)
	return numbers
}

func ParseFloats(text string) []float64 {
	numbers := []float64{}
	json.Unmarshal([]byte(text), &numbers)
	return numbers
}

func StringifyNetwork(network *LFModel) string {
	n := PNetwork{[]PNode{}, []PLink{}}
	for _, nd := range network.Nodes {
		n.Nodes = append(n.Nodes, PNode{
			Id:       nd.ID,
			Capacity: nd.Capacity,
			Power:    nd.Power,
			Type:     nd.Type,
		})
	}
	for _, lnk := range network.Links {
		n.Links = append(n.Links, PLink{
			Id:   lnk.ID,
			From: lnk.From,
			To:   lnk.To,
			X:    lnk.X,
			Max:  lnk.Max,
		})
	}
	data, _ := json.Marshal(n)
	return string(data)
}

func StringifyFlags(flags []bool) string {
	data, _ := json.Marshal(flags)
	return string(data)
}

func StringifyInts(numbers []int) string {
	data, _ := json.Marshal(numbers)
	return string(data)
}
