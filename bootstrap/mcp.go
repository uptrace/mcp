package bootstrap

import (
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/uptrace/mcp/appconf"
)

func NewServer(conf *appconf.Config) *mcp.Server {
	return mcp.NewServer(
		&mcp.Implementation{
			Name:    AppName,
			Version: AppVersion,
		},
		&mcp.ServerOptions{
			Instructions: "Uptrace is an open-source observability platform for distributed tracing, metrics, and logs. " +
				"For comprehensive documentation optimized for LLMs, see https://uptrace.dev/llms.txt",
		},
	)
}
