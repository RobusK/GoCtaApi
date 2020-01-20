package services

import "GoCtaApi/api"

type RoutesService struct {
	routes []api.Route
	client *api.Client
}

func NewRoutesService(client *api.Client) *RoutesService {
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
