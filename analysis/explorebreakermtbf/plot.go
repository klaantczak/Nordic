package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"strconv"

	"github.com/gonum/plot"
	"github.com/gonum/plot/plotter"
	"github.com/gonum/plot/plotutil"
	"github.com/gonum/plot/vg"
)

func main() {
	drawLoad()
	drawTime()
	drawEvents()
	drawTimePerEvent()
}

func drawLoad() {
	recs := readData("data/report.csv")

	loads := make(plotter.XYs, len(recs))
	for i, r := range recs {
		loads[i].Y = r.Load
		loads[i].X = r.Mtbf
	}

	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	p.Title.Text = "Average Load"
	p.X.Label.Text = "Breaker MTBF (fraction of year)"
	p.Y.Label.Text = "Load"
	p.X.Scale = plot.LogScale{}
	p.X.Tick.Marker = plot.LogTicks{}

	plotutil.AddLines(p, []interface{}{loads}...)

	// Save the plot to a PNG file.
	if err := p.Save(6*vg.Inch, 4*vg.Inch, "charts/load.png"); err != nil {
		panic(err)
	}

}

func drawTime() {
	recs := readData("data/report.csv")

	loads := make(plotter.XYs, len(recs))
	for i, r := range recs {
		loads[i].Y = r.Time
		loads[i].X = r.Mtbf
	}

	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	p.Title.Text = "Average Time for Simulation"
	p.X.Label.Text = "Breaker MTBF (fraction of year)"
	p.Y.Label.Text = "Time"
	p.X.Scale = plot.LogScale{}
	p.X.Tick.Marker = plot.LogTicks{}

	plotutil.AddLines(p, []interface{}{loads}...)

	// Save the plot to a PNG file.
	if err := p.Save(6*vg.Inch, 4*vg.Inch, "charts/time.png"); err != nil {
		panic(err)
	}

}

func drawEvents() {
	recs := readData("data/report.csv")

	loads := make(plotter.XYs, len(recs))
	for i, r := range recs {
		loads[i].Y = r.Events
		loads[i].X = r.Mtbf
	}

	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	p.Title.Text = "Average Events per Simulation"
	p.X.Label.Text = "Breaker MTBF (fraction of year)"
	p.Y.Label.Text = "Events"
	p.X.Scale = plot.LogScale{}
	p.X.Tick.Marker = plot.LogTicks{}

	plotutil.AddLines(p, []interface{}{loads}...)

	// Save the plot to a PNG file.
	if err := p.Save(6*vg.Inch, 4*vg.Inch, "charts/events.png"); err != nil {
		panic(err)
	}

}

func drawTimePerEvent() {
	recs := readData("data/report.csv")

	loads := make(plotter.XYs, len(recs))
	for i, r := range recs {
		loads[i].Y = r.Time / r.Events
		loads[i].X = r.Mtbf
	}

	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	p.Title.Text = "Average Timeper Event"
	p.X.Label.Text = "Breaker MTBF (fraction of year)"
	p.Y.Label.Text = "Time"
	p.X.Scale = plot.LogScale{}
	p.X.Tick.Marker = plot.LogTicks{}

	plotutil.AddLines(p, []interface{}{loads}...)

	// Save the plot to a PNG file.
	if err := p.Save(6*vg.Inch, 4*vg.Inch, "charts/event.time.png"); err != nil {
		panic(err)
	}

}

type Rec struct {
	Mtbf   float64
	Load   float64
	Events float64
	Time   float64
}

func readData(file string) []Rec {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println("Reading content of the file " + file + " failed with error \"" + err.Error() + "\"")
		os.Exit(1)
	}

	lines := strings.Split(string(content), "\n")

	recs := []Rec{}
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}

		values := strings.Split(line, ";")
		mtbf, _ := strconv.ParseFloat(values[0], 64)
		load, _ := strconv.ParseFloat(values[1], 64)
		events, _ := strconv.ParseFloat(values[2], 64)
		time, _ := strconv.ParseFloat(values[3], 64)
		recs = append(recs, Rec{
			mtbf,
			load / 109400,
			events,
			time,
		})
	}

	return recs
}
