package main

import (
	"log"
	"net/http"
)

func main() {
	fs := http.FileServer(http.Dir("./styles"))
	http.Handle("/styles/", http.StripPrefix("/styles", fs))
	fs = http.FileServer(http.Dir("./node_modules"))
	http.Handle("/node_modules/", http.StripPrefix("/node_modules", fs))
	fs = http.FileServer(http.Dir("./images"))
	http.Handle("/images/", http.StripPrefix("/images", fs))
	fs = http.FileServer(http.Dir("./" + GLOBAL_LATEX_FOLDER))
	http.Handle("/"+GLOBAL_LATEX_FOLDER+"/", http.StripPrefix("/"+GLOBAL_LATEX_FOLDER, fs))

	http.HandleFunc("/error", makeHandler(errorHandler))
	http.HandleFunc("/tool", makeHandler(toolHandler))
	http.HandleFunc("/home", makeHandler(homepageHandler))
	http.HandleFunc("/homepage", makeHandler(homepageHandler))
	http.HandleFunc("/", makeHandler(homepageHandler))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
