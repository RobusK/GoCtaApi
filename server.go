package main

import (
	"GoCtaApi/api"
	"GoCtaApi/services"
	"encoding/json"
	"fmt"
	"github.com/friendsofgo/graphiql"
	"github.com/graphql-go/graphql"
	"net/http"
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
						directionsService.GetOrCreateDirections(route.RouteID, channel)
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

type reqBody struct {
	Query string `json:"query"`
}

func main() {
	graphiqlHandler, err := graphiql.NewGraphiqlHandler("/graphql")
	schema := gqlSchema()
	if err != nil {
		panic(err)
	}

	http.Handle("/graphql", gqlHandler(&schema))
	http.Handle("/graphiql", graphiqlHandler)
	http.ListenAndServe(":3000", nil)
}

var (
	client            = api.NewAPIClient(CtaAPIKey)
	directionsService = services.NewDirectionsService(client)
	routesService     = services.NewRoutesService(client)
	stopService       = services.NewStopsService(client)
)

func gqlHandler(schema *graphql.Schema) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Body == nil {
			http.Error(w, "No query data", 400)
			return
		}

		var rBody reqBody
		err := json.NewDecoder(r.Body).Decode(&rBody)
		if err != nil {
			http.Error(w, "Error parsing JSON request body", 400)
		}

		fmt.Fprintf(w, "%s", processQuery(rBody.Query, schema))

	})
}

func processQuery(query string, schema *graphql.Schema) (result string) {

	params := graphql.Params{Schema: *schema, RequestString: query}
	r := graphql.Do(params)
	if len(r.Errors) > 0 {
		fmt.Printf("failed to execute graphql operation, errors: %+v", r.Errors)
	}
	rJSON, _ := json.Marshal(r)

	return fmt.Sprintf("%s", rJSON)

}

// Define the GraphQL Schema
func gqlSchema() graphql.Schema {
	fields := graphql.Fields{
		"routes": &graphql.Field{
			Name: "",
			Type: graphql.NewList(routeType),
			Args: nil,
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				return routesService.GetOrCreateRoutes(), nil
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
					for _, route := range routesService.GetOrCreateRoutes() {
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
