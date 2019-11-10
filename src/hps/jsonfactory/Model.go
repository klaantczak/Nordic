package jsonfactory

import (
	"encoding/json"
	"io/ioutil"
	"path"
	"path/filepath"
)

// PropertyDefinition JSON unmarshalling struct.
//
// JSON:
//
//     { "type" : "boolean", "comment" : "..." }
type PropertyDefinition struct {
	Name string

	// Property type name. Can be one of (NOT DEFINED YET).
	Type string

	Comment string
}

// UnmarshalJSON unmarshals the definition to the struct.
func (p *PropertyDefinition) UnmarshalJSON(data []byte) error {
	record := struct {
		Type    string `json:"type"`
		Comment string `json:"comment"`
	}{}
	err := json.Unmarshal(data, &record)
	if err != nil {
		return err
	}
	p.Type = record.Type
	p.Comment = record.Comment
	return nil
}

// MachineDefinition JSON unmarshalling struct.
//
// JSON:
//
//     {
//       "name" : "machine",
//       "type" : "state-machine",
//       "properties" : { ... },
//       "structure" : { ... }
//     }
type MachineDefinition struct {
	// Machine name.
	Name string `json:"name"`

	// Machine type. Currently supported: "state-machine", "network-machine".
	Type string `json:"type"`

	// Properties that can be set for this machine.
	Properties map[string]PropertyDefinition `json:"properties"`

	// Unparsed machine structure. Should be parsed with respect to machine type.
	Structure json.RawMessage `json:"structure"`
}

// UnmarshalJSON unmarshals the definition to the struct.
func (p *MachineDefinition) UnmarshalJSON(data []byte) error {
	record := struct {
		Name       string          `json:"name"`
		Type       string          `json:"type"`
		Properties json.RawMessage `json:"properties"`
		Structure  json.RawMessage `json:"structure"`
	}{}

	err := json.Unmarshal(data, &record)
	if err != nil {
		return err
	}

	p.Name = record.Name

	p.Type = record.Type

	definitions := map[string]PropertyDefinition{}

	properties := ([]byte)(record.Properties)
	if len(properties) > 0 {
		err = json.Unmarshal(([]byte)(record.Properties), &definitions)
		if err != nil {
			return err
		}
	}

	p.Properties = definitions

	p.Structure = record.Structure

	return nil
}

// NetworkDefinition JSON unmarshalling struct.
//
// JSON:
//
//     {
//       "name": "main",
//       "machines": [ ... ]
//     }
type NetworkDefinition struct {
	// Network name.
	Name string `json:"name"`

	// Instances of the machines defined in the machines section.
	Machines []InstanceDefinition `json:"machines"`
}

// InstanceDefinition JSON unmarshalling struct.
//
// JSON:
//
//     {
//       "name": "Machine1",
//       "machine": "SuperMachine",
//       "properties": {
//         "superpower": "activated",
//         ...
//       }
//     }
type InstanceDefinition struct {
	// Machine instance name
	Name string `json:"name"`

	// Machine type, from the list of machines.
	Machine string `json:"machine"`

	// Properties of the machine instance.
	Properties map[string]json.RawMessage `json:"properties"`
}

// ProjectDefinition JSON unmarshalling struct.
//
// JSON:
//
//     { machines : [ ... ], networks : [ ... ] }
type ProjectDefinition struct {
	// Project machines.
	Machines []MachineDefinition `json:"machines"`

	// Project networks.
	Networks []NetworkDefinition `json:"networks"`

	// Path to the project file, if loaded from the file.
	Path string
}

// UnmarshalJSON unmarshals the definition to the struct.
func (p *ProjectDefinition) UnmarshalJSON(data []byte) error {
	record := struct {
		Machines []json.RawMessage `json:"machines"`
		Networks []json.RawMessage `json:"networks"`
	}{}

	err := json.Unmarshal(data, &record)
	if err != nil {
		return err
	}

	for _, md := range record.Machines {
		include := struct {
			Path string `json:"include"`
		}{}
		err = json.Unmarshal(([]byte)(md), &include)
		if err != nil {
			return err
		}

		machine := MachineDefinition{}
		var content []byte

		if include.Path != "" {
			path := p.resolvePath(include.Path)
			content, err = ioutil.ReadFile(path)
			if err != nil {
				return err
			}
		} else {
			content = ([]byte)(md)
		}

		err = json.Unmarshal(content, &machine)
		if err != nil {
			return err
		}

		p.Machines = append(p.Machines, machine)
	}

	for _, md := range record.Networks {
		include := struct {
			Path string `json:"include"`
		}{}
		err = json.Unmarshal(([]byte)(md), &include)
		if err != nil {
			return err
		}

		network := NetworkDefinition{}
		var content []byte

		if include.Path != "" {
			path := p.resolvePath(include.Path)
			content, err = ioutil.ReadFile(path)
			if err != nil {
				return err
			}
		} else {
			content = ([]byte)(md)
		}

		err = json.Unmarshal(content, &network)
		if err != nil {
			return err
		}

		p.Networks = append(p.Networks, network)
	}

	return nil
}

func (pd *ProjectDefinition) resolvePath(filePath string) string {
	dir := filepath.Dir(pd.Path)
	return path.Join(dir, filePath)
}

func (pd *ProjectDefinition) FindMachineDefinition(name string) (*MachineDefinition, bool) {
	for _, m := range pd.Machines {
		if m.Name != name {
			continue
		}

		return &m, true
	}
	return nil, false
}

func (pd *ProjectDefinition) FindNetworkDefinition(name string) (*NetworkDefinition, bool) {
	for _, n := range pd.Networks {
		if n.Name != name {
			continue
		}

		return &n, true
	}
	return nil, false
}
