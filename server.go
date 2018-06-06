package main

import (
	"fmt"
	"net/http"

	"google.golang.org/appengine" // Required external App Engine library
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	fmt.Fprintln(w, "Hello, Gopher Network!")
}

func photoHandler(w http.ResponseWriter, r *http.Request) {

}

func main() {
	http.HandleFunc("/tasks/getphoto", photoHandler)
	http.HandleFunc("/", indexHandler)
	appengine.Main() // Starts the server to receive requests
}
