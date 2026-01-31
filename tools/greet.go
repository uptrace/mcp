package tools

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type GreetArgs struct {
	Name string `json:"name" mcp:"Name of the person to greet"`
}

func registerGreetTool(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "greet",
		Description: "Say hello to someone",
	}, handleGreet)
}

func handleGreet(
	ctx context.Context,
	ss *mcp.ServerSession,
	params *mcp.CallToolParamsFor[GreetArgs],
) (*mcp.CallToolResultFor[struct{}], error) {
	greeting := fmt.Sprintf("Hello, %s! Welcome to Uptrace MCP.", params.Arguments.Name)

	return &mcp.CallToolResultFor[struct{}]{
		Content: []mcp.Content{
			&mcp.TextContent{Text: greeting},
		},
	}, nil
}
