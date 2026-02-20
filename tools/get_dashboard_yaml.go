package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/uptrace/mcp/appconf"
	"github.com/uptrace/mcp/uptraceapi"
)

type GetDashboardYamlTool struct {
	client *uptraceapi.Client
	conf   *appconf.Config
}

func NewGetDashboardYamlTool(client *uptraceapi.Client, conf *appconf.Config) *GetDashboardYamlTool {
	return &GetDashboardYamlTool{
		client: client,
		conf:   conf,
	}
}

func (t *GetDashboardYamlTool) Register(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name: "get_dashboard_yaml",
		Annotations: &mcp.ToolAnnotations{
			Title:          "Get dashboard YAML",
			ReadOnlyHint:   true,
			IdempotentHint: true,
			OpenWorldHint:  boolPtr(true),
		},
		Description: uptraceapi.Operations["get_dashboard_yaml"].Description,
	}, t.handler)
}

type getDashboardYAMLInput struct {
	ProjectID   int64 `json:"project_id,omitempty" jsonschema:"Uptrace project ID."`
	DashboardID int64 `json:"dashboard_id" jsonschema:"Dashboard ID." validate:"required"`
}

func (t *GetDashboardYamlTool) handler(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input *getDashboardYAMLInput,
) (*mcp.CallToolResult, any, error) {
	projectID := input.ProjectID
	if projectID == 0 {
		projectID = t.conf.Uptrace.ProjectID
	}

	opts := &uptraceapi.GetDashboardYamlRequestOptions{
		PathParams: &uptraceapi.GetDashboardYamlPath{
			ProjectID:   projectID,
			DashboardID: input.DashboardID,
		},
	}

	resp, err := t.client.GetDashboardYaml(ctx, opts)
	if err != nil {
		return nil, nil, err
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(*resp)},
		},
	}, nil, nil
}
