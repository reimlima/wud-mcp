package tools

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/reimlima/wud-mcp/client"
)

func RegisterApp(s *server.MCPServer, c *client.Client) {
	s.AddTool(
		mcp.NewTool("get_app_info",
			mcp.WithDescription("Get WUD application info including name and version"),
		),
		func(ctx context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			data, err := c.GetApp(ctx)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return mcp.NewToolResultText(string(data)), nil
		},
	)
}
