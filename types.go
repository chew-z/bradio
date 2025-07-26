package main

import (
	"strconv"
	"strings"

	"gitlab.com/AgentNemo/goradios"
)

// Station represents a radio station with MCP-friendly structure
type Station struct {
	ID            string   `json:"id"`
	Name          string   `json:"name"`
	URL           string   `json:"url"`
	URLResolved   string   `json:"url_resolved,omitempty"`
	Tags          []string `json:"tags"`
	Country       string   `json:"country"`
	CountryCode   string   `json:"country_code,omitempty"`
	State         string   `json:"state"`
	Language      string   `json:"language"`
	Codec         string   `json:"codec"`
	Bitrate       int      `json:"bitrate"`
	ClickCount    int      `json:"click_count"`
	VoteCount     int      `json:"vote_count"`
	ClickTrend    int      `json:"click_trend"`
	Homepage      string   `json:"homepage"`
	Favicon       string   `json:"favicon"`
	HLS           bool     `json:"hls"`
	LastCheckOk   bool     `json:"last_check_ok"`
	LastCheckTime string   `json:"last_check_time,omitempty"`
}

// PopularStation extends Station with ranking information
type PopularStation struct {
	Station
	Rank int `json:"rank"`
}

// MCPSearchResult represents the result of a search operation for MCP tools
type MCPSearchResult struct {
	Query      string    `json:"query"`
	QueryType  string    `json:"query_type"`
	TotalFound int       `json:"total_found"`
	Stations   []Station `json:"stations"`
}

// MCPPopularResult represents the result of a popular stations query for MCP tools
type MCPPopularResult struct {
	Query      string           `json:"query"`
	QueryType  string           `json:"query_type"`
	TotalFound int              `json:"total_found"`
	Stations   []PopularStation `json:"stations"`
}

// MCPConfig holds MCP server configuration
type MCPConfig struct {
	ServerName string
	Version    string
	Transport  string // "stdio", "http", "sse"
	Port       int    // for HTTP transport
}

// ConvertStation converts goradios.Station to our Station type
func ConvertStation(gs goradios.Station) Station {
	// Parse tags from comma-separated string to slice
	var tags []string
	if gs.Tags != "" {
		// Split by comma and trim whitespace
		tagParts := strings.Split(gs.Tags, ",")
		for _, tag := range tagParts {
			trimmed := strings.TrimSpace(tag)
			if trimmed != "" {
				tags = append(tags, trimmed)
			}
		}
	}

	// Parse vote count from string
	voteCount := 0
	if gs.Votes != "" {
		if parsed, err := strconv.Atoi(gs.Votes); err == nil {
			voteCount = parsed
		}
	}

	return Station{
		ID:            gs.StationUUID,
		Name:          gs.Name,
		URL:           gs.URL,
		URLResolved:   gs.URLResolved,
		Tags:          tags,
		Country:       gs.Country,
		CountryCode:   gs.CountryCode,
		State:         gs.State,
		Language:      gs.Language,
		Codec:         gs.Codec,
		Bitrate:       gs.Bitrate,
		ClickCount:    gs.ClickCount,
		VoteCount:     voteCount,
		ClickTrend:    gs.ClickTrend,
		Homepage:      gs.Homepage,
		Favicon:       gs.Favicon,
		HLS:           gs.HLS,
		LastCheckOk:   gs.LastCheckOk,
		LastCheckTime: gs.LastCheckTime,
	}
}

// ConvertStations converts a slice of goradios.Station to our Station type
func ConvertStations(stations []goradios.Station) []Station {
	result := make([]Station, len(stations))
	for i, station := range stations {
		result[i] = ConvertStation(station)
	}
	return result
}

// DefaultMCPConfig returns default MCP server configuration
func DefaultMCPConfig() *MCPConfig {
	return &MCPConfig{
		ServerName: "bradio-mcp-server",
		Version:    "1.0.0",
		Transport:  "stdio",
		Port:       8080,
	}
}
