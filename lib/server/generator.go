package server

import (
	"html/template"
	"log"
	"os"
	"path/filepath"
	"sort"

	"github.com/PCCSuite/PCCISO/lib/data"
)

func GenerateIndexPage() {
	tmpl, err := template.ParseFiles("web/index.tmpl")
	if err != nil {
		log.Panic("Failed to parse os.tmpl: ", err)
	}
	file, err := os.Create(filepath.Join(data.Conf.DataDir, "index.html"))
	if err != nil {
		log.Panic("Failed to open test.html: ", err)
	}
	defer file.Close()

	err = tmpl.Execute(file, data.OsList)
	if err != nil {
		log.Panic("Failed to exec tmpl: ", err)
	}
}

type OsPageInfo struct {
	Root     string
	OsName   string
	Color    string
	Pinned   []FileInfo
	Latest   []FileInfo
	Versions []*data.Os
	Files    []FileInfo
}

type FileInfo struct {
	data.FileData
	Path string
}

func GenerateOsPage(Os *data.Os) {
	tmpl, err := template.ParseFiles("web/os.tmpl")
	if err != nil {
		log.Panic("Failed to parse os.tmpl: ", err)
	}
	file, err := os.Create(filepath.Join(data.Conf.DataDir, Os.Path, "index.html"))
	if err != nil {
		log.Panic("Failed to open test.html: ", err)
	}
	defer file.Close()
	info := OsPageInfo{
		Root:     data.Conf.HttpRoot,
		OsName:   Os.Name,
		Color:    Os.Color,
		Pinned:   []FileInfo{},
		Latest:   []FileInfo{},
		Versions: Os.Child,
		Files:    make([]FileInfo, 0, len(Os.MetaData.Files)),
	}

	var latest FileInfo
	for _, v := range Os.MetaData.Files {
		fileInfo := FileInfo{
			FileData: *v,
			Path:     v.Name,
		}
		if v.Pinned {
			info.Pinned = append(info.Pinned, fileInfo)
		} else {
			if !v.GetTime.Before(latest.GetTime) {
				latest = fileInfo
			}
		}
		info.Files = append(info.Files, fileInfo)
	}
	if latest.Name != "" {
		info.Latest = append(info.Latest, latest)
	}

	if Os.Category {
		for _, child := range Os.Child {
			var latest FileInfo
			for _, v := range child.MetaData.Files {
				fileInfo := FileInfo{
					FileData: *v,
					Path:     child.Id + "/" + v.Name,
				}
				if v.Pinned {
					info.Pinned = append(info.Pinned, fileInfo)
				} else {
					if !v.GetTime.Before(latest.GetTime) {
						latest = fileInfo
					}
				}
				info.Files = append(info.Files, fileInfo)
			}
			if latest.Name != "" {
				info.Latest = append(info.Latest, latest)
			}
		}
	}

	sort.SliceStable(info.Files, func(i, j int) bool {
		return info.Files[i].Name < info.Files[j].Name
	})

	err = tmpl.Execute(file, info)
	if err != nil {
		log.Panic("Failed to exec tmpl: ", err)
	}
}
