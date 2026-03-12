package main

import "fmt"

// func main() {
// 	fs := http.FileServer(http.Dir("./styles"))
// 	http.Handle("/styles/", http.StripPrefix("/styles", fs))

// 	http.HandleFunc("/home/", makeHandler(homepageHandler))

// 	log.Fatal(http.ListenAndServe(":8080", nil))
// }

func main() {
	res, err := ParseFile("latex_testing/test.tex")
	if err != nil {
		fmt.Printf("error: %s\n", err.Error())
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
}
