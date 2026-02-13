package tools

import (
	"context"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/uptrace/mcp/appconf"
	"github.com/uptrace/mcp/uptraceapi"
)

type ListSpansTool struct {
	client *uptraceapi.Client
	conf   *appconf.Config
}

func NewListSpansTool(client *uptraceapi.Client, conf *appconf.Config) *ListSpansTool {
	return &ListSpansTool{
		client: client,
		conf:   conf,
	}
}

func (t *ListSpansTool) Register(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name: "list_spans",
		Annotations: &mcp.ToolAnnotations{
			Title:          "List spans",
			ReadOnlyHint:   true,
			IdempotentHint: true,
			OpenWorldHint:  boolPtr(true),
		},
		Description: "List individual spans using UQL (Uptrace Query Language). " +
			"Use this to inspect specific span details, search for errors, or browse recent operations. " +
			"Supports WHERE filters (e.g. where service_name = 'myservice', where _status_code = 'error', " +
			"where _dur_ms > 100ms), full-text search (e.g. 'word1|word2 -excluded'), " +
			"system filtering (e.g. httpserver:all, db:postgresql, log:error), " +
			"duration filtering in milliseconds, and sorting by any span field. " +
			"Span fields use underscore prefix: _name, _dur_ms, _status_code, _time, _trace_id, _kind. " +
			"Attributes use dot-to-underscore: service.name becomes service_name. " +
			"Returns individual span objects with attrs, timing, and status. " +
			"Use list_span_groups instead when you need aggregated metrics (count, avg, p99). " +
			"Use list_traces instead when you need to find traces matching multi-span criteria. " +
			"Documentation: https://uptrace.dev/features/querying/spans",
	}, t.handler)
}

func (t *ListSpansTool) handler(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input *uptraceapi.ListSpansRequestOptions,
) (*mcp.CallToolResult, *uptraceapi.ListSpansResponse, error) {
	if input.PathParams.ProjectID == 0 {
		input.PathParams.ProjectID = t.conf.Uptrace.ProjectID
	}

	if input.Query == nil {
		input.Query = &uptraceapi.ListSpansQuery{}
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

	resp, err := t.client.ListSpans(ctx, input)
	if err != nil {
		return nil, nil, err
	}

	return nil, resp, nil
}
