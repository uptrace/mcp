package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/uptrace/mcp/appconf"
	"github.com/uptrace/mcp/uptraceapi"
)

func registerListSpansTool(server *mcp.Server, client *uptraceapi.Client, conf *appconf.Config) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_spans",
		Description: "List spans from Uptrace. Use to search and analyze distributed traces.",
	}, makeListSpansHandler(client, conf))
}

func makeListSpansHandler(
	client *uptraceapi.Client,
	conf *appconf.Config,
) mcp.ToolHandlerFor[*uptraceapi.ListSpansRequestOptions, any] {
	return func(ctx context.Context, req *mcp.CallToolRequest, input *uptraceapi.ListSpansRequestOptions) (*mcp.CallToolResult, any, error) {
		return handleListSpans(ctx, client, conf, input)
	}
}

func handleListSpans(
	ctx context.Context,
	client *uptraceapi.Client,
	conf *appconf.Config,
	input *uptraceapi.ListSpansRequestOptions,
) (*mcp.CallToolResult, any, error) {
	if input.PathParams.ProjectID == 0 {
		input.PathParams.ProjectID = conf.Uptrace.ProjectID
	}

	// Дефолтные значения времени
	var timeStart, timeEnd time.Time
	if input.Query != nil {
		timeStart = input.Query.TimeStart
		timeEnd = input.Query.TimeEnd
	}
	if timeStart.IsZero() {
		timeStart = time.Now().Add(-time.Hour)
	}
	if timeEnd.IsZero() {
		timeEnd = time.Now()
	}

	limit := uptraceapi.Limit(100)
	if input.Query != nil && input.Query.Limit != nil {
		limit = *input.Query.Limit
	}

	opts := &uptraceapi.ListSpansRequestOptions{
		PathParams: input.PathParams,
		Query: &uptraceapi.ListSpansQuery{
			TimeStart: timeStart,
			TimeEnd:   timeEnd,
			TraceID:   input.Query.TraceID,
			Limit:     &limit,
		},
	}

	resp, err := client.ListSpans(ctx, opts)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Error listing spans: %v", err)},
			},
			IsError: true,
		}, nil, nil
	}

	data, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Error marshaling response: %v", err)},
			},
			IsError: true,
		}, nil, nil
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(data)},
		},
	}, nil, nil
}
