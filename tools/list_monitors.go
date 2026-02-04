package tools

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/uptrace/mcp/appconf"
	"github.com/uptrace/mcp/uptraceapi"
)

func registerListMonitorsTool(server *mcp.Server, client *uptraceapi.Client, conf *appconf.Config) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_monitors",
		Description: "List monitors from Uptrace. Use to view configured alerts and monitoring rules.",
	}, makeListMonitorsHandler(client, conf))
}

func makeListMonitorsHandler(
	client *uptraceapi.Client,
	conf *appconf.Config,
) mcp.ToolHandlerFor[*uptraceapi.ListMonitorsPath, any] {
	return func(ctx context.Context, req *mcp.CallToolRequest, input *uptraceapi.ListMonitorsPath) (*mcp.CallToolResult, any, error) {
		return handleListMonitors(ctx, client, conf, input)
	}
}

func handleListMonitors(
	ctx context.Context,
	client *uptraceapi.Client,
	conf *appconf.Config,
	input *uptraceapi.ListMonitorsPath,
) (*mcp.CallToolResult, any, error) {
	if input.ProjectID == 0 {
		input.ProjectID = conf.Uptrace.ProjectID
	}
	opts := &uptraceapi.ListMonitorsRequestOptions{
		PathParams: input,
	}
	resp, err := client.ListMonitors(ctx, opts)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Error: %v", err)},
			},
			IsError: true,
		}, nil, nil
	}
	return nil, resp, nil
}
