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
				fmt.Printf("%#v\t%#v\t%#v\n", pe.pos, cerr.Error(), cerr.LocateError())
			} else {
				fmt.Printf("%#v\n", pe.Inner.Error())
			}
		}
	} else {
		fmt.Printf("%#v\n", res)
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
