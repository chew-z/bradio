package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	// Check for MCP mode first
	for _, arg := range os.Args[1:] {
		if arg == "--mcp" {
			if err := RunMCPServer(); err != nil {
				log.Fatalf("MCP server error: %v", err)
			}
			return
		}
		if arg == "--help" || arg == "-h" {
			showMainUsage()
			os.Exit(0)
		}
	}

	// Default to CLI mode
	if err := RunCLI(); err != nil {
		log.Fatalf("CLI error: %v", err)
	}
}

// showMainUsage displays the main usage information
func showMainUsage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [mode] [options]\n\n", os.Args[0])
	fmt.Fprintln(os.Stderr, "A CLI tool and MCP server for searching internet radio stations via radio-browser.")
	fmt.Fprintln(os.Stderr, "\nModes:")
	fmt.Fprintln(os.Stderr, "  Default (no flags)    Run as CLI tool")
	fmt.Fprintln(os.Stderr, "  --mcp                Run as MCP server")
	fmt.Fprintln(os.Stderr, "  --help               Show this help")
	fmt.Fprintln(os.Stderr, "\nCLI Examples:")
	fmt.Fprintln(os.Stderr, "  bradio --name 'Milano Lounge'")
	fmt.Fprintln(os.Stderr, "  bradio --tag 'ambient' --limit 30")
	fmt.Fprintln(os.Stderr, "\nMCP Server Examples:")
	fmt.Fprintln(os.Stderr, "  bradio --mcp                    # Start MCP server with stdio transport")
	fmt.Fprintln(os.Stderr, "\nFor mode-specific help:")
	fmt.Fprintln(os.Stderr, "  bradio --help                   # CLI help")
}
