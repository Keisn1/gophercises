package main

import (
	"encoding/xml"
	"io/ioutil"
	"log"
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

func main() {

	//*************************************
	// struct to XML file
	//*************************************
	e1 := &Url{Loc: "Ich.org"}
	e2 := &Url{Loc: "Du.org"}

	var urls []*Url = []*Url{e1, e2}
	var uSet urlset = urlset{Xmlns: "Wir.org", Urls: urls}

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

	err = ioutil.WriteFile("emp.xml", data, 0666)
	if err != nil {
		log.Fatal(err)
	}
}
