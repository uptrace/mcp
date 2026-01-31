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
│   └── greet.go            # Example greet tool
└── openapi/                # Git submodule: Uptrace OpenAPI specification
    └── openapi.yaml        # OpenAPI 3.x spec for Uptrace API
```

## OpenAPI Specification

The Uptrace OpenAPI specification is available at `openapi/openapi.yaml` (git submodule).

Use this spec to generate API clients or reference available endpoints.

## Dependencies

- `github.com/modelcontextprotocol/go-sdk` - MCP protocol implementation
- `go.uber.org/fx` - Dependency injection framework
- `github.com/urfave/cli/v3` - CLI argument parsing
- `gopkg.in/yaml.v3` - YAML config parsing

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

### Structs

| Struct | Package | Fields |
|--------|---------|--------|
| `Config` | `appconf` | `Uptrace UptraceConfig` |
| `UptraceConfig` | `appconf` | `DSN string` |
| `GreetArgs` | `tools` | `Name string` |

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
| `greet` | Say hello to someone (example tool) |

---

## Configuration

Config file (`config.yaml`):

```yaml
uptrace:
  dsn: "https://<token>@api.uptrace.dev/<project_id>"
```

See `config.yaml.example` for reference.

---

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
