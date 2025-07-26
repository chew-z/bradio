package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"gitlab.com/AgentNemo/goradios"
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
	args := request.GetArguments()

	name, ok := args["name"].(string)
	if !ok || name == "" {
		return mcp.NewToolResultError("name parameter is required and must be a string"), nil
	}

	limit, err := m.extractAndValidateLimit(args)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	stations, err := m.radioService.SearchByName(name, limit)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to search stations: %v", err)), nil
	}

	return m.createSearchResultResponse(name, "name", stations)
}

// handleSearchByTag handles the search_radio_by_tag tool
func (m *MCPServer) handleSearchByTag(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()

	tag, ok := args["tag"].(string)
	if !ok || tag == "" {
		return mcp.NewToolResultError("tag parameter is required and must be a string"), nil
	}

	limit, err := m.extractAndValidateLimit(args)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	stations, err := m.radioService.SearchByTag(tag, limit)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to search stations: %v", err)), nil
	}

	return m.createSearchResultResponse(tag, "tag", stations)
}

// handleGetPopularStations handles the get_popular_stations tool
func (m *MCPServer) handleGetPopularStations(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()

	limit, err := m.extractAndValidateLimit(args)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	stations, err := m.radioService.GetPopularStations(limit)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get popular stations: %v", err)), nil
	}

	return m.createPopularResultResponse(stations)
}

// extractAndValidateLimit extracts and validates the limit parameter from MCP tool arguments
func (m *MCPServer) extractAndValidateLimit(args map[string]any) (int, error) {
	limit := 12 // default
	if limitVal, exists := args["limit"]; exists {
		if limitFloat, ok := limitVal.(float64); ok {
			limit = int(limitFloat)
		}
	}

	if limit <= 0 || limit > 1000 {
		return 0, fmt.Errorf("limit must be between 1 and 1000")
	}

	return limit, nil
}

// createSearchResultResponse creates a JSON response for search operations
func (m *MCPServer) createSearchResultResponse(query, queryType string, stations []goradios.Station) (*mcp.CallToolResult, error) {
	result := MCPSearchResult{
		Query:      query,
		QueryType:  queryType,
		TotalFound: len(stations),
		Stations:   ConvertStations(stations),
	}

	jsonData, err := json.Marshal(result)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal result: %v", err)), nil
	}

	return mcp.NewToolResultText(string(jsonData)), nil
}

// createPopularResultResponse creates a JSON response for popular stations
func (m *MCPServer) createPopularResultResponse(stations []goradios.Station) (*mcp.CallToolResult, error) {
	convertedStations := ConvertStations(stations)
	popularStations := make([]PopularStation, len(convertedStations))
	for i, station := range convertedStations {
		popularStations[i] = PopularStation{
			Station: station,
			Rank:    i + 1,
		}
	}

	result := MCPPopularResult{
		Query:      "popular_stations",
		QueryType:  "popular",
		TotalFound: len(stations),
		Stations:   popularStations,
	}

	jsonData, err := json.Marshal(result)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal result: %v", err)), nil
	}

	return mcp.NewToolResultText(string(jsonData)), nil
}
