package ui

import (
	"embed"
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
)

type SpaHandler struct{}

func NewSpaHandler() SpaHandler {
	return SpaHandler{}
}

func (ah SpaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
