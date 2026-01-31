package tools

import "github.com/modelcontextprotocol/go-sdk/mcp"

// Register registers all MCP tools with the server.
func Register(server *mcp.Server) {
	registerGreetTool(server)
}
