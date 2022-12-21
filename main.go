package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"

	"github.com/boltdb/bolt"
	"github.com/plastikov/urlshort/handler"
)

func main() {

	mux := defaultMux()
	flag.Parse()
	ext := filepath.Ext(*filePath)

	pathUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	handler.MapHandler(pathUrls, mux)

	var hndlr http.Handler
	var err error
	if ext == ".yaml" {
		hndlr, err = handler.YAMLHandler(printBytes(*filePath), mux)
		if err != nil {
			panic(err)
		}
	} else if ext == ".json" {
		hndlr, err = handler.JSONHandler(printBytes(*filePath), mux)
		if err != nil {
			panic(err)
		}
	} else if ext == "db" {
		hndlr = handler.DBHandler((createDB(pathUrls)), mux)
	}
	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", hndlr)
}

var (
	filePath = flag.String("File", "urls.yaml", "Pass into this flag a file holding data of url paths")
)

func printBytes(filepath string) []byte {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		panic(err)
	}
	return []byte(data)
}

func createDB(paths map[string]string) bolt.DB {
	db, err := bolt.Open("urls.db", 0600, nil)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	for path, url := range paths {
		err = db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("pathurls"))
			err := b.Put([]byte(path), []byte(url))
			return err
		})
		if err != nil {
			panic(err)
		}
	}
	return *db
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
