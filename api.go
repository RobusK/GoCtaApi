package main

type BusTimeRoutesResponse struct {
	RoutesList RoutesList `json:"bustime-response"`
}

type BusTimeDirectionsResponse struct {
	DirectionList DirectionList `json:"bustime-response"`
}

type RoutesList struct {
	Routes []Route
}

type DirectionList struct {
	Directions []Direction `json:"directions"`
}

type Route struct {
	RouteId    string `json:"rt"`
	CommonName string `json:"rtnm"`
	Color      string `json:"rtclr"`
	Rtdd       string
}

type Direction struct {
	Value string `json:"dir"`
}
