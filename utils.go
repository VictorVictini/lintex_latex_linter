package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func CompileLaTeXFile(fileName string) (string, CreateError) {
	// get working directory
	workingDir, errFn := GetWorkingDirectory()
	if errFn != nil {
		return "", errFn
	}

	// run lualatex command (basing it on latex folder)
	cmd := exec.Command("lualatex", "--halt-on-error", "--output-directory="+filepath.Join(workingDir, GLOBAL_LATEX_FOLDER), filepath.Join(workingDir, GLOBAL_LATEX_FOLDER, fileName+".tex"))
	_, err := cmd.Output()
	if err != nil {
		return "", CANNOT_CONVERT_TO_PDF
	}

	// return the created file
	return filepath.Join(workingDir, fileName+".pdf"), nil
}

// helper function to get the working directory
func GetWorkingDirectory() (string, CreateError) {
	ex, err := os.Executable()
	if err != nil {
		return "", SERVER_RESPONSIBLE_CANNOT_FIND_WORKING_DIRECTORY
	}
	exPath := filepath.Dir(ex)
	return exPath, nil
}

// helper function to structure an err list's contents into a fixed format
func FormatErrors(list errList) []CustomError {
	res := make([]CustomError, 0)
	for _, err := range list {
		var cerr CustomError
		pe, ok := err.(*parserError)
		if ok {
			cerr, ok = pe.Inner.(CustomError)
		} else {
			cerr, ok = err.(CustomError)
		}
		if ok {
			res = append(res, cerr)
		} else {
			res = append(res, DummyError(SERVER_RESPONSIBLE_UNIDENTIFIED_ERROR))
		}
	}
	return res
}

// helper function to calculate the integer offset to reach the coordinate in a given file
func CalculateOffset(fileBytes []byte, coord Coordinate) int {
	fileContents := string(fileBytes)
	lines := strings.Split(fileContents, "\n")
	position := coord.charPos - 1
	fmt.Printf("len(lines) %#v\ncoord.line - 1 %#v\n", len(lines), coord.line-1)
	for _, line := range lines[0 : coord.line-1] {
		position += len(line) + 1
	}
	return position
}

// helper function for returning a slice of lines contains a given name
func FindLinesWithName(lines []Line, name string) []Line {
	res := make([]Line, 0)
	for _, line := range lines {
		if line.GetName() == name {
			res = append(res, line)
		}
	}
	return res
}

// helper function to return the start and end coordinates when provided the starting coordinate information and text
func GetCoordinates(line int, column int, text []byte) (*Coordinate, *Coordinate, CreateError) {
	// read file contents
	fileBytes, err := GetFileBytes(GLOBAL_FILE_NAME)
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
// fileName should not include extensions (assumed .tex)
func GetFileBytes(fileName string) ([]byte, CreateError) {
	// get working directory
	workingDir, errFn := GetWorkingDirectory()
	if errFn != nil {
		return []byte{}, errFn
	}

	// use cache if possible
	if len(GLOBAL_FILE_DATA_CACHE) > 0 {
		return GLOBAL_FILE_DATA_CACHE, nil
	}

	// read from the file
	path := filepath.Join(workingDir, GLOBAL_LATEX_FOLDER, fileName+".tex")
	fileBytes, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("read file error: %#v\n", err.Error())
		return nil, SERVER_RESPONSIBLE_READ_FILE_ERROR
	}
	return fileBytes, nil
}

// fileName should not include extensions (assumed .tex)
func WriteFileBytes(fileName string, data []byte) CreateError {
	// get working directory
	workingDir, errFn := GetWorkingDirectory()
	if errFn != nil {
		return errFn
	}

	// could not write to file
	path := filepath.Join(workingDir, GLOBAL_LATEX_FOLDER, fileName+".tex")
	err := os.WriteFile(path, data, 0644)
	if err != nil {
		fmt.Printf("write file error: %#v\n", err.Error())
		return SERVER_RESPONSIBLE_WRITE_FILE_ERROR
	}
	return nil
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
func RemoveNilElements[T comparable](slice []T, nilValue T) []T {
	if slice == nil {
		return slice
	}
	var res []T
	for _, v := range slice {
		if v != nilValue {
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
