package main

import (
	"regexp"
	"strings"
)

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
	value = strings.Trim(value, WHITESPACE)
	// ignore empty strings
	if value == "" {
		return &KeyValueArgument{
			value: make(map[string]string, 0),
		}, nil
	}

	// validate overall format
	matched, err := regexp.MatchString(KEY_VALUE_ARGUMENT_REGEX, value)
	if err != nil {
		return nil, SERVER_RESPONSIBLE_INVALID_REGEX_ERROR
	}
	if !matched {
		return nil, INVALID_ARGUMENT_CONTENT_KEY_VALUE_FORMAT
	}

	// parse into key-value pairs
	re := regexp.MustCompile(KEY_VALUE_PAIR_REGEX)
	matches := re.FindAllStringSubmatch(value, -1)
	mapping := make(map[string]string)
	for _, match := range matches {
		if len(match) < 3 {
			return nil, SERVER_RESPONSIBLE_INVALID_REGEX_ERROR
		}
		key := match[1]
		value := match[2]
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
