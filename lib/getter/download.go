package getter

import (
	"log"
	"net/url"
	"time"

	"github.com/PCCSuite/PCCISO/lib/data"
	"github.com/cavaliergopher/grab/v3"
	"github.com/docker/go-units"
)

func DownloadFile(filePath string, url *url.URL) int64 {
	client := grab.DefaultClient
	req, err := grab.NewRequest(filePath, url.String())
	if err != nil {
		log.Panicf("Failed to make request to %s: %v", url, err)
	}
	resp := client.Do(req)
	if data.Conf.Debug {
		log.Printf("Downloading %s from %s", filePath, url)
		log.Printf("Download size: %d(%s)", resp.Size(), units.HumanSize(float64(resp.Size())))
		log.Printf("Download resumable: %t, resumed: %t", resp.CanResume, resp.DidResume)
		ticker := time.NewTicker(30 * time.Second)
	progress_loop:
		for {
			select {
			case <-ticker.C:
				log.Printf("Download progress: %f ETA: %s", resp.Progress()*100, time.Until(resp.ETA()).String())
			case <-resp.Done:
				break progress_loop
			}
		}
		ticker.Stop()
		log.Printf("Downloading done")
	}
	if resp.Err() != nil {
		log.Panicf("Failed to get %s: %v", url, resp.Err())
	}
	return resp.Size()
}
