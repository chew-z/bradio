package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"gitlab.com/AgentNemo/goradios"
)

// CLIConfig holds CLI configuration
type CLIConfig struct {
	SearchName string
	SearchTag  string
	Limit      int
	ShowHelp   bool
}

// RunCLI runs the command-line interface
func RunCLI() error {
	// Reset flag package for CLI use
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	config := parseCLIFlags()

	if config.ShowHelp {
		showUsage()
		os.Exit(0)
	}

	// Validate input
	if err := validateCLIFlags(config); err != nil {
		log.Printf("Error: %v\n", err)
		showUsage()
		os.Exit(1)
	}

	// Create radio service and fetch stations
	radioService := NewRadioService()
	return fetchAndDisplayStations(radioService, config)
}

// parseCLIFlags parses command-line flags
func parseCLIFlags() *CLIConfig {
	config := &CLIConfig{}

	flag.StringVar(&config.SearchName, "name", "", "Search stations by name")
	flag.StringVar(&config.SearchTag, "tag", "", "Search stations by tag")
	flag.IntVar(&config.Limit, "limit", 12, "Limit the number of results (default: 12)")
	flag.BoolVar(&config.ShowHelp, "help", false, "Show usage information")

	// Custom usage function
	flag.Usage = showUsage

	flag.Parse()
	return config
}

// showUsage displays usage information
func showUsage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\n", os.Args[0])
	fmt.Fprintln(os.Stderr, "A CLI tool to search for internet radio stations via radio-browser.")
	fmt.Fprintln(os.Stderr, "Sorts results by popularity (click count for --name, trend for --tag).")
	fmt.Fprintln(os.Stderr, "Options:")
	flag.PrintDefaults()
	fmt.Fprintln(os.Stderr, "\nExamples:")
	fmt.Fprintln(os.Stderr, "  bradio --name 'Milano Lounge'")
	fmt.Fprintln(os.Stderr, "  bradio --tag 'ambient'")
	fmt.Fprintln(os.Stderr, "  bradio --tag 'chillout' --limit 30")
	fmt.Fprintln(os.Stderr, "\nOutput format: (clicks) name; tags; codec[bitrate]; url")
}

// validateCLIFlags validates the provided CLI flags
func validateCLIFlags(config *CLIConfig) error {
	if config.SearchName == "" && config.SearchTag == "" {
		return fmt.Errorf("you must specify either --name or --tag")
	}

	if config.SearchName != "" && config.SearchTag != "" {
		return fmt.Errorf("please specify either --name or --tag, not both")
	}

	if config.Limit <= 0 {
		return fmt.Errorf("limit must be a positive number, got %d", config.Limit)
	}

	if config.Limit > 1000 {
		return fmt.Errorf("limit too high (max 1000), got %d", config.Limit)
	}

	return nil
}

// fetchAndDisplayStations retrieves stations and displays them
func fetchAndDisplayStations(radioService *RadioService, config *CLIConfig) error {
	var stations []goradios.Station
	var err error

	if config.SearchName != "" {
		stations, err = radioService.SearchByName(config.SearchName, config.Limit)
	} else if config.SearchTag != "" {
		stations, err = radioService.SearchByTag(config.SearchTag, config.Limit)
	}

	if err != nil {
		return fmt.Errorf("failed to fetch stations: %v", err)
	}

	if len(stations) == 0 {
		fmt.Printf("No stations found for search criteria.\n")
		return nil
	}

	// Display results
	for _, station := range stations {
		fmt.Println(radioService.FormatStationOutput(station))
	}

	return nil
}
