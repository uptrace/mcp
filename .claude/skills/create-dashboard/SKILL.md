---
name: create-dashboard
description: Create an Uptrace monitoring dashboard by discovering metrics, analyzing attributes, and generating YAML. Use when the user wants to create a dashboard, visualize metrics, or monitor services.
argument-hint: [description of what to monitor]
---

# Uptrace Dashboard Creation

This skill creates monitoring dashboards in Uptrace by discovering available metrics, analyzing their attributes, and generating valid YAML definitions. It handles the full workflow from metric discovery to dashboard creation and verification.

IMPORTANT: Call all Uptrace MCP tools directly. Do NOT delegate to subagents or Task tools — they do not have access to MCP tools.

## When to Use This Skill

- Creating a new monitoring dashboard for a service or library
- Visualizing metrics from OpenTelemetry instrumentation
- Building overview dashboards with charts, gauges, and tables
- Monitoring specific subsystems (HTTP, database, Redis, runtime, etc.)

## What This Skill Does

1. **Discovers metrics** using `explore_metrics` — finds available metrics with instrument types, units, and attributes
2. **Analyzes attributes** using `list_metric_attributes` / `list_metric_attribute_values` — identifies grouping dimensions
3. **Groups by library** — metrics from the same `libraryName` share attributes and belong on one dashboard
4. **Learns from templates** using `list_dashboard_templates` / `get_dashboard_template` — uses curated built-in templates as reference for correct metric combinations and query patterns
5. **Generates YAML** — produces valid dashboard YAML following the strict schema
6. **Creates and verifies** — submits via `create_dashboard` and verifies with `get_dashboard_yaml`

## How to Use

Create a monitoring dashboard in Uptrace for: **$ARGUMENTS**

### Workflow

1. `explore_metrics` — discover metrics, group by `libraryName`
2. `list_metric_attributes` / `list_metric_attribute_values` — find grouping attributes
3. `list_dashboard_templates` — find a matching built-in template to use as reference
4. `get_dashboard_template` — fetch the template YAML to learn structure, metric combinations, and query patterns
5. `list_dashboard_tags` — get available tags
6. `create_dashboard` — submit YAML
7. `get_dashboard_yaml` — verify result

## Instructions

### Metric Discovery

