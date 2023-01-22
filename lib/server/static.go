package server

import "net/http"

var staticHandler http.Handler = http.FileServer(http.Dir("data"))

func serveStatic(w http.ResponseWriter, r *http.Request) {
	staticHandler.ServeHTTP(w, r)
}
