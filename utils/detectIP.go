package utils

import (
	"GO-DUC/helpers"
	"io/ioutil"
	"log"
	"net/http"
)

//GetIP Polls api.ipify.org for ip address and triggers a change if required.
func GetIP() {

	resp, err := http.Get("https://api.ipify.org")
	if err != nil {
		log.Println(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}

	lastip, err := ioutil.ReadFile(helpers.GetExePath() + "\\DUCip.txt")
	if err != nil {
		log.Fatalln(err)
	}

	if string(lastip) != string(body) {
		err = ioutil.WriteFile(helpers.GetExePath()+"\\DUCip.txt", body, 0644)
		log.Println("IP change detected new IP: " + string(body))
		if err != nil {
			panic(err)
		}
		UpdateIP(string(body))
	}

}
