package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"urlshort"
)

func readFile(file string) []byte {
	content, err := os.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}
	return content
}

func main() {
	yamlFlag := flag.String("yml", "", "The yaml file to be parsed")
	jsonFlag := flag.String("json", "", "The json file to be parsed")
	dbFlag := flag.String("db", "", "Expects the path to the db")
	flag.Parse()

	var data []byte
	defaultFlag := "yaml"
	flag := &defaultFlag
	if *yamlFlag != "" {
		data = readFile(*yamlFlag)
		flag = yamlFlag
	} else if *jsonFlag != "" {
		data = readFile(*jsonFlag)
		flag = jsonFlag
	} else if *dbFlag != "" {
		flag = dbFlag
	} else {
		log.Println("No input file given, defaulting to data.yaml")
		data = readFile("data.yaml")
	}

	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
		"/spiegel":        "https://spiegel.de",
	}

	mapHandler := urlshort.MapHandler(pathsToUrls, mux)
	// Build the YAMLHandler using the mapHandler as the
	// fallback
	handler, err := urlshort.Handler(data, mapHandler, flag)
	if err != nil {
		panic(err)
	}
	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", handler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
