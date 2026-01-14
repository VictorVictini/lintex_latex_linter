package main

import (
	"html/template"
	"net/http"
	"regexp"
)

type Page struct {
	Title     string
	Body      []byte
	Responses []string
}

var templates = template.Must(template.ParseFiles("webpages/home.html"))
var validPath = regexp.MustCompile("^/home/?$")

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
		fn(w, r, m[0])
	}
}

func homepageHandler(w http.ResponseWriter, r *http.Request, title string) {
	p := &Page{
		Title:     "LinTeX: Promoting Accessibility in LaTeX",
		Responses: make([]string, 0)}
	renderTemplate(w, "home", p)
}
