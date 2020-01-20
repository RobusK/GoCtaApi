package api

type busTimeRoutesResponse struct {
	RoutesList RoutesList `json:"bustime-response"`
}

type busTimeStopsResponse struct {
	RoutesList StopList `json:"bustime-response"`
}

type busTimePredictionsResponse struct {
	Predictions PredictionList `json:"bustime-response"`
}

// PredictionList is the response we get from CTA for arrivals
type PredictionList struct {
	Predictions []Prediction    `json:"prd"`
	Error       []PredictionError `json:"error"`
}

// PredictionError describes an error in a Prediction request
type PredictionError struct {
	RouteID string `json:"rt"`
	StopID  string `json:"stpid"`
	Message string `json:"msg"`
}
// Prediction describes a prediction from the predictions API
type Prediction struct {
	PrdGenerated  string `json:"tmstmp"`
	Type          string `json:"typ"`
	StopName      string `json:"stpnm"`
	StopID        string `json:"stpid"`
	VehicleID     string `json:"vid"`
	Distance      int    `json:"dstp"`
	RouteID       string `json:"rt"`
	Direction     string `json:"rtdir"`
	Destination   string `json:"des"`
	PredictedTime string `json:"prdtm"`
	Delayed       bool   `json:"dly"`
	TimeLeft      string `json:"prdctdn"`
}

// BusTimeDirectionsResponse is the top level response of Directions API
type BusTimeDirectionsResponse struct {
	DirectionList DirectionList `json:"bustime-response"`
}

// StopList is a successful response from stop API
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
// Stop is a representation of a CTA Stop
type Stop struct {
	StopID     string  `json:"stpid"`
	CommonName string  `json:"stpnm"`
	Lat        float64 `json:"lat"`
	Lon        float64 `json:"lon"`
}

// Direction is a representation of a CTA Route direction
type Direction struct {
	Value string `json:"dir"`
}
