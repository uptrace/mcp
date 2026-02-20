package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/uptrace/mcp/appconf"
	"github.com/uptrace/mcp/uptraceapi"
)

type ListMetricAttributesTool struct {
	client *uptraceapi.Client
	conf   *appconf.Config
}

func NewListMetricAttributesTool(client *uptraceapi.Client, conf *appconf.Config) *ListMetricAttributesTool {
	return &ListMetricAttributesTool{
		client: client,
		conf:   conf,
	}
}

func (t *ListMetricAttributesTool) Register(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name: "list_metric_attributes",
		Annotations: &mcp.ToolAnnotations{
			Title:          "List metric attributes",
			ReadOnlyHint:   true,
			IdempotentHint: true,
			OpenWorldHint:  boolPtr(true),
		},
		Description: uptraceapi.Operations["list_metric_attributes"].Description,
	}, t.handler)
}

func (t *ListMetricAttributesTool) handler(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input *uptraceapi.ListMetricAttributesRequestOptions,
) (*mcp.CallToolResult, *uptraceapi.ListMetricAttributesResponseJSON, error) {
	if input.PathParams == nil || input.PathParams.ProjectID == 0 {
		input.PathParams = &uptraceapi.ListMetricAttributesPath{
			ProjectID: t.conf.Uptrace.ProjectID,
		}
	}
	resp, err := t.client.ListMetricAttributes(ctx, input)
	if err != nil {
		return nil, nil, err
	}
	return nil, resp, nil
}
