package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
)

type Patch struct {
	Time  float64
	Delta interface{}
	Raw   string
}

type State struct {
	time      float64
	load      float64
	attacking bool
	action    bool
	critical  float64
}

func (s *State) Update(time float64, load float64, attacking bool, idle bool, action bool) {
	p := *s

	s.time = time
	if attacking {
		s.attacking = true
	}
	if idle {
		s.attacking = false
	}
	s.load = load
	s.action = action

	if p.load == s.load {
		if p.attacking != s.attacking {
			s.print()
		}
	} else {
		s.print()
	}
}

func (s *State) print() {
	iif := func(v bool, then string, otherwise string) string {
		if v {
			return then
		}
		return otherwise
	}
	fmt.Printf("%g %.0f %s %s %s\n", s.time, s.load,
		iif(s.attacking, "attacking", "idle"),
		iif(s.action, "malicious", "failure"),
		iif(s.load < s.critical, "critical", "normal"))
}

func main() {
	var file string
	var critical float64

	flag.StringVar(&file, "file", "", "jslog file")
	flag.Float64Var(&critical, "critical", 0.0, "critical load threshold")

	flag.Parse()

	if file == "" {
		fmt.Println("reward.sm.simple.go reads the .jslog file for experiment " +
			"without attacker or with attacker that disconnects the substation " +
			"elements after penetrating the firewall and prints only the states " +
			"changes between attacking or idle attacker's state, malicious or " +
			"failure event, and between critical and not critical load shed.")
		fmt.Println("")
		fmt.Println("Arguments:")
		fmt.Println("  --file=<path to the file> - the jslog file, required")
		fmt.Println("  --critical=<number> - critical load threshold, optional, defauilt is 0")
		os.Exit(0)
	}

	handler, err := os.Open(file)
	if err != nil {
		panic(err)
	}

	reader := bufio.NewReader(handler)

	// seed
	readLongString(reader)

	// init
	eventTime := 0.0
	eventPatches := []*Patch{}
	state := &State{}
	state.critical = critical
	model := initModel(readLongString(reader))

	// main loop
	for {
		patch := initPatch(readLongString(reader))
		if patch == nil {
			break
		}
		if patch.Time == eventTime {
			eventPatches = append(eventPatches, patch)
		} else {
			processPatches(state, model, eventPatches)
			eventTime = patch.Time
			eventPatches = []*Patch{patch}
		}
	}

	processPatches(state, model, eventPatches)
}

func processPatches(state *State, model interface{}, patches []*Patch) {
	if len(patches) == 0 {
		return
	}

	attacking := false
	idle := false
	action := false

	for _, p := range patches {
		if machines, ok := p.Delta.(map[string]interface{})["machines"]; ok {
			if attacker, ok := machines.(map[string]interface{})["Attacker"]; ok {
				if _, ok := attacker.(map[string]interface{})["state"]; ok {
					state := attacker.(map[string]interface{})["state"].(string)
					attacking = state == "attack"
					idle = state == "idle"
					action = strings.HasPrefix(state, "disconnect")
				}
			}
		}
	}

	for _, p := range patches {
		applyPatch(model, p)
		state.Update(p.Time, load(model), attacking, idle, action)
	}
}

func load(model interface{}) float64 {
	load := 0.0
	for _, m := range model.(map[string]interface{})["machines"].(map[string]interface{}) {
		mm := m.(map[string]interface{})
		kind := mm["machine"].(string)
		if strings.HasPrefix(kind, "Substation ") {
			ms := mm["content"].(map[string]interface{})["machines"].(map[string]interface{})
			for _, sm := range ms {
				smm := sm.(map[string]interface{})
				if smm["machine"] != "Load Bay" {
					continue
				}

				props := smm["properties"].(map[string]interface{})
				if !props["connected"].(map[string]interface{})["value"].(bool) {
					continue
				}

				if "ok" != smm["content"].(map[string]interface{})["machines"].(map[string]interface{})["Load"].(map[string]interface{})["state"] {
					continue
				}

				load += props["power"].(map[string]interface{})["value"].(float64)
			}
		}
	}
	return load
}

func initModel(line string) interface{} {
	var model interface{}
	err := json.Unmarshal([]byte(line), &model)
	if err != nil {
		panic(err)
	}
	return model
}

func initPatch(line string) *Patch {
	if line == "" {
		return nil
	}

	var raw interface{}
	err := json.Unmarshal([]byte(line), &raw)
	if err != nil {
		panic(err)
	}

	data := raw.(map[string]interface{})

	patch := &Patch{}
	patch.Raw = line
	patch.Time, _ = data["time"].(float64)
	if delta, ok := data["delta"]; ok {
		patch.Delta = delta
	} else {
		return nil
	}

	if _, ok := patch.Delta.(map[string]interface{}); !ok {
		return nil
	}

	return patch
}

func applyPatch(model interface{}, patch *Patch) {
	apply(model, patch.Delta)
}

func apply(dst interface{}, src interface{}) {
	dt, dOk := dst.(map[string]interface{})
	st, sOk := src.(map[string]interface{})

	if !dOk || !sOk {
		fmt.Println(st)
		panic("Assume that the patched value exists.")
	}

	// iterate over the keys of src map
	for k, v := range st {
		if _, ok := v.(map[string]interface{}); ok {
			// dst = { test : { ok : false } } ; src = { test : { ok : true } }
			// k = test ; v = { ok : true }
			apply(dt[k], st[k])
		} else {
			// dst = { test : false } ; src = { test : true }
			// k = test ; v = true
			dt[k] = v
		}
	}
}

func readLongString(r *bufio.Reader) string {
	s, _ := r.ReadBytes('\n')
	return strings.Trim(string(s), "\r\n")
}
