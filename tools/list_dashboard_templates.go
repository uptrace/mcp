package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/uptrace/mcp/appconf"
	"github.com/uptrace/mcp/uptraceapi"
)

type ListDashboardTemplatesTool struct {
	client *uptraceapi.Client
	conf   *appconf.Config
}

func NewListDashboardTemplatesTool(client *uptraceapi.Client, conf *appconf.Config) *ListDashboardTemplatesTool {
	return &ListDashboardTemplatesTool{
		client: client,
		conf:   conf,
	}
}

func (t *ListDashboardTemplatesTool) Register(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name: "list_dashboard_templates",
		Annotations: &mcp.ToolAnnotations{
			Title:          "List dashboard templates",
			ReadOnlyHint:   true,
			IdempotentHint: true,
			OpenWorldHint:  boolPtr(true),
		},
		Description: uptraceapi.Operations["list_dashboard_templates"].Description,
	}, t.handler)
}

func (t *ListDashboardTemplatesTool) handler(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input *uptraceapi.ListDashboardTemplatesRequestOptions,
) (*mcp.CallToolResult, *uptraceapi.ListDashboardTemplatesResponseJSON, error) {
	if input.PathParams == nil || input.PathParams.ProjectID == 0 {
		input.PathParams = &uptraceapi.ListDashboardTemplatesPath{
			ProjectID: t.conf.Uptrace.ProjectID,
		}
	}
	resp, err := t.client.ListDashboardTemplates(ctx, input)
	if err != nil {
		return nil, nil, err
	}
	return nil, resp, nil
}
