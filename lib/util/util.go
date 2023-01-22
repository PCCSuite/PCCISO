package util

import (
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

func GetFilenameFromUrl(url string) string {
	return url[strings.LastIndex(url, "/")+1:]
}

func VerifySum(filePath string, sum string, sumtype string) (bool, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return false, err
	}
	hash := GetSum(file, sumtype)
	return sum == hash, nil
}

func GetSum(data io.Reader, sumtype string) string {
	var hasher hash.Hash
	switch sumtype {
	case "sha1":
		hasher = sha1.New()
	case "sha256":
		hasher = sha256.New()
	case "sha512":
		hasher = sha512.New()
	}
	_, err := io.Copy(hasher, data)
	if err != nil {
		log.Panic("Failed to copy to hasher: ", err)
	}
	return hex.EncodeToString(hasher.Sum(make([]byte, 0, 64)))
}

// return redirected url and filename
func GetUrlInfo(url string) (*url.URL, string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, "", fmt.Errorf("failed to request: %w", err)
	}
	resp.Body.Close()
	contentDisposition := strings.Split(resp.Header.Get("Content-Disposition"), ";")
	if len(contentDisposition) >= 2 && contentDisposition[0] == "attachment" {
		for _, v := range contentDisposition[1:] {
			statement := strings.TrimSpace(v)
			split := strings.SplitN(statement, "=", 2)
			if strings.TrimSpace(split[0]) == "filename" {
				return resp.Request.URL, strings.TrimPrefix(strings.TrimSuffix(split[1], "\""), "\""), nil
			}
		}
	}
	return resp.Request.URL, filepath.Base(resp.Request.URL.Path), nil
}
