package mcp

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Run starts the MCP server over stdio.
func Run() error {
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "hartool",
		Version: "1.0.0",
	}, nil)

	registerDumpTool(server)
	registerSearchTool(server)

	return server.Run(context.Background(), &mcp.StdioTransport{})
}
