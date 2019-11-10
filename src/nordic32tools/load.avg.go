package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

type ReportSummary struct {
	Load float64 `json:"load"`
}

func processSummaries(rs []ReportSummary) {
	totalLoad := 0.0
	for _, v := range rs {
		totalLoad += v.Load
	}
	fmt.Println(totalLoad / float64(len(rs)))
}

func main() {
	var logs string

	flag.StringVar(&logs, "logs", "", "folder with log reports")

	flag.Parse()

	if logs == "" {
		fmt.Println("load.avg.go reads the .json files from the logs folder and")
		fmt.Println("parses simulation run summary of events and reward function.")
		fmt.Println("")
		fmt.Println("The folder with logs should be specified as follows:")
		fmt.Println("  --logs=<path to the folder with logs>")
		os.Exit(0)
	}

	files, err := ioutil.ReadDir(logs)
	if err != nil {
		fmt.Println("Reading the folder " + logs + " failed with error \"" + err.Error() + "\"")
		os.Exit(1)
	}

	list := []ReportSummary{}

	for _, item := range files {
		if item.IsDir() {
			continue
		}

		if filepath.Ext(item.Name()) != ".json" {
			continue
		}

		content, err := ioutil.ReadFile(path.Join(logs, item.Name()))
		if err != nil {
			fmt.Println("Reading content of the file " + item.Name() + " failed with error \"" + err.Error() + "\"")
			os.Exit(1)
		}

		var rs ReportSummary
		err = json.Unmarshal(content, &rs)
		if err != nil {
			fmt.Println("Unmarshaling content of the file " + item.Name() + " failed with error \"" + err.Error() + "\"")
			os.Exit(1)
		}

		list = append(list, ReportSummary{rs.Load})
	}

	processSummaries(list)
}
