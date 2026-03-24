package main

const (
	GLOBAL_LATEX_FOLDER      string = "latex"
	WHITESPACE               string = " \t\n\r"
	NUMBERS_ARGUMENT_REGEX   string = `\A\s*\{\s*\d+\s*\}\s*\z`
	KEY_VALUE_ARGUMENT_REGEX string = `\A(?:` + KEY_VALUE_PAIR_REGEX + `,)*` + KEY_VALUE_PAIR_REGEX + `\z`
	KEY_VALUE_PAIR_REGEX     string = `\s*(\S+)\s*=\s*({[^}]*}|[^,{]+)\s*`
	DEFAULT_FILE_NAME        string = "document"
)

// golang does not support structs as constants, so despite the irony this will be a var
var (
	GLOBAL_ERROR_ID_GENERATOR ErrorIDGenerator = newErrorIDGenerator() // move into singleton later
	GLOBAL_FILE_DATA_CACHE    []byte           = []byte{}              // move into singleton later
	GLOBAL_FILE_NAME          string           = DEFAULT_FILE_NAME     // move into singleton later
)
