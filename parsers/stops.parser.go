package parsers

import (
	"GoCtaApi/api"
	"bufio"
	"encoding/csv"
	"io"
	"log"
	"os"
	"strconv"
)

func GetStopCoordinates() []api.Stop {
	csvFile, _ := os.Open("./data/stops_with_bus_service.csv")
	reader := csv.NewReader(bufio.NewReader(csvFile))
	var stops []api.Stop
	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}
		if line[6] != "0" {
			continue
		}
		lat, _ := strconv.ParseFloat(line[4], 64)
		lon, _ := strconv.ParseFloat(line[5], 64)
		stops = append(stops, api.Stop{
			StopID:     line[0],
			CommonName: line[2],
			Lat:        lat,
			Lon:        lon,
		})
	}
	return stops
}
