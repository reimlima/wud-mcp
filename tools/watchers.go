package tools

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/reimlima/wud-mcp/client"
)

func RegisterWatchers(s *server.MCPServer, c *client.Client) {
	s.AddTool(
		mcp.NewTool("list_watchers",
			mcp.WithDescription("List all configured watchers"),
		),
		func(ctx context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			data, err := c.ListWatchers(ctx)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return mcp.NewToolResultText(string(data)), nil
		},
	)

	s.AddTool(
		mcp.NewTool("get_watcher",
			mcp.WithDescription("Get details of a specific watcher by type and name"),
			mcp.WithString("type",
				mcp.Required(),
				mcp.Description("Watcher type (e.g. docker, kubernetes)"),
			),
			mcp.WithString("name",
				mcp.Required(),
				mcp.Description("Watcher name"),
			),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			args := req.GetArguments()
			watchType, _ := args["type"].(string)
			name, _ := args["name"].(string)
			if watchType == "" || name == "" {
				return mcp.NewToolResultError("type and name are required"), nil
			}
			data, err := c.GetWatcher(ctx, watchType, name)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return mcp.NewToolResultText(string(data)), nil
		},
	)
}
