package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type ApiClient struct {
	tr *http.Transport
}

func NewApiClient() *ApiClient {
	client := &ApiClient{}
	client.tr = &http.Transport{
		MaxIdleConns:        15,
		IdleConnTimeout:     1 * time.Second,
		TLSHandshakeTimeout: 1 * time.Second,
	}
	return client
}

func (c ApiClient) RetrieveRoutes() []Route {
	fmt.Println("Retrieving stops")

	client := &http.Client{Transport: c.tr, Timeout: time.Second * 10,}
	req, _ := http.NewRequest("GET", "http://www.ctabustracker.com/bustime/api/v2/getroutes", nil)

	q := req.URL.Query()
	q.Add("key", CtaApiKey)
	q.Add("format", "json")
	req.URL.RawQuery = q.Encode()

	resp, _ := client.Do(req)
	parsed := BusTimeRoutesResponse{}
	err := json.NewDecoder(resp.Body).Decode(&parsed)

	if err != nil {
		fmt.Println(err)
	}
	return parsed.RoutesList.Routes
}

func (c ApiClient) RetrieveDirectionsForRoute(routeId string) []Direction {
	fmt.Println("Retrieving Directions for routeId ", routeId)

	client := &http.Client{Transport: c.tr, Timeout: time.Second * 10,}
	req, _ := http.NewRequest("GET", "http://www.ctabustracker.com/bustime/api/v2/getdirections", nil)

	q := req.URL.Query()
	q.Add("key", CtaApiKey)
	q.Add("format", "json")
	q.Add("rt", routeId)
	req.URL.RawQuery = q.Encode()

	resp, _ := client.Do(req)
	parsed := BusTimeDirectionsResponse{}
	err := json.NewDecoder(resp.Body).Decode(&parsed)

	if err != nil {
		fmt.Println(err)
	}
	return parsed.DirectionList.Directions
}
