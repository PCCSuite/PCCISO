package getter

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/PCCSuite/PCCISO/lib/data"
	"github.com/docker/go-units"
)

func DownloadFile(filePath string, url *url.URL) int64 {
	file, err := os.Create(filePath)
	if err != nil {
		log.Panic("Failed to create file: ", err)
	}
	defer file.Close()
	resp, err := http.Get(url.String())
	if err != nil {
		log.Panicf("Failed to get %s: %v", url, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		log.Panic("status code not success: ", resp.StatusCode)
	}
	if data.Conf.Debug {
		log.Printf("Downloading %s from %s", filePath, url)
		log.Printf("Download size: %d(%s)", resp.ContentLength, units.HumanSize(float64(resp.ContentLength)))
		defer log.Printf("Downloading done")
	}
	size, err := io.Copy(file, resp.Body)
	if err != nil {
		log.Panic("Failed to write file: ", err)
	}
	return size
}
