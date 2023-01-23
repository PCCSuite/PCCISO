package server

import (
	"net/http"

	"github.com/PCCSuite/PCCISO/lib/data"
)

var staticHandler http.Handler = http.StripPrefix(data.Conf.HttpRoot, http.FileServer(http.Dir("data")))

func serveStatic(w http.ResponseWriter, r *http.Request) {
	staticHandler.ServeHTTP(w, r)
}
