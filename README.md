# Bract - MCP Browser Automation Server

Bract is a Go implementation of a Model Context Protocol (MCP) server that enables browser automation through a Chrome extension. It provides a standardized interface for AI assistants and automation tools to control web browsers programmatically.

## Features

- 🌐 **WebSocket-based Communication**: Real-time bidirectional communication
- 🔧 **Comprehensive Browser Control**: Tab management, content interaction, and capture capabilities
- 🔒 **Secure by Default**: Localhost-only connections with input validation
- 🚀 **High Performance**: Optimized for concurrent operations
- 🔄 **Automatic Reconnection**: Robust connection handling with retry logic
- 📝 **MCP Compliant**: Full implementation of the Model Context Protocol

## Architecture

```
MCP Client ←→ WebSocket (ws://localhost:8765) ←→ Bract Server ←→ Chrome Extension
```

## Quick Start

```bash
# Clone the repository
git clone https://github.com/periplon/bract.git
cd bract

# Install dependencies
go mod download

# Build the server
go build -o bin/bract cmd/mcp-server/main.go

# Run the server
./bin/bract
```

## Documentation

- [Product Requirements Document](docs/PRD-mcp-browser-automation-server.md)
- [Development Guidelines](CLAUDE.md)

## Project Status

🚧 **In Development** - This project is currently in the planning phase. See the PRD for implementation timeline and roadmap.

## Contributing

Please read our [development guidelines](CLAUDE.md) before contributing. We follow conventional commits and semantic versioning.

## License

TBD