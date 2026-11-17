package main

import (
	"fmt"
)

func readFileContents(fileName string) CreateError {
	// read from the file
	fileBytes, errFn := GetFileBytes(fileName)
	if errFn != nil {
		return errFn
	}

	// use string(fileBytes)
	fmt.Println(string(fileBytes))
	return nil
}
