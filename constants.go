package main

const (
	GLOBAL_LATEX_FOLDER    string = "latex"
	WHITESPACE             string = " \t\n\r"
	NUMBERS_ARGUMENT_REGEX string = `\A\s*\{\s*\d+\s*\}\s*\z`
	DEFAULT_FILE_NAME      string = "document"
)

// golang does not support structs as constants, so despite the irony this will be a var
var (
	GLOBAL_ERROR_ID_GENERATOR ErrorIDGenerator = newErrorIDGenerator() // move into singleton later
	GLOBAL_FILE_DATA_CACHE    []byte           = []byte{}              // move into singleton later
	GLOBAL_FILE_NAME          string           = DEFAULT_FILE_NAME     // move into singleton later
)
