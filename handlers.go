package main

import (
	"html/template"
	"net/http"
	"regexp"
	"strings"
)

type Page struct {
	Body      string
	Responses []string
}

var funcs = template.FuncMap{"join": strings.Join}
var templates = template.Must(template.New("").Funcs(funcs).ParseFiles("webpages/home.html"))
var validPath = regexp.MustCompile("^/(home)/[a-zA-Z0-9/]*$")

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

func homepageHandler(w http.ResponseWriter, r *http.Request, title string) {
	// parse form
	page := ProcessForm(r)
	page.Responses = make([]string, 0)

	// retrieve error outputs
	GLOBAL_FILE_DATA_CACHE = []byte(page.Body)
	document, err := Parse("user_data", []byte(page.Body))
	if err != nil {
		page.Responses = append(page.Responses, FormatErrors(err.(errList))...)
	} else {
		// get base document
		var doc IDocument = document.(IDocument)

		// add a decorator for accessibility
		accDoc := newAccessibleDocument(doc).(*AccessibleDocument)
		err := accDoc.VerifyAccessibility()
		if err != nil {
			page.Responses = append(page.Responses, FormatErrors(err)...)
		}
	}

	// output webpage
	renderTemplate(w, "home", page)
}

func ProcessForm(r *http.Request) *Page {
	// default for first load
	defaultPage := &Page{}

	// parse form data
	err := r.ParseForm()
	if err != nil {
		return defaultPage
	}

	return &Page{
		Body:      r.FormValue("body"),
		Responses: strings.Split(r.FormValue("responses"), ","),
	}
}
