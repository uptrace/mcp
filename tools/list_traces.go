package tools

import (
	"context"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/uptrace/mcp/appconf"
	"github.com/uptrace/mcp/uptraceapi"
)

type ListTracesTool struct {
	client *uptraceapi.Client
	conf   *appconf.Config
}

func NewListTracesTool(client *uptraceapi.Client, conf *appconf.Config) *ListTracesTool {
	return &ListTracesTool{
		client: client,
		conf:   conf,
	}
}

func (t *ListTracesTool) Register(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name: "list_traces",
		Annotations: &mcp.ToolAnnotations{
			Title:          "List traces",
			ReadOnlyHint:   true,
			IdempotentHint: true,
			OpenWorldHint:  boolPtr(true),
		},
		Description: "List individual traces using correlated sub-queries. " +
			"Use this to find specific traces matching complex multi-span criteria. " +
			"Requires parallel arrays: query[], alias[], system[] with matching lengths. " +
			"One alias must be 'root' to identify the root span query. " +
			"Additional sub-queries filter traces where child spans match specific criteria. " +
			"Systems: spans:all, httpserver:all, db:postgresql, log:error, etc. " +
			"Returns root spans for matching traces sorted by time (DESC by default). " +
			"Use list_trace_groups instead when you need aggregated trace metrics. " +
			"Use list_spans instead when you don't need cross-span trace correlation. " +
			"Documentation: https://uptrace.dev/features/querying/spans",
	}, t.handler)
}

func (t *ListTracesTool) handler(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input *uptraceapi.ListTracesRequestOptions,
) (*mcp.CallToolResult, *uptraceapi.ListTracesResponse, error) {
	if input.PathParams.ProjectID == 0 {
		input.PathParams.ProjectID = t.conf.Uptrace.ProjectID
	}

	if input.Query == nil {
		input.Query = &uptraceapi.ListTracesQuery{}
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

	resp, err := t.client.ListTraces(ctx, input)
	if err != nil {
		return nil, nil, err
	}

	return nil, resp, nil
}
