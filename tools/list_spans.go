// tools/list_spans.go
package tools

import (
	"context"
	"fmt"
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
		Description: "List spans from Uptrace for distributed tracing analysis. " +
			"Spans represent individual operations in a trace. " +
			"Use to search traces by trace_id, filter by time range, analyze service performance. " +
			"Documentation: https://uptrace.dev/llms.txt#features > 'Querying Spans and Logs'",
	}, t.handler)
}

func (t *ListSpansTool) handler(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input *uptraceapi.ListSpansRequestOptions,
) (*mcp.CallToolResult, any, error) {
	if input.PathParams.ProjectID == 0 {
		input.PathParams.ProjectID = t.conf.Uptrace.ProjectID
	}

	if input.Query == nil {
		input.Query = &uptraceapi.ListSpansQuery{}
	}
	if input.Query.TimeStart.IsZero() {
		input.Query.TimeStart = time.Now().Add(-time.Hour)
	}
	if input.Query.TimeEnd.IsZero() {
		input.Query.TimeEnd = time.Now()
	}
	if input.Query.Limit == nil {
		defaultLimit := uptraceapi.Limit(t.conf.Default.Limit)
		input.Query.Limit = &defaultLimit
	}

	resp, err := t.client.ListSpans(ctx, input)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Error listing spans: %v", err)},
			},
			IsError: true,
		}, nil, nil
	}

	return nil, resp, nil
}
