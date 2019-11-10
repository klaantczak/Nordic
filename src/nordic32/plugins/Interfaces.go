package plugins

import (
	"hps"
	"hps/engine"
)

// IPlugin defines methods for the Nordic32 simulation plugin.
type IPlugin interface {
	Name() string
	Init(environment hps.IEnvironment) error
	Done(r *engine.SimulationResult)
}
