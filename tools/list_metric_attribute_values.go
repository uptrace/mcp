package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/uptrace/mcp/appconf"
	"github.com/uptrace/mcp/uptraceapi"
)

type ListMetricAttributeValuesTool struct {
	client *uptraceapi.Client
	conf   *appconf.Config
}

func NewListMetricAttributeValuesTool(client *uptraceapi.Client, conf *appconf.Config) *ListMetricAttributeValuesTool {
	return &ListMetricAttributeValuesTool{
		client: client,
		conf:   conf,
	}
}

func (t *ListMetricAttributeValuesTool) Register(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name: "list_metric_attribute_values",
		Annotations: &mcp.ToolAnnotations{
			Title:          "List attribute values",
			ReadOnlyHint:   true,
			IdempotentHint: true,
			OpenWorldHint:  boolPtr(true),
		},
		Description: uptraceapi.Operations["listMetricAttributeValues"].Description,
	}, t.handler)
}

func (t *ListMetricAttributeValuesTool) handler(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input *uptraceapi.ListMetricAttributeValuesRequestOptions,
) (*mcp.CallToolResult, *uptraceapi.ListMetricAttributeValuesResponseJSON, error) {
	if input.PathParams == nil || input.PathParams.ProjectID == 0 {
		if input.PathParams == nil {
			input.PathParams = &uptraceapi.ListMetricAttributeValuesPath{}
		}
		input.PathParams.ProjectID = t.conf.Uptrace.ProjectID
	}
	resp, err := t.client.ListMetricAttributeValues(ctx, input)
	if err != nil {
		return nil, nil, err
	}
	return nil, resp, nil
}
