package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const baseURL = "http://www.ctabustracker.com/bustime/api/v2"

// Client is a HTTP client making requests to the CTA API
type Client struct {
	APIKey     string
	httpclient *http.Client
}

// NewAPIClient creates a new instance of a API Client
func NewAPIClient(key string, httpclient *http.Client) *Client {
	client := &Client{
		APIKey:     key,
		httpclient: httpclient,
	}

	return client
}

// RetrieveRoutes retrieves all routes from the CTA API
func (c Client) RetrieveRoutes() []Route {
	fmt.Println("Retrieving stops")

	req, _ := http.NewRequest("GET", baseURL+"/getroutes", nil)

	q := req.URL.Query()
	q.Add("key", c.APIKey)
	q.Add("format", "json")
	req.URL.RawQuery = q.Encode()

	resp, httpDoErr := c.httpclient.Do(req)
	if httpDoErr != nil {
		fmt.Println(httpDoErr)
		return []Route{}
	}
	defer resp.Body.Close()

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

	req, _ := http.NewRequest("GET", baseURL+"/getdirections", nil)

	q := req.URL.Query()
	q.Add("key", c.APIKey)
	q.Add("format", "json")
	q.Add("rt", routeID)
	req.URL.RawQuery = q.Encode()

	resp, httpDoErr := c.httpclient.Do(req)
	if httpDoErr != nil {
		fmt.Println(httpDoErr)
		return []Direction{}
	}
	defer resp.Body.Close()

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

	req, _ := http.NewRequest("GET", baseURL+"/getstops", nil)

	q := req.URL.Query()
	q.Add("key", c.APIKey)
	q.Add("format", "json")
	q.Add("rt", routeID)
	q.Add("dir", direction)
	req.URL.RawQuery = q.Encode()

	resp, httpDoErr := c.httpclient.Do(req)

	if httpDoErr != nil {
		fmt.Println(httpDoErr)
		return []Stop{}
	}
	defer resp.Body.Close()

	parsed := busTimeStopsResponse{}
	err := json.NewDecoder(resp.Body).Decode(&parsed)

	if err != nil {
		fmt.Println(err)
	}
	return parsed.RoutesList.Stops
}

//RetrievePredictionsForStopAndRoute retrieves predicted arrival/deprature times for a stop and a route
func (c Client) RetrievePredictionsForStopAndRoute(stopID string, routeID string) (*PredictionList, error) {
	fmt.Println("Retrieving Predictions for StopID, routeID ", stopID, routeID)

	req, _ := http.NewRequest("GET", baseURL+"/getpredictions", nil)

	q := req.URL.Query()
	q.Add("key", c.APIKey)
	q.Add("format", "json")
	q.Add("rt", routeID)
	q.Add("stpid", stopID)
	req.URL.RawQuery = q.Encode()

	resp, httpDoErr := c.httpclient.Do(req)
	if httpDoErr != nil {
		fmt.Println(httpDoErr)
		return nil, httpDoErr
	}
	defer resp.Body.Close()

	parsed := busTimePredictionsResponse{}
	err := json.NewDecoder(resp.Body).Decode(&parsed)

	if err != nil {
		fmt.Println(err)
	}
	return &parsed.Predictions, nil
}
