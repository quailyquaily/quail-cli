package mcp

import (
	"fmt"
	"log/slog"

	"github.com/quailyquaily/quail-cli/client"
	"github.com/quailyquaily/quail-cli/cmd/common"
	"github.com/quailyquaily/quail-cli/mcp"
	"github.com/quailyquaily/quail-cli/mcp/resources"
	"github.com/spf13/cobra"

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
			ctx := cmd.Context()
			cl := ctx.Value(common.CTX_CLIENT{}).(*client.Client)
			version := ctx.Value(common.CTX_VERSION{}).(string)

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
			if err := mcp.AddTools(ctx, s, cl); err != nil {
				slog.Error("failed to add tools", "error", err)
				return
			}

			// Start the server
			if useSSE {
				sseServer := mcps.NewSSEServer(s, fmt.Sprintf("http://localhost:%d", ssePort))
				slog.Info("🚀 SSE server listening", "port", ssePort, "url", fmt.Sprintf("http://localhost:%d/sse", ssePort))
				if err := sseServer.Start(fmt.Sprintf(":%d", ssePort)); err != nil {
					slog.Error("failed to start SSE server", "error", err)
				}
			} else {
				if err := mcps.ServeStdio(s); err != nil {
					slog.Error("failed to serve stdio", "error", err)
				}
			}
		},
	}

	cmd.Flags().BoolVar(&useSSE, "sse", false, "Use SSE for the server")
	cmd.Flags().IntVar(&ssePort, "port", 8083, "Port to listen on for SSE")

	return cmd
}
