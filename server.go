package main

import (
	"fmt"
	mathrand "math/rand"
	"net/http"
	"os"
	"os/exec"

	"google.golang.org/appengine" // Required external App Engine library
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	fmt.Fprintln(w, "Hello, Gopher Network!")
}

// photoHandler gets a new photo and processes it
func photoHandler(w http.ResponseWriter, r *http.Request) {
	photourl := searchForPhoto()
	uuid, inputFile := downloadPhoto(photourl)
	outputFile := "./img/output/" + uuid + "-out.jpg"
	n := 200 + mathrand.Intn(300)
	modes := []int{1, 4, 6}
	m := 0
	processPhoto(inputFile, outputFile, n, modes[m])
}

// buildHandler rebuilds the site to incorporate new photo
func buildHandler(w http.ResponseWriter, r *http.Request) {
	cmd := "hugo"
	args := []string{}
	if err := exec.Command(cmd, args...).Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println("Site Rebuilt!")
}

func main() {
	http.HandleFunc("/tasks/getphoto", photoHandler)
	http.HandleFunc("/tasks/buildsite", photoHandler)
	http.Handle("/", http.FileServer(http.Dir("public")))
	appengine.Main() // Starts the server to receive requests
}
