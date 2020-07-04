package dataloaders

import (
	"context"
	"github.com/RobusK/GoCtaApi/api"
	"time"
)

type contextKey string

const key = contextKey("dataloaders")

// Loaders holds references to the individual dataloaders.
type Loaders struct {
	PredictionByStopIDs *PredictionLoader
	// individual loaders will be defined here
}

func newLoaders(ctx context.Context, client api.Client) *Loaders {
	return &Loaders{
		PredictionByStopIDs: newPredictionByStopIDs(ctx, client),
	}
}

// Retriever retrieves dataloaders from the request context.
type Retriever interface {
	Retrieve(context.Context) *Loaders
}

type retriever struct {
	key contextKey
}

func (r *retriever) Retrieve(ctx context.Context) *Loaders {
	return ctx.Value(r.key).(*Loaders)
}

// NewRetriever instantiates a new implementation of Retriever.
func NewRetriever() Retriever {
	return &retriever{key: key}
}

func newPredictionByStopIDs(ctx context.Context, client api.Client) *PredictionLoader {
	return NewabcLoader(PredictionLoaderConfig{
		MaxBatch: 10,
		Wait:     5 * time.Millisecond,
		Fetch: func(stopIDs []string) ([][]api.Prediction, []error) {
			// db query
			res, err := client.RetrievePredictionsForStops(stopIDs)
			if err != nil {
				return nil, []error{err}
			}
			// map
			groupByStopID := make(map[string][]api.Prediction, len(stopIDs))
			for _, prediction := range res.Predictions {

				val, _ := groupByStopID[prediction.StopID]
				updatedPredictionList := append(val, prediction)
				groupByStopID[prediction.StopID] = updatedPredictionList
			}
			// order
			result := make([][]api.Prediction, len(stopIDs))
			for i, stopID := range stopIDs {
				result[i] = groupByStopID[stopID]
			}
			return result, nil
		},
	})
}