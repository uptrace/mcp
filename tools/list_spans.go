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
		Description: uptraceapi.Operations["listSpans"].Description,
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
