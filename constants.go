package main

// golang does not support structs as constants, so despite the irony this will be a var
var (
	GLOBAL_ERROR_ID_GENERATOR ErrorIDGenerator = newErrorIDGenerator()
)

const (
	GLOBAL_STORAGE_LOCATION string = "C:/Users/Danyal/lintex"
	WHITESPACE              string = " \t\n\r"
)