Uptrace metrics come from [OpenTelemetry instrumentation](https://uptrace.dev/opentelemetry/metrics) — each metric is produced by an instrument (Counter, Histogram, Gauge, UpDownCounter) and carries attributes (key-value pairs like `http.method=GET`). See [Querying Metrics](https://uptrace.dev/features/querying/metrics) for the full query language reference.

#### Searching metrics with `explore_metrics`

Use `explore_metrics` with the `search` parameter to filter metrics by name. The search is a substring match on the metric name:

- `search="http"` — finds `http_server_duration`, `http_client_duration`, etc.
- `search="redis"` — finds `redis_commands`, `redis_memory_rss`, `redis_db_keys`, etc.
- `search="go_sql"` — finds `go_sql_query_timing`, `go_sql_open_connections`, etc.
- `search=""` (empty or omitted) — returns all available metrics

Each metric in the response contains:

- **`name`** — metric name (e.g. `http_server_duration`)
- **`instrument`** — OpenTelemetry instrument type that determines which aggregate functions are allowed: `histogram`, `counter`, `gauge`, `additive`. Check the "Aggregate Functions by Instrument Type" section to pick the correct functions.
- **`unit`** — metric unit (e.g. `milliseconds`, `bytes`). Use this to set the `unit` property in overrides.
- **`description`** — what the metric measures
- **`attrKeys`** — all available attributes with type suffixes (e.g. `service_name::str`, `http_status_code::int`). Use these for `group by` clauses and `$var{attr="value"}` filters in queries.
- **`libraryName`** — the OpenTelemetry instrumentation library that produces this metric (e.g. `go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp`). This is the OTel [instrumentation scope](https://uptrace.dev/opentelemetry/metrics) name.
- **`numTimeseries`** — number of active timeseries (helps gauge metric cardinality)

**Dashboards are built from metrics that share the same `libraryName`.** Group discovered metrics by `libraryName` — metrics from the same library share the same attributes and belong on one dashboard. For example, all Redis metrics come from `opentelemetry-collector-contrib/receiver/redisreceiver`, all HTTP server metrics from `go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp`.

#### Finding attributes with `list_metric_attributes` and `list_metric_attribute_values`

Use `list_metric_attributes` (with optional `search` to filter by metric name) to see attribute keys with usage counts — helpful for finding common `group by` attributes for dashboard tables.

Use `list_metric_attribute_values` with a specific `attr_key` including the type suffix (e.g. `host_name::str`, `http_status_code::int`) to check actual values for filtering (e.g. `$var{service_name="myservice"}`).

### Aggregate Functions by Instrument Type

- **histogram** — `count()`, `perMin(count())`, `avg()`, `p50()`, `p90()`, `p99()`, `max()`, `min()`
- **counter** — `sum()`, `perMin(sum())`
- **gauge** — `sum()`, `avg()`, `max()`, `min()`
- **additive** — `sum()`, `perMin(sum())`

**CRITICAL:** Do NOT use `count()`, `p50()`, `p90()`, `p99()` on counter/gauge/additive metrics — they will error with "count is not supported for span metrics".

### Learning from Templates

**Always prefer built-in templates over existing project dashboards as examples.** Templates are curated and correct — project dashboards may have incorrect metric combinations, bad query patterns, or nonsensical table/grid layouts.

1. Call `list_dashboard_templates` to see all available templates (id, name, description)
2. Find a template that matches the domain you're building for (e.g. HTTP, Redis, Go runtime, PostgreSQL)
3. Call `get_dashboard_template` with the template ID to fetch its YAML structure
4. Use the template as a reference for metric selection, query patterns, layout, and tags

Only fall back to `get_dashboard_yaml` on existing project dashboards if no matching template exists.

### Dashboard YAML Format

**WARNING:** If you use `get_dashboard_yaml` on existing dashboards, the output contains read-only properties (`sparkline`, `color`, `aggFunc`) that the API **rejects on create/update**. Do NOT copy these properties — only use `unit` in overrides. Templates from `get_dashboard_template` do NOT have this problem.

```yaml
schema: v2
name: "Dashboard Name"
tags:
- otel
- app

table:
- metrics:
  - metric_name as $var
  query:
  - perMin(sum($var{})) as rate
  - group by service_name::str
  overrides:
  - column: rate              # REQUIRED: must match a query alias
    properties:
    - name: unit
      value: bytes

table_grid_items:             # Optional: gauge widgets above the table
- title: Total Rate
  width: 3
  height: 10
  type: gauge
  metrics:
  - metric_name as $var
  query:
  - perMin(sum($var{})) as rate
  overrides:                  # Same column format as table
  - column: rate
    properties:
    - name: unit
      value: req/min

grid_rows:
- title: Section Title
  items:
  - title: Chart Title
    width: 12                 # 12=half, 24=full, 3=gauge
    height: 28                # 28=chart, 10=gauge, 40=heatmap
    type: chart               # chart | gauge | table | heatmap
    metrics:
    - metric_name as $var
    query:
    - perMin(sum($var{}))
    properties:               # Only known: fillOpacity
    - name: fillOpacity
      value: 0.1
    overrides:                # Grid items use matchers format
    - matchers:
      - title: "metric:alias"
        target: metric
        value: alias
      properties:
      - name: unit
        value: bytes
```

### Critical Rules

- **No unknown fields** — the API rejects any field not in the schema
- **Override formats differ by context:**
  - Table / `table_grid_items`: `column` (required, must match query alias) + `properties`
  - Grid items (in `grid_rows`): `matchers` + `properties`
- **Properties** are arrays of `{name, value}` pairs, never bare key-value maps
- **Do NOT copy read-only properties from `get_dashboard_yaml` output.** The following properties appear in read output but are **rejected on create/update**:
  - `sparkline` — never include in overrides
  - `color` — never include (e.g. `name: color, value: ""`)
  - `aggFunc` — never include (e.g. `name: aggFunc, value: ""`)
- **Only include `unit` property** in table/table_grid_items overrides. For grid item overrides, only `unit` is needed.
- **Unit values:** `utilization`, `bytes`, `milliseconds`, `microseconds`, `nanoseconds`, `seconds`, `"1"`, `req/min`, `span/min`, `log/min`, `call/min`, `query/min`

### Layout Reference

- **Two-column:** width `12` + `x_axis: 12`
- **Three-column:** width `8` + `x_axis: 8` + `x_axis: 16`
- **Gauge row:** width `3`, height `10`, `x_axis: 3, 6, 9, 12, 15`
- **Vertical stacking:** `y_axis` increments by item height (28 for charts, 10 for gauges)

## Output

After creating the dashboard, explain:
1. Selected metrics, their instrument types, and what each measures
2. How metrics were grouped (by `libraryName`) and why
3. What each dashboard section shows
4. Which tags were applied
