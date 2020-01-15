package main

import (
	"github.com/friendsofgo/graphiql"
	"net/http"
)

func main() {
	graphiqlHandler, err := graphiql.NewGraphiqlHandler("/graphql")
	schema := gqlSchema()
	if err != nil {
		panic(err)
	}

	http.Handle("/graphql", gqlHandler(&schema))
	http.Handle("/graphiql", graphiqlHandler)
	http.ListenAndServe(serverHostname+":"+serverPort, nil)
}
