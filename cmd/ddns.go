package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"os/exec"
	"strconv"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

var (
	configFile = flag.String("f", "", "config file path")

	usageMessage = "\nUsage:" +
		"\n     ddns -f ${CONFIG_FILE}" +
		"\n" +
		"\n[-f] set config file path. ex) -f /root/tools/config.yaml" +
		"\n"
)

type Config struct {
	Authorization string `yaml:"Authorization"`
	ZoneID        string `yaml:"ZoneID"`
	DomainName    string `yaml:"DomainName"`
}

type Zone struct {
	CurrentVersionID string `json:"current_version_id"`
}

type Record struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Type        string    `json:"type"`
	EnableAlias bool      `json:"enable_alias"`
	TTL         int       `json:"ttl"`
	Records     []Address `json:"records"`
}

type Address struct {
	Address string `json:"address"`
}

type ClientInfo struct {
	Type string `json:"type"`
	Data Data   `json:"data"`
}

type Data struct {
	IPAddress string `json:"iPAddress"`
}

func readConfig(c string) (Config, error) {
	var config Config
	buf, err := ioutil.ReadFile(c)
	if err != nil {
		return config, err
	}
	err = yaml.Unmarshal(buf, &config)
	if err != nil {
		return config, err
	}
	return config, nil
}

func main() {
	flag.Parse()
	if *configFile == "" {
		fmt.Println(usageMessage)
	}

	config, err := readConfig(*configFile)
	if err != nil {
		panic(err)
	}

	apiURL := "https://api.gis.gehirn.jp/dns/v1/zones/" + config.ZoneID
	req, _ := http.NewRequest("GET", apiURL, nil)
	req.Header.Set("Authorization", config.Authorization)
	req.Header.Set("Content-Type", "application/json")
	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil || resp.Status != "200 OK" {
		panic(err)
	}
	defer resp.Body.Close()
	dumpResp, _ := httputil.DumpResponse(resp, true)
	dumpBody := strings.Split(string(dumpResp), "\r\n\r\n")[1]
	var zone Zone
	err = json.Unmarshal([]byte(string(dumpBody)), &zone)
	if err != nil {
		panic(err)
	}
	zoneVersion := zone.CurrentVersionID

	apiURL = apiURL + "/versions/" + zoneVersion + "/records/"
	req, _ = http.NewRequest("GET", apiURL, nil)
	req.Header.Set("Authorization", config.Authorization)
	req.Header.Set("Content-Type", "application/json")
	client = new(http.Client)
	resp, err = client.Do(req)
	if err != nil || resp.Status != "200 OK" {
		panic(err)
	}
	defer resp.Body.Close()
	dumpResp, _ = httputil.DumpResponse(resp, true)
	dumpBody = strings.Split(string(dumpResp), "\r\n\r\n")[1]
	var records []Record
	err = json.Unmarshal([]byte(string(dumpBody)), &records)
	if err != nil {
		panic(err)
	}
	var recordID string
	var addr string
	var rectype string
	var enableAlias bool
	var ttl int
	for _, r := range records {
		if r.Name == config.DomainName+"." {
			recordID = r.ID
			addr = r.Records[0].Address
			rectype = r.Type
			enableAlias = r.EnableAlias
			ttl = r.TTL
		}
	}
	if recordID == "" || addr == "" || rectype == "" {
		panic(err)
	}

	out, err := exec.Command("curl", "cli.fyi/me").Output()
	if err != nil {
		panic(err)
	}
	var info ClientInfo
	err = json.Unmarshal([]byte(string(out)), &info)
	if err != nil || info.Data.IPAddress == "" {
		panic(err)
	}

	if addr == info.Data.IPAddress {
		fmt.Println("INFO:", "No Update")
	} else {
		apiURL = apiURL + recordID
		jsonStr := `{"id":"` + recordID + `","name":"` + config.DomainName + `.","type":"` + rectype + `","enable_alias":` + strconv.FormatBool(enableAlias) + `,"ttl":` + strconv.Itoa(ttl) + `,"records":[{"address":"` + info.Data.IPAddress + `"}]}`

		req, _ = http.NewRequest("PUT", apiURL, bytes.NewBuffer([]byte(jsonStr)))
		req.Header.Set("Authorization", config.Authorization)
		req.Header.Set("Content-Type", "application/json")
		client = new(http.Client)
		resp, err = client.Do(req)
		if err != nil || resp.Status != "200 OK" {
			panic(err)
		}
		defer resp.Body.Close()
		fmt.Println("INFO:", "Update Done", addr, "=>", info.Data.IPAddress)
	}
}
