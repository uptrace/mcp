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

func registerListSpanGroupsTool(
	server *mcp.Server,
	client *uptraceapi.Client,
	conf *appconf.Config,
) {
	mcp.AddTool(server, &mcp.Tool{
		Name: "list_span_groups",
		Description: `Aggregate spans using UQL (Uptrace Query Language).

Use this tool to group and analyze spans by attributes like host_name, service_name, etc.

UQL query examples:
- "group by host_name" - group spans by hostname
- "group by service_name" - group spans by service
- "where _status_code = \"error\" | group by service_name" - group errors by service
- "group by service_name | having count() > 100" - services with >100 spans
- "where _dur_ms > 1s | group by _name" - slow operations
`,
	}, makeListSpanGroupsHandler(client, conf))
}

func makeListSpanGroupsHandler(
	client *uptraceapi.Client,
	conf *appconf.Config,
) mcp.ToolHandlerFor[*uptraceapi.ListSpanGroupsRequestOptions, any] {
	return func(ctx context.Context, req *mcp.CallToolRequest, input *uptraceapi.ListSpanGroupsRequestOptions) (*mcp.CallToolResult, any, error) {
		return handleListSpanGroups(ctx, client, conf, input)
	}
}

func handleListSpanGroups(
	ctx context.Context,
	client *uptraceapi.Client,
	conf *appconf.Config,
	input *uptraceapi.ListSpanGroupsRequestOptions,
) (*mcp.CallToolResult, any, error) {
	timeStart := input.Query.TimeStart
	if timeStart.IsZero() {
		timeStart = time.Now().Add(-time.Hour)
	}
	timeEnd := input.Query.TimeEnd
	if timeEnd.IsZero() {
		timeEnd = time.Now()
	}

	limit := uptraceapi.Limit(100)
	if input.Query.Limit != nil {
		limit = min(*input.Query.Limit, 10000)
	}

	var query string
	if input.Query.Query != nil {
		query = *input.Query.Query
	}

	var search *string
	if input.Query.Search != nil {
		search = input.Query.Search
	}

	var durationGte *int64
	if input.Query.DurationGte != nil && *input.Query.DurationGte > 0 {
		durationGte = input.Query.DurationGte
	}
	var durationLt *int64
	if input.Query.DurationLt != nil && *input.Query.DurationLt > 0 {
		durationLt = input.Query.DurationLt
	}

	opts := &uptraceapi.ListSpanGroupsRequestOptions{
		PathParams: &uptraceapi.ListSpanGroupsPath{
			ProjectID: conf.Uptrace.ProjectID,
		},
		Query: &uptraceapi.ListSpanGroupsQuery{
			TimeStart:   timeStart,
			TimeEnd:     timeEnd,
			Query:       &query,
			Limit:       &limit,
			Search:      search,
			DurationGte: durationGte,
			DurationLt:  durationLt,
		},
	}

	resp, err := client.ListSpanGroups(ctx, opts)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Error listing span groups: %v", err)},
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
