package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"gopkg.in/yaml.v3"
)

func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if dest, ok := pathsToUrls[path]; ok {
			http.Redirect(w, r, dest, http.StatusFound)
			return
		}
		fallback.ServeHTTP(w, r)
	}
}

func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedYaml, err := parseYAML(yml)
	if err != nil {
		return nil, err
	}
	pathMap, err := buildMap(parsedYaml)
	if err != nil {
		fmt.Println("Could not parse YAML")
	}
	return MapHandler(pathMap, fallback), nil
}

// create a JSONHandler function similar to YAMLHandler
func JSONHandler(jsonData []byte, fallback http.Handler) (http.HandlerFunc, error) {
	type jsonPath []struct{
		Path string `json:"path"`
		Url string `json:"url"`
	}

	var pathsToUrl []jsonPath
	if err := json.Unmarshal(jsonData, &pathsToUrl); err != nil{
		return nil, err
	}

	return func(w http.ResponseWriter, r *http.Request){
		for _, pathtourl := range pathsToUrl{
			if pathtourl.Path == r.URL.Path{
				http.Redirect(w, r, pathtourl.Url, http.StatusFound)
			return
			}
		}
		fallback.ServeHTTP(w, r)
	}, nil
}

func parseYAML(yml []byte) (parsedYaml []byte, err error) {
    yamlData, err := ioutil.ReadAll(strings.NewReader(string(yml)))
    if err != nil {
        log.Fatalf("error: %v", err)
    }
	return yamlData, err
}

func buildMap(yamlData []byte) (map[string]string, error) {
	var data map[string]string
	if err := yaml.Unmarshal(yamlData, &data); err != nil {
		return nil, err
	}
	return data, nil
}
