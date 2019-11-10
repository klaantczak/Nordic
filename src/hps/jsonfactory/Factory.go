package jsonfactory

import (
	"encoding/json"
	"fmt"
	"hps"
	nm "hps/networkmachine"
	sm "hps/statemachine"
	"hps/statemachine/triggers"
	"hps/tools"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
)

// FactoryRnd provides basic implementation of the random number generator.
type FactoryRnd struct{}

// Next returns next float64 random number from 0 (inclusive) to 1 (exclusive).
func (f *FactoryRnd) Next() float64 {
	return rand.Float64()
}

// Factory creates machines and newtorks from the loaded JSON definition.
type Factory struct {
	project     *ProjectDefinition
	rnd         hps.IRnd
	typeFactory hps.ITypeFactory
}

// NewFactory creates new Factory struct with given object factory.
func NewFactory(typeFactory hps.ITypeFactory, rnd hps.IRnd) *Factory {
	f := &Factory{}
	f.rnd = &FactoryRnd{}
	f.typeFactory = typeFactory
	_ = hps.IFactory(f)
	return f
}

// LoadJSON loads the project definition from binary-encoded JSON string.
func (f *Factory) LoadJSON(b []byte) error {
	p := &ProjectDefinition{}

	err := json.Unmarshal(b, &p)
	if err != nil {
		return err
	}

	f.project = p

	return nil
}

// Load loads the project definition from JSON file.
func (f *Factory) Load(path string) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	content, err := ioutil.ReadFile(absPath)
	if err != nil {
		return err
	}

	p := &ProjectDefinition{}
	p.Path = absPath

	err = json.Unmarshal(content, &p)
	if err != nil {
		return err
	}

	f.project = p

	return nil
}

// Create creates the machine instance, looking up the machine definition
// from the loaded list of machines.
func (f *Factory) Create(machineTypeName, machineName string) (hps.IMachine, error) {
	m, ok := f.project.FindMachineDefinition(machineTypeName)
	if !ok {
		return nil, fmt.Errorf("machine type %s is unknown", machineTypeName)
	}

	machine, err := f.createMachine(m, machineName)
	if err != nil {
		return nil, fmt.Errorf("cannot create machine '%v' of type '%v', %v", machineName, machineTypeName, err)
	}

	return machine, nil
}

// CreateNetwork creates and initialises all machines defined in the specified
// network.
func (f *Factory) CreateNetwork(name string) ([]hps.IMachine, error) {
	n, ok := f.project.FindNetworkDefinition(name)
	if !ok {
		return nil, fmt.Errorf("network %s is unknown", name)
	}
	r := make([]hps.IMachine, 0)

	for _, m := range n.Machines {
		instance, err := f.Create(m.Machine, m.Name)
		if err != nil {
			return nil, err
		}

		// load properties
		for pname, pvalue := range m.Properties {
			p, ok := instance.Property(pname)

			if !ok {
				return nil, fmt.Errorf("trying to set unknown property %s for the machine %s", pname, m.Name)
			}

			err := f.setPropertyValue(p, pvalue)

			if err != nil {
				return nil, fmt.Errorf("cannot create machine '%v' of type '%v' "+
					"defined in the network '%v'. Cannot set value of the property '%v', %v",
					m.Name, m.Machine, name, p.Name(), err)
			}
		}

		r = append(r, instance)
	}
	return r, nil
}

// GetNetworkNames returns the list of network names that are loaded
func (f *Factory) GetNetworkNames() []string {
	result := []string{}
	for _, n := range f.project.Networks {
		result = append(result, n.Name)
	}
	return result
}

func (f *Factory) setPropertyValue(p *hps.Property, b []byte) error {
	switch p.Type() {
	case "String", "Lookup", "string":
		v := ""
		err := json.Unmarshal(b, &v)
		if err != nil {
			return fmt.Errorf("Cannot parse property value. %v", err)
		}
		p.SetValue(v)
	case "Number", "number":
		v := 0.0
		err := json.Unmarshal(b, &v)
		if err != nil {
			return fmt.Errorf("Cannot parse property value. %v", err)
		}
		p.SetValue(v)
	case "Boolean", "bool":
		v := false
		err := json.Unmarshal(b, &v)
		if err != nil {
			return fmt.Errorf("Cannot parse property value. %v", err)
		}
		p.SetValue(v)
	case "Activation":
		v, err := f.loadTriggerFromJSON(b)
		if err != nil {
			return fmt.Errorf("Cannot parse property value. %v", err)
		}
		p.SetValue(v)
	default:
		return fmt.Errorf("Property type '%v' is unknown.", p.Type())
	}
	return nil
}

