package tools

import (
	"context"
	"io"
	"net/http"
	"strings"

	"github.com/uptrace/oapi-codegen-dd/v3/pkg/runtime"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/uptrace/mcp/appconf"
	"github.com/uptrace/mcp/uptraceapi"
)

type UpdateDashboardYAMLTool struct {
	client *uptraceapi.Client
	conf   *appconf.Config
}

func NewUpdateDashboardYAMLTool(client *uptraceapi.Client, conf *appconf.Config) *UpdateDashboardYAMLTool {
	return &UpdateDashboardYAMLTool{
		client: client,
		conf:   conf,
	}
}

func (t *UpdateDashboardYAMLTool) Register(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name: "update_dashboard_yaml",
		Annotations: &mcp.ToolAnnotations{
			Title:           "Update dashboard from YAML",
			DestructiveHint: boolPtr(true),
			IdempotentHint:  true,
			OpenWorldHint:   boolPtr(true),
		},
		Description: uptraceapi.Operations["updateDashboardFromYAML"].Description,
	}, t.handler)
}

type updateDashboardYAMLInput struct {
	ProjectID   int64  `json:"project_id,omitempty" jsonschema:"Uptrace project ID."`
	DashboardID int64  `json:"dashboard_id" jsonschema:"Dashboard ID." validate:"required"`
	Body        string `json:"body" jsonschema:"YAML dashboard definition." validate:"required"`
}

func (t *UpdateDashboardYAMLTool) handler(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input *updateDashboardYAMLInput,
) (*mcp.CallToolResult, *uptraceapi.UpdateDashboardFromYAMLResponse, error) {
	projectID := input.ProjectID
	if projectID == 0 {
		projectID = t.conf.Uptrace.ProjectID
	}

	opts := &uptraceapi.UpdateDashboardFromYAMLRequestOptions{
		PathParams: &uptraceapi.UpdateDashboardFromYAMLPath{
			ProjectID:   projectID,
			DashboardID: input.DashboardID,
		},
	}

	yamlBody := input.Body
	setBody := func(_ context.Context, req *http.Request) error {
		req.Body = io.NopCloser(strings.NewReader(yamlBody))
		req.ContentLength = int64(len(yamlBody))
		req.Header.Set("Content-Type", "application/yaml")
		return nil
	}

	resp, err := t.client.UpdateDashboardFromYAML(ctx, opts, runtime.RequestEditorFn(setBody))
	if err != nil {
		return nil, nil, err
	}

	return nil, resp, nil
}
