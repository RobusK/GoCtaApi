package api_test

import (
	api2 "GoCtaApi/api"
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
)

// RoundTripFunc .
type RoundTripFunc func(req *http.Request) *http.Response

// RoundTrip .
func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

//NewTestClient returns *http.Client with Transport replaced to avoid making real calls
func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(fn),
	}
}

func TestDoStuffWithRoundTripper(t *testing.T) {
	assert := assert.New(t)
	client := NewTestClient(func(req *http.Request) *http.Response {
		assert.Equal(req.URL.Query()["stpid"][0], "123")
		fmt.Println(req.Form.Get("stopID"))
		return &http.Response{
			StatusCode: 200,
			// Send response to be tested
			Body: ioutil.NopCloser(bytes.NewBufferString(`{}`)),
			// Must be set to non-nil value or it panics
			Header: make(http.Header),
		}
	})

	api := api2.NewAPIClient("123", client)
	resp, _ := api.RetrievePredictionsForStopAndRoute("123", "456")
	fmt.Print(resp.Predictions)

}
