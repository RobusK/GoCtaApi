package main

import (
	"github.com/RobusK/GoCtaApi/api"
	"github.com/RobusK/GoCtaApi/dataloaders"
	"github.com/RobusK/GoCtaApi/services"
	"github.com/graphql-go/handler"
	"net/http"
	"os"
	"time"
)

var httpClient = &http.Client{
	Transport: &http.Transport{
		MaxIdleConns:        15,
		IdleConnTimeout:     1 * time.Second,
		TLSHandshakeTimeout: 1 * time.Second,
	},
	Timeout: time.Second * 10,
}

var (
	client            = api.NewAPIClient(getAPIKey(), httpClient)
	directionService  = services.NewDirectionsService(client)
	routeService      = services.NewRoutesService(client)
	stopService       = services.NewStopsService(client)
	//predictionService = services.NewPredictionsService(client)
)

func getAPIKey() string {
	if len(CtaAPIKey) == 0 {
		return os.Getenv("CTA_KEY")
	}
	return CtaAPIKey
}

func main() {
	schema := gqlSchema(dataloaders.NewRetriever())

	h := handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true,
	})
	staticFileHandler := http.FileServer(http.Dir("./static/"))
	dlMiddleware := dataloaders.Middleware(*client)
	http.Handle("/graphql", dlMiddleware(disableCors(h)))
	http.Handle("/", staticFileHandler)
	http.ListenAndServe(serverHostname+":"+serverPort, nil)
}

func disableCors(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, Content-Length, Accept-Encoding")

		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Max-Age", "86400")
			w.WriteHeader(http.StatusOK)
			return
		}

		h.ServeHTTP(w, r)
	})
}
