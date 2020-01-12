package main

type Database struct {
	routes     []Route
	directions map[string][]string
	client     *ApiClient
}

func NewDatabase(client *ApiClient) *Database {
	return &Database{client: client, directions: make(map[string][]string)}
}

func (f *Database) getOrCreateRoutes() []Route {
	if f.routes == nil {
		f.routes = f.client.RetrieveRoutes()
	}
	return f.routes
}

func (f *Database) getOrCreateDirections(routeId string) []string {
	if f.directions[routeId] == nil {
		directions := f.client.RetrieveDirectionsForRoute(routeId)
		directionsArray := make([]string, len(directions))
		for index, element := range directions {
			directionsArray[index] = element.Value
		}
		f.directions[routeId] = directionsArray
	}
	return f.directions[routeId]
}
