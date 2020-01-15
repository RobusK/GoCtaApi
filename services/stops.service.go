package services

import (
	"GoCtaApi/api"
)

type stopKey struct {
	routeID   string
	direction string
}

type StopsService struct {
	stops  map[stopKey][]api.Stop
	client *api.APIClient
}

func NewStopsService(client *api.APIClient) *StopsService {
	return &StopsService{
		client: client,
		stops:  make(map[stopKey][]api.Stop),
	}
}

func (service *StopsService) GetOrCreateStops(routeID string, direction string) []api.Stop {

	key := stopKey{routeID: routeID, direction: direction}

	if service.stops[key] == nil {
		service.stops[key] = service.client.RetrieveStopsForRoute(routeID, direction)
	}
	return service.stops[key]
}
