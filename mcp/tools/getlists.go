package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	mcps "github.com/mark3labs/mcp-go/server"
	"github.com/quail-ink/quail-cli/client"
	"github.com/spf13/viper"
)

func handleListsTool(cl *client.Client) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		userID := viper.GetInt64("app.user.id")

		lists, err := cl.GetUserLists(uint64(userID))
		if err != nil {
			slog.Error("failed to get user lists", "error", err)
			return nil, err
		}
		fmt.Printf("lists: %v\n", lists)
		res := make([]string, 0)
		for _, list := range lists {
			jsonItem := make(map[string]any)
			jsonItem["id"] = list.ID
			jsonItem["slug"] = list.Slug
			jsonItem["title"] = list.Title
			jsonItem["tagline"] = list.Tagline
			jsonItem["description"] = list.Description
			buf, err := json.Marshal(jsonItem)
			if err != nil {
				continue
			}
			res = append(res, string(buf))
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: strings.Join(res, "\n\n"),
				},
			},
		}, nil
	}
}

func GetListsTool(cl *client.Client) (mcp.Tool, mcps.ToolHandlerFunc, error) {
	tool := mcp.NewTool("get_my_channels",
		mcp.WithDescription("Return the channels list of me. The channels list items are json objects."),
	)

	return tool, handleListsTool(cl), nil
}
