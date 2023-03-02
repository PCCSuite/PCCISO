package task

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/PCCSuite/PCCISO/lib/data"
)

type Webhook struct {
	Embeds []Embed `json:"embeds,omitempty"`
}

type Embed struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Thumbnail   struct {
		Url string `json:"url,omitempty"`
	} `json:"thumbnail,omitempty"`
	Fields []EmbedField `json:"fields,omitempty"`
}

type EmbedField struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func NotifyResult(scanResult data.ScanFilesResult, update UpdateResults, limit data.LimitFilesResult, pages bool) {
	if data.Conf.Webhook == "" {
		return
	}
	embed := Embed{
		Title: "PCCISO reports",
	}
	ok := pages
	if scanResult.Corrupted != 0 || scanResult.Failed != 0 || scanResult.NotFound != 0 || scanResult.Unknown != 0 {
		ok = false
	}
	if update.Failed != 0 {
		ok = false
	}
	if limit.Failed != 0 {
		ok = false
	}
	if ok {
		embed.Thumbnail.Url = "http://icons.iconarchive.com/icons/paomedia/small-n-flat/256/sign-check-icon.png"
	} else {
		embed.Thumbnail.Url = "https://icons.iconarchive.com/icons/paomedia/small-n-flat/256/sign-warning-icon.png"
	}
	embed.Fields = append(embed.Fields, EmbedField{
		Name:  "Data Sanity Check",
		Value: fmt.Sprintf("Passed: %d\nUnknown file: %d\nFile Not Found: %d\nCorrupted: %d\nFailed: %d", scanResult.Pass, scanResult.Unknown, scanResult.NotFound, scanResult.Corrupted, scanResult.Failed),
	})
	embed.Fields = append(embed.Fields, EmbedField{
		Name:  "Update Data",
		Value: fmt.Sprintf("Done: %d\nSkipped with matched sum: %d\nSkipped with matched name: %d\nFailed: %d\nSkipped with no update setting: %d", update.Done, update.SkippedSum, update.SkippedName, update.Failed, update.SkippedNoProcess),
	})
	embed.Fields = append(embed.Fields, EmbedField{
		Name:  "Limiting Number of Version",
		Value: fmt.Sprintf("Passed: %d\nDeleted: %d\nFailed: %d", limit.Pass, limit.Deleted, limit.Failed),
	})
	if pages {
		embed.Fields = append(embed.Fields, EmbedField{
			Name:  "Index page generate",
			Value: "Successful",
		})
	} else {
		embed.Fields = append(embed.Fields, EmbedField{
			Name:  "Index page generate",
			Value: "Failed",
		})
	}
	raw, err := json.Marshal(Webhook{
		Embeds: []Embed{embed},
	})
	if err != nil {
		log.Panic("Failed to marshal Webhook: ", err)
	}
	resp, err := http.Post(data.Conf.Webhook, "application/json", bytes.NewReader(raw))
	if err != nil {
		log.Print("Failed to send webhook: ", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 204 {
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Print("Failed to read responce: ", err)
		}
		log.Printf("Webhook return unexpected code %d: %s", resp.StatusCode, string(data))
	}
}
