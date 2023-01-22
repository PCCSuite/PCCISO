package getter

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/PCCSuite/PCCISO/lib/data"
	"github.com/PCCSuite/PCCISO/lib/util"
)

type UpdateResult string

const (
	UpdateResultSkippedSum       UpdateResult = "skipped due to sum matched"
	UpdateResultSkippedName      UpdateResult = "skipped due to filename matched"
	UpdateResultSkippedNoProcess UpdateResult = "skipped due to update not specified"
	UpdateResultDone             UpdateResult = "done"
	UpdateResultFailed           UpdateResult = "failed"
)

func ExecUpdate(Os *data.Os) (result UpdateResult) {
	defer func() {
		err := recover()
		if err != nil {
			log.Print("Failed to update: ", err)
			result = UpdateResultFailed
		}
	}()
	if Os.Update == nil {
		return UpdateResultSkippedNoProcess
	}
	url, filename := GetUpdateURL(Os)
	filePath := filepath.Join(data.Conf.DataDir, Os.Path, filename)
	var sum *data.CheckSum
	for i, v := range Os.MetaData.Files {
		if v.Name == filename {
			// filename matched
			if Os.UpdateSum != nil && v.OfficialSum != nil {
				sum = GetChecksum(Os, url)
				if sum.Sumtype == v.OfficialSum.Sumtype && sum.Value == v.OfficialSum.Value {
					return UpdateResultSkippedSum
				} else {
					log.Printf("Removing file %s due to sum not matched", filePath)
					err := os.Remove(filePath)
					if err != nil {
						log.Printf("Failed to remove %s: %v", filePath, err)
					}
					Os.MetaData.Files = append(Os.MetaData.Files[:i], Os.MetaData.Files[i+1:]...)
					break
				}
			} else {
				return UpdateResultSkippedName
			}
		}
	}
	size := DownloadFile(filePath, url)
	fileinfo, err := os.Stat(filePath)
	if err != nil {
		log.Printf("Failed to stat %s: %v", filePath, err)
	}
	if fileinfo.Size() != size {
		log.Printf("Download size is not match: %s", filePath)
	}
	if sum == nil && Os.UpdateSum != nil {
		sum = GetChecksum(Os, url)
	}
	if sum != nil {
		if data.Conf.Debug {
			log.Print("Checking sum: ", *sum)
		}
		valid, err := util.VerifySum(filePath, sum.Value, sum.Sumtype)
		if err != nil {
			log.Printf("Failed to verify sum %s: %v", filePath, err)
			return UpdateResultFailed
		}
		if !valid {
			log.Printf("Downloaded file sum not matched: %s", filePath)
			return UpdateResultFailed
		}
	}
	var sha512 string
	if sum != nil && sum.Sumtype == "sha512" {
		sha512 = sum.Value
	} else {
		file, err := os.Open(filePath)
		if err != nil {
			log.Printf("Failed to open %s: %v", filePath, err)
			return UpdateResultFailed
		}
		defer file.Close()
		sha512 = util.GetSum(file, "sha512")
	}
	Os.MetaData.Files = append(Os.MetaData.Files, &data.FileData{
		Name:        filename,
		Size:        size,
		GetTime:     time.Now(),
		SHA512:      sha512,
		OfficialSum: sum,
		Pinned:      false,
		AddedBy:     "auto_updater",
	})
	data.SaveOsMeta(Os)
	return UpdateResultDone
}
