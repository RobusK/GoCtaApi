package services

import (
	"github.com/RobusK/GoCtaApi/api"
	"sync"
)

type DirectionsService struct {
	directions      map[string][]string
	directionsMutex sync.RWMutex
	client          *api.Client
}

func NewDirectionsService(client *api.Client) *DirectionsService {
	return &DirectionsService{
		client:          client,
		directions:      make(map[string][]string),
		directionsMutex: sync.RWMutex{},
	}
}

func (service *DirectionsService) GetOrCreateDirections(routeID string, ch chan<- []string) {
	if service.directions[routeID] == nil {
		directions := service.client.RetrieveDirectionsForRoute(routeID)
		directionsArray := make([]string, len(directions))
		for index, element := range directions {
			directionsArray[index] = element.Value
		}
		service.directionsMutex.Lock()
		service.directions[routeID] = directionsArray
		service.directionsMutex.Unlock()
	}
	ch <- service.directions[routeID]
}
