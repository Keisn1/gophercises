package main

import (
	"flag"
	"sitemap"
)

func main() {
	urlFlag := flag.String("url", "https://www.calhoun.io/", "Website from which the sitemap is built")
	flag.Parse()
	sitemap.Sitemap(*urlFlag)
}
