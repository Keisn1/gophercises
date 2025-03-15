package links

import (
	"fmt"
	"testing"
)

func TestExercise(t *testing.T) {
	testCases := []struct {
		htmlFilename string
		wantedLinks  []Link
	}{
		{"ex1.html", []Link{{Href: "/other-page", Text: "A link to another page"}}},
		{"ex2.html", []Link{
			{
				Href: "https://www.twitter.com/joncalhoun",
				Text: "Check me out on twitter",
			},
			{
				Href: "https://github.com/gophercises",
				Text: "Gophercises is on Github!",
			}}},
		{"ex3.html", []Link{
			{
				Href: "https://twitter.com/marcusolsson",
				Text: "@marcusolsson",
			},
			{
				Href: "/lost",
				Text: "Lost? Need help?",
			},
			{
				Href: "#",
				Text: "Login",
			}}},
		{"ex4.html", []Link{
			{
				Href: "/dog-cat",
				Text: "dog cat",
			}}},
	}

	fmt.Println("running tests...")

	for _, tc := range testCases {
		htmlLinks := Exercise(tc.htmlFilename)
		isSame(tc.wantedLinks, htmlLinks, t)
	}
}

func isSame(wantedLinks, htmlLinks []Link, t *testing.T) {
	if len(wantedLinks) != len(htmlLinks) {
		t.Errorf("Not same amount of links")
	}
	for _, wL := range wantedLinks {
		present := false
		var idx int
		for i, L := range htmlLinks {
			if L.Href == wL.Href && L.Text == wL.Text {
				present = true
				idx = i
				break
			}
		}
		if !present {
			t.Errorf("Link with Url '%s' and Text '%s' not present", wL.Href, wL.Text)
			t.Errorf("Link with Url '%s' and Text '%s' not present", wL.Href, wL.Text)
		}
		htmlLinks = append(htmlLinks[:idx], htmlLinks[idx+1:]...)
	}
}
