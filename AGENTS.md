# AGENTS.md

Uptrace MCP Server written in Go using `github.com/modelcontextprotocol/go-sdk`.

## Project Structure

```
mcp/
├── cmd/                    # Application entrypoints
│   └── mcp-server/         # Main server binary
├── internal/               # Private application code
│   ├── config/             # Configuration loading
│   ├── handler/            # MCP request handlers
│   ├── uptrace/            # Uptrace API client
│   └── tools/              # MCP tool implementations
├── openapi/                # Git submodule: Uptrace OpenAPI specification
│   └── openapi.yaml        # OpenAPI 3.x spec for Uptrace API
├── pkg/                    # Public libraries (if any)
└── module.go               # Root fx module
```

## OpenAPI Specification

The Uptrace OpenAPI specification is available at `openapi/openapi.yaml` (git submodule).

Use this spec to generate API clients or reference available endpoints.

## Dependencies

- `github.com/modelcontextprotocol/go-sdk` - MCP protocol implementation
- `go.uber.org/fx` - Dependency injection framework

---

## Code Style Guidelines

### Prefer Early Returns

```go
// Good
func process(data *Data) error {
    if data == nil {
        return ErrNilData
    }
    if !data.IsValid() {
        return ErrInvalidData
    }
    return doProcess(data)
}

// Bad
func process(data *Data) error {
    if data != nil {
        if data.IsValid() {
            return doProcess(data)
        }
        return ErrInvalidData
    }
    return ErrNilData
}
```

### Code Ordering (Top-Down)

1. High-level functions first
2. Helper functions follow
3. Callees appear after their callers

```go
// Good: caller before callee
func HandleRequest(req *Request) (*Response, error) {
    validated := validateRequest(req)
    return buildResponse(validated)
}

func validateRequest(req *Request) *ValidatedRequest { ... }
func buildResponse(v *ValidatedRequest) (*Response, error) { ... }
```

### Memory Management

- **Define clear ownership** - Document which function/struct owns allocated memory
- **Avoid unnecessary copying** - Pass pointers for large structs
- **Minimize allocations** - Reuse buffers, use sync.Pool where appropriate
- **Keep function arguments immutable by default** - Do not mutate input parameters
- **Document when mutation is required** - Use comments to indicate intentional mutation

```go
// Good: clear ownership, no mutation of input
func Transform(input []byte) []byte {
    result := make([]byte, len(input))
    copy(result, input)
    // transform result...
    return result
}

// When mutation is required, document it
// ProcessInPlace modifies data in place for performance.
// Caller must ensure exclusive access.
func ProcessInPlace(data []byte) { ... }
```

### Dependency Injection

Uses `go.uber.org/fx` throughout. Each module has `module.go` with fx definitions.

```go
// internal/handler/module.go
package handler

import "go.uber.org/fx"

var Module = fx.Options(
    fx.Provide(NewHandler),
)

// internal/uptrace/module.go
package uptrace

import "go.uber.org/fx"

var Module = fx.Options(
    fx.Provide(NewClient),
)
```

---

## API Reference

### Functions

<!-- Document all exported functions here -->

| Function | Package | Description |
|----------|---------|-------------|
| - | - | No functions defined yet |

### Structs

<!-- Document all exported structs and their fields here -->

| Struct | Package | Description |
|--------|---------|-------------|
| - | - | No structs defined yet |

### Package-Level Variables

<!-- Document all exported package-level variables here -->

| Variable | Package | Type | Description |
|----------|---------|------|-------------|
| - | - | - | No variables defined yet |

### Constants

<!-- Document all exported constants here -->

| Constant | Package | Value | Description |
|----------|---------|-------|-------------|
| - | - | - | No constants defined yet |

---

## MCP Tools

This server exposes the following MCP tools:

| Tool Name | Description |
|-----------|-------------|
| - | No tools implemented yet |

---

## Configuration

Environment variables:

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `UPTRACE_DSN` | Yes | - | Uptrace DSN for API access |

---

## Development

```bash
# Run the server
go run ./cmd/mcp-server

# Run tests
go test ./...

# Build
go build -o mcp-server ./cmd/mcp-server
```
