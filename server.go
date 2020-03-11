package main

import (
	"GoCtaApi/api"
	"errors"
	"fmt"
	"github.com/graphql-go/graphql"
	"strings"
)

var routeType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Route",
		Fields: graphql.Fields{
			"RouteID": &graphql.Field{
				Type: graphql.String,
			},
			"CommonName": &graphql.Field{
				Type: graphql.String,
			},
			"Color": &graphql.Field{
				Type: graphql.String,
			},
			"Directions": &graphql.Field{
				Type: graphql.NewList(graphql.String),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					route, _ := p.Source.(api.Route)
					channel := make(chan []string)
					go func() {
						defer close(channel)
						directionService.GetOrCreateDirections(route.RouteID, channel)
					}()

					return func() (interface{}, error) {
						result := <-channel
						return result, nil
					}, nil
				},
			},
		},
	})

var stopType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Stop",
		Fields: graphql.Fields{
			"StopID": &graphql.Field{
				Type: graphql.String,
			},
			"CommonName": &graphql.Field{
				Type: graphql.String,
			},
			"Lat": &graphql.Field{
				Type: graphql.Float,
			},
			"Lon": &graphql.Field{
				Type: graphql.Float,
			},
			"Predictions": &graphql.Field{
				Type: graphql.NewList(predictionType),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					stop, _ := p.Source.(api.Stop)
					channel := make(chan []api.Prediction)
					go func() {
						defer close(channel)
						predictionService.GetOrCreatePredictionsForStop(stop.StopID, channel)
					}()

					return func() (interface{}, error) {
						result := <-channel
						return result, nil
					}, nil

				},
			},
		},
	})

var predictionType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Prediction",
		Fields: graphql.Fields{
			"StopID": &graphql.Field{
				Type: graphql.String,
			},
			"Type": &graphql.Field{
				Type: graphql.String,
			},
			"StopName": &graphql.Field{
				Type: graphql.String,
			},
			"VehicleID": &graphql.Field{
				Type: graphql.String,
			},
			"Distance": &graphql.Field{
				Type: graphql.Int,
			},
			"RouteID": &graphql.Field{
				Type: graphql.String,
			},
			"Direction": &graphql.Field{
				Type: graphql.String,
			},
			"Destination": &graphql.Field{
				Type: graphql.String,
			},
			"PredictedTime": &graphql.Field{
				Type: graphql.String,
			},
			"Delayed": &graphql.Field{
				Type: graphql.Boolean,
			},
			"TimeLeft": &graphql.Field{
				Type: graphql.String,
			},
		},
	})

// Define the GraphQL Schema
func gqlSchema() graphql.Schema {
	fields := graphql.Fields{
		"routes": &graphql.Field{
			Type: graphql.NewList(routeType),
			Args: nil,
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				return routeService.GetOrCreateRoutes(), nil
			},
			Description: "All routes",
		},
		"route": &graphql.Field{
			Type:        routeType,
			Description: "Get Route by ID",
			Args: graphql.FieldConfigArgument{
				"RouteID": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				id, success := params.Args["RouteID"].(string)
				if !success {
					return nil, errors.New("missing or invalid arguments")
				}
				for _, route := range routeService.GetOrCreateRoutes() {
					if route.RouteID == id {
						return route, nil
					}
				}
				return nil, errors.New("route id not found")
			},
		},
		"stops": &graphql.Field{
			Type:        graphql.NewList(stopType),
			Description: "Get Stops by Route ID and Direction",
			Args: graphql.FieldConfigArgument{
				"RouteID": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"Direction": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				routeID, success := params.Args["RouteID"].(string)
				direction, success2 := params.Args["Direction"].(string)
				if success && success2 {
					stops := stopService.GetOrCreateStops(routeID, direction)
					return stops, nil
				}
				return nil, errors.New("missing or invalid arguments")
			},
		},
		"closestStops": &graphql.Field{
			Type:        graphql.NewList(stopType),
			Description: "Get Stops by coordinates",
			Args: graphql.FieldConfigArgument{
				"Lat": &graphql.ArgumentConfig{
					Type: graphql.Float,
				},
				"Lon": &graphql.ArgumentConfig{
					Type: graphql.Float,
				},
				"Limit": &graphql.ArgumentConfig{
					Type:         graphql.Int,
					DefaultValue: 10,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				lat, success := params.Args["Lat"].(float64)
				lon, success2 := params.Args["Lon"].(float64)
				limit, success3 := params.Args["Limit"].(int)
				if success && success2 && success3 {
					if limit > 10 {
						limit = 10
					}
					if limit <= 0 {
						limit = 1
					}
					stops := stopService.GetClosest(lat, lon, limit)
					return stops, nil
				}
				return nil, errors.New("missing or invalid arguments")
			},
		},
		"predictions": &graphql.Field{
			Type:        graphql.NewList(predictionType),
			Description: "Get predictions by RouteID and StopID.",
			Args: graphql.FieldConfigArgument{
				"RouteID": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"StopID": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				routeID, success := params.Args["RouteID"].(string)
				stopID, success2 := params.Args["StopID"].(string)
				if success && success2 {
					predictions := predictionService.GetOrCreatePredictionsForStopAndRoute(stopID, routeID)
					if len(predictions.Error) > 0 {
						message := ""
						for _, value := range predictions.Error {
							message += value.Message + ". "
						}
						return nil, errors.New(strings.Trim(message, " "))
					}
					return predictions.Predictions, nil
				}
				return nil, errors.New("missing or invalid arguments")
			},
		},
	}
	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: fields}
	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery)}
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		fmt.Printf("failed to create new schema, error: %v", err)
	}

	return schema

}
