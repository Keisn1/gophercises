package urlshort

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"database/sql"

	"github.com/mattn/go-sqlite3"
	"gopkg.in/yaml.v3"
)

var db *sql.DB

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	//	TODO: Implement this...
	mapHandlerFunc := func(w http.ResponseWriter, r *http.Request) {
		if url, ok := pathsToUrls[r.URL.Path]; ok {
			http.Redirect(w, r, url, http.StatusFound)
			return
		}
		fallback.ServeHTTP(w, r)
	}
	return mapHandlerFunc
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//   - path: /some-path
//     url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func Handler(input []byte, fallback http.Handler, flag *string) (http.HandlerFunc, error) {
	// TODO: Implement this...

	var parsedInput []PathAndUrl
	var err error
	if strings.Contains(*flag, "yaml") {
		parsedInput, err = parseYaml(input)
	} else if strings.Contains(*flag, "json") {
		parsedInput, err = parseJSON(input)
	} else if strings.Contains(*flag, "sqlite") {
		parsedInput, err = readFromDB(*flag)
	} else {
		log.Fatal("Format not supported")
	}

	if err != nil {
		return nil, err
	}

	pathMap := buildMap(parsedInput)
	return MapHandler(pathMap, fallback), nil
}

func buildMap(pathAndUrls []PathAndUrl) map[string]string {
	pathsToUrl := make(map[string]string)
	for _, entry := range pathAndUrls {
		pathsToUrl[entry.Path] = entry.Url
	}
	return pathsToUrl
}

type PathAndUrl struct {
	Path string
	Url  string
}

func parseYaml(yml []byte) ([]PathAndUrl, error) {
	var pathsAndUrls []PathAndUrl
	err := yaml.Unmarshal([]byte(yml), &pathsAndUrls)
	if err != nil {
		log.Fatal(err)
	}
	return pathsAndUrls, nil
}

func parseJSON(yml []byte) ([]PathAndUrl, error) {
	var pathsAndUrls []PathAndUrl
	err := json.Unmarshal([]byte(yml), &pathsAndUrls)
	if err != nil {
		log.Fatal(err)
	}
	return pathsAndUrls, nil
}

func readFromDB(dbPath string) ([]PathAndUrl, error) {
	if false {
		log.Println(sqlite3.ErrBusy)
	}
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	log.Println("Connected!")

	rows, err := db.Query("select path, url from pathToUrls")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var pathsAndUrls []PathAndUrl
	for rows.Next() {
		var pathAndUrl PathAndUrl
		if err := rows.Scan(&pathAndUrl.Path, &pathAndUrl.Url); err != nil {
			return nil, err
		}
		pathsAndUrls = append(pathsAndUrls, pathAndUrl)
	}
	return pathsAndUrls, nil
}
