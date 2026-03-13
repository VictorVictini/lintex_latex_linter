package main

// a data structure to ensure that the IDs assigned to a given error function remains unique
type ErrorIDGenerator struct {
	idList  map[int]bool
	current int
}

// constructor for the Error ID generator
func newErrorIDGenerator() ErrorIDGenerator {
	return ErrorIDGenerator{
		idList:  make(map[int]bool, 0),
		current: 1,
	}
}

// ensures the data structure returns a unique ID
func (generator *ErrorIDGenerator) CreateUniqueID() int {
	for generator.idList[generator.current] {
		generator.current++
	}
	return generator.current
}

// to ensure the data (ID, ErrorInitial, Descriptions) is unchanged by the function creating the error, we use a wrapper to ensure that data is fixed
type CreateError func(Coordinate, Coordinate) CustomError

// A wrapper function to create functions to create errors with
func CustomErrorWrapper(shortDesc string, longDesc string, newCustomError func(int, Coordinate, Coordinate, string, string) CustomError) CreateError {
	id := GLOBAL_ERROR_ID_GENERATOR.CreateUniqueID()
	return func(startCoord Coordinate, endCoord Coordinate) CustomError {
		return newCustomError(id, startCoord, endCoord, shortDesc, longDesc)
	}
}

// A wrapper function to return an error's information by itself, by placing dummy data where a specific error's details would be provided
func DummyError(newCustomErrorWrapper CreateError) CustomError {
	return newCustomErrorWrapper(newCoordinate(0, 0), newCoordinate(0, 0))
}

// Acts as a shorthand list of error functions
// Restructuring the ordering at this point may change an error's ID,
// However, the code is set up such that it can support such a change
var (
	/**
	 * Parsing related error functions.
	 */
	// Argument parsing errors
	OPTIONS_ARGUMENT_START_MISSING = CustomErrorWrapper("short_desc", "long desc", newParseError)
	OPTIONS_ARGUMENT_END_MISSING   = CustomErrorWrapper("short desc", "long desc", newParseError)
	CLASS_ARGUMENT_START_MISSING   = CustomErrorWrapper("short desc", "long desc", newParseError)
	CLASS_ARGUMENT_END_MISSING     = CustomErrorWrapper("short desc", "long desc", newParseError)

	// Group parsing errors
	GROUP_NAME_MISMATCH = CustomErrorWrapper("short desc", "long desc", newParseError)
	GROUP_BEGIN_MISSING = CustomErrorWrapper("short desc", "long desc", newParseError)
	GROUP_END_MISSING   = CustomErrorWrapper("short desc", "long desc", newParseError)

	// Server-responsible parsing errors
	SERVER_RESPONSIBLE_STACK_EMPTY                      = CustomErrorWrapper("short desc", "long desc", newParseError)
	SERVER_RESPONSIBLE_STACK_CONTAINS_NON_GROUP_ELEMENT = CustomErrorWrapper("short desc", "long desc", newParseError)

	/**
	 * Structure related error functions
	 */
	// Structure errors :
	// missing essential arguments to something (e.g. first class arg to \begin, \end, and \documentclass), doc content not being a begin group containing document
	// mostly handled within IDocumentBuilder.buildDocument()

	// Missing essential element structure errors
	DOCUMENT_CLASS_NOT_FOUND = CustomErrorWrapper("short desc", "long desc", newStructureError)
	DOCUMENT_GROUP_NOT_FOUND = CustomErrorWrapper("short desc", "long desc", newStructureError)

	// Containing group constructs within non-grouping contexts
	PREAMBLE_CONTAINS_GROUP     = CustomErrorWrapper("short desc", "long desc", newStructureError)
	PREREQUISITE_CONTAINS_GROUP = CustomErrorWrapper("short desc", "long desc", newStructureError)

	// Extra latex code where it shouldn't belong
	DOCUMENT_CONTENT_ALREADY_EXISTS = CustomErrorWrapper("short desc", "long desc", newStructureError)

	// further examples to be confirmed

	/**
	 * Accessibility related error functions
	 */
	// tbd
)
