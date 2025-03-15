package cyoa

import (
	"log"
	"net/http"
	"strings"
	"text/template"
)

type Story map[string]Chapter

type Chapter struct {
	Title      string   `json:"title"`
	Paragraphs []string `json:"story"`
	Options    []Option `json:"options"`
}

type Option struct {
	Text    string `json:"text"`
	Chapter string `json:"arc"`
}

type handler struct {
	Story  Story
	tpl    *template.Template
	pathFn func(*http.Request) string
}

type HandlerOption func(h *handler)

func WithTemplate(tpl *template.Template) HandlerOption {
	return func(h *handler) {
		h.tpl = tpl
	}
}

var defaultTplName string = "templates/chapter1.html"
var defaultTpl *template.Template = template.Must(template.ParseFiles(defaultTplName))

func WithPath(pathFn func(*http.Request) string) HandlerOption {
	return func(h *handler) {
		h.pathFn = pathFn
	}
}

func defaultPathFn(r *http.Request) string {
	path := strings.TrimSpace(r.URL.Path)
	if path == "/" || path == "" {
		path = "/intro"
	}
	return path[1:]
}

func NewHandler(story Story, options ...HandlerOption) http.Handler {
	h := handler{
		Story:  story,
		tpl:    defaultTpl,
		pathFn: defaultPathFn,
	}

	for _, opt := range options {
		opt(&h)
	}

	return h
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := h.pathFn(r)
	if chapter, ok := h.Story[path]; ok {
		err := h.tpl.Execute(w, chapter)
		if err != nil {
			log.Printf("%v", err)
			http.Error(w, "Something went wrong ...", http.StatusInternalServerError)
		}
		return
	}
	http.Error(w, "Chapter not found", http.StatusNotFound)
}