func (f *Factory) createMachine(def *MachineDefinition, name string) (hps.IMachine, error) {
	var m hps.IMachine
	switch def.Type {
	case "state-machine":
		var definition struct {
			States      []string                   `json:"states"`
			Initial     string                     `json:"initial"`
			Properties  map[string]json.RawMessage `json:"properties"`
			Transitions map[string]map[string][]json.RawMessage
		}

		if err := json.Unmarshal(def.Structure, &definition); err != nil {
			return nil, fmt.Errorf("cannot unmarshal state machine %s definition, %v", name, err)
		}

		stm, err := sm.NewMachine(name, def.Name, definition.States, definition.Initial)
		if err != nil {
			return nil, err
		}

		for pname, pdef := range def.Properties {
			p := stm.AddProperty(pname, pdef.Type)

			if pdef.Comment != "" {
				p.SetComment(pdef.Comment)
			}
		}

		for from, toTr := range definition.Transitions {
			for to, trs := range toTr {
				for _, tr := range trs {
					trigger, err := f.loadTriggerFromJSON(tr)
					if err != nil {
						return nil, fmt.Errorf("cannot read transition from '%v' to '%v', %v", from, to, err)
					}

					s, _ := stm.GetState(from)
					s.AddTransition(to, trigger)
				}
			}
		}

		m = stm
	case "network-machine":
		var definition struct {
			Network string `json:"network"`
		}

		if err := json.Unmarshal(def.Structure, &definition); err != nil {
			return nil, err
		}

		nwt, err := f.CreateNetwork(definition.Network)
		if err != nil {
			return nil, err
		}

		ntm := nm.NewMachine(name, def.Name)
		for pname, pdef := range def.Properties {
			ntm.AddProperty(pname, pdef.Type)
		}

		for _, m := range nwt {
			ntm.AddMachine(m)
		}
		m = ntm
	}

	if f.typeFactory != nil {
		m = f.typeFactory(m)
	}

	if m == nil {
		return nil, fmt.Errorf("Unknown machine type %v", def.Name)
	}

	return m, nil
}

func (f *Factory) loadTriggerFromJSON(b []byte) (sm.ITrigger, error) {
	var t sm.ITrigger
	var err error

	data := make(map[string]interface{})

	err = json.Unmarshal(b, &data)
	if err != nil {
		return nil, err
	}

	jm := tools.JsonMap(data)

	if transitionType, ok := jm.GetString("type"); ok {
		comment, _ := jm.GetString("comment")

		switch transitionType {
		case "idle":
			t = triggers.NewIdleTrigger()

		case "deterministic":
			if parameter, ok := jm.GetString("parameter"); ok {
				if parameter == "never" {
					fmt.Fprintf(os.Stderr, "Deprecated: trigger deterministic/never trigger should be replaced with idle trigger type.")
					t = triggers.NewIdleTrigger()
				} else if value, err := tools.ParseDuration(parameter); err == nil {
					fmt.Fprintf(os.Stderr, fmt.Sprintf("Deprecated: deterministic trigger parameter should be numeric, replace %s with %f.", parameter, value))
					t = triggers.NewDeterministicTrigger(value)
				} else {
					t = triggers.NewIdleTrigger()
					err = fmt.Errorf("Cannot parse transition's deterministic parameter.")
				}
			} else if parameter, ok := jm.GetFloat("parameter"); ok {
				t = triggers.NewDeterministicTrigger(parameter)
			} else if parameter, ok := jm.GetInt("parameter"); ok {
				t = triggers.NewDeterministicTrigger(float64(parameter))
			} else {
				t = triggers.NewIdleTrigger()
				err = fmt.Errorf("Cannot parse transition's deterministic parameter.")
			}

		case "probabilistic":
			if distribution, ok := jm.GetString("distribution"); ok {
				switch distribution {
				case "exponential":
					if param, ok := jm.GetString("parameter"); ok {
						if value, err := tools.ParseDuration(param); err == nil {
							fmt.Fprintf(os.Stderr, "Deprecated: exponential trigger parameter should be numeric.")
							t = triggers.NewProbabilisticTrigger(value)
							t.SetComment(comment)
						} else {
							t = triggers.NewIdleTrigger()
							err = fmt.Errorf("Cannot parse transition's exponential distribution parameter.")
						}
					} else {
						if value, ok := jm.GetFloat("parameter"); ok {
							t = triggers.NewProbabilisticTrigger(value)
						} else if value, ok := jm.GetInt("parameter"); ok {
							t = triggers.NewProbabilisticTrigger(float64(value))
						} else {
							t = triggers.NewIdleTrigger()
							err = fmt.Errorf("Cannot parse transition's exponential distribution parameter.")
						}
					}
				default:
					t = triggers.NewIdleTrigger()
					err = fmt.Errorf("This distribution is not supported.")
				}
			} else {
				t = triggers.NewIdleTrigger()
				err = fmt.Errorf("Cannot parse transition's distribution type.")
			}

		case "property":
			if property, ok := jm.GetString("property"); ok {
				t = triggers.NewPropertyTrigger(property)
			} else {
				t = triggers.NewIdleTrigger()
				err = fmt.Errorf("Cannot parse transition.")
			}
		}

		t.SetComment(comment)
	} else {
		t = triggers.NewIdleTrigger()
		err = fmt.Errorf("Cannot parse transition type.")
	}

	return t, err
}
