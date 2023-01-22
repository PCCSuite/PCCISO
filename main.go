package main

import (
	"log"
	"os"

	"github.com/PCCSuite/PCCISO/lib/server"
	"github.com/PCCSuite/PCCISO/lib/task"
)

func main() {
	server.PlaceDefaultFiles()
	if len(os.Args) == 1 {
		go task.TaskScheduler()
	} else {
		switch os.Args[1] {
		case "now":
			task.ExecPeriodic()
			return
		default:
			log.Fatal("Invalid argument")
		}
	}
	server.StartServer()
}
