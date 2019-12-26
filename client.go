package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func RetrieveStops() BusTimeResponse {

	tr := &http.Transport{
		MaxIdleConns:    10,
		IdleConnTimeout: 10 * time.Second,
	}

	client := &http.Client{Transport: tr}
	req, _ := http.NewRequest("GET", "http://www.ctabustracker.com/bustime/api/v2/getroutes", nil)

	q := req.URL.Query()
	q.Add("key", CtaApiKey)
	q.Add("format", "json")
	req.URL.RawQuery = q.Encode()

	resp, _ := client.Do(req)
	parsed := BusTimeResponse{}
	err := json.NewDecoder(resp.Body).Decode(&parsed)

	if err != nil {
		fmt.Println(err)
	}
	return parsed
}
