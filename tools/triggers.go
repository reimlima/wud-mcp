package tools

import (
	"context"
	"io"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/reimlima/wud-mcp/client"
)

func RegisterTriggers(s *server.MCPServer, c *client.Client) {
	s.AddTool(
		mcp.NewTool("list_triggers",
			mcp.WithDescription("List all configured triggers"),
		),
		func(ctx context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			data, err := c.ListTriggers(ctx)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return mcp.NewToolResultText(string(data)), nil
		},
	)

	s.AddTool(
		mcp.NewTool("get_trigger",
			mcp.WithDescription("Get details of a specific trigger by type and name"),
			mcp.WithString("type",
				mcp.Required(),
				mcp.Description("Trigger type (e.g. smtp, slack, webhook)"),
			),
			mcp.WithString("name",
				mcp.Required(),
				mcp.Description("Trigger name"),
			),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			args := req.GetArguments()
			trigType, _ := args["type"].(string)
			name, _ := args["name"].(string)
			if trigType == "" || name == "" {
				return mcp.NewToolResultError("type and name are required"), nil
			}
			data, err := c.GetTrigger(ctx, trigType, name)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return mcp.NewToolResultText(string(data)), nil
		},
	)

	s.AddTool(
		mcp.NewTool("run_trigger",
			mcp.WithDescription("Run a trigger with optional simulated container data"),
			mcp.WithString("type",
				mcp.Required(),
				mcp.Description("Trigger type (e.g. smtp, slack, webhook)"),
			),
			mcp.WithString("name",
				mcp.Required(),
				mcp.Description("Trigger name"),
			),
			mcp.WithString("container_json",
				mcp.Description("Optional JSON object representing a simulated container"),
			),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			args := req.GetArguments()
			trigType, _ := args["type"].(string)
			name, _ := args["name"].(string)
			if trigType == "" || name == "" {
				return mcp.NewToolResultError("type and name are required"), nil
			}
			containerJSON, _ := args["container_json"].(string)
			var body io.Reader
			if containerJSON != "" {
				body = strings.NewReader(containerJSON)
			}
			data, err := c.RunTrigger(ctx, trigType, name, body)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return mcp.NewToolResultText(string(data)), nil
		},
	)
}
