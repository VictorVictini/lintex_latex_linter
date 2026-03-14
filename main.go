package main

import (
	"fmt"
)

// func main() {
// 	fs := http.FileServer(http.Dir("./styles"))
// 	http.Handle("/styles/", http.StripPrefix("/styles", fs))

// 	http.HandleFunc("/home/", makeHandler(homepageHandler))

// 	log.Fatal(http.ListenAndServe(":8080", nil))
// }

func main() {
	res, err := ParseFile("latex_testing/test.tex")
	if err != nil {
		list := err.(errList)
		for _, err := range list {
			pe := err.(*parserError)
			cerr, ok := pe.Inner.(CustomError)
			if ok {
				errStr, errFn := cerr.LocateError()
				if errFn != nil {
					fmt.Println(DummyError(errFn).Error())
					return
				}
				fmt.Printf("%#v\t%#v\t%#v\n", pe.pos, cerr.Error(), errStr)
			} else {
				fmt.Printf("%s: %#v\n", pe.prefix, pe.Inner.Error())
			}
		}
	} else {
		fmt.Println("BEFORE DOCUMENT CLASS:")
		for _, line := range res.(Document).prerequisiteContent {
			line.PrintTree(1)
		}
		fmt.Println("\nDOCUMENT CLASS:")
		res.(Document).documentClass.PrintTree(1)
		fmt.Println("\nPREAMBLE:")
		for _, line := range res.(Document).preamble {
			line.PrintTree(1)
		}
		fmt.Println("\nDOCUMENT CONTENT:")
		res.(Document).content.PrintTree(1)
		// for i, comp := range AnyInterfaceToTSlice[*Line](res) {
		// 	fmt.Print("A")
		// 	if i > 0 {
		// 		fmt.Println()
		// 	}
		// 	comp.PrintTree(0)
		// }
	}
	// var err error = MISSING_END(newCoordinate(0, 0), newCoordinate(0, 0))
	// fmt.Printf("%#v\n", err)
	// var customErr CustomError = err.(CustomError)
	// fmt.Println(customErr.Error())
	// fmt.Println(customErr.LongError())
	// fmt.Println(customErr.LocateError())
	// fmt.Println(customErr.RetrieveID())
}
