package services

import (
	"GoCtaApi/api"
	"sync"
)

type DirectionsService struct {
	directions      map[string][]string
	directionsMutex sync.RWMutex
	client          *api.APIClient
}

func NewDirectionsService(client *api.APIClient) *DirectionsService {
	return &DirectionsService{
		client:          client,
		directions:      make(map[string][]string),
		directionsMutex: sync.RWMutex{},
	}
}

func (f *DirectionsService) GetOrCreateDirections(routeID string, ch chan<- []string) {
	if f.directions[routeID] == nil {
		directions := f.client.RetrieveDirectionsForRoute(routeID)
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
