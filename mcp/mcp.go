package mcp

import (
	"context"
	"log/slog"

	mcps "github.com/mark3labs/mcp-go/server"
	"github.com/quailyquaily/quail-cli/client"
	"github.com/quailyquaily/quail-cli/mcp/tools"
)

func AddTools(ctx context.Context, s *mcps.MCPServer, cl *client.Client) error {
	listsTool, listsToolHandler, err := tools.GetListsTool(cl)
	if err != nil {
		slog.Error("failed to get lists tool", "error", err)
		return err
	}
	s.AddTool(listsTool, listsToolHandler)

	publishPostTool, publishPostToolHandler, err := tools.GetPublishPostTool(cl)
	if err != nil {
		slog.Error("failed to get publish post tool", "error", err)
		return err
	}
	s.AddTool(publishPostTool, publishPostToolHandler)

	searchTool, searchToolHandler, err := tools.GetSearchTool(cl)
	if err != nil {
		slog.Error("failed to get search tool", "error", err)
		return err
	}
	s.AddTool(searchTool, searchToolHandler)

	getListPostsTool, getListPostsToolHandler, err := tools.GetListPostsTool(cl)
	if err != nil {
		slog.Error("failed to get list posts tool", "error", err)
		return err
	}
	s.AddTool(getListPostsTool, getListPostsToolHandler)

	getURLTool, getURLToolHandler, err := tools.GetURLTool(cl)
	if err != nil {
		slog.Error("failed to get url tool", "error", err)
		return err
	}
	s.AddTool(getURLTool, getURLToolHandler)

	loginTool, loginToolHandler, err := tools.LoginTool(ctx, cl)
	if err != nil {
		slog.Error("failed to get login tool", "error", err)
		return err
	}
	s.AddTool(loginTool, loginToolHandler)

	getPostTool, getPostToolHandler, err := tools.GetPostTool(cl)
	if err != nil {
		slog.Error("failed to get get post tool", "error", err)
		return err
	}
	s.AddTool(getPostTool, getPostToolHandler)

	return nil
}
