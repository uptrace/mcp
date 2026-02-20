package tools

import (
	"context"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/uptrace/mcp/appconf"
	"github.com/uptrace/mcp/uptraceapi"
)

type PublicListSpanGroupsTool struct {
	client *uptraceapi.Client
	conf   *appconf.Config
}

func NewPublicListSpanGroupsTool(client *uptraceapi.Client, conf *appconf.Config) *PublicListSpanGroupsTool {
	return &PublicListSpanGroupsTool{
		client: client,
		conf:   conf,
	}
}

func (t *PublicListSpanGroupsTool) Register(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name: "public_list_span_groups",
		Annotations: &mcp.ToolAnnotations{
			Title:          "List span groups (public API)",
			ReadOnlyHint:   true,
			IdempotentHint: true,
			OpenWorldHint:  boolPtr(true),
		},
		Description: uptraceapi.Operations["public_list_span_groups"].Description,
	}, t.handler)
}

func (t *PublicListSpanGroupsTool) handler(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input *uptraceapi.PublicListSpanGroupsRequestOptions,
) (*mcp.CallToolResult, *uptraceapi.PublicListSpanGroupsResponse, error) {
	if input.PathParams.ProjectID == 0 {
		input.PathParams.ProjectID = t.conf.Uptrace.ProjectID
	}

	if input.Query == nil {
		input.Query = &uptraceapi.PublicListSpanGroupsQuery{}
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
	if input.Query.Query == nil {
		input.Query.Query = &t.conf.Default.Query
	}

	resp, err := t.client.PublicListSpanGroups(ctx, input)
	if err != nil {
		return nil, nil, err
	}

	return nil, resp, nil
}
