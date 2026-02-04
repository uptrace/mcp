package tools

import (
	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/uptrace/mcp/appconf"
	"github.com/uptrace/mcp/uptraceapi"
)

// Register registers all MCP tools with the server.
func Register(server *mcp.Server, client *uptraceapi.Client, conf *appconf.Config) {
	registerListSpanGroupsTool(server, client, conf)
	registerListSpansTool(server, client, conf)
	registerListMonitorsTool(server, client, conf)
}
