---
name: create-dashboard
description: Create an Uptrace monitoring dashboard by discovering metrics, analyzing attributes, and generating YAML. Use when the user wants to create a dashboard, visualize metrics, or monitor services.
argument-hint: [description of what to monitor]
---

# Create Uptrace Dashboard

Create a monitoring dashboard in Uptrace for: **$ARGUMENTS**

IMPORTANT: Call all Uptrace MCP tools directly. Do NOT delegate to subagents or Task tools — they do not have access to MCP tools.

## Workflow

Follow these steps in order:

### 1. Discover Metrics

Use `list_span_groups` to find relevant metrics:
- Use the `search` parameter to find metrics matching the user's request (e.g. search="http", search="redis")
- Use `group by _name` to list metric names
- Set a system filter if the domain is known (e.g. `where system = 'db:postgresql'`)
- To find metrics from a specific library, filter by `library_name`

### 2. Analyze Attributes

For each discovered metric:
- Use `list_span_groups` with `group by` to discover available attributes (e.g. `group by host_name, service_name`)
- Identify metrics that share the same attributes — these can be combined in table views
- Use `query_timeseries` to preview metric data if needed

### 3. Fetch Available Tags

Call `list_dashboard_tags` to get all available tags for the project.
Select the most relevant tags for this dashboard.

### 4. Study Existing Dashboards

1. Call `list_dashboards` to find similar existing dashboards
2. Call `get_dashboard_yaml` on 1-2 relevant dashboards to learn the exact YAML format
3. Use these as templates for structure, query syntax, and layout patterns

### 5. Generate Dashboard YAML

Build the YAML using the patterns learned from existing dashboards. Key format:

```yaml
schema: v2
name: "Dashboard Name"
tags:
- tag1
- tag2
version: v25.04.20
table:
- metrics:
  - metric_name as $var
  query:
  - perMin(count($var{})) as rate
  - p90($var{}) as dur_p90
  - group by service_name::str
grid_rows:
  - title: Section Title
    items:
      - title: Chart Title
        width: 6          # 12=full, 6=half, 4=third, 3=quarter
        height: 28         # 28 for charts, 10 for gauges
        type: chart        # chart, table, or gauge
        metrics:
          - metric_name as $var
        query:
          - perMin(count($var{}))
```

Syntax rules:
- `metrics`: "metric_name as $variable"
- `query`: aggregations — sum, avg, min, max, count, uniq, perMin, p50, p75, p90, p95, p99
- `group by`: separate query line — "group by attr_name::str"
- `filter`: curly braces — $var{status_code>=400}

**Important — count vs sum:**
- Use `count($var{})` for **duration/histogram metrics** (e.g. `http_server_duration`, `rpc_server_duration`) — counts the number of recorded measurements
- Use `sum($var{})` for **counter/span metrics** (e.g. `uptrace_tracing_spans`, `uptrace_tracing_logs`, `redis_keyspace_hits`) — sums pre-aggregated counter values
- `count` is NOT supported for span/counter metrics; `uniq` requires a specific attribute argument (e.g. `uniq($var{}, attr_name)`)

### 6. Create the Dashboard

Call `create_dashboard` with the generated YAML body.

### 7. Verify

Call `get_dashboard_yaml` with the returned dashboard ID to verify it was created correctly.
If issues are found, use `update_dashboard_yaml` to fix them.

## Output

After creating the dashboard, explain:
1. The list of selected metrics and what each one measures
2. How metrics were grouped and why
3. What each dashboard section/visualization shows
4. Which tags were applied and why
