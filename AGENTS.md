# AGENTS.md

Uptrace MCP Server written in Go using `github.com/modelcontextprotocol/go-sdk`.

## Project Structure

```
mcp/
├── cmd/mcp-server/         # Application entrypoint
│   └── main.go
├── appconf/                # Configuration loading
│   └── config.go
├── bootstrap/              # fx app bootstrap
│   └── bootstrap.go
├── tools/                  # MCP tool implementations
│   ├── tools.go            # Register() for all tools
│   ├── list_groups.go      # list_span_groups tool (UQL aggregation)
│   ├── list_spans.go       # list_spans tool
│   └── list_monitors.go    # list_monitors tool
├── uptraceapi/             # Generated Uptrace API client (do not edit)
│   ├── client.go           # HTTP client methods
│   ├── types.go            # API types and models
│   └── ...                 # Other generated files
└── openapi/                # Git submodule: Uptrace OpenAPI specification
    └── openapi.yaml        # OpenAPI 3.1 spec for Uptrace API
```

## OpenAPI Specification

The Uptrace OpenAPI specification is available at `openapi/openapi.yaml` (git submodule).

Use this spec to generate API clients or reference available endpoints.

## Dependencies

- `github.com/modelcontextprotocol/go-sdk` - MCP protocol implementation
- `go.uber.org/fx` - Dependency injection framework
- `github.com/urfave/cli/v3` - CLI argument parsing
- `github.com/goccy/go-yaml` - YAML config parsing
- `github.com/doordash-oss/oapi-codegen-dd/v3` - Generated API client runtime

## Code Generation

The `uptraceapi/` package is generated from the OpenAPI spec. Do not edit manually.

```bash
task generate  # Regenerate API client from openapi/openapi.yaml
```

Config: `oapi-codegen.yaml`

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

Uses `go.uber.org/fx` throughout. Bootstrap provides core dependencies, additional modules passed via `fx.Option`.

```go
// cmd/mcp-server/main.go
bootstrap.Run(
    ctx,
    cmd,
    fx.Invoke(runServer),
)

// bootstrap/bootstrap.go - provides *appconf.Config, *slog.Logger
func New(conf *appconf.Config, options ...fx.Option) *fx.App
func Run(ctx context.Context, cmd *cli.Command, options ...fx.Option) error
```

---

## API Reference

### Functions

| Function | Package | Description |
|----------|---------|-------------|
| `Load(path string)` | `appconf` | Load config from YAML file |
| `Parse(data []byte)` | `appconf` | Parse YAML bytes into Config |
| `New(conf, ...fx.Option)` | `bootstrap` | Create fx.App with config |
| `Run(ctx, cmd, ...fx.Option)` | `bootstrap` | Load config and run fx app |
| `Register(server)` | `tools` | Register all MCP tools |
| `NewClient(server, opts)` | `uptraceapi` | Create Uptrace API client (generated) |

### Structs

| Struct | Package | Fields |
|--------|---------|--------|
| `Config` | `appconf` | `Service ServiceConfig`, `Logging LoggingConfig`, `Uptrace UptraceConfig` |
| `ServiceConfig` | `appconf` | `StartTimeout time.Duration`, `StopTimeout time.Duration` |
| `LoggingConfig` | `appconf` | `Level string`, `MaxBodySize int` |
| `UptraceConfig` | `appconf` | `DSN string`, `APIURL string`, `APIToken string`, `ProjectID int64` |

### Package-Level Variables

| Variable | Package | Type | Description |
|----------|---------|------|-------------|
| - | - | - | None |

### Constants

| Constant | Package | Value | Description |
|----------|---------|-------|-------------|
| - | - | - | None |

---

## MCP Tools

This server exposes the following MCP tools:

| Tool Name | Description |
|-----------|-------------|
| `list_span_groups` | Aggregate spans using UQL (Uptrace Query Language) |
| `list_spans` | List spans from Uptrace for analyzing distributed traces |
| `list_monitors` | List monitors (alerts) from Uptrace |

---

## Configuration

Config file (`config.yaml`):

```yaml
service:
  start_timeout: 15s  # optional, default: 15s
  stop_timeout: 15s   # optional, default: 15s

logging:
  level: info  # optional, default: info (debug, info, warn, error)

uptrace:
  dsn: "https://<token>@api.uptrace.dev/<project_id>"
  api_url: "https://api.uptrace.dev"
  api_token: "<your-api-token>"
default:
  project_id: 1
  query: "<your-query>"
  limit: 10
```

See `config.yaml.example` for reference.

---
## Creating Dashboards

  Example YAML structure:
  ```yaml
  schema: v2
  name: My Dashboard
  tags: []
  version: v25.04.20
  grid_rows:
    - title: General
      items:
        - title: Metric name
          width: 12
          height: 28
          type: chart
          metrics:
            - metric_name as $var
          query:
            - sum($var)
## Development

```bash
# Copy example config
cp config.yaml.example config.yaml

# Run with hot-reload (requires watchexec)
task dev

# Run manually
go run ./cmd/mcp-server --config config.yaml

# Run tests
go test ./...

# Build
go build -o mcp-server ./cmd/mcp-server
```
