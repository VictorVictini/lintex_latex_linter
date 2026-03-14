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
	OPTIONS_ARGUMENT_START_MISSING = CustomErrorWrapper("The opening square bracket '[' is missing from the argument.", "long desc", newParseError)
	OPTIONS_ARGUMENT_END_MISSING   = CustomErrorWrapper("The closing square bracket ']' is missing from the argument.", "long desc", newParseError)
	CLASS_ARGUMENT_START_MISSING   = CustomErrorWrapper("The opening curly braces '{' is missing from the argument.", "long desc", newParseError)
	CLASS_ARGUMENT_END_MISSING     = CustomErrorWrapper("The closing curly braces '}' is missing from the argument.", "long desc", newParseError)

	// Group parsing errors
	GROUP_NAME_MISMATCH = CustomErrorWrapper("The class argument within \\begin line did not match the class argument within \\end line.", "long desc", newParseError)
	GROUP_BEGIN_MISSING = CustomErrorWrapper("The matching \\begin line could not be found for the \\end line.", "long desc", newParseError)
	GROUP_END_MISSING   = CustomErrorWrapper("The matching \\end line could not be found for the \\begin line.", "long desc", newParseError)

	// Server-responsible parsing errors
	SERVER_RESPONSIBLE_STACK_EMPTY                      = CustomErrorWrapper("Something went wrong on the server's end. The grouping stack is empty. Please report this to the server owner.", "long desc", newParseError)
	SERVER_RESPONSIBLE_STACK_CONTAINS_NON_GROUP_ELEMENT = CustomErrorWrapper("Something went wrong on the server's end. The grouping stack contained a non-group element. Please report this to the server owner.", "long desc", newParseError)

	/**
	 * Structure related error functions
	 */
	// Structure errors :
	// missing essential arguments to something (e.g. first class arg to \begin, \end, and \documentclass), doc content not being a begin group containing document
	// mostly handled within IDocumentBuilder.buildDocument()

	// Missing essential element structure errors
	DOCUMENT_CLASS_NOT_FOUND         = CustomErrorWrapper("Missing \\documentclass line.", "long desc", newStructureError)
	DOCUMENT_CONTENT_GROUP_NOT_FOUND = CustomErrorWrapper("Missing \\begin{document} and \\end{document} lines.", "long desc", newStructureError)

	// Containing group constructs within non-grouping contexts
	PREAMBLE_CONTAINS_GROUP     = CustomErrorWrapper("The preamble cannot contain any grouping constructs such as \\begin or \\end.", "long desc", newStructureError)
	PREREQUISITE_CONTAINS_GROUP = CustomErrorWrapper("The content before the document class cannot contain any grouping constructs such as \\begin or \\end.", "long desc", newStructureError)

	// Extra latex code where it shouldn't belong
	DOCUMENT_CONTENT_ALREADY_EXISTS = CustomErrorWrapper("You cannot have more than one document content component.", "long desc", newStructureError)

	// further examples to be confirmed

	/**
	 * Accessibility related error functions
	 */
	// tbd
	// https://docs.overleaf.com/writing-and-editing/creating-accessible-pdfs
)
