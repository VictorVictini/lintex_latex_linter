package main

import (
	"regexp"
	"strconv"
)

// a data structure to ensure that the IDs assigned to a given error function remains unique
type ErrorIDGenerator struct {
	idList  map[int]CreateError
	current int
}

// constructor for the Error ID generator
func newErrorIDGenerator() ErrorIDGenerator {
	return ErrorIDGenerator{
		idList:  make(map[int]CreateError, 0),
		current: 1,
	}
}

// helper function to retrieve an error function from its id
func (generator *ErrorIDGenerator) GetErrorFromID(idStr string) (CreateError, bool) {
	// extract the id
	re := regexp.MustCompile(`\d+$`)
	str := re.FindString(idStr)
	if str == "" {
		return SERVER_RESPONSIBLE_INVALID_STRING_ERROR, false
	}
	id, err := strconv.Atoi(str)
	if err != nil {
		return SERVER_RESPONSIBLE_INVALID_STRING_ERROR, false
	}

	// return the error
	errFn, ok := generator.idList[id]
	return errFn, ok
}

// ensures the data structure returns a unique ID
func (generator *ErrorIDGenerator) CreateUniqueID() int {
	for true {
		_, ok := generator.idList[generator.current]
		if !ok {
			break
		}
		generator.current++
	}
	return generator.current
}

func (generator *ErrorIDGenerator) SetIDLink(id int, errFn CreateError) {
	generator.idList[id] = errFn
}

// to ensure the data (ID, ErrorInitial, Descriptions) is unchanged by the function creating the error, we use a wrapper to ensure that data is fixed
type CreateError func(Coordinate, Coordinate) CustomError

