package main

import (
	"log"
	"time"

	"golang.org/x/sys/windows/svc/eventlog"
)

func startTicker(eventLog *eventlog.Log) {
	ticker := time.NewTicker(1 * time.Minute) // Set your interval
	defer ticker.Stop()

	for range ticker.C {
		periodicFunction(eventLog)
	}
}

func periodicFunction(eventLog *eventlog.Log) {
	// Log an informational message
	log.Println("Executing periodic function")
	if err := eventLog.Info(1, "Executing periodic function"); err != nil {
		log.Printf("Failed to write to event log: %v", err)
	}
	// Your function logic here
}
