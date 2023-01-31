package getter

import (
	"errors"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"net/url"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/PCCSuite/PCCISO/lib/data"
	"github.com/PCCSuite/PCCISO/lib/util"
)

// Return url and filename
func GetUpdateURL(Os *data.Os) (*url.URL, string) {
	if Os.Update == nil {
		return nil, ""
	}
	nowURL := ""
	for _, v := range Os.Update {
		if v.URL != "" {
			nowURL = v.URL
		} else if v.Regex != "" {
			if v.RegexCache == nil {
				v.RegexCache = regexp.MustCompile(v.Regex)
			}
			nowPage, err := getRequest(nowURL)
			if err != nil {
				log.Panicf("Failed to get page %s: %v", nowURL, err)
			}
			allMatch := v.RegexCache.FindAllSubmatch(nowPage, -1)
			var match [][]byte
			if v.RegexSelect >= 0 {
				match = allMatch[v.RegexSelect]
			} else {
				match = allMatch[len(allMatch)+v.RegexSelect]
			}
			result := html.UnescapeString(string(match[1]))
			if strings.HasPrefix(result, "http://") || strings.HasPrefix(result, "https://") {
				nowURL = result
			} else if strings.HasPrefix(result, "/") {
				nowURL = nowURL[:8] + strings.SplitN(nowURL[8:], "/", 2)[0] + result
			} else {
				nowURL = nowURL[:strings.LastIndex(nowURL, "/")+1] + result
			}
		}
	}
	url, filename, err := util.GetUrlInfo(nowURL)
	if err != nil {
		log.Panic("Failed to get url info: ", err)
	}
	return url, filename
}

func GetChecksum(Os *data.Os, fileURL *url.URL) *data.CheckSum {
	if Os.UpdateSum == nil {
		return nil
	}
	nowURL := ""
	for _, v := range Os.UpdateSum {
		if v.URL != "" {
			nowURL = v.URL
			nowURL = strings.ReplaceAll(nowURL, "$FILENAME", filepath.Base(fileURL.Path))
			nowURL = strings.ReplaceAll(nowURL, "$FILEURL", fileURL.String())
		} else if v.Regex != "" {
			if v.RegexCache == nil {
				v.RegexCache = regexp.MustCompile(v.Regex)
			}
			nowPage, err := getRequest(nowURL)
			if err != nil {
				log.Panicf("Failed to get page %s: %v", nowURL, err)
			}
			allMatch := v.RegexCache.FindAllSubmatch(nowPage, -1)
			var match [][]byte
			if v.RegexSelect >= 0 {
				match = allMatch[v.RegexSelect]
			} else {
				match = allMatch[len(allMatch)+v.RegexSelect]
			}
			result := html.UnescapeString(string(match[1]))
			if strings.HasPrefix(result, "http://") || strings.HasPrefix(result, "https://") {
				nowURL = result
			} else if strings.HasPrefix(result, "/") {
				nowURL = nowURL[:8] + strings.SplitN(nowURL[8:], "/", 2)[0] + result
			} else {
				nowURL = nowURL[:strings.LastIndex(nowURL, "/")+1] + result
			}
		} else if v.SumRegex != "" {
			nowPage, err := getRequest(nowURL)
			if err != nil {
				log.Panicf("Failed to get page %s: %v", nowURL, err)
			}
			rawRegex := strings.ReplaceAll(v.SumRegex, "$FILENAME", filepath.Base(fileURL.Path))
			rawRegex = strings.ReplaceAll(rawRegex, "$FILEURL", fileURL.String())
			regex := regexp.MustCompile(rawRegex)
			allMatch := regex.FindAllSubmatch(nowPage, -1)
			var match [][]byte
			if v.RegexSelect >= 0 {
				match = allMatch[v.RegexSelect]
			} else {
				match = allMatch[len(allMatch)+v.RegexSelect]
			}
			return &data.CheckSum{
				Sumtype: v.SumType,
				Value:   string(match[1]),
			}
		} else {
			log.Panic("Nothing to do")
		}
	}
	log.Panic("No sum")
	return nil
}

func getRequest(url string) ([]byte, error) {
	var lastErr error
	for i := 0; i < data.Conf.Retry; i++ {
		if i != 0 {
			log.Print("Error, retrying...: ", lastErr)
			time.Sleep(1 * time.Millisecond)
		}
		resp, err := http.Get(url)
		if err != nil {
			lastErr = fmt.Errorf("failed to get: %w", err)
			continue
		}
		defer resp.Body.Close()
		raw, err := io.ReadAll(resp.Body)
		if err != nil {
			lastErr = fmt.Errorf("failed read responce: %w", err)
			continue
		}
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			lastErr = errors.New(fmt.Sprint("status code not success: ", resp.StatusCode))
			continue
		}
		return raw, nil
	}
	return nil, lastErr
}
