// this file implements the data structure that the parsed files are inserted into
// using the Composite design pattern, to treat individual lines and groupings (using begin/similar mechanisms) as the same

package main

type Component interface {
	CheckGroupings(groups map[string]bool) bool // provides a list of all current groupings to the nested layers
	GetName() string                            // provides the name of the relevant component
}

// a single LaTeX line
// acting as the Leaf
type Line struct {
	name      string
	arguments []Argument

	requiredGroups []string
}

func newLine(name string, arguments []Argument) *Line {
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

// a grouping of LaTeX lines
// acting as the Composite
type Group struct {
	name      string
	arguments []Argument

	requiredGroups []string

	components []Component
}

func newGroup(name string, arguments []Argument) *Group {
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

// data structure to handle each argument of a line/group
type Argument struct {
	argumentType ArgumentType
	value        string
}

func newArgument(argumentType ArgumentType, value string) Argument {
	return Argument{
		argumentType: argumentType,
		value:        value,
	}
}
