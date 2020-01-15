package services

import "GoCtaApi/api"

type RoutesService struct {
	routes []api.Route
	client *api.APIClient
}

func NewRoutesService(client *api.APIClient) *RoutesService {
	return &RoutesService{
		client: client,
	}
}

func (service *RoutesService) GetOrCreateRoutes() []api.Route {
	if service.routes == nil {
		service.routes = service.client.RetrieveRoutes()
	}
	return service.routes
}
