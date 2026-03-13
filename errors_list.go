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
	MISSING_END CreateError = CustomErrorWrapper("short desc", "long desc", newParseError)
)
