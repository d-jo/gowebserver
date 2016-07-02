package main

import (
	"fmt"
	"github.com/d-jo/webserver/io-ops"
	"github.com/d-jo/webserver/structs"
	"html/template"
	"net/http"
	"regexp"
)

var templates = template.Must(template.ParseGlob("web/templates/*.html"))
var validPath = regexp.MustCompile("^/(e|w|s|p)/([a-zA-Z0-9]+)$")

func viewSnippit(w http.ResponseWriter, r *http.Request, id string) {
	snip, err := io_ops.GetCodeSnipFromDB(id)
	if err != nil {
		// custom 404 TODO
		http.NotFound(w, r)
		return
	}
	executeTemplate(w, "viewcodesnip", snip)
}

func makeHandler(fn func(w http.ResponseWriter, r *http.Request, id string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

func executeTemplate(w http.ResponseWriter, templateName string, cs structs.CodeSnip) {
	err := templates.ExecuteTemplate(w, fmt.Sprintf("%s.html", templateName), cs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/s/", makeHandler(viewSnippit))
	//http.HandleFunc("/write/", makeHandler(saveSnippit))
	//http.HandleFunc("/e/", editSnippit)
	http.ListenAndServe(":8080", nil)
}
