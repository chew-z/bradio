# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview
bradio is a CLI tool for searching radio stations on radio-browser by name or tag, with sorting by popularity. It's a simple Go application that uses the `gitlab.com/AgentNemo/goradios` library to interact with the radio-browser API.

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
This is a Go application (`main.go`) with a robust command-line interface:
- Uses Go's `flag` package for proper argument parsing with validation
- Takes command-line arguments for search type (`--name` or `--tag`) and optional `--limit`
- Includes comprehensive error handling and user-friendly help messages
- Modular design with separate functions for validation and station fetching
- Uses the goradios library to fetch station data from radio-browser API
- Outputs formatted station information including click count, name, tags, codec, bitrate, and URL
- Output format is designed to work with shell pipelines (especially fzf + mpv for radio playback)

## Key Dependencies
- Go 1.24.4 (specified in go.mod)
- `gitlab.com/AgentNemo/goradios` for radio station API access

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

## Usage Patterns
The tool is designed for integration with shell workflows, particularly:
- Piping output to `fzf` for interactive selection
- Extracting URLs for media players like `mpv`
- The output format supports both human reading and programmatic parsing
- Proper error handling prevents crashes and provides helpful user feedback