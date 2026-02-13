package tools

import (
	"context"
	"fmt"

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
		Description: "List dashboards configured in Uptrace. " +
			"Use to view available dashboards for monitoring and visualization. " +
			"Documentation: https://uptrace.dev/llms.txt#features > 'Dashboards'",
	}, t.handler)
}

func (t *ListDashboardsTool) handler(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input *uptraceapi.ListDashboardsRequestOptions,
) (*mcp.CallToolResult, any, error) {
	if input.PathParams == nil || input.PathParams.ProjectID == 0 {
		input.PathParams = &uptraceapi.ListDashboardsPath{
			ProjectID: t.conf.Uptrace.ProjectID,
		}
	}
	resp, err := t.client.ListDashboards(ctx, input)
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
