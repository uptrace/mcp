package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/uptrace/mcp/appconf"
	"github.com/uptrace/mcp/uptraceapi"
)

type DeleteDashboardTool struct {
	client *uptraceapi.Client
	conf   *appconf.Config
}

func NewDeleteDashboardTool(client *uptraceapi.Client, conf *appconf.Config) *DeleteDashboardTool {
	return &DeleteDashboardTool{
		client: client,
		conf:   conf,
	}
}

func (t *DeleteDashboardTool) Register(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name: "delete_dashboard",
		Annotations: &mcp.ToolAnnotations{
			Title:           "Delete dashboard",
			DestructiveHint: boolPtr(true),
			IdempotentHint:  true,
			OpenWorldHint:   boolPtr(true),
		},
		Description: uptraceapi.Operations["delete_dashboard"].Description,
	}, t.handler)
}

type deleteDashboardInput struct {
	ProjectID   int64 `json:"project_id,omitempty" jsonschema:"Uptrace project ID."`
	DashboardID int64 `json:"dashboard_id" jsonschema:"Dashboard ID." validate:"required"`
}

func (t *DeleteDashboardTool) handler(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input *deleteDashboardInput,
) (*mcp.CallToolResult, *struct{}, error) {
	projectID := input.ProjectID
	if projectID == 0 {
		projectID = t.conf.Uptrace.ProjectID
	}

	opts := &uptraceapi.DeleteDashboardRequestOptions{
		PathParams: &uptraceapi.DeleteDashboardPath{
			ProjectID:   projectID,
			DashboardID: input.DashboardID,
		},
	}

	resp, err := t.client.DeleteDashboard(ctx, opts)
	if err != nil {
		return nil, nil, err
	}

	return nil, resp, nil
}
