package main

import (
	"encoding/json"
	"flag"
  	"bufio"
  	"strings"
  	"fmt"
  	"os"
)

type Patch struct {
	Time float64
	Delta interface{}
}

func main() {
	var file string

	flag.StringVar(&file, "file", "", "JSLOG file")

	flag.Parse()

	if file == "" {
		fmt.Println("reward.go reads the .jslog file and replays the state changes.")
		fmt.Println("It is a base tool for different log-related tools.")
		fmt.Println("")		
		fmt.Println("The log file should be specified as follows:")
		fmt.Println("  --file=<path to the file>")
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
	model := initModel(readLongString(reader))
	iter(0.0, model)

	// main loop
	for {
		patch := initPatch(readLongString(reader))
		if patch == nil {
			break
		}
		applyPatch(model, patch)
		iter(patch.Time, model)
	}
}

func iter(time float64, model interface{}) {
	load := 0.0
	off := []string{}
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
					off = append(off, mm["name"].(string))
					continue
				}

				if "ok" != smm["content"].(map[string]interface{})["machines"].(map[string]interface{})["Load"].(map[string]interface{})["state"] {
					off = append(off, mm["name"].(string))
					continue
				}

				load += props["power"].(map[string]interface{})["value"].(float64)
			}
		}
	}
	fmt.Printf("%g %.0f %s\n", time, load, strings.Join(off, ","))
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
	patch.Time, _ = data["time"].(float64)
	patch.Delta = data["delta"]

	return patch
}

func applyPatch(model interface{}, patch *Patch) {
	apply(model, patch.Delta)
}

func apply(dst interface{}, src interface{}) {
	dt, dOk := dst.(map[string]interface{})
	st, sOk := src.(map[string]interface{})

	if !dOk || !sOk {
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
