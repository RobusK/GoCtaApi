package main

// BusTimeRoutesResponse is the top level response of Routes API
type BusTimeRoutesResponse struct {
	RoutesList RoutesList `json:"bustime-response"`
}

// BusTimeDirectionsResponse is the top level response of Directions API
type BusTimeDirectionsResponse struct {
	DirectionList DirectionList `json:"bustime-response"`
}

// RoutesList is a successful response from route API
type RoutesList struct {
	Routes []Route
}

// DirectionList is a successful response from Directions API
type DirectionList struct {
	Directions []Direction `json:"directions"`
}
// Route is a representation of a CTA Route
type Route struct {
	RouteID    string `json:"rt"`
	CommonName string `json:"rtnm"`
	Color      string `json:"rtclr"`
	Rtdd       string
}
// Direction is a representation of a CTA Route direction
type Direction struct {
	Value string `json:"dir"`
}
