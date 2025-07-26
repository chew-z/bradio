package main

import (
	"context"
	"fmt"
	"log"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// MCPServer wraps the MCP server functionality
type MCPServer struct {
	server       *server.MCPServer
	radioService *RadioService
}

// NewMCPServer creates a new MCP server instance
func NewMCPServer() *MCPServer {
	config := DefaultMCPConfig()

	mcpServer := server.NewMCPServer(config.ServerName, config.Version)
	radioService := NewRadioService()

	return &MCPServer{
		server:       mcpServer,
		radioService: radioService,
	}
}

// RunMCPServer starts the MCP server
func RunMCPServer() error {
	log.Println("Starting bradio MCP server...")

	mcpServer := NewMCPServer()

	// Initialize the server with tools
	mcpServer.setupTools()

	log.Println("MCP server ready. Listening on stdio...")

	// Start the server with stdio transport
	return server.ServeStdio(mcpServer.server)
}

// setupTools registers all MCP tools
func (m *MCPServer) setupTools() {
	log.Println("Setting up MCP tools...")

	// Add search_radio_by_name tool
	m.server.AddTool(mcp.Tool{
		Name:        "search_radio_by_name",
		Description: "Search for radio stations by name, sorted by popularity (click count)",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]any{
				"name": map[string]any{
					"type":        "string",
					"description": "The name or partial name of the radio station to search for",
				},
				"limit": map[string]any{
					"type":        "integer",
					"description": "Maximum number of results to return (default: 12, max: 1000)",
					"default":     12,
				},
			},
			Required: []string{"name"},
		},
	}, m.handleSearchByName)

	log.Println("Registered tool: search_radio_by_name")

	// Add search_radio_by_tag tool
	m.server.AddTool(mcp.Tool{
		Name:        "search_radio_by_tag",
		Description: "Search for radio stations by tag, sorted by popularity (click trend)",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]any{
				"tag": map[string]any{
					"type":        "string",
					"description": "The tag to search for (e.g., 'jazz', 'rock', 'electronic')",
				},
				"limit": map[string]any{
					"type":        "integer",
					"description": "Maximum number of results to return (default: 12, max: 1000)",
					"default":     12,
				},
			},
			Required: []string{"tag"},
		},
	}, m.handleSearchByTag)

	log.Println("Registered tool: search_radio_by_tag")

	// Add get_popular_stations tool
	m.server.AddTool(mcp.Tool{
		Name:        "get_popular_stations",
		Description: "Get the most popular radio stations globally, sorted by click count",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]any{
				"limit": map[string]any{
					"type":        "integer",
					"description": "Maximum number of results to return (default: 12, max: 1000)",
					"default":     12,
				},
			},
			Required: []string{},
		},
	}, m.handleGetPopularStations)

	log.Println("Registered tool: get_popular_stations")
}

// handleSearchByName handles the search_radio_by_name tool
func (m *MCPServer) handleSearchByName(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Extract parameters
	args := request.GetArguments()

	name, ok := args["name"].(string)
	if !ok || name == "" {
		return mcp.NewToolResultError("name parameter is required and must be a string"), nil
	}

	// Handle limit parameter
	limit := 12 // default
	if limitVal, exists := args["limit"]; exists {
		if limitFloat, ok := limitVal.(float64); ok {
			limit = int(limitFloat)
		}
	}

	// Validate limit
	if limit <= 0 || limit > 1000 {
		return mcp.NewToolResultError("limit must be between 1 and 1000"), nil
	}

	// Search for stations
	stations, err := m.radioService.SearchByName(name, limit)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to search stations: %v", err)), nil
	}

	// Convert to our Station type
	result := SearchResult{
		Stations: ConvertStations(stations),
		Count:    len(stations),
		Query:    name,
	}

	return mcp.NewToolResultText(fmt.Sprintf("Found %d stations matching '%s'", result.Count, name)), nil
}

// handleSearchByTag handles the search_radio_by_tag tool
func (m *MCPServer) handleSearchByTag(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Extract parameters
	args := request.GetArguments()

	tag, ok := args["tag"].(string)
	if !ok || tag == "" {
		return mcp.NewToolResultError("tag parameter is required and must be a string"), nil
	}

	// Handle limit parameter
	limit := 12 // default
	if limitVal, exists := args["limit"]; exists {
		if limitFloat, ok := limitVal.(float64); ok {
			limit = int(limitFloat)
		}
	}

	// Validate limit
	if limit <= 0 || limit > 1000 {
		return mcp.NewToolResultError("limit must be between 1 and 1000"), nil
	}

	// Search for stations
	stations, err := m.radioService.SearchByTag(tag, limit)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to search stations: %v", err)), nil
	}

	// Convert to our Station type and format as text
	convertedStations := ConvertStations(stations)

	// Create a formatted text response
	var response string
	if len(convertedStations) == 0 {
		response = fmt.Sprintf("No stations found with tag '%s'", tag)
	} else {
		response = fmt.Sprintf("Found %d stations with tag '%s':\n\n", len(convertedStations), tag)
		for i, station := range convertedStations {
			response += fmt.Sprintf("%d. %s (%d clicks)\n   Tags: %s\n   Codec: %s [%d kbps]\n   URL: %s\n\n",
				i+1, station.Name, station.ClickCount, station.Tags, station.Codec, station.Bitrate, station.URL)
		}
	}

	return mcp.NewToolResultText(response), nil
}

// handleGetPopularStations handles the get_popular_stations tool
func (m *MCPServer) handleGetPopularStations(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Extract parameters
	args := request.GetArguments()

	// Handle limit parameter
	limit := 12 // default
	if limitVal, exists := args["limit"]; exists {
		if limitFloat, ok := limitVal.(float64); ok {
			limit = int(limitFloat)
		}
	}

	// Validate limit
	if limit <= 0 || limit > 1000 {
		return mcp.NewToolResultError("limit must be between 1 and 1000"), nil
	}

	// Get popular stations
	stations, err := m.radioService.GetPopularStations(limit)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get popular stations: %v", err)), nil
	}

	// Convert to our Station type and format as text
	convertedStations := ConvertStations(stations)

	// Create a formatted text response
	var response string
	if len(convertedStations) == 0 {
		response = "No popular stations found"
	} else {
		response = fmt.Sprintf("Top %d most popular radio stations:\n\n", len(convertedStations))
		for i, station := range convertedStations {
			response += fmt.Sprintf("%d. %s (%d clicks)\n   Tags: %s\n   Codec: %s [%d kbps]\n   URL: %s\n\n",
				i+1, station.Name, station.ClickCount, station.Tags, station.Codec, station.Bitrate, station.URL)
		}
	}

	return mcp.NewToolResultText(response), nil
}
