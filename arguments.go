package main

import "strings"

// data structures to handle each argument of a line/group
// conforming to an Argument interface to ensure its maintainability in the future as more argument types are added or if they are to be further refined
type Argument interface {
	GetValue() any // retrieves the value associated with the relevant argument
}

type OptionArgument struct {
	value string
}

func newOptionArgument(value string) Argument {
	return &OptionArgument{
		value: value,
	}
}

func (arg *OptionArgument) GetValue() any {
	return arg.value
}

type ClassArgument struct {
	value string
}

func newClassArgument(value string) Argument {
	return &ClassArgument{
		value: value,
	}
}

func (arg *ClassArgument) GetValue() any {
	return arg.value
}

type KeyValueArgument struct {
	value map[string]string
}

func newKeyValueArgument(value string) (Argument, CreateError) {
	mapping := make(map[string]string)
	pairs := strings.Split(value, ",")
	for _, pair := range pairs {
		data := strings.Split(pair, "=")
		if len(data) != 2 {
			return nil, INVALID_ARGUMENT_CONTENT_KEY_VALUE_FORMAT
		}
		key := strings.Trim(data[0], WHITESPACE)
		value := strings.Trim(data[1], WHITESPACE)
		if key == "" {
			return nil, KEY_EMPTY
		}
		if value == "" {
			return nil, VALUE_EMPTY
		}
		mapping[key] = value
	}

	return &KeyValueArgument{
		value: mapping,
	}, nil
}

func (arg *KeyValueArgument) GetValue() any {
	return arg.value
}

func (arg *KeyValueArgument) GetSelectedValue(key string) (string, bool) {
	res, ok := arg.value[key]
	return res, ok
}
