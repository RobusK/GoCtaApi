package services

import (
	"github.com/RobusK/GoCtaApi/api"
	"github.com/RobusK/GoCtaApi/parsers"
	"github.com/wangjohn/quickselect"
	"pault.ag/go/haversine"
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

func (service *StopsService) GetClosest(lat float64, lon float64, k int) []api.Stop {
	stops := make([]api.Stop, len(service.staticStops))
	copy(stops, service.staticStops)

	data := newByDistance(stops, lat, lon)

	quickselect.QuickSelect(data, k)
	insertionSort(data, 0, k)
	return stops[:k]
}

func insertionSort(data quickselect.Interface, a, b int) {
	for i := a + 1; i < b; i++ {
		for j := i; j > a && data.Less(j, j-1); j-- {
			data.Swap(j, j-1)
		}
	}
}
