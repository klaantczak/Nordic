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
	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	p.Title.Text = "Load CDF"
	p.X.Label.Text = "Load"
	p.Y.Label.Text = "P"

	lines := []interface{}{}

	imagePath := os.Args[1]

	args := os.Args[2:]
	for i := 0; i < len(args); i += 2 {
		fmt.Println(args[i], args[i+1])
		lines = append(lines, args[i])
		lines = append(lines, readData(args[i+1]))
	}

	err = plotutil.AddLines(p, lines...)
	if err != nil {
		panic(err)
	}

	p.Y.Min = 0.0
	p.Y.Max = 1.0

	// Save the plot to a PNG file.
	if err := p.Save(6*vg.Inch, 4*vg.Inch, imagePath); err != nil {
		panic(err)
	}
}

func readData(file string) plotter.XYs {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println("Reading content of the file " + file + " failed with error \"" + err.Error() + "\"")
		os.Exit(1)
	}

	lines := strings.Split(string(content), "\n")
	ys := []float64{}
	xs := []float64{}
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}

		values := strings.Split(line, ";")

		x, _ := strconv.ParseFloat(values[0], 64)
		xs = append(xs, x)

		y, _ := strconv.ParseFloat(values[1], 64)
		ys = append(ys, y)
	}

	pts := make(plotter.XYs, len(xs))
	for i, x := range xs {
		pts[i].Y = ys[i]
		pts[i].X = x
	}
	return pts
}
