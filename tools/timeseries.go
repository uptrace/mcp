package tools

import (
	"context"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/uptrace/mcp/appconf"
	"github.com/uptrace/mcp/uptraceapi"
)

type QueryTimeseriesTool struct {
	client *uptraceapi.Client
	conf   *appconf.Config
}

func NewQueryTimeseriesTool(client *uptraceapi.Client, conf *appconf.Config) *QueryTimeseriesTool {
	return &QueryTimeseriesTool{
		client: client,
		conf:   conf,
	}
}

func (t *QueryTimeseriesTool) Register(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name: "timeseries",
		Annotations: &mcp.ToolAnnotations{
			Title:          "Query timeseries",
			ReadOnlyHint:   true,
			IdempotentHint: true,
			OpenWorldHint:  boolPtr(true),
		},
		Description: uptraceapi.Operations["queryTimeseries"].Description,
	}, t.handler)
}

func (t *QueryTimeseriesTool) handler(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input *uptraceapi.QueryTimeseriesRequestOptions,
) (*mcp.CallToolResult, *uptraceapi.QueryTimeseriesResponse, error) {
	if input.PathParams.ProjectID == 0 {
		input.PathParams.ProjectID = t.conf.Uptrace.ProjectID
	}

	if input.Query == nil {
		input.Query = &uptraceapi.QueryTimeseriesQuery{}
	}
	if input.Query.TimeGte.IsZero() {
		input.Query.TimeGte = time.Now().Add(-time.Hour)
	}
	if input.Query.TimeLt.IsZero() {
		input.Query.TimeLt = time.Now()
	}
	if input.Query.Limit == nil {
		defaultLimit := uptraceapi.Limit(t.conf.Default.Limit)
		input.Query.Limit = &defaultLimit
	}
	if input.Query.Query == nil {
		input.Query.Query = &t.conf.Default.Query
	}

	resp, err := t.client.QueryTimeseries(ctx, input)
	if err != nil {
		return nil, nil, err
	}

	return nil, resp, nil
}
