package main

import (
	"GoCtaApi/api"
	"GoCtaApi/services"
	"github.com/graphql-go/handler"
	"net/http"
	"os"
)

var (
	client            = api.NewAPIClient(getAPIKey())
	directionService  = services.NewDirectionsService(client)
	routeService      = services.NewRoutesService(client)
	stopService       = services.NewStopsService(client)
	predictionService = services.NewPredictionsService(client)
)

func getAPIKey() string {
	if len(CtaAPIKey) == 0 {
		return os.Getenv("CTA_KEY")
	}
	return CtaAPIKey
}

func main() {
	schema := gqlSchema()

	h := handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true,
	})

	staticFileHandler := http.FileServer(http.Dir("./static/"))

	http.Handle("/graphql", disableCors(h))
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
