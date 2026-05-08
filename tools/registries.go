package tools

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/reimlima/wud-mcp/client"
)

func RegisterRegistries(s *server.MCPServer, c *client.Client) {
	s.AddTool(
		mcp.NewTool("list_registries",
			mcp.WithDescription("List all configured container registries"),
		),
		func(ctx context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			data, err := c.ListRegistries(ctx)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return mcp.NewToolResultText(string(data)), nil
		},
	)

	s.AddTool(
		mcp.NewTool("get_registry",
			mcp.WithDescription("Get details of a specific registry by type and name"),
			mcp.WithString("type",
				mcp.Required(),
				mcp.Description("Registry type (e.g. hub, ghcr, ecr)"),
			),
			mcp.WithString("name",
				mcp.Required(),
				mcp.Description("Registry name"),
			),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			args := req.GetArguments()
			regType, _ := args["type"].(string)
			name, _ := args["name"].(string)
			if regType == "" || name == "" {
				return mcp.NewToolResultError("type and name are required"), nil
			}
			data, err := c.GetRegistry(ctx, regType, name)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return mcp.NewToolResultText(string(data)), nil
		},
	)
}
