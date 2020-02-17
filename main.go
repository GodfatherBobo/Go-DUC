package main

import (
	"GO-DUC/helpers"
	"GO-DUC/utils"
	"flag"
	"log"
	"os"
	"time"

	"github.com/kardianos/service"
)

var logger service.Logger

type program struct{}

func (p *program) Start(s service.Service) error {
	// Start should not block. Do the actual work async.
	go p.run()
	return nil
}
func (p *program) run() {

	// Creates Log file IF MISSING.
	logFile, err := os.Create(helpers.GetExePath() + "\\DUCLog.txt")
	if err != nil {
		log.Fatalln("Cant Create Log")
	}
	// Outputs Log to file
	log.SetOutput(logFile)

	// Checked for Conig file exists. Exits if not.
	if helpers.FileExists(helpers.GetExePath() + "\\DUCConfig.json") {
		log.Println("Config Detected. Starting Application.")
		runloop()

	} else {
		log.Println("Please run with -setup arg with the email and password and hostname assosated with the account Example: duc.exe -setup myemail mypassword myhostname")
		os.Exit(1)
	}

}
func (p *program) Stop(s service.Service) error {
	// Stop should not block. Return with a few seconds.
	return nil
}

// Main function controls flag parsing and service requests.
func main() {

	svcFlag := flag.String("service", "", "Control the system service.")
	setupFlag := flag.String("setup", "", "Starts setup function of application.")
	flag.Parse()

	options := make(service.KeyValue)
	options["Restart"] = "on-success"
	options["SuccessExitStatus"] = "1 2 8 SIGKILL"
	svcConfig := &service.Config{
		Name:        "GoDUC",
		DisplayName: "Dynamic DNS Update Client",
		Description: "Dynamic DNS update client for No-IP.",
		Option:      options,
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}
	errs := make(chan error, 5)
	logger, err = s.Logger(errs)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			err := <-errs
			if err != nil {
				log.Print(err)
			}
		}
	}()

	if len(*svcFlag) != 0 {
		err := service.Control(s, *svcFlag)
		if err != nil {
			log.Printf("Valid actions: %q\n", service.ControlAction)
			log.Fatal(err)
		}
		return
	}

	if len(*setupFlag) != 0 {
		utils.Setup()
	}

	err = s.Run()
	if err != nil {
		logger.Error(err)
	}

}

//Main loop triggers ip check / update
func runloop() {
	pollInterval := 1000

	timerCh := time.Tick(time.Duration(pollInterval) * time.Millisecond)

	for range timerCh {
		utils.GetIP()
	}
}
