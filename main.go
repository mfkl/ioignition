package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
)

const Port = "8080"

func main() {
	r := chi.NewRouter()

	// Create a route along /files that will serve contents from
	// the ./data/ folder.
	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "public"))
	fileServer(r, "/", filesDir)

	server := http.Server{
		Addr:    ":" + Port,
		Handler: r,
	}

	fmt.Printf("Server listing on port: %s\n", Port)
	log.Fatal(server.ListenAndServe())
}

func fileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("fileserver does not permit any parameters")
	}

	// add trailing '/' if not existing
	if path != "/" && path[len(path)-1] != '/' {
		// letting the caller know that resource has moved from path to path/
		r.Get(path, http.RedirectHandler(path+"/", http.StatusMovedPermanently).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		// ex: public/*
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		// removes public from /public/*
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))

		fs.ServeHTTP(w, r)
	})
}
