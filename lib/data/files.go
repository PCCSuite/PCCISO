package data

import (
	"log"
	"os"
	"path/filepath"

	"github.com/PCCSuite/PCCISO/lib/util"
)

var skip = [...]string{
	"icon.png",
	"PCCISO_meta.json",
	"index.html",
}

type ScanFilesResult struct {
	Pass      int
	Corrupted int
	NotFound  int
	Unknown   int
	Failed    int
}

func ScanFiles(Os *Os, result *ScanFilesResult) {
	path := filepath.Join(Conf.DataDir, Os.Path)
	realFilesRaw, err := os.ReadDir(path)
	if err != nil {
		log.Print("Failed to read dir: ", err)
		result.Failed++
		return
	}
	realFiles := make(map[string]os.FileInfo)
file_loop:
	for _, v := range realFilesRaw {
		if v.IsDir() {
			continue
		}
		for _, name := range skip {
			if v.Name() == name {
				continue file_loop
			}
		}
		realFiles[v.Name()], err = v.Info()
		if err != nil {
			log.Print("Failed to get info: ", err)
			result.Failed++
			continue
		}
	}
	files := make([]*FileData, 0, len(Os.MetaData.Files))
	for _, v := range Os.MetaData.Files {
		realFile, ok := realFiles[v.Name]
		if !ok {
			log.Print("File not found: ", filepath.Join(path, v.Name))
			result.NotFound++
			continue
		}
		delete(realFiles, v.Name)
		if v.Size != realFile.Size() {
			result.Corrupted++
			log.Print("Size verification failed: ", filepath.Join(path, v.Name))
			err := os.Remove(filepath.Join(path, v.Name))
			if err != nil {
				log.Print("Failed to remove corrupted file: ", filepath.Join(path, v.Name))
			}
			continue
		}
		if Conf.Strict {
			if Conf.Debug {
				log.Print("Checking sum: ", v.Name)
				log.Printf("Size: %d(%s)", v.Size, v.SizeString())
			}
			valid, err := util.VerifySum(filepath.Join(path, v.Name), v.SHA512, "sha512")
			if err != nil {
				log.Print("Failed to verify sum: ", err)
				result.Failed++
			}
			if !valid {
				result.Corrupted++
				log.Print("Sum verification failed: ", filepath.Join(path, v.Name))
				err := os.Remove(filepath.Join(path, v.Name))
				if err != nil {
					log.Print("Failed to remove corrupted file: ", filepath.Join(path, v.Name))
				}
				continue
			}
		}
		files = append(files, v)
		result.Pass++
	}
	Os.MetaData.Files = files
	for k := range realFiles {
		log.Print("Unknown file detected: ", filepath.Join(Os.Path, k))
		result.Unknown++
	}
	SaveOsMeta(Os)
}
