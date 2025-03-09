package mcp

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/quail-ink/quail-cli/client"
	"github.com/quail-ink/quail-cli/cmd/common"
	"github.com/spf13/cobra"

	"github.com/mark3labs/mcp-go/mcp"
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

func handleSearchTool(cl *client.Client) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		query := request.Params.Arguments["q"].(string)
		fmt.Printf("query: %v\n", query)

		results, err := cl.Search(query)
		if err != nil {
			slog.Error("failed to search", "error", err)
			return nil, err
		}
		ret := make([]string, 0)
		posts := results.Data.Items
		for _, post := range posts {
			ret = append(ret, fmt.Sprintf("title: %s\nsummary: %s\npublished at: %s\nurl: %s",
				post.Title, post.Summary, post.PublishedAt.Format(time.RFC3339),
				fmt.Sprintf("https://quaily.com/%s/p/%s", post.List.Slug, post.Slug)))
		}
		result := strings.Join(ret, "\n\n")
		fmt.Printf("result: %v\n", result)

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: result,
				},
			},
		}, nil
	}
}

func anyDatetimeToRFC3339(datetimeStr string) string {
	if datetimeStr == "" {
		return time.Now().Format(time.RFC3339)
	}
	// try to parse the datetimeStr in common formats
	formats := []string{
		time.RFC3339,
		time.RFC3339Nano,
		time.RFC1123,
		time.RFC1123Z,
		"2006-01-02 15:04:05",
		"2006-01-02 15:04",
		"2006-01-02 15",
		"2006-01-02",
	}
	for _, format := range formats {
		tt, err := time.Parse(format, datetimeStr)
		if err == nil {
			return tt.Format(time.RFC3339)
		}
	}
	return time.Now().Format(time.RFC3339)
}

func handlePublishPostTool(cl *client.Client) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		fmt.Printf("request.Params.Arguments: %+v\n", request.Params.Arguments)
		var (
			ok          bool
			title       string
			slug        string
			summary     string
			channelSlug string
			datetimeStr string
			content     string
			tags        string
		)

		if _, ok = request.Params.Arguments["title"]; !ok {
			return nil, fmt.Errorf("title is required")
		} else {
			title = request.Params.Arguments["title"].(string)
		}
		if _, ok = request.Params.Arguments["slug"]; !ok {
			slug = ""
		} else {
			slug = request.Params.Arguments["slug"].(string)
		}
		if _, ok = request.Params.Arguments["summary"]; !ok {
			summary = ""
		} else {
			summary = request.Params.Arguments["summary"].(string)
		}
		if _, ok = request.Params.Arguments["channel"]; !ok {
			return nil, fmt.Errorf("channel is required")
		} else {
			channelSlug = request.Params.Arguments["channel"].(string)
		}
		if _, ok = request.Params.Arguments["datetime"]; ok {
			datetimeStr = request.Params.Arguments["datetime"].(string)
		}
		datetimeStr = anyDatetimeToRFC3339(datetimeStr)

		if _, ok = request.Params.Arguments["content"]; !ok {
			return nil, fmt.Errorf("content is required")
		} else {
			content = request.Params.Arguments["content"].(string)
		}
		payload := map[string]any{
			"slug":               slug,
			"title":              title,
			"summary":            summary,
			"channel":            channelSlug,
			"content":            content,
			"cover_image_url":    "https://quaily.com/images/default-cover.png",
			"datetime":           datetimeStr,
			"first_published_at": datetimeStr,
			"tags":               tags,
		}

		result := ""
		ret, err := cl.CreatePost(channelSlug, payload)
		if err != nil {
			slog.Error("failed to create post", "error", err)
			result = err.Error()
		} else {
			result = fmt.Sprintf(`Post created successfully!
			slug=%s, url=%s
			Please update the slug in markdown file frontmatter`,
				ret.Data.Slug, fmt.Sprintf("https://quaily.com/%s/p/%s", channelSlug, ret.Data.Slug))
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: result,
				},
			},
		}, nil
	}
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

			// Add tools
			s.AddTool(mcp.NewTool("publish_post",
				mcp.WithDescription(`Publish a post to quaily according to the given payload.
				The payload could be found at the current context.
				It could be a markdown file which contains the post content and the frontmatter.
				The tool will try to parse the markdown file and extract the title, summary, channel slug, and content from it.
				The slug could be empty or omitted, in which case the tool will generate a slug for the post. The slug should only contain lowercase letters, numbers, and hyphens.
				The datetime could be empty or omitted, in which case the tool will use an empty string.
				The tags could be empty or omitted, in which case the tool will use an empty string.
				If the post is published successfully, the tool will return the metadata of the post.
				The tool must update the slug, datetime to the file, more specifically, the frontmatter.
				`),
				mcp.WithString("title", mcp.Description("Title of the post"), mcp.Required()),
				mcp.WithString("channel", mcp.Description("Channel slug to publish the post to"), mcp.Required()),
				mcp.WithString("content", mcp.Description("Content of the post"), mcp.Required()),
				mcp.WithString("slug", mcp.Description("Slug of the post"), mcp.DefaultString("")),
				mcp.WithString("summary", mcp.Description("Summary of the post"), mcp.DefaultString("")),
				mcp.WithString("datetime", mcp.Description("Datetime of the post"), mcp.DefaultString("")),
				mcp.WithString("tags", mcp.Description("Tags of the post"), mcp.DefaultString("")),
			), handlePublishPostTool(cl))

			s.AddTool(mcp.NewTool("search",
				mcp.WithDescription("Search quaily.com for a given query"),
				mcp.WithString("q",
					mcp.Description("Query to search quaily.com for a given topic"),
					mcp.Required(),
				),
			), handleSearchTool(cl))

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
