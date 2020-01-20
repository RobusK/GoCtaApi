package services

import "GoCtaApi/api"

type PredictionsService struct {
	client *api.Client
}

func NewPredictionsService(client *api.Client) *PredictionsService {
	return &PredictionsService{client: client}
}

func (service *PredictionsService) GetOrCreatePredictionsForStopAndRoute(stopID string, routeID string) api.PredictionList {
	return service.client.RetrievePredictionsForStopAndRoute(stopID, routeID)
}

func (service *PredictionsService) GetOrCreatePredictionsForStop(stopID string, ch chan<- []api.Prediction) {
	ch <- service.client.RetrievePredictionsForStopAndRoute(stopID, "").Predictions
}
