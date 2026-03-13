package main

import "strconv"

type ErrorInitial string

const (
	ParseErrorInitial         ErrorInitial = "PE"
	StructureErrorInitial     ErrorInitial = "SE"
	AccessibilityErrorInitial ErrorInitial = "AE"
)

type CustomError interface {
	RetrieveID() string  // Retrieves the error message's ID to make it easier to identify for the user
	Error() string       // Retrieves the short error message to be used when it occurs
	LongError() string   // Provides a longer description of the error to be used
	LocateError() string // Returns the full start and end location of the LaTeX component that caused the error (including its line and character count information)
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

func (err *ParseError) LocateError() string {
	// read file contents
	// get the line at the relevant position
	// return it
	return "..."
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

func (err *StructureError) LocateError() string {
	// read file contents
	// get the line at the relevant position
	// return it
	return "..."
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

func (err *AccessibilityError) LocateError() string {
	// read file contents
	// get the line at the relevant position
	// return it
	return "..."
}

// a data structure to handle locating at what point in the file the error had occurred
type Coordinate struct {
	line, charPos int
}

func newCoordinate(line int, charPos int) Coordinate {
	return Coordinate{
		line:    line,
		charPos: charPos,
	}
}
