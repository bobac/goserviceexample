package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/kardianos/service"
	"golang.org/x/sys/windows/registry"
	"golang.org/x/sys/windows/svc/eventlog"
)

type program struct {
	eventLog *eventlog.Log
}

var logFile *os.File

func (p *program) Start(s service.Service) error {
	// Start logic
	log.Println("Starting service...")
	var err error
	p.eventLog, err = eventlog.Open("GoServiceExample")
	if err != nil {
		return err
	}
	p.eventLog.Info(100, "Service starting")
	go p.run()
	return nil
}

func (p *program) Stop(s service.Service) error {
	// Stop logic
	log.Println("Stopping service...")
	p.eventLog.Info(101, "Service stopping")
	p.eventLog.Close()
	return nil
}

func (p *program) run() {
	go startTicker(p.eventLog)
	go startHttpServer()
}

func main() {
	var err error
	logPathFileName := filepath.Join(filepath.Dir(os.Args[0]), "goserviceexample.log")
	logFile, err = os.OpenFile(logPathFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Printf("Failed to open log file: %v", err)
		return
	}
	defer logFile.Close()

	log.SetOutput(logFile)
	log.Println("Running main()...")

	svcConfig := &service.Config{
		Name:        "GoServiceExample",
		DisplayName: "Go Service Example",
		Description: "This is an example Go service that logs to the Windows event log.",
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "install":
			log.Println("Installing service...")
			err := installService(s)
			if err != nil {
				log.Fatal("Failed:", err)
			}
			return
		case "uninstall":
			log.Println("Uninstalling service...")
			err := uninstallService(s)
			if err != nil {
				log.Fatal("Failed:", err)
			}
			return
		case "start":
			err = s.Start()
			if err != nil {
				log.Println(err)
			}
			return
		case "stop":
			err = s.Stop()
			if err != nil {
				log.Println(err)
			}
			return
		}
		fmt.Printf("Valid commands: %q\n", []string{"install", "uninstall", "start", "stop"})
		return
	}

	if err = s.Run(); err != nil {
		log.Fatal(err)
	}
}

func installService(s service.Service) error {
	if err := s.Install(); err != nil {
		return err
	}
	return createEventSource("GoServiceExample", "Application")
}

func uninstallService(s service.Service) error {
	if err := s.Uninstall(); err != nil {
		log.Printf("Failed to uninstall service: %v", err)
		return err
	}
	return nil
}

func createEventSource(source, logName string) error {
	k, _, err := registry.CreateKey(registry.LOCAL_MACHINE, `SYSTEM\CurrentControlSet\Services\EventLog\`+logName+`\`+source, registry.SET_VALUE)
	if err != nil {
		log.Printf("Failed to create event source in Registry: %v", err)
		return err
	}
	defer k.Close()

	// Set the path to your service's executable as the event message file
	exePath, err := os.Executable()
	if err != nil {
		log.Printf("Failed to get executable path: %v", err)
		return err
	}
	err = k.SetStringValue("EventMessageFile", exePath)
	if err != nil {
		log.Printf("Failed to set event message file: %v", err)
		return err
	}

	// Set the types of events your service can log
	err = k.SetDWordValue("TypesSupported", 7)
	if err != nil {
		log.Printf("Failed to set supported types: %v", err)
		return err
	}

	return nil
}
