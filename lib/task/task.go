package task

import (
	"log"

	"github.com/PCCSuite/PCCISO/lib/data"
	"github.com/PCCSuite/PCCISO/lib/getter"
	"github.com/PCCSuite/PCCISO/lib/server"
)

func ExecPeriodic() {
	scanResult := ScanAllFiles()
	log.Print("Scan Result: ", scanResult)
	updateResult := UpdateAllOs()
	log.Print("Update Result: ", updateResult)
	limitResult := LimitAllFiles()
	log.Print("Limit file num result: ", limitResult)
	generateResult := GenerateAllPages()
	log.Print("Page generation result: ", generateResult)
	NotifyResult(scanResult, updateResult, limitResult, generateResult)
}

func ScanAllFiles() (scanResult data.ScanFilesResult) {
	for _, v := range data.OsList {
		scanFiles(v, &scanResult)
	}
	return
}

func scanFiles(Os *data.Os, result *data.ScanFilesResult) {
	data.ScanFiles(Os, result)
	for _, v := range Os.Child {
		scanFiles(v, result)
	}
}

func LimitAllFiles() (limitResult data.LimitFilesResult) {
	for _, v := range data.OsList {
		limitFiles(v, &limitResult)
	}
	return
}

func limitFiles(Os *data.Os, result *data.LimitFilesResult) {
	data.LimitFiles(Os, result)
	for _, v := range Os.Child {
		limitFiles(v, result)
	}
}

type UpdateResults struct {
	SkippedSum       int
	SkippedName      int
	SkippedNoProcess int
	Done             int
	Failed           int
}

func UpdateAllOs() (result UpdateResults) {
	for _, v := range data.OsList {
		updateOs(v, &result)
	}
	return
}

func updateOs(Os *data.Os, result *UpdateResults) {
	if data.Conf.Debug {
		log.Print("Updating: ", Os.Name)
	}
	res := getter.ExecUpdate(Os)
	if data.Conf.Debug {
		log.Print("Update result: ", res)
	}
	switch res {
	case getter.UpdateResultSkippedSum:
		result.SkippedSum++
	case getter.UpdateResultSkippedName:
		result.SkippedName++
	case getter.UpdateResultSkippedNoProcess:
		result.SkippedNoProcess++
	case getter.UpdateResultDone:
		result.Done++
	case getter.UpdateResultFailed:
		result.Failed++
	}
	for _, v := range Os.Child {
		updateOs(v, result)
	}
}

func GenerateAllPages() (success bool) {
	defer func() {
		err := recover()
		if err != nil {
			log.Print("Failed to generate pages: ", err)
			success = false
		}
	}()
	server.GenerateIndexPage()
	for _, v := range data.OsList {
		generatePage(v)
	}
	success = true
	return
}

func generatePage(Os *data.Os) {
	server.GenerateOsPage(Os)
	for _, v := range Os.Child {
		generatePage(v)
	}
}
