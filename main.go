package main

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/cloudflare/cloudflare-go"
)

const IP_API_URL = "http://ip-api.com/json/"

var CLOUDFLARE_API_TOKEN = os.Getenv("IPSTER_CLOUDFLARE_API_TOKEN")
var ZONE_NAME = os.Getenv("IPSTER_CLOUDFLARE_ZONE_NAME")
var DNS_RECORD_NAME = os.Getenv("IPSTER_CLOUDFLARE_DNS_RECORD_NAME")

type Result struct {
	Result string
	Error  error
}

type CFResult struct {
	Result cloudflare.DNSRecord
	Error  error
}

type IP struct {
	Query string
}

func main() {
	verifyEnvVars()
	ticker := time.NewTicker(60 * time.Second)
	for ; true; <-ticker.C {
		log.Println("Verifying IPs")

		ipCh, cfCh := fetchIP(), fetchCF()
		ipRes, cfRes := <-ipCh, <-cfCh

		if cfRes.Error != nil {
			log.Println(cfRes.Error)
			continue
		}

		if ipRes.Error != nil {
			log.Println(ipRes.Error)
			continue
		}

		if cfRes.Result.Content != ipRes.Result {
			log.Println("IPs do not match. Updating...")
			cfRes.Result.Content = ipRes.Result
			err := fixIp(cfRes.Result)
			if err != nil {
				log.Println(err)
			} else {
				log.Println("Record updated!")
			}
			continue
		}
		log.Println("No change")
	}
}

func verifyEnvVars() {
	const msg = `ipster keeps your CloudFlare DNS record in sync with this machine's IP 🤝

Please set the following environmental variables:
	* IPSTER_CLOUDFLARE_API_TOKEN - your CloudFlare API_TOKEN https://dash.cloudflare.com/profile/api-tokens (Use the Edit zone DNS template)
	* IPSTER_CLOUDFLARE_ZONE_NAME - your CloudFlare zone name. Usually your domain name e.g. example.com
	* IPSTER_CLOUDFLARE_DNS_RECORD_NAME - the CloudFlare dns record that you want to keep in sync e.g. home.example.com

Example call:
	IPSTER_CLOUDFLARE_API_TOKEN=xxxxxxxxx_yyyyyyyyyyyyyyyyyyyyyyyyyyyyyy IPSTER_CLOUDFLARE_ZONE_NAME=example.com IPSTER_CLOUDFLARE_DNS_RECORD_NAME=home.example.com ipster`
	if CLOUDFLARE_API_TOKEN == "" {
		log.Fatalln(msg)
		os.Exit(1)
	}
	if ZONE_NAME == "" {
		log.Fatalln(msg)
		os.Exit(2)
	}
	if DNS_RECORD_NAME == "" {
		log.Fatalln(msg)
		os.Exit(3)
	}
}

func fetchIP() <-chan Result {
	ch := make(chan Result)
	go func() {
		defer close(ch)

		client := http.Client{
			Timeout: 5 * time.Second,
		}

		req, err := client.Get(IP_API_URL)
		if err != nil {
			ch <- Result{Result: "", Error: err}
			return
		}

		if req.StatusCode != 200 {
			ch <- Result{Result: "", Error: errors.New("connection failed")}
		}

		defer req.Body.Close()

		body, err := io.ReadAll(req.Body)
		if err != nil {
			ch <- Result{Result: "", Error: err}
		}

		var ip IP
		json.Unmarshal(body, &ip)
		ch <- Result{Result: ip.Query, Error: nil}
	}()

	return ch
}

func fetchCF() <-chan CFResult {
	ch := make(chan CFResult)
	go func() {
		defer close(ch)

		api, err := cloudflare.NewWithAPIToken(CLOUDFLARE_API_TOKEN)
		if err != nil {
			ch <- CFResult{Result: cloudflare.DNSRecord{}, Error: err}
			return
		}

		zoneID, err := api.ZoneIDByName(ZONE_NAME)
		if err != nil {
			ch <- CFResult{Result: cloudflare.DNSRecord{}, Error: err}
			return
		}

		records, err := api.DNSRecords(context.Background(), zoneID, cloudflare.DNSRecord{Name: DNS_RECORD_NAME})
		if err != nil {
			ch <- CFResult{Result: cloudflare.DNSRecord{}, Error: err}
			return
		}
		for _, record := range records {
			if record.Name == DNS_RECORD_NAME {
				ch <- CFResult{Result: record, Error: nil}
			}
		}

		ch <- CFResult{Result: cloudflare.DNSRecord{}, Error: errors.New("record not found")}
	}()

	return ch
}

func fixIp(record cloudflare.DNSRecord) error {
	api, err := cloudflare.NewWithAPIToken(CLOUDFLARE_API_TOKEN)
	if err != nil {
		return err
	}
	zoneID, err := api.ZoneIDByName(ZONE_NAME)
	if err != nil {
		return err
	}
	err = api.UpdateDNSRecord(context.Background(), zoneID, record.ID, record)
	if err != nil {
		return err
	}
	return nil

}
