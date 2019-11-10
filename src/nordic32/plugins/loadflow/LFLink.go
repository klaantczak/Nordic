package loadflow

import "nordic32/model"

type LFLink struct {
	ID      int
	From    int
	To      int
	X       float64
	Max     float64
	Enabled bool
	Flow    float64
	Machine *model.Link
	Status  int
}
