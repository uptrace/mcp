package tools

import (
	"context"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/uptrace/mcp/appconf"
	"github.com/uptrace/mcp/uptraceapi"
)

type ListSpanGroupsTool struct {
	client *uptraceapi.Client
	conf   *appconf.Config
}

func NewListSpanGroupsTool(client *uptraceapi.Client, conf *appconf.Config) *ListSpanGroupsTool {
	return &ListSpanGroupsTool{
		client: client,
		conf:   conf,
	}
}

func (t *ListSpanGroupsTool) Register(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name: "list_span_groups",
		Annotations: &mcp.ToolAnnotations{
			Title:          "List span groups",
			ReadOnlyHint:   true,
			IdempotentHint: true,
			OpenWorldHint:  boolPtr(true),
		},
		Description: "Aggregate spans into groups using UQL queries. " +
			"Use this to get aggregated metrics like request count, error rate, or latency percentiles. " +
			"Aggregate functions: count(), avg(), sum(), min(), max(), p50(), p75(), p90(), p99(), uniq(), apdex(). " +
			"Supports GROUP BY (e.g. group by service_name), HAVING (e.g. having p50(_dur_ms) > 100ms), " +
			"WHERE filters, full-text search, system filtering (e.g. httpserver:all, db:postgresql), " +
			"and duration filtering. Example query: 'perMin(count()) | group by host_name'. " +
			"Returns grouped rows with dynamic columns based on the query. " +
			"Use list_spans instead when you need individual span details. " +
			"Use timeseries instead when you need time-bucketed data for charts. " +
			"Documentation: https://uptrace.dev/features/querying/grouping",
	}, t.handler)
}

func (t *ListSpanGroupsTool) handler(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input *uptraceapi.ListSpanGroupsRequestOptions,
) (*mcp.CallToolResult, *uptraceapi.ListSpanGroupsResponse, error) {
	if input.PathParams.ProjectID == 0 {
		input.PathParams.ProjectID = t.conf.Uptrace.ProjectID
	}

	if input.Query == nil {
		input.Query = &uptraceapi.ListSpanGroupsQuery{}
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

	resp, err := t.client.ListSpanGroups(ctx, input)
	if err != nil {
		return nil, nil, err
	}

	return nil, resp, nil
}
