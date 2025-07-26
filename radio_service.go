package main

import (
	"fmt"

	"gitlab.com/AgentNemo/goradios"
)

// RadioService handles all radio station operations
type RadioService struct{}

// NewRadioService creates a new radio service instance
func NewRadioService() *RadioService {
	return &RadioService{}
}

// SearchByName searches radio stations by name
func (rs *RadioService) SearchByName(name string, limit int) ([]goradios.Station, error) {
	if name == "" {
		return nil, fmt.Errorf("name cannot be empty")
	}
	if limit <= 0 {
		return nil, fmt.Errorf("limit must be positive, got %d", limit)
	}
	if limit > 1000 {
		return nil, fmt.Errorf("limit too high (max 1000), got %d", limit)
	}

	stations := goradios.FetchStationsDetailed(
		goradios.StationsByName,
		name,
		goradios.StationsOrderClickCount,
		true, // reverse (highest first)
		0,    // offset
		uint(limit),
		true, // hideBroken
	)

	return stations, nil
}

// SearchByTag searches radio stations by tag
func (rs *RadioService) SearchByTag(tag string, limit int) ([]goradios.Station, error) {
	if tag == "" {
		return nil, fmt.Errorf("tag cannot be empty")
	}
	if limit <= 0 {
		return nil, fmt.Errorf("limit must be positive, got %d", limit)
	}
	if limit > 1000 {
		return nil, fmt.Errorf("limit too high (max 1000), got %d", limit)
	}

	stations := goradios.FetchStationsDetailed(
		goradios.StationsByTagExact,
		tag,
		goradios.StationsOrderClickTrend,
		true, // reverse (highest first)
		0,    // offset
		uint(limit),
		true, // hideBroken
	)

	return stations, nil
}

// GetPopularStations retrieves the most popular stations globally
func (rs *RadioService) GetPopularStations(limit int) ([]goradios.Station, error) {
	if limit <= 0 {
		return nil, fmt.Errorf("limit must be positive, got %d", limit)
	}
	if limit > 1000 {
		return nil, fmt.Errorf("limit too high (max 1000), got %d", limit)
	}

	stations := goradios.FetchAllStationsDetailed(
		goradios.StationsOrderClickCount,
		true, // reverse (highest first)
		0,    // offset
		uint(limit),
		true, // hideBroken
	)

	return stations, nil
}

// FormatStationOutput formats a station for CLI output
func (rs *RadioService) FormatStationOutput(station goradios.Station) string {
	return fmt.Sprintf("(%d) %s; %s; %s[%d]; %s",
		station.ClickCount,
		station.Name,
		station.Tags,
		station.Codec,
		station.Bitrate,
		station.URL)
}
