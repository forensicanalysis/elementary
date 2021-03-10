package pluginlib

import (
	"encoding/csv"
	"fmt"
	"strings"
)

type ParameterType int

func (t ParameterType) IsList() bool {
	return t == StringArray || t == PathArray
}

const (
	_ ParameterType = iota
	Bool
	String
	StringArray
	Path
	PathArray
)

type Parameter struct {
	Name        string
	Description string
	Type        ParameterType
	Value       interface{}
	Required    bool
	Argument    bool
}

func (p *Parameter) BoolValue() bool {
	if p.Value == nil {
		return false
	}
	if s, ok := p.Value.(bool); ok {
		return s
	}
	panic(fmt.Errorf("parameter %s is not a bool: %T", p.Name, p.Value))
}

func (p *Parameter) StringValue() string {
	if p.Value == nil {
		return ""
	}
	if s, ok := p.Value.(string); ok {
		return s
	}
	panic(fmt.Errorf("parameter %s is not a string: %T", p.Name, p.Value))
}

func (p *Parameter) StringArray() []string {
	if p.Value == nil {
		return nil
	}
	if s, ok := p.Value.([]string); ok {
		return s
	}
	panic(fmt.Errorf("parameter %s is not a string array: %T", p.Name, p.Value))
}

type ParameterList []*Parameter

func (pl ParameterList) Get(name string) (*Parameter, error) {
	for _, p := range pl {
		if p.Name == name {
			return p, nil
		}
	}
	return nil, fmt.Errorf("parameter %s not found", name)
}

func (pl ParameterList) BoolValue(name string) bool {
	p, err := pl.Get(name)
	if err != nil {
		panic(err)
	}
	return p.BoolValue()
}

func (pl ParameterList) StringValue(name string) string {
	p, err := pl.Get(name)
	if err != nil {
		panic(err)
	}
	return p.StringValue()
}

func (pl ParameterList) GetStringArrayValue(name string) []string {
	p, err := pl.Get(name)
	if err != nil {
		panic(err)
	}
	return p.StringArray()
}

func (pl ParameterList) Set(name string, value interface{}) {
	p, err := pl.Get(name)
	if err != nil {
		panic(fmt.Sprintf("parameter %s not found", name))
	}
	p.Value = value
}

func (pl ParameterList) ToCommandlineArgs() []string {
	var cmdArgs []string
	for _, p := range pl {
		if p.Argument {
			cmdArgs = append(cmdArgs, p.Name)
		}

		value := fmt.Sprint(p.Value)
		if p.Type.IsList() && strings.HasPrefix(value, "[") && strings.HasSuffix(value, "]") {
			slice, err := readAsCSV(strings.TrimSuffix(strings.TrimPrefix(value, "["), "]"))
			if err == nil {
				for _, value := range slice {
					cmdArgs = append(cmdArgs, fmt.Sprintf("--%s=%s", p.Name, value))
				}
				continue
			}
		}
		if p.Type == Bool {
			if p.BoolValue() {
				cmdArgs = append(cmdArgs, fmt.Sprintf("--%s", p.Name))
			}
			continue
		}
		cmdArgs = append(cmdArgs, fmt.Sprintf("--%s=%s", p.Name, value))
	}
	return cmdArgs
}

func readAsCSV(val string) ([]string, error) {
	if val == "" {
		return []string{}, nil
	}
	stringReader := strings.NewReader(val)
	csvReader := csv.NewReader(stringReader)
	return csvReader.Read()
}

type Property struct {
	Type        string      `json:"type,omitempty"`
	Description string      `json:"description,omitempty"`
	Default     interface{} `json:"default,omitempty"`
}

type JSONSchema struct {
	Properties map[string]Property `json:"properties,omitempty"`
	Required   []string            `json:"required,omitempty"`
}

func JsonschemaToParameter(schema JSONSchema) []*Parameter {
	var parameters []*Parameter
	for name, property := range schema.Properties {
		p := &Parameter{Name: name, Description: property.Description}
		switch property.Type {
		case "string":
			p.Type = String
			if defaultValue, ok := property.Default.(string); ok {
				p.Value = defaultValue
			} else {
				p.Value = ""
			}
		case "boolean":
			p.Type = Bool
			if defaultValue, ok := property.Default.(bool); ok {
				p.Value = defaultValue
			} else {
				p.Value = false
			}
		default:
			panic(fmt.Sprintf("unknown jsonschema type %s", property.Type))
		}
		if contains(schema.Required, name) {
			p.Required = true
		}
		parameters = append(parameters, p)
	}
	return parameters
}

func contains(list []string, elem string) bool {
	for _, i := range list {
		if i == elem {
			return true
		}
	}
	return false
}
