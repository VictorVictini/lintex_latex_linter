package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
)

type Page struct {
	FileName  string
	Body      string
	Responses []string
}

func newPage(fileName string, body string, responses []string) *Page {
	return &Page{
		FileName:  fileName,
		Body:      body,
		Responses: responses,
	}
}

var funcs = template.FuncMap{"join": strings.Join}
var templates = template.Must(template.New("").Funcs(funcs).ParseFiles("webpages/tool.html", "webpages/homepage.html"))
var validPath = regexp.MustCompile("^/(?:(homepage|tool)/)?[a-zA-Z0-9/]*$")

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[1])
	}
}

func toolHandler(w http.ResponseWriter, r *http.Request, title string) {
	// parse form
	page := ProcessForm(r)

	// // Initialize error messages slice
	// var serverMessages []string

	// // Parse the multipart form, 10 MB max upload size
	// r.ParseMultipartForm(10 << 20)

	// Retrieve the file from form data
	file, header, err := r.FormFile("tex_file")
	if err != nil {
		fmt.Printf("uh oh spaghetios %#v\n", err)
		return
		// if err == http.ErrMissingFile {
		// 	serverMessages = append(serverMessages, "No file submitted")
		// } else {
		// 	serverMessages = append(serverMessages, "Error retrieving the file")
		// }

		// if len(serverMessages) > 0 {
		// 	templates.ExecuteTemplate(w, "messages", serverMessages)
		// 	return
		// }

	}
	fmt.Printf("File Data %#v\n", file)
	defer file.Close()

	// create file
	out, err := os.Create("C:\\Users\\Danyal\\lintex\\latex\\example.tex")
	if err != nil {
		fmt.Printf("%#v\n", err)
		fmt.Fprintf(w, "Unable to create the file for writing. Check your write access privilege")
		return
	}
	defer out.Close()

	// write the content from POST to the file
	_, err = io.Copy(out, file)
	if err != nil {
		fmt.Fprintln(w, err)
	}

	fmt.Fprintf(w, "File uploaded successfully : ")
	fmt.Fprintf(w, header.Filename)

	// render tool
	renderTemplate(w, "tool", page)
}

func homepageHandler(w http.ResponseWriter, r *http.Request, title string) {
	// parse form
	page := ProcessForm(r)

	// read file if any
	fmt.Printf("home %#v\n", page)

	// render homepage
	renderTemplate(w, "homepage", page)

	// // writing to a file
	// GLOBAL_FILE_DATA_CACHE = []byte(page.Body)
	// errFn := WriteFileBytes("example", GLOBAL_FILE_DATA_CACHE)
	// if errFn != nil {
	// 	err := make(errList, 1)
	// 	err = errList{DummyError(errFn)}
	// 	page.Responses = append(page.Responses, FormatErrors(err)...)
	// }

	// // retrieve error outputs
	// document, errLs := Parse("example", []byte(page.Body))
	// if errLs != nil {
	// 	page.Responses = append(page.Responses, FormatErrors(errLs.(errList))...)
	// } else {
	// 	// get base document
	// 	var doc IDocument = document.(IDocument)

	// 	// add a decorator for accessibility
	// 	accDoc := newAccessibleDocument(doc).(*AccessibleDocument)
	// 	err := accDoc.VerifyAccessibility()
	// 	if len(err) != 0 {
	// 		page.Responses = append(page.Responses, FormatErrors(err)...)
	// 	} else {
	// 		_, errFn := CompileLaTeXFile("example")
	// 		if errFn != nil {
	// 			err := make(errList, 1)
	// 			err = errList{DummyError(errFn)}
	// 			page.Responses = append(page.Responses, FormatErrors(err)...)
	// 		}
	// 	}
	// }
	//
	// // output webpage
	// renderTemplate(w, "homepage", page)
}

func ProcessForm(r *http.Request) *Page {
	// default for first load
	defaultPage := newPage("", "", make([]string, 0))

	// parse form data
	err := r.ParseForm()
	if err != nil {
		return defaultPage
	}

	return newPage(r.FormValue("tex_file"), r.FormValue("body"), strings.Split(r.FormValue("responses"), ","))
}
