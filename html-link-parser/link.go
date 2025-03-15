package links

import (
	"io"
	"log"
	"os"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type Link struct {
	Href string
	Text string
}

var l Link = Link{
	Href: "hello",
	Text: "world",
}

func Exercise(htmlFilename string) []Link {
	r, err := os.Open(htmlFilename)
	if err != nil {
		log.Fatal(err)
	}
	links, _ := Parse(r)
	return links
}

// Parse will take an html document and will return a slice of Link
func Parse(r io.Reader) ([]Link, error) {
	doc, err := html.Parse(r)
	if err != nil {
		log.Fatal(err)
	}

	nodes := linkNodes(doc)
	var htmlLinks []Link
	for _, node := range nodes {
		htmlLinks = append(htmlLinks, buildLink(node))
	}

	return htmlLinks, nil
}

func linkNodes(n *html.Node) []*html.Node {
	if n.Type == html.ElementNode && n.DataAtom == atom.A {
		return []*html.Node{n}
	}
	var ret []*html.Node
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		ret = append(ret, linkNodes(c)...)
	}
	return ret
}

func buildLink(n *html.Node) Link {
	var ret Link
	for _, attr := range n.Attr {
		if attr.Key == "href" {
			ret.Href = strings.TrimSpace(attr.Val)
			break
		}
	}
	ret.Text = getTextNode(n)
	return ret
}

func getTextNode(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}

	if n.Type != html.ElementNode {
		return ""
	}

	var ret string
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		ret += getTextNode(c)
	}
	return strings.Join(strings.Fields(ret), " ")
}
