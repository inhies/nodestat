package main

import (
	"html/template"
	"net/http"
)

var (
	t *template.Template
)

// Serve HTTP
func Serve() {
	// Start the HTTP JSON API if enabled.
	if SystemConfig.Access.JSONApi.Enabled {
		l.Infoln("Starting HTTP JSON API")
		http.HandleFunc("/peers/", peerStatsHandler)
		http.HandleFunc("/node/", nodeStatsHandler)
		http.HandleFunc("/all/", allStatsHandler)

		// If we set the config option to true, show the
		// front-end.
		if SystemConfig.Web.EnableFrontEnd {
			l.Infoln("Starting HTTP front-end")
			http.HandleFunc("/static/", assetsHandler)
			http.HandleFunc("/", rootHandler)
			// Compile templates
			t = template.Must(template.ParseGlob("templates/*.html"))
		}

		// Listen and serve, bitches!
		http.ListenAndServe(SystemConfig.Access.JSONApi.Addr, nil)
	}
}

// AssetsHandler is a static file server that serves everything in
// static directory
func assetsHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, r.URL.Path[1:])
}

// RootHandler handles the "/" connections
func rootHandler(w http.ResponseWriter, r *http.Request) {
	t.ExecuteTemplate(w, "index", nil)
}
