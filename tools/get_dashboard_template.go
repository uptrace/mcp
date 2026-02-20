package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/uptrace/mcp/appconf"
	"github.com/uptrace/mcp/uptraceapi"
)

type GetDashboardTemplateTool struct {
	client *uptraceapi.Client
	conf   *appconf.Config
}

func NewGetDashboardTemplateTool(client *uptraceapi.Client, conf *appconf.Config) *GetDashboardTemplateTool {
	return &GetDashboardTemplateTool{
		client: client,
		conf:   conf,
	}
}

func (t *GetDashboardTemplateTool) Register(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name: "get_dashboard_template",
		Annotations: &mcp.ToolAnnotations{
			Title:          "Get dashboard template",
			ReadOnlyHint:   true,
			IdempotentHint: true,
			OpenWorldHint:  boolPtr(true),
		},
		Description: uptraceapi.Operations["get_dashboard_template"].Description,
	}, t.handler)
}

func (t *GetDashboardTemplateTool) handler(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input *uptraceapi.GetDashboardTemplateRequestOptions,
) (*mcp.CallToolResult, *uptraceapi.GetDashboardTemplateResponseJSON, error) {
	if input.PathParams == nil {
		input.PathParams = &uptraceapi.GetDashboardTemplatePath{
			ProjectID: t.conf.Uptrace.ProjectID,
		}
	} else if input.PathParams.ProjectID == 0 {
		input.PathParams.ProjectID = t.conf.Uptrace.ProjectID
	}
	resp, err := t.client.GetDashboardTemplate(ctx, input)
	if err != nil {
		return nil, nil, err
	}
	return nil, resp, nil
}
