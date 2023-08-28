package main

import (
	"fmt"
	"os"
	"strconv"

	"gitlab.com/AgentNemo/goradios"
)

func main() {
	var stations []goradios.Station
	limit := 12
	if len(os.Args) > 4 && os.Args[3] == "--limit" {
		limit, _ = strconv.Atoi(os.Args[4])
	}
	if os.Args[1] == "--name" {
		stations = goradios.FetchStationsDetailed(goradios.StationsByName, os.Args[2], goradios.StationsOrderClickCount, true, uint(0), uint(limit), true)
	} else if os.Args[1] == "--tag" {
		stations = goradios.FetchStationsDetailed(goradios.StationsByTagExact, os.Args[2], goradios.StationsOrderClickTrend, true, uint(0), uint(limit), true)
	} else {
		stations = goradios.FetchStations(goradios.StationsByName, os.Args[1])
	}
	for _, station := range stations {
		s := fmt.Sprintf("(%d) %s; %s; %s[%d]; %s", station.ClickCount, station.Name, station.Tags, station.Codec, station.Bitrate, station.URL)
		fmt.Println(s)
	}
}
