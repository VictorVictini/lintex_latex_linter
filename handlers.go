package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
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
	// get working directory
	workingDir, errFn := GetWorkingDirectory()
	if errFn != nil {
		fmt.Fprintf(w, "%s\n", DummyError(errFn).Error())
		return
	}

	// parse form
	page := ProcessForm(r)

	// read and upload the form's file
	file, header, err := r.FormFile("tex_file")
	if err == nil {
		// ensures the file is closed after the function executes
		defer file.Close()

		// update page information
		filename := header.Filename
		extension := filepath.Ext(filename)
		page.FileName = filename
		filename = filename[0 : len(filename)-len(extension)]

		// create file
		out, err := os.Create(filepath.Join(workingDir, GLOBAL_LATEX_FOLDER, page.FileName+".tex"))
		if err != nil {
			fmt.Fprintf(w, "Unable to create the file for writing.")
			return
		}
		defer out.Close()

		// write the content from POST to the file
		_, err = io.Copy(out, file)
		if err != nil {
			fmt.Fprintln(w, err)
		}

		// update page body to include file contents
		fileBytes, errFn := GetFileBytes(filename)
		if errFn != nil {
			fmt.Fprintf(w, "%s\n", DummyError(errFn).Error())
			return
		}
		page.Body = string(fileBytes)

		// assume 'create file' option was selected
	} else {
		page.FileName = "document"
		errFn = WriteFileBytes(page.FileName, []byte{})
		if errFn != nil {
			fmt.Fprintf(w, "%s\n", DummyError(errFn).Error())
		}
	}

	// render tool
	fmt.Printf("send off %#v\n", page)
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
