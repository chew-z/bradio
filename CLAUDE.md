# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview
bradio is a dual-mode application for searching radio stations from radio-browser.info. It operates as both a traditional CLI tool and an MCP (Model Context Protocol) server for AI assistant integration. The application uses the `gitlab.com/AgentNemo/goradios` library to interact with the radio-browser API and provides structured JSON responses for programmatic consumption.

## Development Commands

### Building
```bash
go build -o bin/bradio .
```

### Testing
```bash
./run_test.sh
# or directly:
go test -v ./...
```

### Linting
```bash
./run_lint.sh
# Uses golangci-lint with --fix flag
```

### Formatting
```bash
./run_format.sh
# Uses gofmt to format all Go files
```

## Architecture
This is a modular Go application with dual-mode operation and clean separation of concerns:

### File Structure
- `main.go` - Entry point implementing dual-mode operation logic with early argument detection
- `cli.go` - CLI-specific functionality with flag parsing and MCP flag filtering
- `mcp_server.go` - Complete MCP server implementation with three tools and JSON responses
- `radio_service.go` - Core radio search logic extracted for reuse between CLI and MCP modes
- `types.go` - Shared data structures for Station, SearchParams, SearchResult, and MCP response types

### Key Features
- **Dual-mode operation**: Automatic detection between CLI mode (default) and MCP server mode (`--mcp` flag)
- **Flag filtering**: Robust handling prevents CLI crashes when MCP-specific flags are present ([commit 5cf3fdb](https://github.com/user/repo/commit/5cf3fdb))
- **Structured JSON responses**: MCP tools return properly formatted JSON instead of text ([commit b318651](https://github.com/user/repo/commit/b318651))
- **Comprehensive error handling**: Input validation and user-friendly error messages
- **Shell pipeline compatibility**: CLI output format designed for fzf + mpv workflows

## Key Dependencies
- Go 1.24.4 (specified in go.mod)
- `gitlab.com/AgentNemo/goradios` for radio station API access
- `github.com/mark3labs/mcp-go` v0.35.0 for Model Context Protocol server implementation

## CLI Usage
The tool supports the following flags:
- `--name "string"`: Search stations by name
- `--tag "string"`: Search stations by tag  
- `--limit N`: Limit number of results (default: 12, max: 1000)
- `--help`: Show usage information

Examples:
```bash
bradio --name "Milano Lounge"
bradio --tag "ambient" 
bradio --tag "chillout" --limit 30
bradio --help
```

## MCP Server Mode
The tool can run as an MCP (Model Context Protocol) server to integrate with AI assistants:

### Starting MCP Server
```bash
bradio --mcp                    # Start with stdio transport
```

### Available MCP Tools
1. **`search_radio_by_name`**
   - Parameters: `name` (string), `limit` (optional int, default: 12, max: 1000)
   - Searches stations by name, sorted by click count
   - Returns: `MCPSearchResult` with structured JSON

2. **`search_radio_by_tag`**
   - Parameters: `tag` (string), `limit` (optional int, default: 12, max: 1000)  
   - Searches stations by tag, sorted by click trend
   - Returns: `MCPSearchResult` with structured JSON

3. **`get_popular_stations`**
   - Parameters: `limit` (optional int, default: 12, max: 1000)
   - Returns most popular stations globally
   - Returns: `MCPPopularResult` with ranking information

### JSON Response Format
All MCP tools return structured JSON responses instead of formatted text:

**Search Tools Response (`MCPSearchResult`):**
```json
{
  "query": "jazz",
  "query_type": "tag",
  "total_found": 15,
  "stations": [
    {
      "id": "station_uuid",
      "name": "101 SMOOTH JAZZ",
      "url": "http://www.101smoothjazz.com/101-smoothjazz.m3u",
      "url_resolved": "http://resolved-stream-url",
      "tags": ["easy listening", "jazz", "smooth jazz"],
      "country": "United States",
      "country_code": "US",
      "state": "California",
      "language": "english",
      "codec": "MP3",
      "bitrate": 128,
      "click_count": 19832,
      "vote_count": 45,
      "click_trend": 150,
      "homepage": "http://www.101smoothjazz.com",
      "favicon": "http://www.101smoothjazz.com/favicon.ico",
      "hls": false,
      "last_check_ok": true,
      "last_check_time": "2025-07-26T10:30:00Z"
    }
  ]
}
```

**Popular Stations Response (`MCPPopularResult`):**
```json
{
  "query": "popular_stations",
  "query_type": "popular",
  "total_found": 12,
  "stations": [
    {
      "rank": 1,
      "id": "station_uuid",
      "name": "Classic Vinyl HD",
      "url": "https://icecast.walmradio.com:8443/classic",
      "click_count": 28593,
      "vote_count": 234,
      "tags": ["1930", "1940", "1950", "beautiful music", "big band"],
      "country": "United States",
      "codec": "MP3",
      "bitrate": 320
      // ... additional station fields
    }
  ]
}
```

### MCP Integration
The MCP server enables AI assistants to:
- **Programmatic Search**: Search radio stations by name or tag with full parameter validation
- **Structured Data**: Access comprehensive station metadata including geolocation, codec details, and popularity metrics
- **Ranking Information**: Get popular stations with explicit ranking for playlist generation
- **Rich Metadata**: Access resolved URLs, homepage links, favicon URLs, and health check status
- **Conversational Workflows**: Integrate radio discovery into smart home systems and automation workflows
- **JSON-First Design**: All responses are structured JSON for easy parsing and integration

## Usage Patterns
The tool is designed for integration with:

**Shell workflows:**
- Piping output to `fzf` for interactive selection
- Extracting URLs for media players like `mpv`
- The output format supports both human reading and programmatic parsing

**AI Assistant workflows:**
- **LLM Integration**: Claude, GPT, and other LLMs can use MCP tools to search stations programmatically
- **Structured Data Processing**: JSON responses enable automated playlist generation and data analysis
- **Smart Home Integration**: Direct integration with home automation systems through MCP protocol
- **Enhanced Metadata**: Access to comprehensive station information including health status, codecs, and geographic data
- **Ranking and Popularity**: Explicit ranking information for creating curated station lists

## Recent Improvements

### JSON Response Standardization ([commit b318651](https://github.com/user/repo/commit/b318651))
- **Enhanced Data Structures**: Introduced `MCPSearchResult`, `MCPPopularResult`, and `PopularStation` types
- **Structured JSON Responses**: All MCP tools now return properly formatted JSON instead of formatted text
- **Rich Station Metadata**: Added comprehensive fields including `url_resolved`, `country_code`, `click_trend`, `hls`, `last_check_ok`, and `last_check_time`
- **Tags as Arrays**: Convert comma-separated tag strings to proper string arrays for easier programmatic access

### Flag Filtering Fix ([commit 5cf3fdb](https://github.com/user/repo/commit/5cf3fdb))
- **Robust CLI Parsing**: Implemented `filterMCPFlags()` function to prevent CLI crashes when MCP-specific flags are present
- **Dual-Mode Compatibility**: Ensures clean separation between CLI and MCP mode argument processing
- **Error Prevention**: Eliminates "flag provided but not defined" errors when switching between modes

### Modular Architecture Refactoring ([commit b107940](https://github.com/user/repo/commit/b107940))
- **Code Organization**: Split monolithic `main.go` into focused modules (`cli.go`, `mcp_server.go`, `radio_service.go`, `types.go`)
- **Reusable Components**: Extracted core radio search logic for use in both CLI and MCP modes
- **Clean Separation**: Each module has clear responsibilities and well-defined interfaces
- **Helper Functions**: Added `extractAndValidateLimit()`, `createSearchResultResponse()`, and `createPopularResultResponse()` to eliminate code duplication

## Development Status

### Current State
- ✅ **Core MCP Functionality**: Complete with three fully functional tools
- ✅ **Structured JSON Responses**: All MCP tools return properly formatted JSON
- ✅ **Dual-Mode Operation**: Robust CLI and MCP server modes with flag filtering
- ✅ **Zero Linting Issues**: Clean codebase with comprehensive error handling
- ✅ **Modular Architecture**: Well-organized code structure with clear separation of concerns

### Pending Enhancements (See MCP_INTEGRATION_PLAN.md)
- 🔄 **Unit Testing**: Add comprehensive unit tests for all MCP tools (HIGH priority)
- 🔄 **HTTP Transport**: Add HTTP transport support for web-based MCP clients (MEDIUM priority)
- 🔄 **MCP Resources**: Implement MCP resources for browsable data access (LOW priority)

The core MCP functionality is production-ready and provides a robust foundation for radio station discovery and integration workflows.