package main

import (
	"sync"
)

type Database struct {
	routes          []Route
	directions      map[string][]string
	directionsMutex sync.RWMutex
	client          *ApiClient
}

func NewDatabase(client *ApiClient) *Database {
	return &Database{client: client, directions: make(map[string][]string), directionsMutex: sync.RWMutex{}}
}

func (f *Database) getOrCreateRoutes() []Route {
	if f.routes == nil {
		f.routes = f.client.RetrieveRoutes()
	}
	return f.routes
}

func (f *Database) getOrCreateDirections(routeId string, ch chan<- []string) {
	if f.directions[routeId] == nil {
		directions := f.client.RetrieveDirectionsForRoute(routeId)
		directionsArray := make([]string, len(directions))
		for index, element := range directions {
			directionsArray[index] = element.Value
		}
		f.directionsMutex.Lock()
		f.directions[routeId] = directionsArray
		f.directionsMutex.Unlock()
	}
	ch <- f.directions[routeId]
}
