package data

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

const configFile = "./config.yaml"
const osFile = "./os.yaml"

var Conf Config

var OsList []*Os

func init() {
	file, err := os.ReadFile(configFile)
	if err != nil {
		log.Fatal("Config load failed:", err)
	}
	err = yaml.Unmarshal(file, &Conf)
	if err != nil {
		log.Fatal("Config parse failed:", err)
	}

	file, err = os.ReadFile(osFile)
	if err != nil {
		log.Fatal("OS list load failed:", err)
	}
	err = yaml.Unmarshal(file, &OsList)
	if err != nil {
		log.Fatal("Config parse failed:", err)
	}

	if Conf.Debug {
		log.Print("Debug is enabled")
	}

	for _, v := range OsList {
		processOsData(v, "")
	}
}

func processOsData(Os *Os, path string) {
	Os.Path = path + "/" + Os.Id
	Os.MetaData = loadOsMeta(Os)
	err := os.Mkdir(filepath.Join(Conf.DataDir, Os.Path), os.ModePerm)
	if err != nil && !errors.Is(err, os.ErrExist) {
		log.Panicf("Failed to make dir %s: %v", filepath.Join(Conf.DataDir, Os.Path), err)
	}
	if Os.Child == nil {
		return
	}
	for _, v := range Os.Child {
		processOsData(v, Os.Path)
	}
}

func loadOsMeta(Os *Os) (meta OsMetaData) {
	file, err := os.ReadFile(filepath.Join(Conf.DataDir, Os.Path, "PCCISO_meta.json"))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return
		}
		log.Panic("OS Metadata load failed:", err)
	}
	err = json.Unmarshal(file, &meta)
	if err != nil {
		log.Panic("OS Metadata unmarshal failed:", err)
	}
	return
}

func SaveOsMeta(Os *Os) {
	raw, err := json.MarshalIndent(Os.MetaData, "", "	")
	if err != nil {
		log.Panic("OS Metadata marshal failed:", err)
	}
	file := filepath.Join(Conf.DataDir, Os.Path, "PCCISO_meta.json")
	backupFile := filepath.Join(Conf.DataDir, Os.Path, "PCCISO_meta.backup.json")
	err = os.Rename(file, backupFile)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			log.Panic("Backup PCCISO_meta.json failed:", err)
		}
	}
	err = os.WriteFile(file, raw, os.ModePerm)
	if err != nil {
		log.Panic("Write PCCISO_meta.json failed:", err)
	}
	err = os.Remove(backupFile)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			log.Panic("Remove PCCISO_meta_backup.json failed:", err)
		}
	}
}
