package tools

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/reimlima/wud-mcp/client"
)

func RegisterStore(s *server.MCPServer, c *client.Client) {
	s.AddTool(
		mcp.NewTool("get_store_config",
			mcp.WithDescription("Get WUD store configuration including path and file"),
		),
		func(ctx context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			data, err := c.GetStore(ctx)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return mcp.NewToolResultText(string(data)), nil
		},
	)

	s.AddTool(
		mcp.NewTool("get_log_config",
			mcp.WithDescription("Get WUD logger configuration including log level"),
		),
		func(ctx context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			data, err := c.GetLog(ctx)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return mcp.NewToolResultText(string(data)), nil
		},
	)
}
