package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/kardianos/service"
)

// Credentials : Stores login detail and hostname for HTTPS Request.
type Credentials struct {
	Encodedcred, Hostname string
}

var logger service.Logger

type program struct{}

func (p *program) Start(s service.Service) error {
	// Start should not block. Do the actual work async.
	go p.run()
	return nil
}
func (p *program) run() {

	// Creates Log file IF MISSING.
	logFile, err := os.Create(getExePath() + "\\DUCLog.txt")
	if err != nil {
		log.Fatalln("Cant Create Log")
	}
	// Outputs Log to file
	log.SetOutput(logFile)

	// Checked for Conig file exists. Exits if not.
	if fileExists(getExePath() + "\\DUCConfig.json") {
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
		setup()
	}

	err = s.Run()
	if err != nil {
		logger.Error(err)
	}

}

// Set up funtion triggred by setup flag. Created config file and log files.
func setup() {

	data := Credentials{
		Encodedcred: base64.StdEncoding.EncodeToString([]byte(os.Args[2] + ":" + os.Args[3])),
		Hostname:    os.Args[4],
	}

	file, _ := json.MarshalIndent(data, "", " ")
	_ = ioutil.WriteFile("DUCConfig.json", file, 0644)
	_ = ioutil.WriteFile("DUCip.txt", nil, 0644)
	_ = ioutil.WriteFile("DUClog.txt", nil, 0644)

	log.Println("Setup Complete. Please install as service. Run : Go-DUC.exe -service install in command line.")
	os.Exit(3)
}

// Sends GET request to NO-IP based on intrgration requirements. https://www.noip.com/integrate/request
func updateIP(newip string) {

	client := &http.Client{}
	jsonFile, _ := os.Open(getExePath() + "\\DUCConfig.json")
	byteData, _ := ioutil.ReadAll(jsonFile)

	var credentials Credentials
	json.Unmarshal(byteData, &credentials)

	var url = "https://dynupdate.no-ip.com/nic/update?hostname=" + credentials.Hostname + "&myip=" + newip
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Host", string(credentials.Hostname))
	req.Header.Add("Authorization", "Basic "+credentials.Encodedcred)
	req.Header.Set("User-Agent", "Barajas-net Go-DUC-1.0/WIN maintainer-goduc@barajas-net.com")
	if err != nil {
		log.Fatalln(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("No-IP Responce: " + string(body))

	if string(body) == "badauth" {
		os.Exit(4)
	}
}

// Polls ipify for ip address.
func getIP() {

	resp, err := http.Get("https://api.ipify.org")
	if err != nil {
		log.Println(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}

	lastip, err := ioutil.ReadFile(getExePath() + "\\DUCip.txt")
	if err != nil {
		log.Fatalln(err)
	}

	if string(lastip) != string(body) {
		err = ioutil.WriteFile(getExePath()+"\\DUCip.txt", body, 0644)
		log.Println("IP change detected new IP: " + string(body))
		if err != nil {
			panic(err)
		}
		updateIP(string(body))
	}

}

// Main loop triggers ip check / update
func runloop() {
	pollInterval := 1000

	timerCh := time.Tick(time.Duration(pollInterval) * time.Millisecond)

	for range timerCh {
		getIP()
	}

}

// Helper Functions

func getExePath() (path string) {
	exe, _ := os.Executable()
	exPath := filepath.Dir(exe)
	return exPath
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
