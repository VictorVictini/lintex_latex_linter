package main

// data structures to handle each argument of a line/group
// conforming to an Argument interface to ensure its maintainability in the future as more argument types are added or if they are to be further refined
type Argument interface {
	GetValue() string // retrieves the value associated with the relevant argument
}

type OptionArgument struct {
	value string
}

func newOptionArgument(value string) Argument {
	return &OptionArgument{
		value: value,
	}
}

func (arg *OptionArgument) GetValue() string {
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

func (arg *ClassArgument) GetValue() string {
	return arg.value
}
