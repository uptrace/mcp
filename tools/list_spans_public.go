package tools

import (
	"context"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/uptrace/mcp/appconf"
	"github.com/uptrace/mcp/uptraceapi"
)

type PublicListSpansTool struct {
	client *uptraceapi.Client
	conf   *appconf.Config
}

func NewPublicListSpansTool(client *uptraceapi.Client, conf *appconf.Config) *PublicListSpansTool {
	return &PublicListSpansTool{
		client: client,
		conf:   conf,
	}
}

func (t *PublicListSpansTool) Register(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name: "public_list_spans",
		Annotations: &mcp.ToolAnnotations{
			Title:          "List spans (public API)",
			ReadOnlyHint:   true,
			IdempotentHint: true,
			OpenWorldHint:  boolPtr(true),
		},
		Description: "List spans using the stable public API. " +
			"Use this for simple span lookups by trace_id, span ID, or parent_id. " +
			"Best for retrieving known spans when you already have an ID. " +
			"For advanced filtering with UQL queries (WHERE, search, system filtering), " +
			"use list_spans instead. " +
			"Documentation: https://uptrace.dev/features/querying/spans",
	}, t.handler)
}

func (t *PublicListSpansTool) handler(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input *uptraceapi.PublicListSpansRequestOptions,
) (*mcp.CallToolResult, *uptraceapi.PublicListSpansResponse, error) {
	if input.PathParams.ProjectID == 0 {
		input.PathParams.ProjectID = t.conf.Uptrace.ProjectID
	}

	if input.Query == nil {
		input.Query = &uptraceapi.PublicListSpansQuery{}
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

	resp, err := t.client.PublicListSpans(ctx, input)
	if err != nil {
		return nil, nil, err
	}

	return nil, resp, nil
}
