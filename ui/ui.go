package ui

import (
	"embed"
	"log"
	"net/http"
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
