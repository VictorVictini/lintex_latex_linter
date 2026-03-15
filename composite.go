// this file implements the data structure that the parsed files are inserted into
// using the Composite design pattern, to treat individual lines and groupings (using begin/similar mechanisms) as the same

package main

import (
	"fmt"
	"strings"
)

type Component interface {
	GetName() string                        // provides the name of the relevant component
	PrintTree(depth int)                    // prints the full data structure recursively
	FindAnyLine(name string) *Line          // Finds any line with a given name
	FindAllLines(name string) []Line        // Finds all lines with a given name
	GetStartCoordinate() Coordinate         // Returns the starting coordinate of the component
	GetEndCoordinate() Coordinate           // Returns the ending coordinate of the component
	SetEndCoordinate(coordinate Coordinate) // Assigns a new end coordinate to the component
}

// a single LaTeX line
// acting as the Leaf
type Line struct {
	name      string
	arguments []Argument

	startCoordinate Coordinate
	endCoordinate   Coordinate
}

func newLine(name string, arguments []Argument, startCoord Coordinate, endCoord Coordinate) Component {
	return &Line{
		name:            name,
		arguments:       arguments,
		startCoordinate: startCoord,
		endCoordinate:   endCoord,
	}
}

func (line *Line) GetName() string {
	return line.name
}

func (line *Line) PrintTree(depth int) {
	fmt.Printf("%s%s\n", strings.Repeat("\t", depth), line.GetName())
}

func (line *Line) FindAnyLine(name string) *Line {
	if line.GetName() == name {
		return line
	}
	return nil
}

func (line *Line) FindAllLines(name string) []Line {
	if line.GetName() == name {
		return []Line{*line}
	}
	return nil
}

func (line *Line) GetStartCoordinate() Coordinate {
	return line.startCoordinate
}

func (line *Line) GetEndCoordinate() Coordinate {
	return line.endCoordinate
}

func (line *Line) SetEndCoordinate(coord Coordinate) {
	line.endCoordinate = coord
}

// a grouping of LaTeX lines
// acting as the Composite
type Group struct {
	name      string
	arguments []Argument

	startCoordinate Coordinate
	endCoordinate   Coordinate

	components []Component
}

func newGroup(name string, arguments []Argument, startCoord Coordinate, endCoord Coordinate) Component {
	return &Group{
		name:            name,
		arguments:       arguments,
		startCoordinate: startCoord,
		endCoordinate:   endCoord,
	}
}

func (group *Group) GetName() string {
	return group.name
}

func (group *Group) FindAnyLine(name string) *Line {
	for _, component := range group.components {
		found := component.FindAnyLine(name)
		if found != nil {
			return found
		}
	}
	return nil
}

func (group *Group) FindAllLines(name string) []Line {
	res := make([]Line, 0)
	for _, component := range group.components {
		addSlice := component.FindAllLines(name)
		if addSlice != nil {
			res = append(res, addSlice...)
		}
	}
	return res
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

func (group *Group) GetStartCoordinate() Coordinate {
	return group.startCoordinate
}

func (group *Group) GetEndCoordinate() Coordinate {
	return group.endCoordinate
}

func (group *Group) SetEndCoordinate(coord Coordinate) {
	group.endCoordinate = coord
}
