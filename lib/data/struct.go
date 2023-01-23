package data

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/docker/go-units"
)

type Config struct {
	Debug    bool   `yaml:"debug"`
	DataDir  string `yaml:"data_dir"`
	HttpRoot string `yaml:"http_root"`
	Webhook  string `yaml:"webhook"`
	Strict   bool   `yaml:"strict"`
	RunTime  string `yaml:"run_time"`
	Os       []*Os  `yaml:"os"`
}

type Os struct {
	Path string `yaml:"-"`
	Id   string `yaml:"id"`
	Name string `yaml:"name"`
	// OS theme color
	Color string `yaml:"color"`
	// Is it category
	Category  bool             `yaml:"category"`
	Update    []FlowStep       `yaml:"update"`
	UpdateSum []FlowStep       `yaml:"update_sum"`
	Download  *DownloadOptions `yaml:"download"`
	Child     []*Os            `yaml:"child"`
	MetaData  OsMetaData       `yaml:"-"`
}

func (Os *Os) IconPath() string {
	iconDir := Os.Path
	for {
		iconPath := filepath.Join(Conf.DataDir, iconDir, "icon.png")
		_, err := os.Stat(iconPath)
		if err == nil {
			break
		}
		cutIndex := strings.LastIndex(iconDir, "/")
		iconDir = iconDir[:cutIndex]
	}
	return strings.TrimSuffix(Conf.HttpRoot, "/") + iconDir + "/icon.png"
}

type FlowStep struct {
	URL         string         `yaml:"url"`
	Regex       string         `yaml:"regex"`
	RegexCache  *regexp.Regexp `yaml:"-"`
	RegexSelect int            `yaml:"regex_select"`
	SumRegex    string         `yaml:"sum_regex"`
	SumType     string         `yaml:"sum_type"`
}

type DownloadOptions struct {
	ValidRegex string `yaml:"valid_regex"`
}

type OsMetaData struct {
	Files []*FileData `json:"file"`
}

type FileData struct {
	Name        string    `json:"name"`
	Size        int64     `json:"size"`
	GetTime     time.Time `json:"get_time"`
	SHA512      string    `json:"sha512"`
	OfficialSum *CheckSum `json:"official_sum"`
	Pinned      bool      `json:"pinned"`
	AddedBy     string    `json:"added_by"`
}

type CheckSum struct {
	Sumtype string
	Value   string
}

func (f FileData) SizeString() string {
	return units.HumanSize(float64(f.Size))
}

func (f FileData) GetTimeString() string {
	return f.GetTime.Format("2006/01/02 15:04")
}
