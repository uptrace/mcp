package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/uptrace/mcp/appconf"
	"github.com/uptrace/mcp/uptraceapi"
)

// WARNING: More information about creating dashboards is needed; the current documentation is insufficient.
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
		Description: "Create a new dashboard from YAML definition. " +
			"Supports grid-based and table-based dashboards with metrics queries. " +
			"Use PromQL-style expressions to visualize spans, events, logs, and metrics. " +
			"Full YAML format guide: https://uptrace.dev/raw/features/dashboards.md, is necessary to reference the documentation to use this tool correctly.",
	}, t.handler)
}

func (t *CreateDashboardTool) handler(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input *uptraceapi.CreateDashboardFromYAMLRequestOptions,
) (*mcp.CallToolResult, any, error) {
	if input.PathParams.ProjectID == nil {
		input.PathParams.ProjectID = &t.conf.Uptrace.ProjectID
	}

	if input.Body == nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: "YAML dashboard definition is required"},
			},
			IsError: true,
		}, nil, nil
	}
	if !strings.Contains(*input.Body, "schema:") ||
		!strings.Contains(*input.Body, "name") ||
		!strings.Contains(*input.Body, "version") ||
		!strings.Contains(*input.Body, "tags") {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: "Missing required field 'schema: v2' or 'schema: v3', 'name', 'version', or 'tags' in the YAML definition",
				},
			},
			IsError: true,
		}, nil, nil
	}
	resp, err := t.client.CreateDashboardFromYAML(ctx, input)
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
