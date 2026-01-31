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

The repository uses git submodules for the OpenAPI specification:

```bash
git submodule update --init --recursive
```

Or clone with submodules in one step:

```bash
git clone --recurse-submodules git@github.com:uptrace/mcp.git
```

### Build

```bash
go build -o mcp-server ./cmd/mcp-server
```

### Run

```bash
export UPTRACE_DSN="https://<token>@api.uptrace.dev/<project_id>"
./mcp-server
```

## Configuration

| Environment Variable | Required | Description |
|---------------------|----------|-------------|
| `UPTRACE_DSN` | Yes | Uptrace DSN for API access |

## Development

See [AGENTS.md](AGENTS.md) for coding guidelines and project documentation.

## License

See [LICENSE](LICENSE) for details.
