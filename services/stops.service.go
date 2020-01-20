package services

import (
	"GoCtaApi/api"
	"GoCtaApi/parsers"
	"pault.ag/go/haversine"
	"sort"
)

type byDistance struct {
	data   []api.Stop
	origin haversine.Point
}

func newByDistance(data []api.Stop, lat float64, lon float64) *byDistance {
	return &byDistance{data: data, origin: haversine.Point{Lat: lat, Lon: lon}}
}

func (a byDistance) Len() int { return len(a.data) }
func (a byDistance) Less(i, j int) bool {
	iPoint := haversine.Point{Lat: a.data[i].Lat, Lon: a.data[i].Lon}
	jPoint := haversine.Point{Lat: a.data[j].Lat, Lon: a.data[j].Lon}

	return a.origin.MetresTo(iPoint) < a.origin.MetresTo(jPoint)
}
func (a byDistance) Swap(i, j int) { a.data[i], a.data[j] = a.data[j], a.data[i] }

type stopKey struct {
	routeID   string
	direction string
}

type StopsService struct {
	stops       map[stopKey][]api.Stop
	client      *api.Client
	staticStops []api.Stop
}

func NewStopsService(client *api.Client) *StopsService {
	return &StopsService{
		client:      client,
		stops:       make(map[stopKey][]api.Stop),
		staticStops: parsers.GetStopCoordinates(),
	}
}

func (service *StopsService) GetOrCreateStops(routeID string, direction string) []api.Stop {

	key := stopKey{routeID: routeID, direction: direction}

	if service.stops[key] == nil {
		service.stops[key] = service.client.RetrieveStopsForRoute(routeID, direction)
	}
	return service.stops[key]
}

func (service *StopsService) GetClosest(lat float64, lon float64) *[]api.Stop {
	coordinates := service.staticStops
	sort.Sort(newByDistance(coordinates, lat, lon))
	return &coordinates
}
