package utils

import (
	"GO-DUC/helpers"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

//UpdateIP Sends GET request to NO-IP based on intrgration requirements. https://www.noip.com/integrate/request
func UpdateIP(newip string) {
	client := &http.Client{}
	jsonFile, _ := os.Open(helpers.GetExePath() + "\\DUCConfig.json")
	byteData, _ := ioutil.ReadAll(jsonFile)
	var credentials = helpers.Credentials{}
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
