package tools

import (
	"context"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/uptrace/mcp/appconf"
	"github.com/uptrace/mcp/uptraceapi"
)

type QueryQuantilesTool struct {
	client *uptraceapi.Client
	conf   *appconf.Config
}

func NewQueryQuantilesTool(client *uptraceapi.Client, conf *appconf.Config) *QueryQuantilesTool {
	return &QueryQuantilesTool{
		client: client,
		conf:   conf,
	}
}

func (t *QueryQuantilesTool) Register(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name: "quantiles",
		Annotations: &mcp.ToolAnnotations{
			Title:          "Query quantiles",
			ReadOnlyHint:   true,
			IdempotentHint: true,
			OpenWorldHint:  boolPtr(true),
		},
		Description: "Query duration percentiles (p50, p90, p99) and count/error rate over time. " +
			"Use this to identify latency outliers, track performance degradation, or compare SLO compliance. " +
			"Returns named timeseries: count, countPerMin, errorCount, errorCountPerMin, " +
			"durationP50, durationP90, durationP99, durationMax. " +
			"Supports WHERE filters (e.g. where service_name = 'myservice'), " +
			"full-text search, system filtering (e.g. httpserver:all), and duration filtering. " +
			"Use timeseries instead when you need custom aggregation queries with GROUP BY. " +
			"Use list_span_groups instead when you need a single aggregated snapshot. " +
			"Documentation: https://uptrace.dev/features/querying/spans",
	}, t.handler)
}

func (t *QueryQuantilesTool) handler(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input *uptraceapi.QueryQuantilesRequestOptions,
) (*mcp.CallToolResult, *uptraceapi.QueryQuantilesResponse, error) {
	if input.PathParams.ProjectID == 0 {
		input.PathParams.ProjectID = t.conf.Uptrace.ProjectID
	}

	if input.Query == nil {
		input.Query = &uptraceapi.QueryQuantilesQuery{}
	}
	if input.Query.TimeGte.IsZero() {
		input.Query.TimeGte = time.Now().Add(-time.Hour)
	}
	if input.Query.TimeLt.IsZero() {
		input.Query.TimeLt = time.Now()
	}
	if input.Query.Limit == nil {
		defaultLimit := uptraceapi.Limit(t.conf.Default.Limit)
		input.Query.Limit = &defaultLimit
	}

	resp, err := t.client.QueryQuantiles(ctx, input)
	if err != nil {
		return nil, nil, err
	}

	return nil, resp, nil
}
