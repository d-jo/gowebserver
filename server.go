package main

import (
	//"github.com/d-jo/webserver/structs"
	"github.com/d-jo/webserver/io-ops"
	"net/http"
	"regexp"
)

func viewSnippit(w http.ResponseWriter, r *http.Request, id string) {

}

var validPath = regexp.MustCompile("^/(e|w|s)/([a-zA-Z0-9]+)$")

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

func main() {
	//http.HandleFunc("/s/", makeHandler(viewSnippit))
	//http.HandleFunc("/w/", makeHandler(saveSnippit))
	//http.HandleFunc("/e/", makeHandler(editSnippit))
	io_ops.Test()
}
