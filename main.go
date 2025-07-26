package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"gitlab.com/AgentNemo/goradios"
)

func main() {
	// Define flags for command-line arguments
	var (
		searchName string
		searchTag  string
		limit      int
		showHelp   bool
	)

	flag.StringVar(&searchName, "name", "", "Search stations by name")
	flag.StringVar(&searchTag, "tag", "", "Search stations by tag")
	flag.IntVar(&limit, "limit", 12, "Limit the number of results (default: 12)")
	flag.BoolVar(&showHelp, "help", false, "Show usage information")

	// Custom usage function
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "A CLI tool to search for internet radio stations via radio-browser.")
		fmt.Fprintln(os.Stderr, "Sorts results by popularity (click count for --name, trend for --tag).\n")
		fmt.Fprintln(os.Stderr, "Options:")
		flag.PrintDefaults()
		fmt.Fprintln(os.Stderr, "\nExamples:")
		fmt.Fprintln(os.Stderr, "  bradio --name 'Milano Lounge'")
		fmt.Fprintln(os.Stderr, "  bradio --tag 'ambient'")
		fmt.Fprintln(os.Stderr, "  bradio --tag 'chillout' --limit 30")
		fmt.Fprintln(os.Stderr, "\nOutput format: (clicks) name; tags; codec[bitrate]; url")
	}

	flag.Parse()

	// Handle help flag
	if showHelp {
		flag.Usage()
		os.Exit(0)
	}

	// Validate input
	if err := validateFlags(searchName, searchTag, limit); err != nil {
		log.Printf("Error: %v\n", err)
		flag.Usage()
		os.Exit(1)
	}

	// Fetch and display stations
	if err := fetchAndDisplayStations(searchName, searchTag, limit); err != nil {
		log.Fatalf("Error: %v", err)
	}
}

// validateFlags checks that the provided flags are valid
func validateFlags(searchName, searchTag string, limit int) error {
	if searchName == "" && searchTag == "" {
		return fmt.Errorf("you must specify either --name or --tag")
	}

	if searchName != "" && searchTag != "" {
		return fmt.Errorf("please specify either --name or --tag, not both")
	}

	if limit <= 0 {
		return fmt.Errorf("limit must be a positive number, got %d", limit)
	}

	if limit > 1000 {
		return fmt.Errorf("limit too high (max 1000), got %d", limit)
	}

	return nil
}

// fetchAndDisplayStations retrieves stations based on search criteria and displays them
func fetchAndDisplayStations(searchName, searchTag string, limit int) error {
	var stations []goradios.Station

	if searchName != "" {
		stations = goradios.FetchStationsDetailed(
			goradios.StationsByName,
			searchName,
			goradios.StationsOrderClickCount,
			true, // reverse (highest first)
			0,    // offset
			uint(limit),
			true, // hideBroken
		)
	} else if searchTag != "" {
		stations = goradios.FetchStationsDetailed(
			goradios.StationsByTagExact,
			searchTag,
			goradios.StationsOrderClickTrend,
			true, // reverse (highest first)
			0,    // offset
			uint(limit),
			true, // hideBroken
		)
	}

	if len(stations) == 0 {
		fmt.Printf("No stations found for search criteria.\n")
		return nil
	}

	// Display results
	for _, station := range stations {
		output := fmt.Sprintf("(%d) %s; %s; %s[%d]; %s",
			station.ClickCount,
			station.Name,
			station.Tags,
			station.Codec,
			station.Bitrate,
			station.URL)
		fmt.Println(output)
	}

	return nil
}
