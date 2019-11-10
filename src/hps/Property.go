package hps

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type PropertyValueJSONUnmarshaller func([]byte) (interface{}, error)

var PropertyValueJSONUnmarshallers = map[string]PropertyValueJSONUnmarshaller{
	"String": func(data []byte) (interface{}, error) {
		var value string
		err := json.Unmarshal(data, &value)
		return interface{}(value), err
	},
	"Integer": func(data []byte) (interface{}, error) {
		var value int64
		err := json.Unmarshal(data, &value)
		return interface{}(value), err
	},
	"Float": func(data []byte) (interface{}, error) {
		var value float64
		err := json.Unmarshal(data, &value)
		return interface{}(value), err
	},
}

func RegisterPropertyValueJSONUnmarshaller(typeName string, unmarshaller PropertyValueJSONUnmarshaller) {
	PropertyValueJSONUnmarshallers[typeName] = unmarshaller
}

func PropertyValueUnmarshalJSON(typeName string, data []byte) (interface{}, error) {
	unmarshaller, ok := PropertyValueJSONUnmarshallers[typeName]
	if !ok {
		return nil, fmt.Errorf("The specified type cannot be unmarshalled.")
	}
	return unmarshaller(data)
}

type PropertyEventHandler func(p *Property, oldValue interface{})

type Property struct {
	name     string
	dataType string
	required bool
	value    interface{}
	comment  string
	changed  []PropertyEventHandler
}

func NewProperty(name string, dataType string) *Property {
	p := &Property{}
	p.name = name
	p.dataType = dataType
	p.required = true
	p.changed = []PropertyEventHandler{}
	return p
}

func (p *Property) Name() string {
	return p.name
}

func (p *Property) Type() string {
	return p.dataType
}

func (p *Property) Value() interface{} {
	return p.value
}

func (p *Property) GetBool() (bool, bool) {
	v, ok := p.value.(bool)
	return v, ok
}

func (p *Property) GetFloat() (float64, bool) {
	v, ok := p.value.(float64)
	return v, ok
}

func (p *Property) GetInt() (int, bool) {
	v, ok := p.value.(int)
	return v, ok
}

func (p *Property) GetString() (string, bool) {
	v, ok := p.value.(string)
	return v, ok
}

func (p *Property) SetValue(value interface{}) {
	if p.value == value {
		return
	}

	oldValue := p.value
	p.value = value
	for _, h := range p.changed {
		h(p, oldValue)
	}
}

func (p *Property) SetComment(comment string) {
	p.comment = comment
}

func (p *Property) Changed(handler PropertyEventHandler) {
	p.changed = append(p.changed, handler)
}

func (p *Property) MarshalJSON() ([]byte, error) {
	buf := bytes.Buffer{}

	buf.WriteString("{")

	buf.WriteString("\"name\":")
	name, _ := json.Marshal(p.Name())
	buf.Write(name)

	var value interface{}
	switch p.value.(type) {
	case bool, int, float64, string:
		value = p.value
	default:
		if js, ok := p.value.(IJSON); ok {
			value = js.ToJSON()
		} else {
			value = p.value
		}
	}
	data, _ := json.Marshal(value)
	buf.WriteString(",\"value\":")
	buf.Write(data)

	buf.WriteString(",\"required\":")
	required, _ := json.Marshal(p.required)
	buf.Write(required)

	if p.comment != "" {
		buf.WriteString(",\"comment\":")
		comment, _ := json.Marshal(p.comment)
		buf.Write(comment)
	}

	buf.WriteString("}")

	return buf.Bytes(), nil
}
