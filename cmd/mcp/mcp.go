package mcp

import (
	"fmt"
	"log/slog"

	"github.com/quail-ink/quail-cli/client"
	"github.com/quail-ink/quail-cli/cmd/common"
	"github.com/quail-ink/quail-cli/mcp/resources"
	"github.com/quail-ink/quail-cli/mcp/tools"
	"github.com/spf13/cobra"

	"github.com/mark3labs/mcp-go/server"
	mcps "github.com/mark3labs/mcp-go/server"
)

var (
	useSSE  bool
	ssePort int
)

func ServeSSE(mcpServer *mcps.MCPServer, addr string) *mcps.SSEServer {
	return mcps.NewSSEServer(mcpServer, fmt.Sprintf("http://%s", addr))
}

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mcp",
		Short: "Start a MCP server",
		Run: func(cmd *cobra.Command, args []string) {
			cl := cmd.Context().Value(common.CTX_CLIENT{}).(*client.Client)
			version := cmd.Context().Value(common.CTX_VERSION{}).(string)

			s := mcps.NewMCPServer(
				"Quail CLI MCP Server",
				version,
				mcps.WithResourceCapabilities(true, true),
				mcps.WithPromptCapabilities(true),
				mcps.WithLogging(),
			)

			// add resources
			listsRes, listsResHandler, err := resources.GetListsResource(cl)
			if err != nil {
				slog.Error("failed to get lists resource", "error", err)
				return
			}
			s.AddResource(listsRes, listsResHandler)

			// add tools
			listsTool, listsToolHandler, err := tools.GetListsTool(cl)
			if err != nil {
				slog.Error("failed to get lists tool", "error", err)
				return
			}
			s.AddTool(listsTool, listsToolHandler)

			publishPostTool, publishPostToolHandler, err := tools.GetPublishPostTool(cl)
			if err != nil {
				slog.Error("failed to get publish post tool", "error", err)
				return
			}
			s.AddTool(publishPostTool, publishPostToolHandler)

			searchTool, searchToolHandler, err := tools.GetSearchTool(cl)
			if err != nil {
				slog.Error("failed to get search tool", "error", err)
				return
			}
			s.AddTool(searchTool, searchToolHandler)

			getListPostsTool, getListPostsToolHandler, err := tools.GetListPostsTool(cl)
			if err != nil {
				slog.Error("failed to get list posts tool", "error", err)
				return
			}
			s.AddTool(getListPostsTool, getListPostsToolHandler)

			getURLTool, getURLToolHandler, err := tools.GetURLTool(cl)
			if err != nil {
				slog.Error("failed to get url tool", "error", err)
				return
			}
			s.AddTool(getURLTool, getURLToolHandler)

			// Start the server
			if useSSE {
				sseServer := mcps.NewSSEServer(s, fmt.Sprintf("http://localhost:%d", ssePort))
				slog.Info("🚀 SSE server listening", "port", ssePort, "url", fmt.Sprintf("http://localhost:%d/sse", ssePort))
				if err := sseServer.Start(fmt.Sprintf(":%d", ssePort)); err != nil {
					slog.Error("failed to start SSE server", "error", err)
				}
			} else {
				if err := server.ServeStdio(s); err != nil {
					slog.Error("failed to serve stdio", "error", err)
				}
			}
		},
	}

	cmd.Flags().BoolVar(&useSSE, "sse", false, "Use SSE for the server")
	cmd.Flags().IntVar(&ssePort, "port", 8083, "Port to listen on for SSE")

	return cmd
}
