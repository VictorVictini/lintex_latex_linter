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
	Responses []CustomError
	Error     CustomError
}

func newPage(fileName string, body string, responses []CustomError, err CustomError) *Page {
	return &Page{
		FileName:  fileName,
		Body:      body,
		Responses: responses,
		Error:     err,
	}
}

var funcs = template.FuncMap{"join": strings.Join}
var templates = template.Must(template.New("").Funcs(funcs).ParseFiles("webpages/tool.html", "webpages/homepage.html", "webpages/error.html"))
var validPath = regexp.MustCompile("^/$|^/(home|homepage|tool|error)/[a-zA-Z0-9/]*$")

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
		fmt.Printf("\n\nfile uploading\n\n")
		// ensures the file is closed after the function executes
		defer file.Close()

		// update page information
		filename := header.Filename
		extension := filepath.Ext(filename)
		filename = filename[0 : len(filename)-len(extension)]
		page.FileName = filename
		GLOBAL_FILE_NAME = page.FileName

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

		// re-upload to same page
	} else if page.FileName != "" {
		fmt.Printf("\n\nno file name provided\n\n")
		// reset error output
		page.FileName = GLOBAL_FILE_NAME
		page.Responses = make([]CustomError, 0)

		// writing to a file
		GLOBAL_FILE_DATA_CACHE = []byte(page.Body)
		errFn := WriteFileBytes(page.FileName, GLOBAL_FILE_DATA_CACHE)
		if errFn != nil {
			page.Responses = append(page.Responses, DummyError(errFn))
		}

		// retrieve error outputs
		document, errLs := Parse(page.FileName, []byte(page.Body))
		if errLs != nil {
			page.Responses = append(page.Responses, FormatErrors(errLs.(errList))...)
		} else {
			// get base document
			var doc IDocument = document.(IDocument)

			// add a decorator for accessibility
			accDoc := newAccessibleDocument(doc).(*AccessibleDocument)
			err := accDoc.VerifyAccessibility()
			if len(err) != 0 {
				page.Responses = append(page.Responses, FormatErrors(err)...)
			} else {
				_, errFn := CompileLaTeXFile(page.FileName)
				if errFn != nil {
					page.Responses = append(page.Responses, DummyError(errFn))
				}
			}
		}

		// assume 'create file' option was selected
	} else {
		fmt.Printf("\n\nnew document creation\n\n")
		page.FileName = "document"
		GLOBAL_FILE_NAME = page.FileName
		errFn = WriteFileBytes(page.FileName, []byte{})
		if errFn != nil {
			fmt.Fprintf(w, "%s\n", DummyError(errFn).Error())
		}
	}

	// render tool
	renderTemplate(w, "tool", page)
}

func homepageHandler(w http.ResponseWriter, r *http.Request, title string) {
	// parse form
	page := ProcessForm(r)

	// render homepage
	renderTemplate(w, "homepage", page)
}

func errorHandler(w http.ResponseWriter, r *http.Request, title string) {
	// parse form
	page := ProcessForm(r)

	// render homepage
	renderTemplate(w, "error", page)
}

func ProcessForm(r *http.Request) *Page {
	// default for first load
	defaultPage := newPage("", "", make([]CustomError, 0), nil)

	// parse form data
	err := r.ParseForm()
	if err != nil {
		return defaultPage
	}

	errFn, _ := GLOBAL_ERROR_ID_GENERATOR.GetErrorFromID(r.FormValue("id"))

	// manage errors
	return newPage(r.FormValue("tex_file"), r.FormValue("body"), make([]CustomError, 0), DummyError(errFn))
}
