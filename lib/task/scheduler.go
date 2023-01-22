package task

import (
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/PCCSuite/PCCISO/lib/data"
)

func TaskScheduler() {
	splitTime := strings.Split(data.Conf.RunTime, ":")
	hour, err := strconv.Atoi(splitTime[0])
	if err != nil {
		log.Panic("Invalid run_time")
	}
	min, err := strconv.Atoi(splitTime[1])
	if err != nil {
		log.Panic("Invalid run_time")
	}
	for {
		next := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), hour, min, 0, 0, time.Local)
		if next.Before(time.Now()) {
			next = time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day()+1, hour, min, 0, 0, time.Local)
		}
		if data.Conf.Debug {
			log.Print("Sleep until: ", next.Format(time.RFC3339))
			log.Print("Sleep for: ", time.Until(next))
		}
		time.Sleep(time.Until(next))
		log.Print("Periodic Task executing...")
		ExecPeriodic()
	}
}
