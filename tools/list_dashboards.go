package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/uptrace/mcp/appconf"
	"github.com/uptrace/mcp/uptraceapi"
)

type ListDashboardsTool struct {
	client *uptraceapi.Client
	conf   *appconf.Config
}

func NewListDashboardsTool(client *uptraceapi.Client, conf *appconf.Config) *ListDashboardsTool {
	return &ListDashboardsTool{
		client: client,
		conf:   conf,
	}
}

func (t *ListDashboardsTool) Register(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name: "list_dashboards",
		Annotations: &mcp.ToolAnnotations{
			Title:          "List dashboards",
			ReadOnlyHint:   true,
			IdempotentHint: true,
			OpenWorldHint:  boolPtr(true),
		},
		Description: uptraceapi.Operations["list_dashboards"].Description,
	}, t.handler)
}

func (t *ListDashboardsTool) handler(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input *uptraceapi.ListDashboardsRequestOptions,
) (*mcp.CallToolResult, *uptraceapi.ListDashboardsResponse, error) {
	if input.PathParams == nil || input.PathParams.ProjectID == 0 {
		input.PathParams = &uptraceapi.ListDashboardsPath{
			ProjectID: t.conf.Uptrace.ProjectID,
		}
	}
	resp, err := t.client.ListDashboards(ctx, input)
	if err != nil {
		return nil, nil, err
	}
	return nil, resp, nil
}
