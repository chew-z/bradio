# 📻 Bradio

A powerful dual-mode CLI tool and MCP server for discovering internet radio stations. Built with Go and powered by the fantastic [Radio Browser](https://www.radio-browser.info/) database.

## ✨ Features

- **🔍 Smart Search**: Find stations by name or tag with popularity-based sorting
- **🤖 AI Integration**: MCP server mode for seamless AI assistant integration  
- **⚡ Fast & Lightweight**: Single binary with no external dependencies
- **🎵 Shell-Friendly**: Perfect for integration with `fzf`, `mpv`, and other tools
- **📊 Rich Output**: Detailed station info including click count, codec, bitrate, and streaming URL

## 🚀 Quick Start

### Installation

```bash
# Clone and build
git clone <repository-url>
cd bradio
go build -o bin/bradio .
```

### Basic Usage

```bash
# Search by station name
bradio --name "Milano Lounge"

# Search by tag/genre
bradio --tag "ambient"

# Limit results
bradio --tag "jazz" --limit 10

# Get help
bradio --help
```

## 🎛️ Operating Modes

### 1. CLI Mode (Default)

Traditional command-line interface for interactive use:

```bash
# Search examples
bradio --name "BBC Radio"
bradio --tag "classical" --limit 20
bradio --tag "electronic"
```

**Output Format:**
```
(15420) BBC Radio 1; pop,hits,uk,bbc; MP3[128]; http://stream.live.vc.bbcmedia.co.uk/bbc_radio_one
(8965) Jazz FM; jazz,smooth jazz,uk; AAC[128]; http://edge-audio-01-gos2.sharp-stream.com/jazzmmp3
```

### 2. MCP Server Mode

Model Context Protocol server for AI assistant integration:

```bash
# Start MCP server
bradio --mcp
```

**Available MCP Tools:**
- `search_radio_by_name` - Search stations by name
- `search_radio_by_tag` - Search stations by tag/genre  
- `get_popular_stations` - Get most popular stations globally

**Example MCP Usage:**
```bash
# List available tools
echo '{"jsonrpc":"2.0","id":1,"method":"tools/list","params":{}}' | bradio --mcp

# Search for jazz stations
echo '{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"search_radio_by_tag","arguments":{"tag":"jazz","limit":5}}}' | bradio --mcp
```

## 🎵 Integration Examples

### With fzf + mpv (Interactive Radio Player)

```bash
# One-liner for quick radio browsing
bradio --tag 'lounge' | fzf --nth 1 --preview="echo {}" --preview-window=bottom:2:nohidden | awk -F $'; ' '{print $4}' | mpv --playlist=-

# Create a handy function
function br() { 
    bradio "$@" | fzf --nth 1 --preview='echo {}' --preview-window=bottom:2:nohidden | awk -F $'; ' '{print $4}' | mpv --playlist=- 
}

# Use the function
br --tag 'chillout' --limit 50
```

### With AI Assistants (MCP Mode)

```bash
# Configure Claude/GPT to use bradio as MCP server
# Add to your AI assistant's MCP configuration:
{
  "servers": {
    "bradio": {
      "command": "/path/to/bradio",
      "args": ["--mcp"]
    }
  }
}
```

## 🛠️ Development

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
```

### Formatting
```bash
./run_format.sh
```

## 📚 Command Reference

### CLI Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--name` | Search stations by name | - |
| `--tag` | Search stations by tag/genre | - |
| `--limit` | Maximum number of results | 12 |
| `--help` | Show help information | - |
| `--mcp` | Run as MCP server | - |

### Search Tips

- **Popular genres**: `jazz`, `classical`, `rock`, `electronic`, `ambient`, `news`, `talk`
- **Language tags**: `english`, `french`, `german`, `spanish`, `italian`
- **Country codes**: `uk`, `usa`, `france`, `germany`, `italy`
- **Combine terms**: Use specific tags like `smooth jazz`, `deep house`, `indie rock`

## 🌐 Data Source

Bradio is powered by the excellent **[Radio Browser](https://www.radio-browser.info/)** - a community-driven database of internet radio stations. Radio Browser provides:

- 🌍 **Global Coverage**: Thousands of stations worldwide
- 🆓 **Free & Open**: No API keys or registration required  
- 📊 **Rich Metadata**: Detailed station information and statistics
- 🔄 **Live Updates**: Community-maintained, always current
- 🚀 **High Performance**: Fast, reliable API infrastructure

*Special thanks to the Radio Browser team and community for maintaining this fantastic resource!*

## 🏗️ Architecture

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

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- **[Radio Browser](https://www.radio-browser.info/)** - For providing the incredible radio station database
- **[goradios](https://gitlab.com/AgentNemo/goradios)** - Go library for Radio Browser API
- **[mcp-go](https://github.com/mark3labs/mcp-go)** - Go implementation of Model Context Protocol
