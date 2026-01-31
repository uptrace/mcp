package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/uptrace/mcp/appconf"
	"github.com/uptrace/mcp/uptraceapi"
)

type ListMonitorsArgs struct{}

func registerListMonitorsTool(server *mcp.Server, client *uptraceapi.Client, conf *appconf.Config) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_monitors",
		Description: "List monitors from Uptrace. Use to view configured alerts and monitoring rules.",
	}, makeListMonitorsHandler(client, conf))
}

func makeListMonitorsHandler(client *uptraceapi.Client, conf *appconf.Config) func(
	ctx context.Context,
	ss *mcp.ServerSession,
	params *mcp.CallToolParamsFor[ListMonitorsArgs],
) (*mcp.CallToolResultFor[struct{}], error) {
	return func(
		ctx context.Context,
		ss *mcp.ServerSession,
		params *mcp.CallToolParamsFor[ListMonitorsArgs],
	) (*mcp.CallToolResultFor[struct{}], error) {
		return handleListMonitors(ctx, client, conf)
	}
}

func handleListMonitors(
	ctx context.Context,
	client *uptraceapi.Client,
	conf *appconf.Config,
) (*mcp.CallToolResultFor[struct{}], error) {
	opts := &uptraceapi.ListMonitorsRequestOptions{
		PathParams: &uptraceapi.ListMonitorsPath{
			ProjectID: conf.Uptrace.ProjectID,
		},
	}

	resp, err := client.ListMonitors(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("list monitors: %w", err)
	}

	data, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal response: %w", err)
	}

	return &mcp.CallToolResultFor[struct{}]{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(data)},
		},
	}, nil
}
