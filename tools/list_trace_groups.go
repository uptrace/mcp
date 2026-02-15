package tools

import (
	"context"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/uptrace/mcp/appconf"
	"github.com/uptrace/mcp/uptraceapi"
)

type ListTraceGroupsTool struct {
	client *uptraceapi.Client
	conf   *appconf.Config
}

func NewListTraceGroupsTool(client *uptraceapi.Client, conf *appconf.Config) *ListTraceGroupsTool {
	return &ListTraceGroupsTool{
		client: client,
		conf:   conf,
	}
}

func (t *ListTraceGroupsTool) Register(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name: "list_trace_groups",
		Annotations: &mcp.ToolAnnotations{
			Title:          "List trace groups",
			ReadOnlyHint:   true,
			IdempotentHint: true,
			OpenWorldHint:  boolPtr(true),
		},
		Description: "Aggregate traces into groups using correlated sub-queries. " +
			"Use this to find trace patterns and get aggregated trace metrics. " +
			"Requires parallel arrays: query[], alias[], system[] with matching lengths. " +
			"One alias must be 'root' to identify the root span query. " +
			"Additional sub-queries filter traces where child spans match specific criteria. " +
			"Systems: spans:all, httpserver:all, db:postgresql, log:error, etc. " +
			"Returns grouped rows with dynamic columns. " +
			"Use list_traces instead when you need individual trace details. " +
			"Use list_span_groups instead when you don't need cross-span trace correlation. " +
			"Documentation: https://uptrace.dev/features/querying/spans",
	}, t.handler)
}

func (t *ListTraceGroupsTool) handler(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input *uptraceapi.ListTraceGroupsRequestOptions,
) (*mcp.CallToolResult, *uptraceapi.ListTraceGroupsResponse, error) {
	if input.PathParams.ProjectID == 0 {
		input.PathParams.ProjectID = t.conf.Uptrace.ProjectID
	}

	if input.Query == nil {
		input.Query = &uptraceapi.ListTraceGroupsQuery{}
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

	resp, err := t.client.ListTraceGroups(ctx, input)
	if err != nil {
		return nil, nil, err
	}

	return nil, resp, nil
}
