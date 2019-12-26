package main

import (
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
			"RouteId": &graphql.Field{
				Type: graphql.String,
			},
			"CommonName": &graphql.Field{
				Type: graphql.String,
			},
			"Color": &graphql.Field{
				Type: graphql.String,
			},
		},
	})

type reqBody struct {
	Query string `json:"query"`
}

func main() {
	graphiqlHandler, err := graphiql.NewGraphiqlHandler("/graphql")
	if err != nil {
		panic(err)
	}

	http.Handle("/graphql", gqlHandler())
	http.Handle("/graphiql", graphiqlHandler)
	http.ListenAndServe(":3000", nil)
}

func gqlHandler() http.Handler {
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

		fmt.Fprintf(w, "%s", processQuery(rBody.Query))

	})
}

func processQuery(query string) (result string) {

	params := graphql.Params{Schema: gqlSchema(), RequestString: query}
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
				return RetrieveStops().RoutesList.Routes, nil
			},
			DeprecationReason: "",
			Description:       "All routes",
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
