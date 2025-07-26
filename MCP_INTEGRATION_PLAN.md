# MCP Server Integration Plan for Bradio

## Research Summary

### MCP-Go v0.35.0 Key Features (Latest Version)
- **HTTP headers support** - Enhanced HTTP transport capabilities
- **Race condition fixes** - Improved stability for concurrent operations  
- **Protocol version negotiation** - Better compatibility handling
- **Enhanced session management** - Client capabilities embedded into sessions

### Core MCP-Go Capabilities
- **Multiple transports**: Stdio, HTTP, Server-Sent Events (SSE)
- **Three capability types**: 
  - **Tools** (function execution)
  - **Resources** (data exposure) 
  - **Prompts** (interaction templates)
- **Extensible hooks system** for request lifecycle customization
- **Session management** with client capability tracking

---

## Phase 1: Basic MCP Server Structure
### **Goal**: Create minimal MCP server that wraps existing CLI functionality

**Tasks:**
1. **Add MCP dependencies** to `go.mod`
   - `github.com/mark3labs/mcp-go/server v0.35.0`
   - `github.com/mark3labs/mcp-go/mcp v0.35.0`

2. **Create MCP server wrapper**
   - New `mcp_server.go` file 
   - Initialize MCP server with name/version
   - Choose transport (start with Stdio for development)

3. **Refactor existing code**
   - Extract radio search logic into reusable functions
   - Separate CLI handler from core business logic

## Phase 2: Tool Implementation  
### **Goal**: Expose radio search as MCP tools

**Tools to implement:**
1. **`search_radio_by_name`**
   - Parameters: `name` (string), `limit` (optional int)
   - Returns: Array of station objects

2. **`search_radio_by_tag`** 
   - Parameters: `tag` (string), `limit` (optional int)
   - Returns: Array of station objects

3. **`get_popular_stations`**
   - Parameters: `limit` (optional int)
   - Returns: Most popular stations globally

## Phase 3: Resource Implementation
### **Goal**: Expose radio data as queryable resources

**Resources to implement:**
1. **`station://{id}`** - Individual station details
2. **`stations://popular`** - Popular stations list
3. **`stations://recent`** - Recently added stations
4. **`tags://list`** - Available tags for searching

## Phase 4: Advanced Features
### **Goal**: Enhanced functionality and multiple transports

**Enhancements:**
1. **Caching layer** - Cache API responses for better performance
2. **HTTP transport** - Enable web-based access
3. **Station favorites** - Save/retrieve favorite stations
4. **Playlist generation** - Create playlists from search results

## Phase 5: Integration & Testing
### **Goal**: Seamless dual-mode operation

**Features:**
1. **Dual mode support**
   - CLI mode (existing functionality)
   - MCP server mode (new `--mcp` flag)

2. **Configuration**
   - Transport selection (stdio/http/sse)
   - Port configuration for HTTP
   - API rate limiting

3. **Testing & Documentation**
   - Unit tests for MCP tools
   - Integration tests with MCP clients
   - Usage documentation

---

## Architecture Overview

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   CLI Mode      │    │   Core Logic     │    │  MCP Server     │
│                 │    │                  │    │     Mode        │
│ • Flag parsing  │───▶│ • Radio search   │◀───│ • Tools         │
│ • Direct output │    │ • Data formatting│    │ • Resources     │
│ • Error handling│    │ • API calls      │    │ • Transport     │
└─────────────────┘    └──────────────────┘    └─────────────────┘
                              │
                              ▼
                    ┌──────────────────┐
                    │ Radio-Browser    │
                    │      API         │
                    └──────────────────┘
```

## File Structure Plan
```
bradio/
├── main.go              # Entry point with mode selection
├── cli.go               # CLI-specific logic  
├── mcp_server.go        # MCP server implementation
├── radio_service.go     # Core radio search logic
├── types.go             # Shared data structures
├── go.mod
├── go.sum
└── README.md           # Updated with MCP usage
```

## Implementation Strategy

### Key Strategy
- **Dual-mode operation**: Maintain existing CLI functionality while adding MCP server mode
- **Incremental development**: Start simple with stdio transport, expand to HTTP/SSE
- **Modular architecture**: Separate CLI, MCP, and core radio logic for maintainability

### MCP Integration Benefits
1. **AI Assistant Integration**: Enable Claude, GPT, and other LLMs to search radio stations
2. **Programmable Interface**: Structured tool calls instead of text parsing
3. **Resource Exposure**: Queryable station data, tags, and popular lists
4. **Multiple Transports**: Stdio, HTTP, and SSE for different use cases

### Backward Compatibility
- All existing CLI functionality preserved
- Same output format for shell pipelines
- No breaking changes to current usage patterns

## Development Notes

### MCP-Go v0.35.0 Specific Features to Leverage
- Use HTTP headers support for enhanced web integration
- Implement proper session management with embedded client capabilities  
- Utilize protocol version negotiation for better compatibility
- Benefit from race condition fixes for stable concurrent operations

### Testing Strategy
- Unit tests for individual MCP tools
- Integration tests with real MCP clients
- CLI regression tests to ensure no breaking changes
- Performance tests for API rate limiting

### Documentation Updates
- Update README.md with MCP usage examples
- Add MCP client configuration examples
- Document tool and resource schemas
- Provide troubleshooting guide

---

*Plan created: July 26, 2025*
*Target MCP-Go version: v0.35.0*
*Estimated implementation time: 2-3 weeks*