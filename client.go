package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type ApiClient struct {
}

func (c ApiClient) RetrieveRoutes() []Route {
	fmt.Println("Retrieving stops")
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
	parsed := BusTimeRoutesResponse{}
	err := json.NewDecoder(resp.Body).Decode(&parsed)

	if err != nil {
		fmt.Println(err)
	}
	return parsed.RoutesList.Routes
}

func (c ApiClient) RetrieveDirectionsForRoute(routeId string) []Direction {
	fmt.Println("Retrieving Directions")
	tr := &http.Transport{
		MaxIdleConns:    10,
		IdleConnTimeout: 10 * time.Second,
	}

	client := &http.Client{Transport: tr}
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
