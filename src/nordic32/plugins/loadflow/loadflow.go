package loadflow

import (
	"bytes"
	"encoding/binary"

	"github.com/gonum/matrix/mat64"
)

type CacheItem struct {
	Key    []byte
	Result []float64
	Hits   int
}

var cache = make(map[uint32][]*CacheItem)
var cacheCount = 0

func getCached(hash uint32, key []byte) ([]float64, bool) {
	if items, ok := cache[hash]; ok {
		for _, item := range items {
			if len(key) == len(item.Key) {
				for i, v := range key {
					if v != item.Key[i] {
						return []float64{}, false
					}
				}
				item.Hits++
				return item.Result, true
			}
		}
	}
	return []float64{}, false
}

func putCache(hash uint32, key []byte, result []float64) {
	if len(cache) > 1000 {
		compactCache()
	}
	cacheCount++
	if items, ok := cache[hash]; ok {
		cache[hash] = append(items, &CacheItem{key, result, 0})
	} else {
		cache[hash] = []*CacheItem{&CacheItem{key, result, 0}}
	}
}

func compactCache() {
	compacted := make(map[uint32][]*CacheItem)
	compactedCount := 0
	for key, segment := range cache {
		compactedSegment := []*CacheItem{}
		for _, item := range segment {
			if item.Hits > 1 {
				compactedSegment = append(compactedSegment,
					&CacheItem{item.Key, item.Result, item.Hits})
				compactedCount++
			}
		}
		if len(compactedSegment) > 0 {
			compacted[key] = compactedSegment
		}
	}
	cache = compacted
	cacheCount = compactedCount
}

func networkHash(n Network) (uint32, []byte) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, len(n.Nodes))
	for _, n := range n.Nodes {
		binary.Write(buf, binary.LittleEndian, n.Power)
	}
	binary.Write(buf, binary.LittleEndian, len(n.Links))
	for _, lnk := range n.Links {
		binary.Write(buf, binary.LittleEndian, lnk.From)
		binary.Write(buf, binary.LittleEndian, lnk.To)
		binary.Write(buf, binary.LittleEndian, lnk.X)
	}
	arr := buf.Bytes()
	var hash uint32
	for i, j, length := 0, 0, len(arr); i < length; i, j = i+1, ((j + 1) % 4) {
		hash = (hash * 397) ^ (uint32(arr[i]) << uint32(j*8))
	}
	return hash, arr
}

type Node struct {
	Power float64
	Data  interface{}
}

type Link struct {
	From int
	To   int
	X    float64
	Flow float64
	Data interface{}
}

type Network struct {
	Nodes []Node
	Links []Link
}

func CalculateLoadFlow(n Network) {
	hash, key := networkHash(n)
	if r, ok := getCached(hash, key); ok {
		for i := 0; i < len(n.Links); i++ {
			n.Links[i].Flow = r[i]
		}
		return
	}

	nodes := n.Nodes
	links := n.Links

	ba := mat64.NewDense(len(links), len(nodes), nil)

	bm := mat64.NewDense(len(nodes), len(nodes), nil)

	bl := mat64.NewDense(len(nodes), 1, nil)

	for i, link := range links {
		xm1 := 1.0 / link.X

		ba.Set(i, link.From, xm1)
		ba.Set(i, link.To, -xm1)

		bm.Set(link.From, link.To, bm.At(link.From, link.To)-xm1)
		bm.Set(link.To, link.From, bm.At(link.To, link.From)-xm1)

		bl.Set(link.From, 0, bl.At(link.From, 0)+xm1)
		bl.Set(link.To, 0, bl.At(link.To, 0)+xm1)
	}

	for i := 0; i < len(nodes); i += 1 {
		bm.Set(i, i, bl.At(i, 0))
	}

	br := mat64.NewDense(len(nodes)-1, len(nodes)-1, nil)
	br.Copy(bm.View(1, 1, len(nodes)-1, len(nodes)-1))

	bri := mat64.NewDense(len(nodes)-1, len(nodes)-1, nil)
	bri.Inverse(br)

	p := mat64.NewDense(len(nodes)-1, 1, nil)
	for i := 0; i < len(nodes)-1; i += 1 {
		node := nodes[i+1]
		p.Set(i, 0, node.Power)
	}

	tr := mat64.NewDense(len(nodes)-1, 1, nil)
	tr.Mul(bri, p)

	te := mat64.NewDense(len(nodes), 1, nil)

	for i := 0; i < len(nodes)-1; i += 1 {
		te.Set(i+1, 0, tr.At(i, 0))
	}

	f := mat64.NewDense(len(links), 1, nil)
	f.Mul(ba, te)

	flow := []float64{}
	for i := 0; i < len(links); i += 1 {
		flow = append(flow, f.At(i, 0))
	}

	putCache(hash, key, flow)

	for i := 0; i < len(links); i += 1 {
		links[i].Flow = flow[i]
	}
}
