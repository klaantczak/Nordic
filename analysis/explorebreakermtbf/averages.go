package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

type Message struct {
	Time float64 `json:"time"`
	Text string  `json:"message"`
}

func parseSummary(str string) (int, time.Duration) {
	p := strings.Index(str, " events in ")

	eventsStr := str[0:p]
	events, _ := strconv.ParseInt(eventsStr, 10, 32)

	timeStr := str[p+len(" events in ") : len(str)]
	diration, _ := time.ParseDuration(timeStr)
	return int(events), diration
}

func main() {
	var jslog string
	var label string

	flag.StringVar(&jslog, "jslog", "", "file with jslog")
	flag.StringVar(&label, "label", "", "label")

	flag.Parse()

	if jslog == "" || label == "" {
		fmt.Println("averages.go reads the .jslog file and calculates average load, simulation duration and number of events.")
		fmt.Println("")
		fmt.Println("Every entry in the log file should contain time and message,")
		fmt.Println("similar to")
		fmt.Println("  {\"time\":10,\"message\":\"Total load: 109034.9\"}")
		fmt.Println("")
		fmt.Println("The file with jslog should be specified as follows:")
		fmt.Println("  --jslogs=<path to the file>")
		fmt.Println("")
		fmt.Println("The label should be specified as follows:")
		fmt.Println("  --label=<label>")
		os.Exit(0)
	}

	content, err := ioutil.ReadFile(jslog)
	if err != nil {
		fmt.Println("Reading content of the file " + jslog + " failed with error \"" + err.Error() + "\"")
		os.Exit(1)
	}

	totalTime := 0.0
	totalEvents := 0
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

		if strings.Contains(msg.Text, " events in ") {
			events, time := parseSummary(msg.Text)
			totalEvents += int(events)
			totalTime += time.Seconds()
		}
	}
	fmt.Printf(
		"%s;%v;%v;%v\n",
		label,
		totalLoad/float64(samples),
		totalEvents/samples,
		totalTime/float64(samples))
}
