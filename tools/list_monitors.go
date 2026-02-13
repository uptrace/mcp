package tools

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/uptrace/mcp/appconf"
	"github.com/uptrace/mcp/uptraceapi"
)

type ListMonitorsTool struct {
	client *uptraceapi.Client
	conf   *appconf.Config
}

func NewListMonitorsTool(client *uptraceapi.Client, conf *appconf.Config) *ListMonitorsTool {
	return &ListMonitorsTool{
		client: client,
		conf:   conf,
	}
}

func (t *ListMonitorsTool) Register(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name: "list_monitors",
		Description: "List monitoring rules and alerts configured in Uptrace. " +
			"Use to view alert configurations, check notification settings, understand monitoring thresholds. " +
			"Documentation: https://uptrace.dev/llms.txt#features > 'Monitoring and Alerts Configuration'",
	}, t.handler)
}

func (t ListMonitorsTool) handler(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input *uptraceapi.ListMonitorsRequestOptions,
) (*mcp.CallToolResult, any, error) {
	if input.PathParams == nil || input.PathParams.ProjectID == 0 {
		input.PathParams = &uptraceapi.ListMonitorsPath{
			ProjectID: t.conf.Uptrace.ProjectID,
		}
	}
	resp, err := t.client.ListMonitors(ctx, input)
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
