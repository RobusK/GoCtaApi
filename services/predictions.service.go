package services

import "GoCtaApi/api"

type PredictionsService struct {
	client *api.Client
}

func NewPredictionsService(client *api.Client) *PredictionsService {
	return &PredictionsService{client: client}
}

func (service *PredictionsService) GetOrCreatePredictionsForStopAndRoute(stopID string, routeID string) api.PredictionList {
	if v, err := service.client.RetrievePredictionsForStopAndRoute(stopID, routeID); err != nil {
		return *v
	}
	return api.PredictionList{}
}

func (service *PredictionsService) GetOrCreatePredictionsForStop(stopID string) (*[]api.Prediction, error) {
	predList, err := service.client.RetrievePredictionsForStopAndRoute(stopID, "")
	if err != nil {
		return nil, err
	}
	return &predList.Predictions, nil
}
