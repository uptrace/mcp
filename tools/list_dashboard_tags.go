package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/uptrace/mcp/appconf"
	"github.com/uptrace/mcp/uptraceapi"
)

type ListDashboardTagsTool struct {
	client *uptraceapi.Client
	conf   *appconf.Config
}

func NewListDashboardTagsTool(client *uptraceapi.Client, conf *appconf.Config) *ListDashboardTagsTool {
	return &ListDashboardTagsTool{
		client: client,
		conf:   conf,
	}
}

func (t *ListDashboardTagsTool) Register(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name: "list_dashboard_tags",
		Annotations: &mcp.ToolAnnotations{
			Title:          "List dashboard tags",
			ReadOnlyHint:   true,
			IdempotentHint: true,
			OpenWorldHint:  boolPtr(true),
		},
		Description: uptraceapi.Operations["list_dashboard_tags"].Description,
	}, t.handler)
}

func (t *ListDashboardTagsTool) handler(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input *uptraceapi.ListDashboardTagsRequestOptions,
) (*mcp.CallToolResult, *uptraceapi.ListDashboardTagsResponse, error) {
	if input.PathParams == nil || input.PathParams.ProjectID == 0 {
		input.PathParams = &uptraceapi.ListDashboardTagsPath{
			ProjectID: t.conf.Uptrace.ProjectID,
		}
	}
	resp, err := t.client.ListDashboardTags(ctx, input)
	if err != nil {
		return nil, nil, err
	}
	return nil, resp, nil
}
