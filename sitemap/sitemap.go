package sitemap

import (
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"sitemap/links"
	"strings"
	// "sitemap/links"
)

type urlset struct {
	XMLName xml.Name `xml:"urlset"`
	Xmlns   string   `xml:"xmlns,attr"`
	Urls    []*Url
}

type Url struct {
	XMLName xml.Name `xml:"url"`
	Loc     string   `xml:"loc"`
}

func encode(root string, urlMap map[string]bool) error {
	var urls []*Url
	for u := range urlMap {
		newUrl := Url{Loc: u}
		urls = append(urls, &newUrl)
	}

	var uSet urlset = urlset{Xmlns: root, Urls: urls}
	data, err := xml.MarshalIndent(uSet, "", " ")
	if err != nil {
		log.Fatal(err)
	}

	const (
		// A generic XML header suitable for use with the output of Marshal.
		// This is not automatically added to any output of this package,
		// it is provided as a convenience.
		Header = `<?xml version="1.0" encoding="UTF-8"?>` + "\n"
	)
	data = []byte(xml.Header + string(data))

	err = os.WriteFile("emp.xml", data, 0666)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func Sitemap(root string) {
	foundLinks := make(map[string]bool)
	foundLinks[root] = true
	linksMap, err := findLinks(root, foundLinks)
	if err != nil {
		log.Fatalf("oh no")
	}

	encode(root, linksMap)
}

func findLinks(l string, foundLinks map[string]bool) (map[string]bool, error) {
	fmt.Printf("%s\n", l)
	resp, err := http.Get(l)
	if err != nil {
		log.Fatalf("Couldn't get website %s", l)
	}
	defer resp.Body.Close()

	htmlLinks, err := links.Parse(resp.Body)
	if err != nil {
		log.Fatal("Couldn't extract links")
	}

	headUrl, err := url.Parse(l)
	if err != nil {
		log.Fatal("couldn't parse url")
	}

	for _, l := range htmlLinks {
		if strings.HasPrefix(l.Href, "/") {
			link := "https://" + headUrl.Host + l.Href
			if !foundLinks[link] {
				foundLinks[link] = true
				findLinks(link, foundLinks)
			}
			continue
		}

		subUrl, err := url.Parse(l.Href)
		if err != nil {
			log.Printf("Couldn't Parse %s\n", l.Href)
			continue
		}

		if subUrl.Host == headUrl.Host && !foundLinks[l.Href] {
			foundLinks[l.Href] = true
			findLinks(l.Href, foundLinks)
		}
	}

	return foundLinks, nil
}
