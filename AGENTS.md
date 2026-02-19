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
│   ├── list_monitors.go    # list_monitors tool
│   ├── create_dashboard.go # create_dashboard tool
│   ├── list_dashboards.go  # list_dashboards tool
│   ├── get_dashboard.go    # get_dashboard tool
│   ├── get_dashboard_yaml.go    # get_dashboard_yaml tool
│   ├── update_dashboard_yaml.go # update_dashboard_yaml tool
│   └── delete_dashboard.go      # delete_dashboard tool
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
| `query_timeseries` | Query timeseries data for metrics |
| `query_quantiles` | Query quantile data for metrics |
| `list_trace_groups` | List trace groups |
| `list_traces` | List traces |
| `list_monitors` | List monitors (alerts) from Uptrace |
| `create_dashboard` | Create a dashboard from YAML definition |
| `list_dashboards` | List all dashboards |
| `list_dashboard_tags` | List available dashboard tags for a project |
| `get_dashboard` | Get dashboard details by ID |
| `get_dashboard_yaml` | Get dashboard YAML definition by ID |
| `update_dashboard_yaml` | Update a dashboard from YAML definition |
| `delete_dashboard` | Delete a dashboard by ID |

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

**IMPORTANT:** Before creating a dashboard, use `get_dashboard_yaml` on an existing dashboard to see the exact format. Only use fields shown below — the API strictly rejects unknown fields.

### Minimal example

```yaml
schema: v2
name: "HTTP: Server"
tags:
- otel
- app
grid_rows:
- title: Overview
  items:
  - title: Request Rate
    width: 12
    height: 28
    type: chart
    metrics:
    - http_server_duration as $dur
    query:
    - perMin(count($dur{}))
```

### Full YAML structure reference

```yaml
schema: v2                    # Required: v2 or v3
name: "Dashboard Name"        # Required
tags:                          # Optional: string labels
- otel
- app
version: v25.04.20            # Optional

# Table section (summary table at top of dashboard)
table:
- metrics:
  - metric_name as $var
  query:
  - perMin(count($var{})) as rate
  - p90($var{}) as p90
  - group by service_name::str
  overrides:                   # Table column overrides
  - column: p90                # Must match a query alias
    properties:
    - name: unit
      value: milliseconds

# Table grid items (optional gauge widgets above the table)
table_grid_items:
- title: Request Rate
  width: 3                     # Gauges typically use width: 3
  height: 10                   # Gauges typically use height: 10
  type: gauge
  metrics:
  - metric_name as $var
  query:
  - perMin(count($var{})) as rate
  overrides:                   # Uses column format (NOT matchers)
  - column: rate               # REQUIRED: must match a query alias
    properties:
    - name: unit
      value: utilization

# Grid section (charts and gauges)
grid_rows:
- title: Section Title
  items:
  - title: Chart Title
    width: 12                  # Standard: 12 (half), 24 (full), 3 (gauge)
    height: 28                 # Standard: 28 (chart), 10 (gauge), 40 (heatmap/table)
    x_axis: 12                 # Optional: horizontal offset (12 = right column)
    y_axis: 28                 # Optional: vertical offset (multiples of height)
    type: chart                # chart | gauge | table | heatmap
    metrics:
    - metric_name as $var
    query:
    - perMin(count($var{}))
    properties:                 # Chart-level properties
    - name: fillOpacity         # Only known property for charts
      value: 0.1
    overrides:                  # Per-metric overrides
    - matchers:                 # Match by metric alias
      - title: "metric:alias"
        target: metric
        value: alias
      properties:
      - name: unit
        value: milliseconds
```

### Key rules

- **No unknown fields.** The API rejects any field not in the schema. Use `get_dashboard_yaml` on an existing dashboard to verify the format.
- **`overrides` format differs** by context:
  - **Table overrides**: `column` (required) + `properties` — match by query alias
  - **`table_grid_items` overrides**: `column` (required) + `properties` — same as table, `column` must match a query alias
  - **Grid item overrides** (in `grid_rows`): `matchers` + `properties` — match by metric alias
- **`properties` and `overrides`** are arrays of `{name, value}` pairs, never bare key-value maps.
- **Sparkline** cannot be set via the YAML create/update API. Do NOT include `sparkline` in overrides.
- **Empty arrays** (`properties: []`, `overrides: []`) are valid and can be omitted.
- **Only known chart property**: `fillOpacity: 0.1` (for stacked/area charts).

### Unit values

`utilization`, `bytes`, `milliseconds`, `microseconds`, `nanoseconds`, `seconds`, `"1"` (dimensionless),
`log/min`, `span/min`, `req/min`, `call/min`, `query/min`, `services`, `hosts`

### Layout conventions

- **Two-column layout**: width 12 + `x_axis: 12` for side-by-side charts
- **Three-column layout**: width 8 + `x_axis: 8` + `x_axis: 16`
- **Gauge row**: width 3 + `x_axis: 3, 6, 9, 12, 15` (up to 8 gauges)
- **Vertical stacking**: `y_axis` increments by item height (28 for charts, 10 for gauges)

### Aggregate functions

- `sum($var{})` — counter/gauge metrics
- `perMin(sum($var{}))` — rate of counters
- `count($var{})` — histogram/duration metrics only
- `perMin(count($var{}))` — rate of histograms
- `avg($var{})` — averages
- `p50()`, `p90()`, `p99()` — latency percentiles (duration/histogram metrics only)
- `max()`, `min()` — extremes
- `uniq($var{}, attr)` — count distinct values of an attribute

### Metric types and allowed functions

**Duration/histogram metrics** (e.g. `http_server_duration`, `go_sql_query_timing`, `rpc_server_duration`):
- `count()`, `perMin(count())`, `avg()`, `p50()`, `p90()`, `p99()`, `max()`, `min()`

**Counter/span metrics** (e.g. `uptrace_tracing_spans`, `uptrace_tracing_logs`, `redis_commands`):
- `sum()`, `perMin(sum())` — for rates and totals
- `uniq($var{}, attr)` — to count distinct values (e.g. `uniq($spans{}, service_name::str)`)
- Do NOT use `count()`, `p50()`, `p90()`, `p99()` on span/counter metrics — they will error with "count is not supported for span metrics, use uniq instead"

**Gauge metrics** (e.g. `system_memory_usage`, `redis_memory_rss`, `redis_db_keys`):
- `sum()`, `avg()`, `max()`, `min()`

### Tags

Simple string labels for categorizing dashboards. Common tags:
`otel`, `app`, `db`, `infra`, `logs`, `tracing`, `network`, `self_monitoring`
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
