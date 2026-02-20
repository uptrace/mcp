package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/uptrace/mcp/appconf"
	"github.com/uptrace/mcp/uptraceapi"
)

type ListMonitorsTool struct {
	client *uptraceapi.Client
	conf   *appconf.Config
}

func NewListMonitorsTool(client *uptraceapi.Client, conf *appconf.Config) *ListMonitorsTool {
	return &ListMonitorsTool{
		client: client,
		conf:   conf,
	}
}

func (t *ListMonitorsTool) Register(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name: "list_monitors",
		Annotations: &mcp.ToolAnnotations{
			Title:          "List monitors",
			ReadOnlyHint:   true,
			IdempotentHint: true,
			OpenWorldHint:  boolPtr(true),
		},
		Description: uptraceapi.Operations["list_monitors"].Description,
	}, t.handler)
}

func (t ListMonitorsTool) handler(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input *uptraceapi.ListMonitorsRequestOptions,
) (*mcp.CallToolResult, *uptraceapi.ListMonitorsResponse, error) {
	if input.PathParams == nil || input.PathParams.ProjectID == 0 {
		input.PathParams = &uptraceapi.ListMonitorsPath{
			ProjectID: t.conf.Uptrace.ProjectID,
		}
	}
	resp, err := t.client.ListMonitors(ctx, input)
	if err != nil {
		return nil, nil, err
	}
	return nil, resp, nil
}
