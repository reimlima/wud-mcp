package tools

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/reimlima/wud-mcp/client"
)

func RegisterContainers(s *server.MCPServer, c *client.Client) {
	s.AddTool(
		mcp.NewTool("list_containers",
			mcp.WithDescription("List all containers currently watched by WUD"),
		),
		func(ctx context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			data, err := c.ListContainers(ctx)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return mcp.NewToolResultText(string(data)), nil
		},
	)

	s.AddTool(
		mcp.NewTool("watch_all_containers",
			mcp.WithDescription("Trigger a manual watch on all containers"),
		),
		func(ctx context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			data, err := c.WatchAllContainers(ctx)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return mcp.NewToolResultText(string(data)), nil
		},
	)

	s.AddTool(
		mcp.NewTool("get_container",
			mcp.WithDescription("Get details of a specific container by ID"),
			mcp.WithString("id",
				mcp.Required(),
				mcp.Description("Container ID"),
			),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			id, _ := req.GetArguments()["id"].(string)
			if id == "" {
				return mcp.NewToolResultError("id is required"), nil
			}
			data, err := c.GetContainer(ctx, id)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return mcp.NewToolResultText(string(data)), nil
		},
	)

	s.AddTool(
		mcp.NewTool("get_container_triggers",
			mcp.WithDescription("List all triggers associated with a specific container"),
			mcp.WithString("id",
				mcp.Required(),
				mcp.Description("Container ID"),
			),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			id, _ := req.GetArguments()["id"].(string)
			if id == "" {
				return mcp.NewToolResultError("id is required"), nil
			}
			data, err := c.GetContainerTriggers(ctx, id)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return mcp.NewToolResultText(string(data)), nil
		},
	)

	s.AddTool(
		mcp.NewTool("watch_container",
			mcp.WithDescription("Trigger a manual watch on a specific container"),
			mcp.WithString("id",
				mcp.Required(),
				mcp.Description("Container ID"),
			),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			id, _ := req.GetArguments()["id"].(string)
			if id == "" {
				return mcp.NewToolResultError("id is required"), nil
			}
			data, err := c.WatchContainer(ctx, id)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return mcp.NewToolResultText(string(data)), nil
		},
	)

	s.AddTool(
		mcp.NewTool("run_container_trigger",
			mcp.WithDescription("Manually run a specific trigger on a container"),
			mcp.WithString("id",
				mcp.Required(),
				mcp.Description("Container ID"),
			),
			mcp.WithString("type",
				mcp.Required(),
				mcp.Description("Trigger type (e.g. smtp, slack)"),
			),
			mcp.WithString("name",
				mcp.Required(),
				mcp.Description("Trigger name"),
			),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			args := req.GetArguments()
			id, _ := args["id"].(string)
			trigType, _ := args["type"].(string)
			name, _ := args["name"].(string)
			if id == "" || trigType == "" || name == "" {
				return mcp.NewToolResultError("id, type, and name are required"), nil
			}
			data, err := c.RunContainerTrigger(ctx, id, trigType, name)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return mcp.NewToolResultText(string(data)), nil
		},
	)

	s.AddTool(
		mcp.NewTool("delete_container",
			mcp.WithDescription("Delete a container from WUD by ID"),
			mcp.WithString("id",
				mcp.Required(),
				mcp.Description("Container ID"),
			),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			id, _ := req.GetArguments()["id"].(string)
			if id == "" {
				return mcp.NewToolResultError("id is required"), nil
			}
			data, err := c.DeleteContainer(ctx, id)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return mcp.NewToolResultText(string(data)), nil
		},
	)
}
