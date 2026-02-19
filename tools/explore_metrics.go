package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/uptrace/mcp/appconf"
	"github.com/uptrace/mcp/uptraceapi"
)

type ExploreMetricsTool struct {
	client *uptraceapi.Client
	conf   *appconf.Config
}

func NewExploreMetricsTool(client *uptraceapi.Client, conf *appconf.Config) *ExploreMetricsTool {
	return &ExploreMetricsTool{
		client: client,
		conf:   conf,
	}
}

func (t *ExploreMetricsTool) Register(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name: "explore_metrics",
		Annotations: &mcp.ToolAnnotations{
			Title:          "Explore metrics",
			ReadOnlyHint:   true,
			IdempotentHint: true,
			OpenWorldHint:  boolPtr(true),
		},
		Description: uptraceapi.Operations["exploreMetrics"].Description,
	}, t.handler)
}

func (t *ExploreMetricsTool) handler(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input *uptraceapi.ExploreMetricsRequestOptions,
) (*mcp.CallToolResult, *uptraceapi.ExploreMetricsResponseJSON, error) {
	if input.PathParams == nil || input.PathParams.ProjectID == 0 {
		input.PathParams = &uptraceapi.ExploreMetricsPath{
			ProjectID: t.conf.Uptrace.ProjectID,
		}
	}
	resp, err := t.client.ExploreMetrics(ctx, input)
	if err != nil {
		return nil, nil, err
	}
	return nil, resp, nil
}
