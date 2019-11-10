package loadflow

import (
	"hps"
	sm "hps/statemachine"
)

type LFLoad struct {
	Machine   *sm.Machine
	Connected *hps.Property
	Value     float64
	Enabled   bool
}
