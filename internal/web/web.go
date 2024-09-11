package web

import (
	"embed"
	"io"
	"io/fs"
	"log"
	"net/http"
	"path/filepath"
)

//go:embed dist
var webAssets embed.FS

func CreateMiddlewareWebFiles() {
	distFS, err := fs.Sub(webAssets, "dist")
	if err != nil {
		log.Fatal(err)
	}

	// Handle requests
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Suppression du '/' initial pour obtenir le chemin relatif
		filePath := r.URL.Path
		if filePath == "/" {
			filePath = "/index.html"
		}

		// Essayer d'ouvrir le fichier demandé
		file, err := distFS.Open(filePath[1:]) // Enlève le '/' initial du chemin
		if err != nil {
			// Si le fichier n'existe pas, servir index.html
			file, err = distFS.Open("index.html")
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		}
		defer file.Close()

		// Lire le contenu du fichier et le servir
		content, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Déterminer le type de contenu en fonction de l'extension du fichier
		contentType := "text/html; charset=utf-8"
		extension := filepath.Ext(filePath)

		// TODO: change that as quick as possible

		if extension == ".js" {
			contentType = "application/javascript"
		}
		if extension == ".css" {
			contentType = "text/css"
		}
		if extension == ".png" {
			contentType = "image/png"
		}

		w.Header().Set("Content-Type", contentType)
		w.WriteHeader(http.StatusOK)
		w.Write(content)
	})

}
