package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const baseURL = "http://www.ctabustracker.com/bustime/api/v2"

// Client is a HTTP client making requests to the CTA API
type Client struct {
	transport *http.Transport
	APIKey    string
}

// NewAPIClient creates a new instance of a API Client
func NewAPIClient(key string) *Client {
	client := &Client{
		APIKey: key,
		transport: &http.Transport{
			MaxIdleConns:        15,
			IdleConnTimeout:     1 * time.Second,
			TLSHandshakeTimeout: 1 * time.Second,
		}}

	return client
}

// RetrieveRoutes retrieves all routes from the CTA API
func (c Client) RetrieveRoutes() []Route {
	fmt.Println("Retrieving stops")

	client := &http.Client{Transport: c.transport, Timeout: time.Second * 10,}
	req, _ := http.NewRequest("GET", baseURL+"/getroutes", nil)

	q := req.URL.Query()
	q.Add("key", c.APIKey)
	q.Add("format", "json")
	req.URL.RawQuery = q.Encode()

	resp, _ := client.Do(req)
	parsed := busTimeRoutesResponse{}
	err := json.NewDecoder(resp.Body).Decode(&parsed)
	if err != nil {
		fmt.Println(err)
	}
	return parsed.RoutesList.Routes
}

//RetrieveDirectionsForRoute retrieves all Directions based on a Route
func (c Client) RetrieveDirectionsForRoute(routeID string) []Direction {
	fmt.Println("Retrieving Directions for routeID ", routeID)

	client := &http.Client{Transport: c.transport, Timeout: time.Second * 10,}
	req, _ := http.NewRequest("GET", baseURL+"/getdirections", nil)

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

//RetrieveStopsForRoute Retrieves all stops for a route
func (c Client) RetrieveStopsForRoute(routeID string, direction string) []Stop {
	fmt.Println("Retrieving Stops for routeID ", routeID, direction)

	client := &http.Client{Transport: c.transport, Timeout: time.Second * 10,}
	req, _ := http.NewRequest("GET", baseURL+"/getstops", nil)

	q := req.URL.Query()
	q.Add("key", c.APIKey)
	q.Add("format", "json")
	q.Add("rt", routeID)
	q.Add("dir", direction)
	req.URL.RawQuery = q.Encode()

	resp, _ := client.Do(req)
	parsed := busTimeStopsResponse{}

	err := json.NewDecoder(resp.Body).Decode(&parsed)

	if err != nil {
		fmt.Println(err)
	}
	return parsed.RoutesList.Stops
}

//RetrievePredictionsForStopAndRoute retrieves predicted arrival/deprature times for a stop and a route
func (c Client) RetrievePredictionsForStopAndRoute(stopID string, routeID string) PredictionList {
	fmt.Println("Retrieving Predictions for StopID, routeID ", stopID, routeID)

	client := &http.Client{Transport: c.transport, Timeout: time.Second * 10,}
	req, _ := http.NewRequest("GET", baseURL+"/getpredictions", nil)

	q := req.URL.Query()
	q.Add("key", c.APIKey)
	q.Add("format", "json")
	q.Add("rt", routeID)
	q.Add("stpid", stopID)
	req.URL.RawQuery = q.Encode()

	resp, _ := client.Do(req)

	parsed := busTimePredictionsResponse{}
	err := json.NewDecoder(resp.Body).Decode(&parsed)

	if err != nil {
		fmt.Println(err)
	}
	return parsed.Predictions
}
