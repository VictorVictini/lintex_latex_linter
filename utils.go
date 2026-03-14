package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// helper function to return the start and end coordinates when provided the starting coordinate information and text
func GetCoordinates(line int, column int, text []byte) (*Coordinate, *Coordinate, CreateError) {
	// read file contents
	fileBytes, err := GetFileBytes("latex_testing/test.tex")
	if err != nil {
		return nil, nil, err
	}

	// useful variables
	fileContents := string(fileBytes)
	lines := strings.Split(fileContents, "\n")
	startCoord := newCoordinate(line, column)

	// calculate where to start extracting the line
	startPosition := startCoord.charPos - 1
	for _, currLine := range lines[0 : startCoord.line-1] {
		startPosition += len(currLine) + 1
	}

	// calculate line where end coordinate finishes at
	endLine := 0
	endPosition := startPosition + len(text) + 1
	for _, line := range lines {
		endLine++
		if endPosition < len(line) && endPosition >= 0 {
			break
		}
		endPosition -= len(line) + 1
	}

	// handling an edge case where endCoord reaches the EOF
	endCoord := newCoordinate(endLine, endPosition)
	if endCoord.line == startCoord.line && endCoord.charPos == 0 && endCoord.charPos < startCoord.charPos {
		endCoord = newCoordinate(endLine, len(lines[len(lines)-1])+1)
	}
	return &startCoord, &endCoord, nil
}

// helper function to retrieve a file's contents
func GetFileBytes(fileLocation string) ([]byte, CreateError) {
	path := filepath.Join(GLOBAL_STORAGE_LOCATION, "latex_testing/test.tex")
	fileBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, SERVER_RESPONSIBLE_READ_FILE_ERROR
	}
	return fileBytes, nil
}

// helper function to convert any interfaces (any -> []interface{} -> [][]uint8) to string a string assuming the underlying interfaces conforms to a string format
func AnyInterfaceToString(interfaceArr any) string {
	interArr := interfaceArr.([]interface{})
	var res string
	for _, v := range interArr {
		for _, val := range v.([]uint8) {
			res += fmt.Sprintf("%c", val)
		}
	}
	return res
}

// helper function to remove nil arguments from a slice
func RemoveNilElements(slice []any) []any {
	if slice == nil {
		return slice
	}
	var res []any
	for _, v := range slice {
		if v != nil {
			res = append(res, v)
		}
	}
	return res
}

// helper function to convert any interfaces (any -> []interface{} -> []Component) to Component slices assuming the underlying interfaces conform to Component slices
func AnyInterfaceToComponentSlice(original any) []Component {
	origin := original.([]interface{})
	var res []Component
	for _, v := range origin {
		if v != nil {
			res = append(res, v.(Component))
		} else {
			res = append(res, nil)
		}
	}
	return res
}

// helper function to convert any interfaces (any -> []interface{} -> []T) to a slice containing values of type T, assuming the parameter provided confirms to a T slice
func AnyInterfaceToTSlice[T any](original any) []T {
	origin := original.([]interface{})
	var res []T
	for _, value := range origin {
		res = append(res, value.(T))
	}
	return res
}

// helper function that checks that every value in a list, required, is contained in the provided map
func AllContains(dict map[string]bool, required []string) bool {
	for _, expected := range required {
		if !dict[expected] {
			return false
		}
	}
	return true
}
