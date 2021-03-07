package daggy

import "fmt"

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

func (p Parameter) BoolValue() bool {
	if p.Value == nil {
		return false
	}
	if s, ok := p.Value.(bool); ok {
		return s
	}
	panic(fmt.Errorf("parameter %s is not a bool", p.Name))
}

func (p Parameter) StringValue() string {
	if p.Value == nil {
		return ""
	}
	if s, ok := p.Value.(string); ok {
		return s
	}
	panic(fmt.Errorf("parameter %s is not a string", p.Name))
}

func (p Parameter) StringArrayValue() []string {
	if p.Value == nil {
		return nil
	}
	if s, ok := p.Value.([]string); ok {
		return s
	}
	panic(fmt.Errorf("parameter %s is not a string array", p.Name))
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
	return p.StringArrayValue()
}

func (pl ParameterList) Set(name string, value interface{}) {
	p, err := pl.Get(name)
	if err != nil {
		panic(fmt.Sprintf("parameter %s not found", name))
	}
	p.Value = value
}
