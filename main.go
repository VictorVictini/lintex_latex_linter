package main

import "fmt"

func main() {
	errFn := readFileContents(GLOBAL_FILE_NAME)
	if errFn != nil {
		fmt.Println(DummyError(errFn))
	} else {
		fmt.Println("Doneso")
	}
}
