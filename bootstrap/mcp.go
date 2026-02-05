package bootstrap

import (
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/uptrace/mcp/appconf"
	"github.com/uptrace/mcp/tools"
	"github.com/uptrace/mcp/uptraceapi"
)

func newServer(client *uptraceapi.Client, conf *appconf.Config) *mcp.Server {
	server := mcp.NewServer(&mcp.Implementation{
		Name:    AppName,
		Version: AppVersion,
	}, &mcp.ServerOptions{
		Instructions: "Uptrace is an open-source observability platform for distributed tracing, metrics, and logs. " +
			"For comprehensive documentation optimized for LLMs, see https://uptrace.dev/llms.txt",
	})
	tools.Register(server, client, conf)
	return server
}
