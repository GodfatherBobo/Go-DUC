package utils

import (
	"GO-DUC/helpers"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

//Setup funtion triggred by setup flag. Created config file and log files.
func Setup() {

	data := helpers.Credentials{
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
