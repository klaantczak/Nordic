package experiments

import (
	"hps/engine"
	"sort"
)

type Experiment struct {
	Environment *engine.Environment
}

type constructor func(environment *engine.Environment) *Experiment

var constructors = make(map[string]constructor)

// List returns array of the available experiment names.
func List() []string {
	list := []string{}
	for k := range constructors {
		list = append(list, k)
	}
	sort.Strings(list)
	return list
}

// CreateExperiment creates experiment by name.
func CreateExperiment(name string, environment *engine.Environment) *Experiment {
	if constructor, ok := constructors[name]; ok {
		return constructor(environment)
	}
	return nil
}
