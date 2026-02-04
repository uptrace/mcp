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
	}, nil)
	tools.Register(server, client, conf)
	return server
}
