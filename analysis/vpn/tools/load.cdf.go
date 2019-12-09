package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Message struct {
	Time float64 `json:"time"`
	Text string  `json:"message"`
}

func main() {
	var jslog string

	flag.StringVar(&jslog, "jslog", "", "file with jslog")

	flag.Parse()

	if jslog == "" {
		fmt.Println("load.cdf.go reads the .jslog file and builds cumulative")
		fmt.Println("distribution function.")
		fmt.Println("")
		fmt.Println("Every entry in the log file should contain time and")
		fmt.Println("message, similar to")
		fmt.Println("  {\"time\":10,\"message\":\"Total load: 109034.9\"}")
		fmt.Println("")
		fmt.Println("The file with jslog should be specified as follows:")
		fmt.Println("  --jslogs=<path to the file>")
		fmt.Println("")
		fmt.Println("Result is ordered loads with counters, similar to")
		fmt.Println("  108879.212815;1")
		fmt.Println("  108913.812992;2")
		os.Exit(0)
	}

	content, err := ioutil.ReadFile(jslog)
	if err != nil {
		fmt.Println("Reading content of the file " + jslog + " failed with error \"" + err.Error() + "\"")
		os.Exit(1)
	}

	loads := []float64{}
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		msg := &Message{}
		err := json.Unmarshal([]byte(line), &msg)
		if err != nil {
			panic(err)
		}

		loadText := ""
		if strings.HasPrefix(msg.Text, "Total load: ") {
			loadText = msg.Text[len("Total load: "):]
		} else if strings.HasPrefix(msg.Text, "[counters] Total load: ") {
			loadText = msg.Text[len("[counters] Total load: "):]
		}

		if loadText != "" {
			load, _ := strconv.ParseFloat(loadText, 64)
			loads = append(loads, load)
		}
	}

	sort.Float64s(loads)

	normalisator := 1.0 / float64(len(loads))
	for i, v := range loads {
		fmt.Printf("%v;%v\n", v/109400, float64(i+1)*normalisator)
	}
}
