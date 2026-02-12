package tools

import (
	"context"
	"fmt"
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
		Description: "List grouped spans to analyze operation patterns and performance bottlenecks. " +
			"Span groups aggregate similar operations together for better analysis. " +
			"Documentation: https://uptrace.dev/llms.txt#features > 'Grouping similar spans and events together'",
	}, t.handler)
}

func (t *ListSpanGroupsTool) handler(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input *uptraceapi.ListSpanGroupsRequestOptions,
) (*mcp.CallToolResult, any, error) {
	if input.PathParams.ProjectID == 0 {
		input.PathParams.ProjectID = t.conf.Uptrace.ProjectID
	}

	if input.Query == nil {
		input.Query = &uptraceapi.ListSpanGroupsQuery{}
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
	if input.Query.Query == nil {
		input.Query.Query = &t.conf.Default.Query
	}

	resp, err := t.client.ListSpanGroups(ctx, input)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Error listing span groups: %v", err)},
			},
			IsError: true,
		}, nil, nil
	}

	return nil, resp, nil
}
