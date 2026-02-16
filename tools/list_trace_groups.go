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
		Description: uptraceapi.Operations["listTraceGroups"].Description,
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
