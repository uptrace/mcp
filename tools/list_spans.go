package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/doordash-oss/oapi-codegen-dd/v3/pkg/runtime"
	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/uptrace/mcp/appconf"
	"github.com/uptrace/mcp/uptraceapi"
)

type ListSpansArgs struct {
	TimeStart string `json:"time_start" mcp:"Start time (RFC3339 or relative like -1h, -30m)"`
	TimeEnd   string `json:"time_end,omitempty" mcp:"End time (RFC3339 or relative, defaults to now)"`
	TraceID   string `json:"trace_id,omitempty" mcp:"Filter by trace ID"`
	Limit     int    `json:"limit,omitempty" mcp:"Maximum number of spans to return (default 100)"`
}

func registerListSpansTool(server *mcp.Server, client *uptraceapi.Client, conf *appconf.Config) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_spans",
		Description: "List spans from Uptrace. Use to search and analyze distributed traces.",
	}, makeListSpansHandler(client, conf))
}

func makeListSpansHandler(client *uptraceapi.Client, conf *appconf.Config) func(
	ctx context.Context,
	ss *mcp.ServerSession,
	params *mcp.CallToolParamsFor[ListSpansArgs],
) (*mcp.CallToolResultFor[struct{}], error) {
	return func(
		ctx context.Context,
		ss *mcp.ServerSession,
		params *mcp.CallToolParamsFor[ListSpansArgs],
	) (*mcp.CallToolResultFor[struct{}], error) {
		return handleListSpans(ctx, client, conf, params.Arguments)
	}
}

func handleListSpans(
	ctx context.Context,
	client *uptraceapi.Client,
	conf *appconf.Config,
	args ListSpansArgs,
) (*mcp.CallToolResultFor[struct{}], error) {
	timeStart, err := parseTime(args.TimeStart)
	if err != nil {
		return nil, fmt.Errorf("invalid time_start: %w", err)
	}

	timeEnd := time.Now()
	if args.TimeEnd != "" {
		timeEnd, err = parseTime(args.TimeEnd)
		if err != nil {
			return nil, fmt.Errorf("invalid time_end: %w", err)
		}
	}

	limit := uptraceapi.Limit(100)
	if args.Limit > 0 {
		limit = uptraceapi.Limit(args.Limit)
	}

	var traceID *string
	if args.TraceID != "" {
		traceID = &args.TraceID
	}

	opts := &uptraceapi.ListSpansRequestOptions{
		PathParams: &uptraceapi.ListSpansPath{
			ProjectID: conf.Uptrace.ProjectID,
		},
		Query: &uptraceapi.ListSpansQuery{
			TimeStart: uptraceapi.TimeStart{
				TimeStart_OneOf: &uptraceapi.TimeStart_OneOf{
					Either: runtime.NewEitherFromB[float32, time.Time](timeStart),
				},
			},
			TimeEnd: uptraceapi.TimeEnd{
				TimeEnd_OneOf: &uptraceapi.TimeEnd_OneOf{
					Either: runtime.NewEitherFromB[float32, time.Time](timeEnd),
				},
			},
			TraceID: traceID,
			Limit:   &limit,
		},
	}

	resp, err := client.ListSpans(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("list spans: %w", err)
	}

	data, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal response: %w", err)
	}

	return &mcp.CallToolResultFor[struct{}]{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(data)},
		},
	}, nil
}

func parseTime(s string) (time.Time, error) {
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t, nil
	}

	if len(s) > 0 && s[0] == '-' {
		d, err := time.ParseDuration(s)
		if err != nil {
			return time.Time{}, err
		}
		return time.Now().Add(d), nil
	}

	return time.Time{}, fmt.Errorf("invalid time format: %s (use RFC3339 or relative like -1h)", s)
}
