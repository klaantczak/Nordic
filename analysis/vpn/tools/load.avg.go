package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
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
		fmt.Println("load.avg.go reads the .jslog file and calculates average load.")
		fmt.Println("")
		fmt.Println("Every entry in the log file should contain time and message,")
		fmt.Println("similar to")
		fmt.Println("  {\"time\":10,\"message\":\"Total load: 109034.9\"}")
		fmt.Println("")
		fmt.Println("The file with jslog should be specified as follows:")
		fmt.Println("  --jslogs=<path to the file>")
		os.Exit(0)
	}

	content, err := ioutil.ReadFile(jslog)
	if err != nil {
		fmt.Println("Reading content of the file " + jslog + " failed with error \"" + err.Error() + "\"")
		os.Exit(1)
	}

	totalLoad := 0.0
	samples := 0
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
			samples++
			totalLoad += load
		}
	}
	fmt.Printf("Average load: %v\n", totalLoad/float64(samples))
}
