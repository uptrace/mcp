package tools

import (
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"go.uber.org/fx"
)

type RegisterParams struct {
	fx.In
	Server *mcp.Server
	Tools  []Tool `group:"tools"`
}

var Module = fx.Module("tools",
	fx.Provide(
		fx.Annotate(NewPublicListSpansTool, fx.As(new(Tool)), fx.ResultTags(`group:"tools"`)),
		fx.Annotate(NewListSpansTool, fx.As(new(Tool)), fx.ResultTags(`group:"tools"`)),
		fx.Annotate(NewPublicListSpanGroupsTool, fx.As(new(Tool)), fx.ResultTags(`group:"tools"`)),
		fx.Annotate(NewListSpanGroupsTool, fx.As(new(Tool)), fx.ResultTags(`group:"tools"`)),
		fx.Annotate(NewQueryTimeseriesTool, fx.As(new(Tool)), fx.ResultTags(`group:"tools"`)),
		fx.Annotate(NewQueryQuantilesTool, fx.As(new(Tool)), fx.ResultTags(`group:"tools"`)),
		fx.Annotate(NewListTraceGroupsTool, fx.As(new(Tool)), fx.ResultTags(`group:"tools"`)),
		fx.Annotate(NewListTracesTool, fx.As(new(Tool)), fx.ResultTags(`group:"tools"`)),
		fx.Annotate(NewListMonitorsTool, fx.As(new(Tool)), fx.ResultTags(`group:"tools"`)),
		fx.Annotate(NewExploreMetricsTool, fx.As(new(Tool)), fx.ResultTags(`group:"tools"`)),
		fx.Annotate(NewListMetricAttributesTool, fx.As(new(Tool)), fx.ResultTags(`group:"tools"`)),
		fx.Annotate(NewListMetricAttributeValuesTool, fx.As(new(Tool)), fx.ResultTags(`group:"tools"`)),
		fx.Annotate(NewCreateDashboardTool, fx.As(new(Tool)), fx.ResultTags(`group:"tools"`)),
		fx.Annotate(NewListDashboardsTool, fx.As(new(Tool)), fx.ResultTags(`group:"tools"`)),
		fx.Annotate(NewGetDashboardTool, fx.As(new(Tool)), fx.ResultTags(`group:"tools"`)),
		fx.Annotate(NewListDashboardTagsTool, fx.As(new(Tool)), fx.ResultTags(`group:"tools"`)),
		fx.Annotate(NewListDashboardTemplatesTool, fx.As(new(Tool)), fx.ResultTags(`group:"tools"`)),
		fx.Annotate(NewGetDashboardTemplateTool, fx.As(new(Tool)), fx.ResultTags(`group:"tools"`)),
		fx.Annotate(NewGetDashboardYamlTool, fx.As(new(Tool)), fx.ResultTags(`group:"tools"`)),
		fx.Annotate(NewUpdateDashboardYamlTool, fx.As(new(Tool)), fx.ResultTags(`group:"tools"`)),
		fx.Annotate(NewDeleteDashboardTool, fx.As(new(Tool)), fx.ResultTags(`group:"tools"`)),
	),
	fx.Invoke(Register),
)

type Tool interface {
	Register(server *mcp.Server)
}

func boolPtr(b bool) *bool { return &b }

func Register(p RegisterParams) {
	for _, t := range p.Tools {
		t.Register(p.Server)
	}
}
