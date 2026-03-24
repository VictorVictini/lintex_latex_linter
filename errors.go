package main

import (
	"strconv"
	"strings"
)

type ErrorInitial string

const (
	ParseErrorInitial         ErrorInitial = "PE"
	StructureErrorInitial     ErrorInitial = "STE"
	AccessibilityErrorInitial ErrorInitial = "AE"
	ServerErrorInitial        ErrorInitial = "SE"
)

type CustomError interface {
	RetrieveID() string             // Retrieves the error message's ID to make it easier to identify for the user
	Error() string                  // Retrieves the short error message to be used when it occurs
	LongError() string              // Provides a longer description of the error to be used
	LocateError() (string, error)   // Returns the full start and end location of the LaTeX component that caused the error (including its line and character count information)
	GetStartCoordinate() Coordinate // Returns the start coordinate for where the error had occured
	GetEndCoordinate() Coordinate   // Returns the end coordinate for where the error had occured
}

// Parsing errors : errors related to incorrect parsing of the document
// such as invalid {} [] or begin\end usage
type ParseError struct {
	// information relevant to where the error had occurred
	fileLoc              string
	startCoord, endCoord Coordinate

	// information relevant to the information provided by the error
	shortDesc, longDesc string

	// information relevant to the identification of the error
	id      int
	initial ErrorInitial
}

func newParseError(id int, startCoord Coordinate, endCoord Coordinate, shortDesc string, longDesc string) CustomError {
	return &ParseError{
		startCoord: startCoord,
		endCoord:   endCoord,
		shortDesc:  shortDesc,
		longDesc:   longDesc,
		id:         id,
		initial:    ParseErrorInitial,
	}
}

func (err *ParseError) RetrieveID() string {
	return string(err.initial) + strconv.Itoa(err.id)
}

func (err *ParseError) Error() string {
	return err.shortDesc
}

func (err *ParseError) LongError() string {
	return err.longDesc
}

func (err *ParseError) LocateError() (string, error) {
	return LocateError(err)
}

func (err *ParseError) GetStartCoordinate() Coordinate {
	return err.startCoord
}

func (err *ParseError) GetEndCoordinate() Coordinate {
	return err.endCoord
}

type StructureError struct {
	// information relevant to where the error had occurred
	fileLoc              string
	startCoord, endCoord Coordinate

	// information relevant to the information provided by the error
	shortDesc, longDesc string

	// information relevant to the identification of the error
	id      int
	initial ErrorInitial
}

// Structure errors : errors such as missing/mistyping a docclass, having grouping contexts in preamble/prerequisite, having more than 1 component in doc content, missing essential arguments to something (e.g. first class arg to \begin, \end, and \documentclass), doc content not being a begin group containing document
// mostly handled within IDocumentBuilder.buildDocument()
func newStructureError(id int, startCoord Coordinate, endCoord Coordinate, shortDesc string, longDesc string) CustomError {
	return &StructureError{
		startCoord: startCoord,
		endCoord:   endCoord,
		shortDesc:  shortDesc,
		longDesc:   longDesc,
		id:         id,
		initial:    StructureErrorInitial,
	}
}

func (err *StructureError) RetrieveID() string {
	return string(err.initial) + strconv.Itoa(err.id)
}

func (err *StructureError) Error() string {
	return err.shortDesc
}

func (err *StructureError) LongError() string {
	return err.longDesc
}

func (err *StructureError) LocateError() (string, error) {
	return LocateError(err)
}

func (err *StructureError) GetStartCoordinate() Coordinate {
	return err.startCoord
}

func (err *StructureError) GetEndCoordinate() Coordinate {
	return err.endCoord
}

type AccessibilityError struct {
	// information relevant to where the error had occurred
	fileLoc              string
	startCoord, endCoord Coordinate

	// information relevant to the information provided by the error
	shortDesc, longDesc string

	// information relevant to the identification of the error
	id      int
	initial ErrorInitial
}

// Accessibility errors : errors related to incorrect accessibility of the document (using outdated methods) or missing accessibility methods (such as not including \Metadata for tagging or relevant code for alt text on figures, images, or tables), may also include some nuances such as table design, text size, colours(?))
func newAccessibilityError(id int, startCoord Coordinate, endCoord Coordinate, shortDesc string, longDesc string) CustomError {
	return &AccessibilityError{
		startCoord: startCoord,
		endCoord:   endCoord,
		shortDesc:  shortDesc,
		longDesc:   longDesc,
		id:         id,
		initial:    AccessibilityErrorInitial,
	}
}

func (err *AccessibilityError) RetrieveID() string {
	return string(err.initial) + strconv.Itoa(err.id)
}

func (err *AccessibilityError) Error() string {
	return err.shortDesc
}

func (err *AccessibilityError) LongError() string {
	return err.longDesc
}

func (err *AccessibilityError) LocateError() (string, error) {
	return LocateError(err)
}

func (err *AccessibilityError) GetStartCoordinate() Coordinate {
	return err.startCoord
}

func (err *AccessibilityError) GetEndCoordinate() Coordinate {
	return err.endCoord
}

type ServerError struct {
	// information relevant to where the error had occurred
	fileLoc              string
	startCoord, endCoord Coordinate

	// information relevant to the information provided by the error
	shortDesc, longDesc string

	// information relevant to the identification of the error
	id      int
	initial ErrorInitial
}

// Accessibility errors : errors related to incorrect accessibility of the document (using outdated methods) or missing accessibility methods (such as not including \Metadata for tagging or relevant code for alt text on figures, images, or tables), may also include some nuances such as table design, text size, colours(?))
func newServerError(id int, startCoord Coordinate, endCoord Coordinate, shortDesc string, longDesc string) CustomError {
	return &AccessibilityError{
		startCoord: startCoord,
		endCoord:   endCoord,
		shortDesc:  shortDesc,
		longDesc:   longDesc,
		id:         id,
		initial:    ServerErrorInitial,
	}
}

func (err *ServerError) RetrieveID() string {
	return string(err.initial) + strconv.Itoa(err.id)
}

func (err *ServerError) Error() string {
	return err.shortDesc
}

func (err *ServerError) LongError() string {
	return err.longDesc
}

func (err *ServerError) LocateError() (string, error) {
	return LocateError(err)
}

func (err *ServerError) GetStartCoordinate() Coordinate {
	return err.startCoord
}

func (err *ServerError) GetEndCoordinate() Coordinate {
	return err.endCoord
}

func LocateError(err CustomError) (string, CustomError) {
	// read file contents
	fileBytes, fileErr := GetFileBytes(GLOBAL_FILE_NAME)
	if fileErr != nil {
		return "", DummyError(SERVER_RESPONSIBLE_READ_FILE_ERROR)
	}

	// useful variables
	fileContents := string(fileBytes)
	startCoord := err.GetStartCoordinate()
	endCoord := err.GetEndCoordinate()

	// calculate where to start extracting the line
	startPosition := CalculateOffset(fileBytes, startCoord)

	// calculate where to end extracting of the line
	endPosition := CalculateOffset(fileBytes, endCoord)

	// extract the line
	if endPosition <= startPosition {
		return "", nil
	}
	extract := string(fileContents[startPosition:endPosition])
	return strings.Trim(extract, WHITESPACE), nil
}

// a data structure to handle locating at what point in the file the error had occurred
type Coordinate struct {
	line    int
	charPos int
}

func newCoordinate(line int, charPos int) Coordinate {
	return Coordinate{
		line:    line,
		charPos: charPos,
	}
}
