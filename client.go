package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// APIClient is a HTTP client making requests to the CTA API
type APIClient struct {
	transport *http.Transport
}
// NewAPIClient creates a new instance of a API Client
func NewAPIClient() *APIClient {
	client := &APIClient{}
	client.transport = &http.Transport{
		MaxIdleConns:        15,
		IdleConnTimeout:     1 * time.Second,
		TLSHandshakeTimeout: 1 * time.Second,
	}
	return client
}

func (c APIClient) retrieveRoutes() []Route {
	fmt.Println("Retrieving stops")

	client := &http.Client{Transport: c.transport, Timeout: time.Second * 10,}
	req, _ := http.NewRequest("GET", "http://www.ctabustracker.com/bustime/api/v2/getroutes", nil)

	q := req.URL.Query()
	q.Add("key", CtaAPIKey)
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

func (c APIClient) retrieveDirectionsForRoute(routeID string) []Direction {
	fmt.Println("Retrieving Directions for routeID ", routeID)

	client := &http.Client{Transport: c.transport, Timeout: time.Second * 10,}
	req, _ := http.NewRequest("GET", "http://www.ctabustracker.com/bustime/api/v2/getdirections", nil)

	q := req.URL.Query()
	q.Add("key", CtaAPIKey)
	q.Add("format", "json")
	q.Add("rt", routeID)
	req.URL.RawQuery = q.Encode()

	resp, _ := client.Do(req)
	parsed := BusTimeDirectionsResponse{}
	err := json.NewDecoder(resp.Body).Decode(&parsed)

	if err != nil {
		fmt.Println(err)
	}
	return parsed.DirectionList.Directions
}
