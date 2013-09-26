package main

import (
	"net/http"
)

// AssetsHandler is a static file server that serves everything in
// static directory
func assetsHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, r.URL.Path[1:])
}