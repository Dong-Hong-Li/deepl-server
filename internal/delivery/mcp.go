package delivery

import (
	"context"
	"deep-map-server/config"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type MCPToolRegistry interface {
	RegisterTool(server *mcp.Server)
}

func NewServertartMCPServer(config *config.Config, ctx context.Context, toolRegistries []MCPToolRegistry) *mcp.Server {
	server := mcp.NewServer(&mcp.Implementation{Name: "deepl", Version: "0.1.0"}, nil)
	for _, toolRegistry := range toolRegistries {
		toolRegistry.RegisterTool(server)
	}
	return server
}
