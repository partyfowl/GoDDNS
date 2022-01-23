package main

import (
	"bytes"
	"encoding/json"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os/user"
	"regexp"
	"strings"
)

type config struct {
	ApiToken string `yaml:"api_token"`
	ZoneId   string `yaml:"zone_id"`
}

type cloudflareResponseResult struct {
	Name     string `json:"name"`
	ZoneName string `json:"zone_name"`
	Content  string `json:"content"`
	Id       string `json:"id"`
}

type cloudflareResponseError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type cloudflareResponse struct {
	Errors []cloudflareResponseError  `json:"errors"`
	Result []cloudflareResponseResult `json:"result"`
}

func getPublicIp() string {
	response, err := http.Get("https://checkip.amazonaws.com/")

	if err != nil {
		log.Fatal(err)
	}
	if response.StatusCode != http.StatusOK {
		log.Fatal(response.StatusCode)
	}

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		log.Fatal(err)
	}

	ip := strings.ReplaceAll(string(body), "\n", "")

	matched, err := regexp.MatchString(`^\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}$`, ip)

	if err != nil {
		log.Fatal(err)
	}
	if !matched {
		log.Fatal("Invalid IP: ", ip)
	}

	return ip
}

func callCloudflareDnsApi(httpRequestType string, apiToken string, zoneId string, data []byte, id string) *http.Response {
	req, err := http.NewRequest(
		httpRequestType,
		"https://api.cloudflare.com/client/v4/zones/"+zoneId+"/dns_records/"+id,
		bytes.NewBuffer(data))

	req.Header.Add("Authorization", "Bearer "+apiToken)
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}

	response, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	return response

}

func decodeToCloudflareResponse(httpBody io.Reader) *cloudflareResponse {
	body, err := ioutil.ReadAll(httpBody)

	if err != nil {
		log.Fatal(err)
	}
	cfResp := &cloudflareResponse{}

	err = json.Unmarshal(body, cfResp)

	if err != nil {
		log.Fatal(err)
	}

	return cfResp
}

func getDnsRecords(apiToken string, zoneId string) *cloudflareResponse {
	response := callCloudflareDnsApi("GET", apiToken, zoneId, nil, "")
	return decodeToCloudflareResponse(response.Body)
}

func updateDns(ip string, apiToken string, zoneId string) {
	recordName := "goddns"

	values := map[string]interface{}{"type": "A", "name": recordName, "content": ip, "proxied": true, "ttl": 60}
	jsonValue, _ := json.Marshal(values)

	var httpResponse *http.Response = nil

	for _, dnsRecord := range getDnsRecords(apiToken, zoneId).Result {
		if dnsRecord.Name == recordName+"."+dnsRecord.ZoneName {
			if dnsRecord.Content == ip {
				log.Print("No changes to make, GoDDNS record already set to ", ip)
				return
			}
			httpResponse = callCloudflareDnsApi("PUT", apiToken, zoneId, jsonValue, dnsRecord.Id)
		}
	}

	if httpResponse == nil {
		httpResponse = callCloudflareDnsApi("POST", apiToken, zoneId, jsonValue, "")
	}

	if httpResponse.StatusCode == http.StatusOK {
		log.Print("Created/updated record for ", ip)
	} else {
		cfResp := decodeToCloudflareResponse(httpResponse.Body)
		log.Fatal(httpResponse.StatusCode, cfResp.Errors[0].Code, cfResp.Errors[0].Message)
	}
}

func loadConfig() *config {
	currentUser, _ := user.Current()
	yamlFile, err := ioutil.ReadFile(currentUser.HomeDir + "/.goddns/config.yaml")
	if err != nil {
		log.Fatal(err)
	}
	conf := &config{}
	err = yaml.Unmarshal(yamlFile, conf)
	if err != nil {
		log.Fatal(err)
	}
	return conf
}

func main() {
	conf := loadConfig()
	ip := getPublicIp()
	updateDns(ip, conf.ApiToken, conf.ZoneId)
}
