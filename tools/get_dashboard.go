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
		Description: uptraceapi.Operations["getDashboard"].Description,
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

	resp, err := t.client.GetDashboard(ctx, input)
	if err != nil {
		return nil, nil, err
	}
	return nil, resp, nil
}
