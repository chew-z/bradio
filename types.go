package main

import "gitlab.com/AgentNemo/goradios"

// Station represents a radio station with MCP-friendly structure
type Station struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	URL        string `json:"url"`
	Tags       string `json:"tags"`
	Country    string `json:"country"`
	Language   string `json:"language"`
	Codec      string `json:"codec"`
	Bitrate    int    `json:"bitrate"`
	ClickCount int    `json:"clickCount"`
	Votes      string `json:"votes"`
}

// SearchParams represents search parameters for MCP tools
type SearchParams struct {
	Query string `json:"query"`
	Limit int    `json:"limit"`
}

// SearchResult represents the result of a search operation
type SearchResult struct {
	Stations []Station `json:"stations"`
	Count    int       `json:"count"`
	Query    string    `json:"query"`
}

// AppMode represents the application's operating mode
type AppMode int

const (
	// ModeCLI represents command-line interface mode
	ModeCLI AppMode = iota
	// ModeMCP represents MCP server mode
	ModeMCP
)

// Config holds application configuration
type Config struct {
	Mode      AppMode
	MCPConfig *MCPConfig
	CLIConfig *CLIConfig
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
	return Station{
		ID:         gs.StationUUID,
		Name:       gs.Name,
		URL:        gs.URL,
		Tags:       gs.Tags,
		Country:    gs.Country,
		Language:   gs.Language,
		Codec:      gs.Codec,
		Bitrate:    gs.Bitrate,
		ClickCount: gs.ClickCount,
		Votes:      gs.Votes,
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

// DefaultCLIConfig returns default CLI configuration
func DefaultCLIConfig() *CLIConfig {
	return &CLIConfig{
		Limit: 12,
	}
}
