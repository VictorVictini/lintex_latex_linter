package main

import (
	"log"
	"net/http"
)

func main() {
	fs := http.FileServer(http.Dir("./styles"))
	http.Handle("/styles/", http.StripPrefix("/styles", fs))

	http.HandleFunc("/home/", makeHandler(homepageHandler))

	log.Fatal(http.ListenAndServe(":8080", nil))
}

// func main() {
// 	res, err := ParseFile("latex_testing/test.tex")
// 	fmt.Println("A")
// 	if err != nil {
// 		fmt.Printf("error: %s\n", err.Error())
// 	} else {
// 		fmt.Printf("res: %s\n", res)
// 	}
// }
