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
