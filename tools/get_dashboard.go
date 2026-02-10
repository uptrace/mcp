package tools

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/uptrace/mcp/appconf"
	"github.com/uptrace/mcp/uptraceapi"
)

type GetDashboardTool struct {
	client *uptraceapi.Client
	conf   *appconf.Config
}

func NewGetDashboardTool(client *uptraceapi.Client, conf *appconf.Config) *GetDashboardTool {
	return &GetDashboardTool{
		client: client,
		conf:   conf,
	}
}

func (t *GetDashboardTool) Register(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name: "list_monitors",
		Description: "List monitoring rules and alerts configured in Uptrace. " +
			"Use to view alert configurations, check notification settings, understand monitoring thresholds. " +
			"Documentation: https://uptrace.dev/llms.txt#features > 'Dashboards'",
	}, t.handler)
}

// NOTE: Essentially useless mcp endpoint, needed for CreateDashboardTool tests.
func (t GetDashboardTool) handler(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input *uptraceapi.GetDashboardRequestOptions,
) (*mcp.CallToolResult, any, error) {
	if input.PathParams.ProjectID == nil {
		input.PathParams.ProjectID = &t.conf.Uptrace.ProjectID
	}

	if input.PathParams.DashboardID == 0 {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: "Error: DashboardID is required"},
			},
			IsError: true,
		}, nil, nil
	}
	resp, err := t.client.GetDashboard(ctx, input)
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