// A wrapper function to create functions to create errors with
func CustomErrorWrapper(shortDesc string, longDesc string, newCustomError func(int, Coordinate, Coordinate, string, string) CustomError) CreateError {
	id := GLOBAL_ERROR_ID_GENERATOR.CreateUniqueID()
	errFn := func(startCoord Coordinate, endCoord Coordinate) CustomError {
		return newCustomError(id, startCoord, endCoord, shortDesc, longDesc)
	}
	GLOBAL_ERROR_ID_GENERATOR.SetIDLink(id, errFn)
	return errFn
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
	OPTIONS_ARGUMENT_START_MISSING            = CustomErrorWrapper("The opening square bracket '[' is missing from the argument.", "long desc", newParseError)
	OPTIONS_ARGUMENT_END_MISSING              = CustomErrorWrapper("The closing square bracket ']' is missing from the argument.", "long desc", newParseError)
	CLASS_ARGUMENT_START_MISSING              = CustomErrorWrapper("The opening curly braces '{' is missing from the argument.", "long desc", newParseError)
	CLASS_ARGUMENT_END_MISSING                = CustomErrorWrapper("The closing curly braces '}' is missing from the argument.", "long desc", newParseError)
	INNER_ARGUMENT_END_MISSING                = CustomErrorWrapper("A closing argument, '}' or ']' is missing from the argument.", "long desc", newParseError)
	INVALID_ARGUMENT_CONTENT_KEY_VALUE_FORMAT = CustomErrorWrapper("This argument is not in key-value pair form. A pair consists of a key and value separated by '=', with each pair separated by ','", "long desc", newParseError)
	KEY_EMPTY                                 = CustomErrorWrapper("The key provided is empty.", "long desc", newParseError)
	VALUE_EMPTY                               = CustomErrorWrapper("The value provided is emtpy.", "long desc", newParseError)

	// Group parsing errors
	GROUP_NAME_MISMATCH = CustomErrorWrapper("The class argument within \\begin line did not match the class argument within \\end line.", "long desc", newParseError)
	GROUP_BEGIN_MISSING = CustomErrorWrapper("The matching \\begin line could not be found for the \\end line.", "long desc", newParseError)
	GROUP_END_MISSING   = CustomErrorWrapper("The matching \\end line could not be found for the \\begin line.", "long desc", newParseError)

	// Server-responsible parsing errors
	SERVER_RESPONSIBLE_STACK_EMPTY                      = CustomErrorWrapper("Something went wrong on the server's end. The grouping stack is empty. Please report this to the server owner.", "long desc", newParseError)
	SERVER_RESPONSIBLE_STACK_CONTAINS_NON_GROUP_ELEMENT = CustomErrorWrapper("Something went wrong on the server's end. The grouping stack contained a non-group element. Please report this to the server owner.", "long desc", newParseError)
	SERVER_RESPONSIBLE_NON_COMPONENT_DATA_STRUCTURE     = CustomErrorWrapper("Something went wrong on the server's end. The data structure for a LaTeX line did not return one that is supported. Please report this to the server owner.", "long desc", newParseError)
	SERVER_RESPONSIBLE_TABLE_NOT_FOUND                  = CustomErrorWrapper("Something went wrong on the server's end. The tabular group could not be found. Please report this to the server owner.", "long desc", newParseError)
	SERVER_RESPONSIBLE_INVALID_DOCUMENT_CLASS           = CustomErrorWrapper("Something went wrong on the server's end. A non-document class argument was provided as one. Please report this to the server owner.", "long desc", newParseError)

	// conversion parsing errors
	CANNOT_CONVERT_TO_PDF = CustomErrorWrapper("The file was not converted into .pdf format. Please compile it on a TeX engine to identify the cause.", "long desc", newParseError)

	/**
	 * Structure related error functions
	 */
	// Missing essential element structure errors
	DOCUMENT_CLASS_NOT_FOUND                 = CustomErrorWrapper("Missing \\documentclass line.", "long desc", newStructureError)
	SEVERAL_DOCUMENT_CLASSES_FOUND           = CustomErrorWrapper("You cannot have more than 1 document class within a single document.", "long desc", newStructureError)
	DOCUMENT_CLASS_REQUIRED_ARGUMENT_MISSING = CustomErrorWrapper("\\documentclass should have at least 1 required argument. For example, '\\documentclass{article}' is a valid example.", "long desc", newStructureError)
	DOCUMENT_CONTENT_GROUP_NOT_FOUND         = CustomErrorWrapper("Missing \\begin{document} and \\end{document} lines.", "long desc", newStructureError)
	DOCUMENT_CONTENT_GROUP_NOT_DOCUMENT_TYPE = CustomErrorWrapper("Incorrect group name used for the document content.", "long desc", newStructureError)

	// Missing required arguments structure errors
	MISSING_BEGIN_GROUP_NAME = CustomErrorWrapper("You must have at least one required argument '{}' with \\begin lines.", "long desc", newStructureError)
	MISSING_END_GROUP_NAME   = CustomErrorWrapper("You must have at least one required argument '{}' with \\end lines.", "long desc", newStructureError)

	// Containing group constructs within non-grouping contexts
	PREAMBLE_CONTAINS_GROUP     = CustomErrorWrapper("The preamble cannot contain any grouping constructs such as \\begin or \\end.", "long desc", newStructureError)
	PREREQUISITE_CONTAINS_GROUP = CustomErrorWrapper("The content before the document class cannot contain any grouping constructs such as \\begin or \\end.", "long desc", newStructureError)

	// Extra latex code where it shouldn't belong
	DOCUMENT_CONTENT_ALREADY_EXISTS = CustomErrorWrapper("You cannot have more than one document content component.", "long desc", newStructureError)
	USEPACKAGE_OUTSIDE_PREAMBLE     = CustomErrorWrapper("\\usepackage should only exist in the preamble.", "long desc", newStructureError)
	USEPACKAGE_LACKS_CLASS_ARGUMENT = CustomErrorWrapper("\\usepackage should have a required argument. For example, '\\usepackage{graphicx}' is a valid example.", "long desc", newStructureError)

	/**
	 * Accessibility related error functions
	 */
	// tagging accessibility errors
	DOCUMENT_METADATA_MISSING               = CustomErrorWrapper("Missing \\DocumentMetadata line before \\documentclass.", "long desc", newAccessibilityError)
	DOCUMENT_METADATA_APPEARED_LATE         = CustomErrorWrapper("\\DocumentMetadata should not appear after \\documentclass.", "long desc", newAccessibilityError)
	DOCUMENT_METADATA_REPEATED              = CustomErrorWrapper("There should only be one \\DocumentMetadata in a single LaTeX file.", "long desc", newAccessibilityError)
	METADATA_SEVERAL_ARGUMENTS              = CustomErrorWrapper("\\DocumentMetadata should only have one required argument. For example, '\\DocumentMetadata{tagging=on}'.", "long desc", newAccessibilityError)
	METADATA_OTHER_ARGUMENTS_UNSUPPORTED    = CustomErrorWrapper("\\DocumentMetadata should not have other argument types, such as those surrounded by square brackets '[]'.", "long desc", newAccessibilityError)
	METADATA_LACKS_CLASS_ARGUMENT           = CustomErrorWrapper("\\DocumentMetadata should have one required argument.", "long desc", newAccessibilityError)
	METADATA_LACKS_ENABLED_TAGGING_ARGUMENT = CustomErrorWrapper("\\DocumentMetadata should have 'tagging=on' within its required argument.", "long desc", newAccessibilityError)
	METADATA_LACKS_TAGGING_SETUP_ARGUMENT   = CustomErrorWrapper("\\DocumentMetadata should have 'tagging-setup={...}' within its required argument. It is recommended to use 'tagging-setup={activate/all}'.", "long desc", newAccessibilityError)
	METADATA_LACKS_PDF_STANDARD_ARGUMENT    = CustomErrorWrapper("\\DocumentMetadata should have 'pdfstandard=...' within its required argument. It is recommended to use 'pdfstandard=ua-2'.", "long desc", newAccessibilityError)
	METADATA_LACKS_LANGUAGE_ARGUMENT        = CustomErrorWrapper("\\DocumentMetadata should have 'lang=...' within its required argument. For example, 'lang=en-GB' is a valid option.", "long desc", newAccessibilityError)

	// image-related alt text accessibility errors
	GRAPHICS_MISSING_ARGUMENTS            = CustomErrorWrapper("\\includegraphics should at least 1 argument.", "long desc", newAccessibilityError)
	GRAPHICS_TOO_MANY_ARGUMENTS           = CustomErrorWrapper("\\includegraphics should have at most 2 arguments.", "long desc", newAccessibilityError)
	GRAPHICS_FIRST_ARGUMENT_NOT_OPTIONAL  = CustomErrorWrapper("\\includegraphics should have its first of both arguments being optional. For example, '\\includegraphics[...]{...}' is a valid format.", "long desc", newAccessibilityError)
	GRAPHICS_SECOND_ARGUMENT_NOT_REQUIRED = CustomErrorWrapper("\\includegraphics should have its second of both arguments being required. For example, '\\includegraphics[...]{...}' is a valid format.", "long desc", newAccessibilityError)
	GRAPHICS_OUTSIDE_DOCUMENT_CONTENT     = CustomErrorWrapper("\\includegraphics should only appear within the document content.", "long desc", newAccessibilityError)
	GRAPHICS_LACKS_ALT_TEXT               = CustomErrorWrapper("\\includegraphics should contain an optional argument containing alternative text ('alt text'). For example, '\\includegraphics[alt={An image of an apple}]{apple.png} is an accessible graphic.", "long desc", newAccessibilityError)
	GRAPHICS_LACKS_SOURCE                 = CustomErrorWrapper("\\includegraphics should contain a required argument for its source. For example, '\\includegraphics{apple.png}' is a valid graphic.", "long desc", newAccessibilityError)
	MISSING_GRAPHICX_PACKAGE              = CustomErrorWrapper("\\includegraphics requires the 'graphicx' package. This can be included by inserting '\\usepackage{graphicx}' into the preamble.", "long desc", newAccessibilityError)

	// table-related accessibility errors
	TABLE_LACKS_REQUIRED_ARGUMENT                   = CustomErrorWrapper("Tabular groups should have an additional required argument.", "long desc", newAccessibilityError)
	TABLE_CANNOT_PARSE_COLUMNS                      = CustomErrorWrapper("Could not parse how many columns the table has.", "long desc", newAccessibilityError)
	TABLE_MISSING_TAG_PDF_SETUP                     = CustomErrorWrapper("Tabular groups should have \\tagpdfsetup{...} the line before the group.", "long desc", newAccessibilityError)
	TAG_PDF_SETUP_REQUIRED_ARGUMENT                 = CustomErrorWrapper("\\tagpdfsetup should have exactly one required argument.", "long desc", newAccessibilityError)
	TAG_PDF_SETUP_NON_REQUIRED_ARGUMENT             = CustomErrorWrapper("\\tagpdfsetup should an argument surrounded by curly braces '{}'.", "long desc", newAccessibilityError)
	TAG_PDF_SETUP_HEADER_ROWS_INVALID_VALUE         = CustomErrorWrapper("\\tagpdfsetup{table/header-rows=...} should have an integer value surrounded by curly braces as the inner argument. For example, \\tagpdfsetup{table/header-rows={1}} indicates the first row is the heading row.", "long desc", newAccessibilityError)
	TAG_PDF_SETUP_LACKS_HEADER_ROWS_OR_PRESENTATION = CustomErrorWrapper("\\tagpdfsetup should have an argument with 'table/header-rows={N}' with N being the row the heading occurs, or 'table/tagging=presentation' to indicate the screen reader is not to read the table.", "long desc", newAccessibilityError)
	TAG_PDF_SETUP_TABLE_TAGGING_INVALID_VALUE       = CustomErrorWrapper("\\tagpdfsetup{table/tagging=...} should have the inner value 'presentation' creating '\\tagpdfsetup{table/tagging=presentation}'.", "long desc", newAccessibilityError)

	/**
	 * Server-created errors. These can be caused by server-related processes failing
	 */
	SERVER_RESPONSIBLE_CANNOT_FIND_WORKING_DIRECTORY = CustomErrorWrapper("Something went wrong on the server's end. The working directory could not be retrieved. Please report this to the server owner.", "long desc", newServerError)
	SERVER_RESPONSIBLE_WRITE_FILE_ERROR              = CustomErrorWrapper("Something went wrong on the server's end. Could not write to the file. Please report this to the server owner.", "long desc", newServerError)
	SERVER_RESPONSIBLE_READ_FILE_ERROR               = CustomErrorWrapper("Something went wrong on the server's end. The file location for the LaTeX file could not be read. Please report this to the server owner.", "long desc", newServerError)
	SERVER_RESPONSIBLE_UNIDENTIFIED_ERROR            = CustomErrorWrapper("Something went wrong on the server's end. There is an unidentified error. Please report this to the server owner.", "long desc", newServerError)
	SERVER_RESPONSIBLE_INVALID_STRING_ERROR          = CustomErrorWrapper("Something went wrong on the server's end. Could not parse string. Please report this to the server owner.", "long desc", newServerError)
	SERVER_RESPONSIBLE_INVALID_REGEX_ERROR           = CustomErrorWrapper("Something went wrong on the server's end. Could not parse regular expression. Please report this to the server owner.", "long desc", newServerError)
)
