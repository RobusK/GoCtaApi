package main

import (
	"GoCtaApi/api"
	"fmt"
	"github.com/graphql-go/graphql"
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
		},
	})

// Define the GraphQL Schema
func gqlSchema() graphql.Schema {
	fields := graphql.Fields{
		"routes": &graphql.Field{
			Name: "",
			Type: graphql.NewList(routeType),
			Args: nil,
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				return routeService.GetOrCreateRoutes(), nil
			},
			DeprecationReason: "",
			Description:       "All routes",
		},
		"route": &graphql.Field{
			Type:        routeType,
			Description: "Get Route by ID",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				id, success := params.Args["id"].(string)
				if success {
					for _, route := range routeService.GetOrCreateRoutes() {
						if route.RouteID == id {
							return route, nil
						}
					}
				}
				return nil, nil
			},
		},
		"stops": &graphql.Field{
			Type:        graphql.NewList(stopType),
			Description: "Get Stops by Route ID and Direction",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"direction": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				id, success := params.Args["id"].(string)
				direction, success2 := params.Args["direction"].(string)
				if success && success2 {
					stops := stopService.GetOrCreateStops(id, direction)
					return stops, nil
				}
				return nil, nil
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
