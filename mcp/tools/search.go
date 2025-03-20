package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	mcps "github.com/mark3labs/mcp-go/server"
	"github.com/quailyquaily/quail-cli/client"
)

func handleSearchTool(cl *client.Client) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		query, ok := request.Params.Arguments["q"].(string)
		if !ok {
			return nil, fmt.Errorf("query is required")
		}

		results, err := cl.Search(query)
		if err != nil {
			slog.Error("failed to search", "error", err)
			return nil, err
		}

		ret := make([]string, 0)
		posts := results.Data.Items
		for _, post := range posts {
			buf, err := json.Marshal(post)
			if err != nil {
				continue
			}
			ret = append(ret, string(buf))
		}
		result := strings.Join(ret, "\n\n")

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

func GetSearchTool(cl *client.Client) (mcp.Tool, mcps.ToolHandlerFunc, error) {
	tool := mcp.NewTool("quaily_search",
		mcp.WithDescription("Search quaily.com for a given query"),
		mcp.WithString("q",
			mcp.Description("Query to search quaily.com for a given topic"),
			mcp.Required(),
		),
	)

	return tool, handleSearchTool(cl), nil
}
