package ui

import (
	"embed"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"

	"github.com/alfreddobradi/go-bb-man/helper"
	"github.com/sirupsen/logrus"
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
	logger.WithFields(logrus.Fields{
		"method": r.Method,
		"url":    r.URL,
	}).Info()

	sub, err := fs.Sub(dist, "dist")
	if err != nil {
		logger.WithError(err).Error("Failed to substitute root path")
		helper.E(w, http.StatusInternalServerError)
		return
	}
	fs := http.FS(sub)

	path := r.URL.Path
	logger.WithField("path", path).Trace("Checking if path exists in embedded fs")
	if _, err := fs.Open(path); os.IsNotExist(err) {
		// For now I don't see any point in making this more universal but in the future we might want to add more exceptions
		if filepath.Ext(path) == ".ico" {
			helper.E(w, http.StatusNotFound)
			return
		}

		logger.WithField("path", path).Trace("Falling back to index.html")
		f, err := fs.Open("index.html")
		if err != nil {
			logger.WithField("path", path).WithError(err).Error("Failed to open index.html for reading")
			helper.E(w, http.StatusInternalServerError)
			return
		}

		indexContents, err := io.ReadAll(f)
		if err != nil {
			logger.WithField("path", path).WithError(err).Error("Failed to read index.html")
			helper.E(w, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html")
		w.Write(indexContents) // nolint
		return
	} else if err != nil {
		logger.WithField("path", path).WithError(err).Error("Failed to open path in the embedded filesystem")
		helper.E(w, http.StatusInternalServerError)
		return
	}

	http.FileServer(fs).ServeHTTP(w, r)
}
