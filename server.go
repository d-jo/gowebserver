package main

import (
	"fmt"
	"github.com/d-jo/webserver/io-ops"
	"github.com/d-jo/webserver/structs"
	"html/template"
	"net/http"
	"regexp"
)

var templates = template.Must(template.New("main").Funcs(template.FuncMap{

}).ParseGlob("web/templates/*.html"))
var validPath = regexp.MustCompile("^/(e|write|s|p)/([a-zA-Z0-9]+)$")

func viewSnippit(w http.ResponseWriter, r *http.Request, id string) {
	snip, err := io_ops.GetCodeSnipFromDB(id)
	if err != nil {
		// custom 404 TODO
		http.NotFound(w, r)
		return
	}
	executeViewTemplate(w, "viewcodesnip", snip)
}

func createSnippit(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, "createcodesnip.html", nil)
	if err != nil {
		http.NotFound(w, r)
	}
}

func saveSnippit(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	body := r.FormValue("body")
	id, err := io_ops.InsertCodeSnipToDB(structs.CodeSnip{Title: title, Content: body})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	http.Redirect(w, r, fmt.Sprintf("/s/%d", id), http.StatusFound)
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

func executeViewTemplate(w http.ResponseWriter, templateName string, cs structs.CodeSnip) {
	err := templates.ExecuteTemplate(w, fmt.Sprintf("%s.html", templateName), cs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/s/", makeHandler(viewSnippit))
	http.HandleFunc("/write/", saveSnippit)
	http.HandleFunc("/c/", createSnippit)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))
	http.ListenAndServe(":8080", nil)
}
