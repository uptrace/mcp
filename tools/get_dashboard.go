package tools

import (
	"context"

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
		Name: "get_dashboard",
		Annotations: &mcp.ToolAnnotations{
			Title:          "Get dashboard",
			ReadOnlyHint:   true,
			IdempotentHint: true,
			OpenWorldHint:  boolPtr(true),
		},
		Description: "Get a dashboard by ID from Uptrace. " +
			"Use this to retrieve full dashboard details including grid rows, items, and metric queries. " +
			"Requires a dashboard_id — use list_dashboards first to find available dashboard IDs. " +
			"Documentation: https://uptrace.dev/llms.txt#features > 'Dashboards'",
	}, t.handler)
}

func (t *GetDashboardTool) handler(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input *uptraceapi.GetDashboardRequestOptions,
) (*mcp.CallToolResult, *uptraceapi.GetDashboardResponse, error) {
	if input.PathParams == nil || input.PathParams.ProjectID == 0 {
		if input.PathParams == nil {
			input.PathParams = &uptraceapi.GetDashboardPath{}
		}
		input.PathParams.ProjectID = t.conf.Uptrace.ProjectID
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
		return nil, nil, err
	}
	return nil, resp, nil
}
