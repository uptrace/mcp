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

type CreateDashboardTool struct {
	client *uptraceapi.Client
	conf   *appconf.Config
}

func NewCreateDashboardTool(client *uptraceapi.Client, conf *appconf.Config) *CreateDashboardTool {
	return &CreateDashboardTool{
		client: client,
		conf:   conf,
	}
}

func (t *CreateDashboardTool) Register(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name: "create_dashboard",
		Annotations: &mcp.ToolAnnotations{
			Title:           "Create dashboard",
			DestructiveHint: boolPtr(false),
			IdempotentHint:  false,
			OpenWorldHint:   boolPtr(true),
		},
		Description: uptraceapi.Operations["createDashboardFromYAML"].Description,
	}, t.handler)
}

type createDashboardInput struct {
	ProjectID int64  `json:"project_id,omitempty" jsonschema:"Uptrace project ID."`
	Body      string `json:"body" jsonschema:"YAML dashboard definition." validate:"required"`
}

func (t *CreateDashboardTool) handler(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input *createDashboardInput,
) (*mcp.CallToolResult, *uptraceapi.CreateDashboardFromYAMLResponse, error) {
	projectID := input.ProjectID
	if projectID == 0 {
		projectID = t.conf.Uptrace.ProjectID
	}

	if input.Body == "" {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: "YAML dashboard definition is required. " +
					"Provide a YAML body with at minimum:\n\n" +
					"schema: v2\nversion: 1\ntags:\n  - env: prod\nname: My Dashboard\n" +
					"grid_rows:\n  - title: General\n    items:\n      - title: Request rate\n" +
					"        width: 12\n        height: 28\n        type: chart\n" +
					"        metrics:\n          - metric_name as $var\n" +
					"        query:\n          - sum($var)\n\n" +
					"See https://uptrace.dev/raw/features/dashboards.md for full format guide."},
			},
			IsError: true,
		}, nil, nil
	}

	opts := &uptraceapi.CreateDashboardFromYAMLRequestOptions{
		PathParams: &uptraceapi.CreateDashboardFromYAMLPath{
			ProjectID: projectID,
		},
	}

	yamlBody := input.Body
	setBody := func(_ context.Context, req *http.Request) error {
		req.Body = io.NopCloser(strings.NewReader(yamlBody))
		req.ContentLength = int64(len(yamlBody))
		req.Header.Set("Content-Type", "application/yaml")
		return nil
	}

	resp, err := t.client.CreateDashboardFromYAML(ctx, opts, runtime.RequestEditorFn(setBody))
	if err != nil {
		return nil, nil, err
	}

	return nil, resp, nil
}
