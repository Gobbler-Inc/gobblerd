package ui

import (
	"embed"
	"html/template"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/alfreddobradi/go-bb-man/helper"
)

var (
	//go:embed dist
	dist embed.FS

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

	sub, err := fs.Sub(dist, "dist")
	if err != nil {
		log.Printf("Failed to substitute root path: %v", err)
		helper.E(w, http.StatusInternalServerError)
		return
	}
	fs := http.FS(sub)

	path := r.URL.Path
	log.Printf("Checking if %s exists in embedded fs...", path)
	if _, err := fs.Open(path); os.IsNotExist(err) {
		log.Printf("Failed to open %s: %v", path, err)

		// For now I don't see any point in making this more universal but in the future we might want to add more exceptions
		if filepath.Ext("path") == ".ico" {
			helper.E(w, http.StatusNotFound)
			return
		}

		log.Printf("Falling back to index.html")
		f, err := fs.Open("index.html")
		if err != nil {
			log.Printf("Failed to open fallback file (index.html): %v", err)
			helper.E(w, http.StatusInternalServerError)
			return
		}

		indexContents, err := io.ReadAll(f)
		if err != nil {
			log.Printf("Failed to read fallback file (index.html): %v", err)
			helper.E(w, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html")
		w.Write(indexContents) // nolint
		return
	} else if err != nil {
		log.Printf("Failed to open %s: %v", path, err)
		helper.E(w, http.StatusInternalServerError)
		return
	}

	http.FileServer(fs).ServeHTTP(w, r)
}

type SpaHandler struct{}

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
