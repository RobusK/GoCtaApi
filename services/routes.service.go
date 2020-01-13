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

func (f *RoutesService) GetOrCreateRoutes() []api.Route {
	if f.routes == nil {
		f.routes = f.client.RetrieveRoutes()
	}
	return f.routes
}
