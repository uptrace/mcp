# Uptrace MCP Server

A Model Context Protocol (MCP) server for [Uptrace](https://uptrace.dev), written in Go.

## Overview

This server implements the [Model Context Protocol](https://modelcontextprotocol.io/) to expose Uptrace observability data to AI assistants and LLM-powered tools.

## Prerequisites

- Go 1.21+
- Access to an Uptrace instance

## Getting Started

### Clone the repository

```bash
git clone git@github.com:uptrace/mcp.git
cd mcp
```

### Initialize submodules

```bash
git submodule update --init --recursive
```

### Build

```bash
go build -o mcp-server ./cmd/mcp-server
```

### Configure

Create a config file:

```bash
cp config.yaml.example config.yaml
```

Edit `config.yaml` with your Uptrace credentials:

```yaml
uptrace:
  api_url: "https://api.uptrace.dev"
  api_token: "<your-api-token>"
  project_id: 1
```

### Run

```bash
./mcp-server --config config.yaml
```

## Adding to Claude Code

Add the following to your Claude Code MCP settings (`~/.claude/claude_desktop_config.json` on macOS or `%APPDATA%\Claude\claude_desktop_config.json` on Windows):

```json
{
  "mcpServers": {
    "uptrace": {
      "command": "/path/to/mcp-server",
      "args": ["--config", "/path/to/config.yaml"]
    }
  }
}
```

Replace `/path/to/mcp-server` and `/path/to/config.yaml` with absolute paths.

## Available Tools

| Tool | Description |
|------|-------------|
| `list_spans` | List spans from Uptrace. Supports time range filtering, trace ID filtering, and pagination. |
| `greet` | Example greeting tool (for testing) |

### list_spans

Fetch spans from Uptrace for analyzing distributed traces.

**Parameters:**
- `time_start` (required): Start time in RFC3339 format or relative (e.g., `-1h`, `-30m`)
- `time_end` (optional): End time, defaults to now
- `trace_id` (optional): Filter by specific trace ID
- `limit` (optional): Maximum number of spans to return (default: 100)

**Example usage in Claude:**
> "Show me spans from the last hour"
> "Find spans for trace ID abc123"

## Configuration

| Field | Required | Description |
|-------|----------|-------------|
| `uptrace.api_url` | Yes | Uptrace API URL |
| `uptrace.api_token` | Yes | API token for authentication |
| `uptrace.project_id` | Yes | Uptrace project ID |
| `logging.level` | No | Log level: debug, info, warn, error (default: info) |

## Development

```bash
# Run with hot-reload
task dev

# Run tests
task test

# Regenerate API client
task generate
```

See [AGENTS.md](AGENTS.md) for coding guidelines and project documentation.

## License

See [LICENSE](LICENSE) for details.
