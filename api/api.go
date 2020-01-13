package api

// BusTimeRoutesResponse is the top level response of Routes API
type BusTimeRoutesResponse struct {
	RoutesList RoutesList `json:"bustime-response"`
}

type BusTimeStopsResponse struct {
	RoutesList StopList `json:"bustime-response"`
}

// BusTimeDirectionsResponse is the top level response of Directions API
type BusTimeDirectionsResponse struct {
	DirectionList DirectionList `json:"bustime-response"`
}

type StopList struct {
	Stops []Stop `json:"stops"`
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

type Stop struct {
	StopID     string  `json:"stpid"`
	CommonName string  `json:"stpnm"`
	Lat        float32 `json:"lat"`
	Lon        float32 `json:"lon"`
}

// Direction is a representation of a CTA Route direction
type Direction struct {
	Value string `json:"dir"`
}
