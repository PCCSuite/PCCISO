package server

import (
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/PCCSuite/PCCISO/lib/data"
)

var defaultFiles = [...]string{
	"copy-icon.png",
	"style.css",
}

func PlaceDefaultFiles() {
	for _, v := range defaultFiles {
		_, err := os.Stat(filepath.Join(data.Conf.DataDir, v))
		if err != nil && errors.Is(err, os.ErrNotExist) {
			src, err := os.Open(filepath.Join("web", v))
			if err != nil {
				log.Fatalf("Failed to open %s: %v", v, err)
			}
			defer src.Close()
			dst, err := os.Create(filepath.Join(data.Conf.DataDir, v))
			if err != nil {
				log.Fatalf("Failed to create %s: %v", v, err)
			}
			defer dst.Close()
			io.Copy(dst, src)
		}
	}

}

func StartServer() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", serveStatic)

	log.Print("Server starting")
	log.Print(http.ListenAndServe("127.0.0.1:8080", mux))
}
