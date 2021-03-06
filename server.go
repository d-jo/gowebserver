package main

import (
	"fmt"
	"github.com/d-jo/webserver/io-ops"
	"github.com/d-jo/webserver/structs"
	"html/template"
	"net/http"
	"regexp"
	"time"
	"strings"
)

var templates = template.Must(template.New("main").Funcs(template.FuncMap{

}).ParseGlob("web/templates/*.html"))
var validPath = regexp.MustCompile("^/(e|write|s|p|good|idiom)/([a-zA-Z0-9]+)$")

func viewSnippit(w http.ResponseWriter, r *http.Request, id string) {
	snip, err := io_ops.GetCodeSnipFromDB(id)
	if err != nil {
		// custom 404 TODO
		http.NotFound(w, r)
		return
	}
	executeViewTemplate(w, "viewcodesnip", *snip)
}

func createSnippit(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, "createcodesnip.html", nil)
	if err != nil {
		http.NotFound(w, r)
	}
}

func saveSnippit(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	author := r.FormValue("author")
	body := r.FormValue("body")
	if len(title) < 3 || len(author) < 3 || len(body) < 3 {
		http.Redirect(w, r, "/c/", http.StatusFound)
		return
	}
	id, err := io_ops.InsertCodeSnipToDB(&structs.CodeSnip{Title: title, Author: author, Content: body})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
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

func handleGood(w http.ResponseWriter, r *http.Request, id string) {
	newPoints, _ := io_ops.UpdatePointsInDB(1, 0, id)
	setCookie(w, r, "good", id)
	fmt.Fprint(w, newPoints)
}

func handleIdiom(w http.ResponseWriter, r *http.Request, id string) {
	_, newPoints := io_ops.UpdatePointsInDB(0, 1, id)
	setCookie(w, r, "idiom", id)
	fmt.Fprint(w, newPoints)
}

func setCookie(w http.ResponseWriter, r *http.Request, dist, id string) {
	name := dist + "votes"
	cookie, err := r.Cookie(name)
	expire := time.Now().Add(24 * time.Hour)
	if err != nil {
		cookie := http.Cookie{Name: name, Value: id, Expires: expire, Path: "/s/"}
		http.SetCookie(w, &cookie)
		return
	}

	if !idInCookie(cookie, id) {
		http.SetCookie(w, &http.Cookie{Name: name, Value: cookie.Value + "," + id, Expires: expire, Path: "/s/"})
	}
}

func idInCookie(cookie *http.Cookie, search string) bool{
	values := strings.Split(cookie.Value, ",")
	for _, val := range values {
		if val == search {
			return true
		}
	}
	return false
}

func main() {
	http.HandleFunc("/", createSnippit)
	http.HandleFunc("/s/", makeHandler(viewSnippit))
	http.HandleFunc("/write/", saveSnippit)
	http.HandleFunc("/c/", createSnippit)
	http.HandleFunc("/good/", makeHandler(handleGood))
	http.HandleFunc("/idiom/", makeHandler(handleIdiom))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))
	http.ListenAndServe(":8080", nil)
}
