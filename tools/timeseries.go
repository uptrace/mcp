package tools

import (
	"context"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/uptrace/mcp/appconf"
	"github.com/uptrace/mcp/uptraceapi"
)

type QueryTimeseriesTool struct {
	client *uptraceapi.Client
	conf   *appconf.Config
}

func NewQueryTimeseriesTool(client *uptraceapi.Client, conf *appconf.Config) *QueryTimeseriesTool {
	return &QueryTimeseriesTool{
		client: client,
		conf:   conf,
	}
}

func (t *QueryTimeseriesTool) Register(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name: "timeseries",
		Annotations: &mcp.ToolAnnotations{
			Title:          "Query timeseries",
			ReadOnlyHint:   true,
			IdempotentHint: true,
			OpenWorldHint:  boolPtr(true),
		},
		Description: "Query time-bucketed aggregation data for spans. " +
			"Use this to analyze trends over time, detect anomalies, or build charts. " +
			"Returns aligned timestamps with auto-computed interval based on the time range. " +
			"Use UQL aggregation query (e.g. 'perMin(count()) | group by service_name') " +
			"with aggregate functions: count(), avg(), sum(), p50(), p90(), p99(), etc. " +
			"Use 'column' parameter to select specific aggregate columns for the timeseries. " +
			"Supports WHERE filters, full-text search, system filtering, and duration filtering. " +
			"Returns groups with arrays of float values aligned with the time array. " +
			"Use list_span_groups instead when you need a single aggregated snapshot (not time-bucketed). " +
			"Use quantiles instead when you only need latency percentiles (p50/p90/p99). " +
			"Documentation: https://uptrace.dev/features/querying/spans",
	}, t.handler)
}

func (t *QueryTimeseriesTool) handler(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input *uptraceapi.QueryTimeseriesRequestOptions,
) (*mcp.CallToolResult, *uptraceapi.QueryTimeseriesResponse, error) {
	if input.PathParams.ProjectID == 0 {
		input.PathParams.ProjectID = t.conf.Uptrace.ProjectID
	}

	if input.Query == nil {
		input.Query = &uptraceapi.QueryTimeseriesQuery{}
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
	if input.Query.Query == nil {
		input.Query.Query = &t.conf.Default.Query
	}

	resp, err := t.client.QueryTimeseries(ctx, input)
	if err != nil {
		return nil, nil, err
	}

	return nil, resp, nil
}
