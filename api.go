package main

type BusTimeResponse struct {
	RoutesList RoutesList `json:"bustime-response"`
}

type RoutesList struct {
	Routes []Route
}

type Route struct {
	RouteId    string `json:"rt"`
	CommonName string `json:"rtnm"`
	Color      string `json:"rtclr"`
	Rtdd       string
}
