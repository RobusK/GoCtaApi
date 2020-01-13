package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// APIClient is a HTTP client making requests to the CTA API
type APIClient struct {
	transport *http.Transport
	APIKey    string
}

// NewAPIClient creates a new instance of a API Client
func NewAPIClient(key string) *APIClient {
	client := &APIClient{
		APIKey: key,
		transport: &http.Transport{
			MaxIdleConns:        15,
			IdleConnTimeout:     1 * time.Second,
			TLSHandshakeTimeout: 1 * time.Second,
		}}

	return client
}

func (c APIClient) RetrieveRoutes() []Route {
	fmt.Println("Retrieving stops")

	client := &http.Client{Transport: c.transport, Timeout: time.Second * 10,}
	req, _ := http.NewRequest("GET", "http://www.ctabustracker.com/bustime/api/v2/getroutes", nil)

	q := req.URL.Query()
	q.Add("key", c.APIKey)
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

func (c APIClient) RetrieveDirectionsForRoute(routeID string) []Direction {
	fmt.Println("Retrieving Directions for routeID ", routeID)

	client := &http.Client{Transport: c.transport, Timeout: time.Second * 10,}
	req, _ := http.NewRequest("GET", "http://www.ctabustracker.com/bustime/api/v2/getdirections", nil)

	q := req.URL.Query()
	q.Add("key", c.APIKey)
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

func (c APIClient) RetrieveStopsForRoute(routeID string, direction string) []Stop {
	fmt.Println("Retrieving Stops for routeID ", routeID, direction)

	client := &http.Client{Transport: c.transport, Timeout: time.Second * 10,}
	req, _ := http.NewRequest("GET", "http://www.ctabustracker.com/bustime/api/v2/getstops", nil)

	q := req.URL.Query()
	q.Add("key", c.APIKey)
	q.Add("format", "json")
	q.Add("rt", routeID)
	q.Add("dir", direction)
	req.URL.RawQuery = q.Encode()

	resp, _ := client.Do(req)
	parsed := BusTimeStopsResponse{}

	err := json.NewDecoder(resp.Body).Decode(&parsed)

	if err != nil {
		fmt.Println(err)
	}
	return parsed.RoutesList.Stops
}
