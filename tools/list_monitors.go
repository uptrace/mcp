package tools

import (
	"context"
	"fmt"
	"strings"

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

	var sb strings.Builder
	fmt.Fprintf(&sb, "Monitors:\n")
	for _, m := range resp.Monitors {
		fmt.Fprintf(
			&sb,
			"- ID: %d\n  Name: %s\n  Type: %s\n  Status: %s\n  NotifyEveryoneByEmail: %v\n  TeamIds: %v\n  ChannelIds: %v\n",
			m.ID,
			m.Name,
			m.Type,
			m.Status,
			m.NotifyEveryoneByEmail,
			m.TeamIds,
			m.ChannelIds,
		)
		if m.RepeatInterval != nil {
			fmt.Fprintf(&sb, "  RepeatInterval: %+v\n", *m.RepeatInterval)
		}
		if len(m.Params) > 0 {
			fmt.Fprintf(&sb, "  Params: %+v\n", m.Params)
		}
		if m.CreatedAt != nil {
			fmt.Fprintf(&sb, "  CreatedAt: %f\n", *m.CreatedAt)
		}
		if m.UpdatedAt != nil {
			fmt.Fprintf(&sb, "  UpdatedAt: %f\n", *m.UpdatedAt)
		}
		fmt.Fprint(&sb, "\n")
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: sb.String()},
		},
	}, nil, nil
}
