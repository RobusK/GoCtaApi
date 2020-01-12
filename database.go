package main

import (
	"sync"
)

// Database is a cache for information coming from the CTA
type Database struct {
	routes          []Route
	directions      map[string][]string
	directionsMutex sync.RWMutex
	client          *APIClient
}
// NewDatabase is the constructor of Database
func NewDatabase(client *APIClient) *Database {
	return &Database{client: client, directions: make(map[string][]string), directionsMutex: sync.RWMutex{}}
}

func (f *Database) getOrCreateRoutes() []Route {
	if f.routes == nil {
		f.routes = f.client.retrieveRoutes()
	}
	return f.routes
}

func (f *Database) getOrCreateDirections(routeID string, ch chan<- []string) {
	if f.directions[routeID] == nil {
		directions := f.client.retrieveDirectionsForRoute(routeID)
		directionsArray := make([]string, len(directions))
		for index, element := range directions {
			directionsArray[index] = element.Value
		}
		f.directionsMutex.Lock()
		f.directions[routeID] = directionsArray
		f.directionsMutex.Unlock()
	}
	ch <- f.directions[routeID]
}
