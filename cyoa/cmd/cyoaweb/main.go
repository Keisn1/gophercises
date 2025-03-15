package main

import (
	"cyoa"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"
)

func readFile(filepath string) []byte {
	content, err := os.ReadFile(filepath)
	if err != nil {
		log.Fatal(err)
	}
	return content
}

func main() {
	filename := flag.String("file", "gopher.json", "the JSON file with the CYOA story")
	flag.Parse()

	// 1. Decode json file
	content := readFile(*filename)

	var story cyoa.Story
	err := json.Unmarshal(content, &story)
	if err != nil {
		log.Fatal(err)
	}

	handler := cyoa.NewHandler(story, cyoa.WithPath(pathFn))
	http.Handle("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

var tplname2 string = "templates/chapter2.html"
var tpl2 *template.Template = template.Must(template.ParseFiles(tplname2))

func pathFn(r *http.Request) string {
	path := strings.TrimSpace(r.URL.Path)
	if path == "/story" || path == "story" {
		path = "/story/intro"
	}
	return path[len("/story/"):]
}
