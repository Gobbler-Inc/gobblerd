package ui

import (
	"embed"
	"html/template"
	"log"
	"net/http"

	"github.com/alfreddobradi/go-bb-man/helper"
)

var (
	//go:embed assets
	assets embed.FS

	//go:embed templates
	templates embed.FS
)

type AssetHandler struct {
	Prefix string
}

func NewAssetHandler(prefix string) AssetHandler {
	return AssetHandler{Prefix: prefix}
}

func (ah AssetHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s : %s", r.Method, r.URL)

	http.FileServer(http.FS(assets)).ServeHTTP(w, r)
}

type SpaHandler struct {
	staticPath string
	indexPath  string
}

func (h SpaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	index, err := templates.ReadFile("templates/index.html")
	if err != nil {
		log.Printf("Failed to read templates/index.html: %v", err)
		helper.E(w, http.StatusBadRequest)
		return
	}
	w.Write(index)
}

func MainPageHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tpl, err := template.ParseFS(templates, "templates/index.html")
		if err != nil {
			log.Printf("Failed to parse template: %v", err)
			helper.E(w, http.StatusInternalServerError)
			return
		}

		if err := tpl.Execute(w, nil); err != nil {
			log.Printf("Failed to execute template: %v", err)
			helper.E(w, http.StatusInternalServerError)
			return
		}
	}
}
