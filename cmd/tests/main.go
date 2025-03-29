package main

import (
	"github.com/ze674/EZLine/internal/services"
	"time"
)

func main() {
	plc := services.NewPlc(500 * time.Millisecond)
	processService := services.NewProcessTaskService(plc)
	processService.ProcessTask()
	for {
		time.Sleep(2 * time.Second)
	}
}
