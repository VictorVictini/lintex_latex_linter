// this file implements the data structure that the parsed files are inserted into
// using the Composite design pattern, to treat individual lines and groupings (using begin/similar mechanisms) as the same

package main

import (
	"fmt"
	"strings"
)

type Component interface {
	CheckGroupings(groups map[string]bool) bool // provides a list of all current groupings to the nested layers
	GetName() string                            // provides the name of the relevant component
	PrintTree(depth int)                        // prints the full data structure recursively
}

// a single LaTeX line
// acting as the Leaf
type Line struct {
	name      string
	arguments []Argument

	requiredGroups []string
}

func newLine(name string, arguments []Argument) Component {
	return &Line{
		name:      name,
		arguments: arguments,
	}
}

// checks if all groups within the line's requiredGroups are within the currently existing groups
func (line *Line) CheckGroupings(groups map[string]bool) bool {
	return AllContains(groups, line.requiredGroups)
}

func (line *Line) GetName() string {
	return line.name
}

func (line *Line) PrintTree(depth int) {
	fmt.Printf("%s%s\n", strings.Repeat("\t", depth), line.GetName())
}

// a grouping of LaTeX lines
// acting as the Composite
type Group struct {
	name      string
	arguments []Argument

	requiredGroups []string

	components []Component
}

func newGroup(name string, arguments []Argument) Component {
	return &Group{
		name:      name,
		arguments: arguments,
	}
}

// adds the current group to the set, and checks if nested lines all contain valid groupings
func (group *Group) CheckGroupings(groups map[string]bool) bool {
	if !AllContains(groups, group.requiredGroups) {
		return false
	}

	groups[group.name] = true
	for _, component := range group.components {
		if !component.CheckGroupings(groups) {
			return false
		}
	}
	return true
}

func (group *Group) GetName() string {
	return group.name
}

func (group *Group) AddComponent(component Component) {
	group.components = append(group.components, component)
}

func (group *Group) PrintTree(depth int) {
	fmt.Printf("%s%s\n", strings.Repeat("\t", depth), group.GetName())
	for _, component := range group.components {
		component.PrintTree(depth + 1)
	}
}

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
