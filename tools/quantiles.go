package tools

import (
	"context"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/uptrace/mcp/appconf"
	"github.com/uptrace/mcp/uptraceapi"
)

type QueryQuantilesTool struct {
	client *uptraceapi.Client
	conf   *appconf.Config
}

func NewQueryQuantilesTool(client *uptraceapi.Client, conf *appconf.Config) *QueryQuantilesTool {
	return &QueryQuantilesTool{
		client: client,
		conf:   conf,
	}
}

func (t *QueryQuantilesTool) Register(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name: "quantiles",
		Annotations: &mcp.ToolAnnotations{
			Title:          "Query quantiles",
			ReadOnlyHint:   true,
			IdempotentHint: true,
			OpenWorldHint:  boolPtr(true),
		},
		Description: uptraceapi.Operations["query_quantiles"].Description,
	}, t.handler)
}

func (t *QueryQuantilesTool) handler(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input *uptraceapi.QueryQuantilesRequestOptions,
) (*mcp.CallToolResult, *uptraceapi.QueryQuantilesResponse, error) {
	if input.PathParams.ProjectID == 0 {
		input.PathParams.ProjectID = t.conf.Uptrace.ProjectID
	}

	if input.Query == nil {
		input.Query = &uptraceapi.QueryQuantilesQuery{}
	}
	if input.Query.TimeStart.IsZero() {
		input.Query.TimeStart = time.Now().Add(-t.conf.Default.TimeDuration)
	}
	if input.Query.TimeEnd.IsZero() {
		input.Query.TimeEnd = time.Now()
	}
	if input.Query.Limit == nil {
		defaultLimit := uptraceapi.Limit(t.conf.Default.Limit)
		input.Query.Limit = &defaultLimit
	}

	resp, err := t.client.QueryQuantiles(ctx, input)
	if err != nil {
		return nil, nil, err
	}

	return nil, resp, nil
}
