package main

import (
	"fmt"
	"net/http"
	"flag"

	"github.com/plastikov/urlshort/handler"
)

func main() {
	yamlFlag := flag.String("File", "urls.yaml", "Pass into this flag a file holding data of yaml type")
	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
		"/my-github": "http://github.com/plastikov",
	}
	mapHandler := handler.MapHandler(pathsToUrls, mux)

	// Build the YAMLHandler using the mapHandler as the
	// fallback

	yamlHandler, err := handler.YAMLHandler([]byte(*yamlFlag), mapHandler)
	if err != nil {
		panic(err)
	}
	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", yamlHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}


func hello(w http.ResponseWriter, r *http.Request) {
	asciiArt := `
	_____
   /    /|_ ___________________________________________
  /    // /|                                          /|
 (====|/ //   An apple a day...            _QP_      / |
  (=====|/     keeps the teacher at bay   (  ' )    / .|
 (====|/                                   \__/    / /||
/_________________________________________________/ / ||
|  _____________________________________________  ||  ||
| ||                                            | ||
| ||                                            | ||
| |                                             | |  pjb
	`
	fmt.Fprintln(w, asciiArt)
}
